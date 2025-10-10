.PHONY: help build test run clean docker-build docker-run docker-stop docker-clean fmt vet lint

# Variables
APP_NAME=signing-service
DOCKER_IMAGE=signing-service:latest
DOCKER_CONTAINER=signing-service

# Default target
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# Go commands
build: ## Build the application binary
	@echo "Building $(APP_NAME)..."
	@go build -o $(APP_NAME) .
	@echo "Build complete: ./$(APP_NAME)"

test: ## Run all tests
	@echo "Running tests..."
	@go test ./... -v

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test ./... -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	@go test ./... -race

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test ./crypto -bench=. -benchmem

run: ## Run the application locally
	@echo "Starting $(APP_NAME)..."
	@go run main.go

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(APP_NAME)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

docker-run: ## Run application in Docker container
	@echo "Starting Docker container..."
	@docker run -d \
		--name $(DOCKER_CONTAINER) \
		-p 8080:8080 \
		$(DOCKER_IMAGE)
	@echo "Container started: $(DOCKER_CONTAINER)"
	@echo "Service available at http://localhost:8080"

docker-stop: ## Stop Docker container
	@echo "Stopping Docker container..."
	@docker stop $(DOCKER_CONTAINER) || true
	@docker rm $(DOCKER_CONTAINER) || true
	@echo "Container stopped"

docker-logs: ## Show Docker container logs
	@docker logs -f $(DOCKER_CONTAINER)

docker-clean: docker-stop ## Remove Docker image and container
	@echo "Removing Docker image..."
	@docker rmi $(DOCKER_IMAGE) || true
	@echo "Docker cleanup complete"

# Docker Compose commands
compose-up: ## Start services with docker-compose
	@echo "Starting services with docker-compose..."
	@docker-compose up -d
	@echo "Services started. Check status with: make compose-logs"

compose-down: ## Stop services with docker-compose
	@echo "Stopping services..."
	@docker-compose down
	@echo "Services stopped"

compose-logs: ## Show docker-compose logs
	@docker-compose logs -f

compose-build: ## Build services with docker-compose
	@echo "Building services..."
	@docker-compose build
	@echo "Build complete"

compose-restart: ## Restart services
	@docker-compose restart

# Combined commands
all: clean fmt vet test build ## Run all checks and build

dev: ## Run in development mode with auto-reload (requires air or similar)
	@echo "Starting development server..."
	@go run main.go

demo: ## Run the demo script (requires service to be running)
	@echo "Running demo..."
	@./run-demo.sh

# Installation
install-deps: ## Install development dependencies
	@echo "Installing dependencies..."
	@go mod download
	@echo "Dependencies installed"

