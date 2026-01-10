include .env
export

.PHONY: help build run build-be run-be build-fe run-fe run-docker migrate-create migrate-up migrate-down clean-be clean-fe test seed

# Variables
APP_NAME := hris-app
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR_BE := bin
BUILD_DIR_FE := dist
PATH_DB_MIGRATIONS := ./backend/migrations

# Default target
help:
	@echo "Available targets:"
	@echo "  help        - Show this help message"
	@echo "  build       - Build both backend and frontend"
	@echo "  run         - Run both backend and frontend locally"
	@echo "  build-be    - Build the application backend"
	@echo "  run-be      - Run the application backend"
	@echo "  build-fe    - Build the application frontend"
	@echo "  run-fe      - Run the application frontend"
	@echo "  run-docker  - Run the application using container"
	@echo "  migrate-create - Create database migrations"
	@echo "  migrate-up  - Run database migrations"
	@echo "  migrate-down- Rollback database migrations"
	@echo "  clean-be    - Clean build backend artifacts"
	@echo "  clean-fe    - Clean build frontend artifacts"
	@echo "  test        - Run tests"
	@echo "  seed        - Run seeds databases"

# Build both services
build: build-be build-fe
	@echo "Build complete!"

# Run both services locally
run:
	@echo "Starting backend and frontend..."
	@make -j2 run-be run-fe

# Build BE
build-be:
	@echo "Building backend..."
	@cd backend && mkdir -p $(BUILD_DIR_BE) && go build -o $(BUILD_DIR_BE)/$(APP_NAME) ./cmd/api

# Run BE
run-be:
	@echo "Running backend..."
	@cd backend && go run ./cmd/api

# Build FE
build-fe:
	@echo "Building frontend..."
	@cd frontend && pnpm build

# Run FE
run-fe:
	@echo "Running frontend..."
	@cd frontend && pnpm dev 

# Run application both services
run-docker:
	@echo "Running $(APP_NAME)..."
	@docker compose up -d --build --force-recreate

# Run database migrations up
migrate-up:
	@echo "Running database migrations..."
	@migrate -path $(PATH_DB_MIGRATIONS) -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)?multiStatements=true" up

# Rollback database migrations
migrate-down:
	@echo "Rolling back database migrations..."
	@migrate -path $(PATH_DB_MIGRATIONS) -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)?multiStatements=true" down

# Create new migration
migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=migration_name"; exit 1; fi
	@migrate create -ext sql -dir $(PATH_DB_MIGRATIONS) -seq $(NAME)

# Clean build artifacts BE
clean-be:
	@echo "Cleaning build backend artifacts..."
	@cd backend && rm -rf $(BUILD_DIR_BE)

# Clean build artifacts FE
clean-fe:
	@echo "Cleaning build frontend artifacts..."
	@cd frontend && rm -rf $(BUILD_DIR_FE)

# Run tests
test:
	@echo "Running tests..."
	@cd backend && go test -v ./...

# Run seeds
seed:
	@echo "Running seeds..."
	@cd backend && go run ./cmd/api/main.go seed