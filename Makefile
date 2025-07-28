PHONY: build run test clean docker-build docker-up docker-down logs help

# Variables
BINARY_NAME=zl0y-billing
DOCKER_COMPOSE=docker-compose

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $1, $2}' $(MAKEFILE_LIST)

# Build the application
build: ## Build the Go binary
	go build -o $(BINARY_NAME) .

# Run the application locally
run: ## Run the application locally
	go run main.go

# Run tests
test: ## Run tests
	go test -v ./...

# Clean build artifacts
clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	go clean

# Build Docker image
docker-build: ## Build Docker image
	$(DOCKER_COMPOSE) build

# Start all services with Docker Compose
docker-up: ## Start all services
	$(DOCKER_COMPOSE) up --build -d

# Stop all services
docker-down: ## Stop all services
	$(DOCKER_COMPOSE) down

# View logs
logs: ## View application logs
	$(DOCKER_COMPOSE) logs -f billing

# View all service logs
logs-all: ## View all service logs
	$(DOCKER_COMPOSE) logs -f

# Start only databases
db-up: ## Start only database services
	$(DOCKER_COMPOSE) up postgres mongodb -d

# Stop only databases
db-down: ## Stop only database services
	$(DOCKER_COMPOSE) stop postgres mongodb

# Install dependencies
deps: ## Install Go dependencies
	go mod download
	go mod tidy

# Format code
fmt: ## Format Go code
	go fmt ./...

# Run linter
lint: ## Run golangci-lint
	golangci-lint run

# Development setup
dev-setup: deps db-up ## Setup development environment
	@echo "Development environment ready!"
	@echo "Run 'make run' to start the application"

# Full restart
restart: docker-down docker-up ## Restart all services