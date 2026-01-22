package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshuarp/kubera/internal/claim/entity"
	claimerrors "github.com/joshuarp/kubera/internal/claim/errors"
	claimrepository "github.com/joshuarp/kubera/internal/claim/repository"
	couponentity "github.com/joshuarp/kubera/internal/coupon/entity"
	couponrepository "github.com/joshuarp/kubera/internal/coupon/repository"
)

// Service handles claim business logic
type Service struct {
	claimRepo  *claimrepository.Repository
	couponRepo *couponrepository.Repository
	pool       *pgxpool.Pool
}

// New creates a new claim service
func New(claimRepo *claimrepository.Repository, couponRepo *couponrepository.Repository, pool *pgxpool.Pool) *Service {
	return &Service{
		claimRepo:  claimRepo,
		couponRepo: couponRepo,
		pool:       pool,
	}
}

// ClaimCoupon handles coupon claiming logic with atomic transaction
// The key operation is:
// 1. Check if user already claimed (via unique constraint)
// 2. Lock coupon row for update (FOR UPDATE)
// 3. Check stock availability
// 4. Insert claim record
// 5. Decrement remaining amount
// All within a database transaction
func (s *Service) ClaimCoupon(ctx context.Context, req *entity.ClaimCouponRequest) (*couponentity.CouponResponse, error) {
	// Start database transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Use named return value for proper defer handling
	var result *couponentity.CouponResponse
	var retErr error

	// Defer rollback in case of error
	defer func() {
		if retErr != nil {
			tx.Rollback(ctx)
		}
	}()

	// Get coupon with row-level lock (FOR UPDATE) within transaction
	coupon, err := s.couponRepo.GetForUpdate(ctx, req.CouponName, tx)
	if err != nil {
		retErr = err
		return result, retErr
	}

	// Check if stock is available
	if coupon.RemainingAmount <= 0 {
		retErr = claimerrors.ErrNoStockAvailable
		return result, retErr
	}

	// Create claim record (will fail if unique constraint violated) within transaction
	claim := &entity.Claim{
		UserID:     req.UserID,
		CouponName: req.CouponName,
	}

	if err := s.claimRepo.Create(ctx, claim, tx); err != nil {
		retErr = err
		return result, retErr
	}

	// Decrement remaining amount within transaction
	newRemainingAmount := coupon.RemainingAmount - 1
	if err := s.couponRepo.UpdateRemainingAmount(ctx, req.CouponName, newRemainingAmount, tx); err != nil {
		retErr = fmt.Errorf("failed to update remaining amount: %w", err)
		return result, retErr
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		retErr = fmt.Errorf("failed to commit transaction: %w", err)
		return result, retErr
	}

	// Get list of users who claimed this coupon (outside transaction for better performance)
	claimedBy, err := s.claimRepo.ListClaimedByCoupon(ctx, req.CouponName)
	if err != nil {
		retErr = fmt.Errorf("failed to get claimed by list: %w", err)
		return result, retErr
	}

	// Return response
	result = &couponentity.CouponResponse{
		Name:            coupon.Name,
		Amount:          coupon.Amount,
		RemainingAmount: newRemainingAmount,
		ClaimedBy:       claimedBy,
	}
	return result, nil
}

func (s *Service) GetClaimedByCoupon(ctx context.Context, couponName string) ([]string, error) {
	return s.claimRepo.ListClaimedByCoupon(ctx, couponName)
}
