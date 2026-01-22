package claim

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshuarp/kubera/internal/claim/handler"
	"github.com/joshuarp/kubera/internal/claim/repository"
	"github.com/joshuarp/kubera/internal/claim/service"
	couponrepo "github.com/joshuarp/kubera/internal/coupon/repository"
	sqlc "github.com/joshuarp/kubera/internal/pkg/sql"
)

// Module represents the claim domain module with all its dependencies
type Module struct {
	Handler *handler.Handler
	Service service.Service
	Repo    repository.Repository
	Queries *sqlc.Queries
	Pool    *pgxpool.Pool
}

// NewModule initializes the claim domain module with all dependencies
func NewModule(pool *pgxpool.Pool, queries *sqlc.Queries, couponRepo couponrepo.Repository) *Module {
	// Initialize claim repository
	claimRepo := repository.New(pool, queries)

	// Initialize service (depends on coupon repo)
	svc := service.New(claimRepo, couponRepo, pool)

	// Initialize handler
	h := handler.New(svc)

	return &Module{
		Handler: h,
		Service: svc,
		Repo:    claimRepo,
		Queries: queries,
		Pool:    pool,
	}
}
