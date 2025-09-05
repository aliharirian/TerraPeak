# TerraPeak Makefile
# Build and test automation for TerraPeak Terraform Registry

.PHONY: help build test test-unit test-integration test-coverage clean fmt lint vet deps run docker-build docker-run

# Default target
help: ## Show this help message
	@echo "TerraPeak - Terraform Peak of Features"
	@echo "===================================="
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the TerraPeak binary
	@echo "🔨 Building TerraPeak..."
	cd registry && go build -ldflags="-s -w" -o terrapeak .
	@echo "✅ Build complete: registry/terrapeak"

build-linux: ## Build for Linux (useful for Docker)
	@echo "🔨 Building TerraPeak for Linux..."
	cd registry && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o terrapeak-linux .
	@echo "✅ Linux build complete: registry/terrapeak-linux"

# Test targets
test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests only
	@echo "🧪 Running unit tests..."
	cd registry && go test -v -race ./...

test-integration: ## Run integration tests
	@echo "🧪 Running integration tests..."
	cd registry && go test -v -tags=integration ./...

test-coverage: ## Run tests with coverage report
	@echo "🧪 Running tests with coverage..."
	cd registry && go test -v -race -coverprofile=coverage.out ./...
	cd registry && go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: registry/coverage.html"

test-benchmark: ## Run benchmark tests
	@echo "🏃 Running benchmark tests..."
	cd registry && go test -bench=. -benchmem ./...

# Code quality targets
fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	cd registry && go fmt ./...

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	cd registry && go vet ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "🔍 Running linter..."
	cd registry && golangci-lint run ./...

# Dependency management
deps: ## Download and tidy dependencies
	@echo "📦 Managing dependencies..."
	cd registry && go mod download
	cd registry && go mod tidy

deps-update: ## Update all dependencies
	@echo "📦 Updating dependencies..."
	cd registry && go get -u ./...
	cd registry && go mod tidy

# Development targets
run: ## Run TerraPeak with default config
	@echo "🚀 Starting TerraPeak..."
	cd registry && ./terrapeak -c .cfg.default.yml

run-dev: build ## Build and run TerraPeak
	@echo "🚀 Building and starting TerraPeak..."
	cd registry && ./terrapeak -c .cfg.default.yml

# Docker targets
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t terrapeak:latest .

docker-run: ## Run TerraPeak in Docker container
	@echo "🐳 Running TerraPeak in Docker..."
	docker run -p 8081:8081 -v $(PWD)/cfg.yml:/app/cfg.yml:ro terrapeak:latest

docker-compose-up: ## Start with docker-compose
	@echo "🐳 Starting with docker-compose..."
	docker-compose up -d

docker-compose-down: ## Stop docker-compose services
	@echo "🐳 Stopping docker-compose services..."
	docker-compose down

# Cleanup targets
clean: ## Clean build artifacts and test files
	@echo "🧹 Cleaning up..."
	cd registry && rm -f terrapeak terrapeak-linux
	cd registry && rm -f coverage.out coverage.html
	cd registry && rm -rf ./registry/ # Test storage directory
	@echo "✅ Cleanup complete"

clean-all: clean ## Clean everything including dependencies
	cd registry && go clean -modcache
	docker system prune -f

# Installation targets
install: build ## Install TerraPeak binary to $GOPATH/bin
	@echo "📦 Installing TerraPeak..."
	cd registry && go install .

# Release targets
release-check: test lint vet ## Run all checks for release
	@echo "🔍 Running release checks..."
	@echo "✅ All release checks passed"

# CI/CD targets
ci: deps fmt vet lint test-coverage ## Run CI pipeline
	@echo "🤖 CI pipeline complete"

# Quick targets for common workflows
quick-test: fmt vet test-unit ## Quick test cycle (format, vet, unit tests)

dev-setup: deps ## Setup development environment
	@echo "🔧 Setting up development environment..."
	@echo "Installing golangci-lint..."
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
	@echo "✅ Development environment ready"

# Watch mode (requires entr)
watch-test: ## Watch files and run tests on change (requires 'entr')
	find registry -name "*.go" | entr -c make test-unit

# Status check
status: ## Check project status
	@echo "📊 TerraPeak Status"
	@echo "=================="
	@echo "Go version: $(shell go version)"
	@echo "Git branch: $(shell git branch --show-current 2>/dev/null || echo 'not a git repo')"
	@echo "Git status: $(shell git status --porcelain 2>/dev/null | wc -l | xargs) files changed"
	@echo "Dependencies: $(shell cd registry && go list -m all | wc -l | xargs) modules"
	@echo "Test files: $(shell find registry -name "*_test.go" | wc -l | xargs) files"
	@echo "Source files: $(shell find registry -name "*.go" -not -name "*_test.go" | wc -l | xargs) files"


