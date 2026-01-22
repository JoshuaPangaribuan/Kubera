package service

import (
	"context"

	"github.com/joshuarp/kubera/internal/coupon/entity"
	"github.com/joshuarp/kubera/internal/coupon/repository"
)

// Service defines the coupon service interface
type Service interface {
	CreateCoupon(ctx context.Context, req *entity.CreateCouponRequest) (*entity.Coupon, error)
	GetCoupon(ctx context.Context, name string) (*entity.Coupon, error)
	ListCoupons(ctx context.Context) ([]*entity.Coupon, error)
}

type service struct {
	repo repository.Repository
}

// New creates a new coupon service
func New(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateCoupon(ctx context.Context, req *entity.CreateCouponRequest) (*entity.Coupon, error) {
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

func (s *service) GetCoupon(ctx context.Context, name string) (*entity.Coupon, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *service) ListCoupons(ctx context.Context) ([]*entity.Coupon, error) {
	return s.repo.List(ctx)
}
