.PHONY: run build test clean docker-up docker-down migrate migrate-down migrate-status migrate-create help

# Variables
APP_NAME=kubera
CMD_DIR=cmd/server/http
BUILD_DIR=bin
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=kubera
DB_SSLMODE=disable

# Database connection string for goose
DATABASE_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	go run $(CMD_DIR)/main.go

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

test-stress: ## Run stress tests
	@echo "Running stress tests..."
	go test -v ./tests/stress/...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean

docker-up: ## Start Docker services (database)
	docker-compose up -d

docker-down: ## Stop Docker services
	docker-compose down

docker-logs: ## Show Docker logs
	docker-compose logs -f

# Goose migration commands
migrate: ## Run database migrations
	@echo "Running migrations..."
	goose -dir db/migrations postgres "$(DATABASE_URL)" up

migrate-down: ## Rollback the last migration
	@echo "Rolling back last migration..."
	goose -dir db/migrations postgres "$(DATABASE_URL)" down

migrate-down-all: ## Rollback all migrations
	@echo "Rolling back all migrations..."
	goose -dir db/migrations postgres "$(DATABASE_URL)" reset

migrate-status: ## Show migration status
	@echo "Migration status..."
	goose -dir db/migrations postgres "$(DATABASE_URL)" status

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@echo "Creating migration: $(name)..."
	goose -dir db/migrations create $(name) sql

sqlc-generate: ## Generate SQLC code
	@echo "Generating SQLC code..."
	sqlc generate

deps: ## Download dependencies
	go mod download
	go mod tidy

lint: ## Run linter
	@golangci-lint run ./...

fmt: ## Format code
	go fmt ./...

.DEFAULT_GOAL := help
