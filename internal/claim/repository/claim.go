package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshuarp/kubera/internal/claim/entity"
	claimerrors "github.com/joshuarp/kubera/internal/claim/errors"
	apperrors "github.com/joshuarp/kubera/internal/pkg/error"
	sqlc "github.com/joshuarp/kubera/internal/pkg/sql"
)

// Repository handles claim data persistence
type Repository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// New creates a new claim repository
func New(pool *pgxpool.Pool, queries *sqlc.Queries) *Repository {
	return &Repository{
		pool:    pool,
		queries: queries,
	}
}

func (r *Repository) Create(ctx context.Context, claim *entity.Claim, tx ...pgx.Tx) error {
	params := sqlc.CreateClaimParams{
		UserID:     claim.UserID,
		CouponName: claim.CouponName,
	}

	var err error

	// If transaction provided, use it; otherwise use pool
	if len(tx) > 0 && tx[0] != nil {
		queries := r.queries.WithTx(tx[0])
		_, err = queries.CreateClaim(ctx, params)
	} else {
		_, err = r.queries.CreateClaim(ctx, params)
	}

	if err != nil {
		// Check for unique constraint violation (user already claimed)
		if isDuplicateKeyError(err) {
			return claimerrors.ErrAlreadyClaimed
		}
		return apperrors.InternalServer("failed to create claim", err)
	}

	return nil
}

func (r *Repository) GetByUserAndCoupon(ctx context.Context, userID, couponName string) (*entity.Claim, error) {
	params := sqlc.GetClaimByUserAndCouponParams{
		UserID:     userID,
		CouponName: couponName,
	}
	dbClaim, err := r.queries.GetClaimByUserAndCoupon(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, apperrors.InternalServer("failed to get claim", err)
	}

	return dbClaimToEntity(&dbClaim), nil
}

func (r *Repository) ListClaimedByCoupon(ctx context.Context, couponName string) ([]string, error) {
	userIDs, err := r.queries.ListClaimedByCoupon(ctx, couponName)
	if err != nil {
		return nil, apperrors.InternalServer("failed to list claimed by", err)
	}

	return userIDs, nil
}

func (r *Repository) ExistsByUserAndCoupon(ctx context.Context, userID, couponName string) (bool, error) {
	claim, err := r.GetByUserAndCoupon(ctx, userID, couponName)
	if err != nil {
		return false, err
	}
	return claim != nil, nil
}

// Helper functions

func dbClaimToEntity(c *sqlc.Claim) *entity.Claim {
	return &entity.Claim{
		ID:         int64(c.ID),
		UserID:     c.UserID,
		CouponName: c.CouponName,
		ClaimedAt:  c.ClaimedAt.Time,
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
