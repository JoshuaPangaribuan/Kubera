package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshuarp/kubera/internal/coupon/entity"
	couponerrors "github.com/joshuarp/kubera/internal/coupon/errors"
	apperrors "github.com/joshuarp/kubera/internal/pkg/error"
	sqlc "github.com/joshuarp/kubera/internal/pkg/sql"
)

// Repository defines coupon repository interface
type Repository interface {
	Create(ctx context.Context, coupon *entity.Coupon) error
	GetByName(ctx context.Context, name string) (*entity.Coupon, error)
	GetForUpdate(ctx context.Context, name string, tx ...pgx.Tx) (*entity.Coupon, error)
	UpdateRemainingAmount(ctx context.Context, name string, amount int, tx ...pgx.Tx) error
	List(ctx context.Context) ([]*entity.Coupon, error)
}

type repository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// New creates a new coupon repository
func New(pool *pgxpool.Pool, queries *sqlc.Queries) Repository {
	return &repository{
		pool:    pool,
		queries: queries,
	}
}

func (r *repository) Create(ctx context.Context, coupon *entity.Coupon) error {
	params := sqlc.CreateCouponParams{
		Name:   coupon.Name,
		Amount: int32(coupon.Amount),
	}

	dbCoupon, err := r.queries.CreateCoupon(ctx, params)
	if err != nil {
		// Check for unique constraint violation
		if isDuplicateKeyError(err) {
			return couponerrors.ErrCouponAlreadyExists
		}
		return apperrors.InternalServer("failed to create coupon", err)
	}

	// Update entity with generated ID
	coupon.ID = int64(dbCoupon.ID)
	return nil
}

func (r *repository) GetByName(ctx context.Context, name string) (*entity.Coupon, error) {
	dbCoupon, err := r.queries.GetCouponByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, couponerrors.ErrCouponNotFound
		}
		return nil, apperrors.InternalServer("failed to get coupon", err)
	}

	return dbCouponToEntity(&dbCoupon), nil
}

func (r *repository) GetForUpdate(ctx context.Context, name string, tx ...pgx.Tx) (*entity.Coupon, error) {
	var dbCoupon sqlc.Coupon
	var err error

	// If transaction provided, use it; otherwise use pool
	if len(tx) > 0 && tx[0] != nil {
		queries := r.queries.WithTx(tx[0])
		dbCoupon, err = queries.GetCouponForUpdate(ctx, name)
	} else {
		dbCoupon, err = r.queries.GetCouponForUpdate(ctx, name)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, couponerrors.ErrCouponNotFound
		}
		return nil, apperrors.InternalServer("failed to get coupon for update", err)
	}

	return dbCouponToEntity(&dbCoupon), nil
}

func (r *repository) UpdateRemainingAmount(ctx context.Context, name string, amount int, tx ...pgx.Tx) error {
	params := sqlc.UpdateRemainingAmountParams{
		RemainingAmount: int32(amount),
		Name:            name,
	}

	var err error

	// If transaction provided, use it; otherwise use pool
	if len(tx) > 0 && tx[0] != nil {
		queries := r.queries.WithTx(tx[0])
		_, err = queries.UpdateRemainingAmount(ctx, params)
	} else {
		_, err = r.queries.UpdateRemainingAmount(ctx, params)
	}

	if err != nil {
		return apperrors.InternalServer("failed to update remaining amount", err)
	}

	return nil
}

func (r *repository) List(ctx context.Context) ([]*entity.Coupon, error) {
	dbCoupons, err := r.queries.ListAllCoupons(ctx)
	if err != nil {
		return nil, apperrors.InternalServer("failed to list coupons", err)
	}

	coupons := make([]*entity.Coupon, len(dbCoupons))
	for i, c := range dbCoupons {
		coupons[i] = dbCouponToEntity(&c)
	}

	return coupons, nil
}

// Helper functions

func dbCouponToEntity(c *sqlc.Coupon) *entity.Coupon {
	return &entity.Coupon{
		ID:              int64(c.ID),
		Name:            c.Name,
		Amount:          int(c.Amount),
		RemainingAmount: int(c.RemainingAmount),
		CreatedAt:       c.CreatedAt.Time,
		UpdatedAt:       c.UpdatedAt.Time,
	}
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// PostgreSQL error code 23505 = unique_violation
		return pgErr.Code == "23505"
	}
	return false
}
