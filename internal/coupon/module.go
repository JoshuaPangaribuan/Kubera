package coupon

import (
	"github.com/jackc/pgx/v5/pgxpool"
	claimservice "github.com/joshuarp/kubera/internal/claim/service"
	"github.com/joshuarp/kubera/internal/coupon/handler"
	"github.com/joshuarp/kubera/internal/coupon/repository"
	"github.com/joshuarp/kubera/internal/coupon/service"
	sqlc "github.com/joshuarp/kubera/internal/pkg/sql"
)

// Module represents the coupon domain module with all its dependencies
type Module struct {
	Handler *handler.Handler
	Service *service.Service
	Repo    *repository.Repository
	Queries *sqlc.Queries
}

// NewModule initializes the coupon domain module with all dependencies
func NewModule(pool *pgxpool.Pool) *Module {
	// Initialize SQLC queries
	queries := sqlc.New(pool)

	// Initialize repository
	repo := repository.New(pool, queries)

	// Initialize service
	svc := service.New(repo)

	// Initialize handler (without claim service initially)
	h := handler.New(svc, nil)

	return &Module{
		Handler: h,
		Service: svc,
		Repo:    repo,
		Queries: queries,
	}
}

// SetClaimService sets the claim service on the handler after initialization
// This is needed because the claim module depends on the coupon repository,
// creating a circular dependency if we try to inject the claim service during initialization
func (m *Module) SetClaimService(claimSvc *claimservice.Service) {
	m.Handler = handler.New(m.Service, claimSvc)
}
