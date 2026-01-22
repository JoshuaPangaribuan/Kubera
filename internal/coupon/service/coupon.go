package service

import (
	"context"

	"github.com/joshuarp/kubera/internal/coupon/entity"
	couponrepository "github.com/joshuarp/kubera/internal/coupon/repository"
)

// Service handles coupon business logic
type Service struct {
	repo *couponrepository.Repository
}

// New creates a new coupon service
func New(repo *couponrepository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateCoupon(ctx context.Context, req *entity.CreateCouponRequest) (*entity.Coupon, error) {
	// Create coupon entity
	coupon := &entity.Coupon{
		Name:            req.Name,
		Amount:          req.Amount,
		RemainingAmount: req.Amount,
	}

	// Save to repository
	if err := s.repo.Create(ctx, coupon); err != nil {
		return nil, err
	}

	return coupon, nil
}

func (s *Service) GetCoupon(ctx context.Context, name string) (*entity.Coupon, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *Service) ListCoupons(ctx context.Context) ([]*entity.Coupon, error) {
	return s.repo.List(ctx)
}
