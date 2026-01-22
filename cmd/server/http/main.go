package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshuarp/kubera/internal/claim"
	"github.com/joshuarp/kubera/internal/coupon"
	"github.com/joshuarp/kubera/internal/middlewares"
	appconfig "github.com/joshuarp/kubera/internal/pkg/config"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := appconfig.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.GetString("server.mode"))

	// Create pgxpool configuration
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.GetString("database.host"),
		cfg.GetString("database.port"),
		cfg.GetString("database.user"),
		cfg.GetString("database.password"),
		cfg.GetString("database.dbname"),
		cfg.GetString("database.sslmode"),
	)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse database config: %v", err)
	}

	// Configure pool settings
	poolConfig.MaxConns = int32(cfg.GetInt("database.max_open_conns"))
	poolConfig.MinConns = int32(cfg.GetInt("database.max_idle_conns"))
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	// Connect to database
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()

	// Test database connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection established")

	// Initialize modules
	couponModule := coupon.NewModule(pool)
	claimModule := claim.NewModule(pool, couponModule.Queries, couponModule.Repo)

	// Set claim service on coupon handler (to include claimed_by in GET /api/coupons/:name)
	couponModule.SetClaimService(claimModule.Service)

	// Setup router
	router := gin.Default()

	// Apply middlewares
	router.Use(middlewares.Logger())
	router.Use(middlewares.Recovery())
	router.Use(middlewares.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Register API routes
	api := router.Group("/api")
	{
		couponModule.Handler.RegisterRoutes(api)
		claimModule.Handler.RegisterRoutes(api)
	}

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.GetString("server.port"),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.GetInt("server.read_timeout")) * time.Second,
		WriteTimeout: time.Duration(cfg.GetInt("server.write_timeout")) * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", cfg.GetString("server.port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
