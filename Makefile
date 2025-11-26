# =============================================================================
# TerraPeak - Unified Makefile
# =============================================================================
# Complete build and deployment automation for TerraPeak project
# - Registry (Go) operations
# - Docker containerization and deployment
# - Development and production workflows
# =============================================================================

.PHONY: help build test test-unit test-integration test-coverage test-api-provider clean fmt lint vet deps run docker-build docker-run
.PHONY: docker-up docker-up-minio docker-down docker-logs docker-restart docker-ps
.PHONY: dev-setup status quick-test watch-test help-registry help-docker
.PHONY: landing-deps landing-dev landing-build landing-docker-build landing-docker-run landing-docker-up landing-docker-down landing-clean

# Default target
help: ## Show this help message
	@echo "TerraPeak - Unified Build System"
	@echo "================================"
	@echo ""
	@echo "Quick Start:"
	@echo "  make deps               Install all dependencies (first time setup)"
	@echo "  make build              Build the TerraPeak binary"
	@echo "  make test               Run all tests"
	@echo ""
	@echo "Registry (Go) Commands:"
	@echo "  make build              Build the TerraPeak binary"
	@echo "  make test               Run all tests"
	@echo "  make run                Run TerraPeak server"
	@echo "  make docker-build       Build registry Docker image"
	@echo ""
	@echo "Landing (Next.js) Commands:"
	@echo "  make landing-deps       Install landing dependencies"
	@echo "  make landing-dev        Run landing in development mode"
	@echo "  make landing-build      Build landing for production"
	@echo "  make landing-docker-build  Build landing Docker image"
	@echo "  make landing-docker-up  Start landing service"
	@echo ""
	@echo "API Testing Commands:"
	@echo "  make test-api           Test core API endpoints"
	@echo "  make test-api-provider  Test Terraform provider discovery & download"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make docker-build       Build Docker image"
	@echo "  make docker-up          Start registry (filesystem storage)"
	@echo "  make docker-up-minio    Start registry with MinIO (S3 storage)"
	@echo "  make docker-down        Stop all services"
	@echo ""
	@echo "For detailed help:"
	@echo "  make help-registry      Show all registry commands"
	@echo "  make help-docker        Show all Docker commands"

# Build targets
build: ## Build the TerraPeak binary
	@echo "Building TerraPeak..."
	cd registry && go build -ldflags="-s -w" -o terrapeak .
	@echo "Build complete: registry/terrapeak"

build-linux: ## Build for Linux (useful for Docker)
	@echo "Building TerraPeak for Linux..."
	cd registry && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o terrapeak-linux .
	@echo "Linux build complete: registry/terrapeak-linux"

# Test targets
test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	cd registry && go test -v -race ./...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	cd registry && go test -v -tags=integration ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	cd registry && go test -v -race -coverprofile=coverage.out ./...
	cd registry && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: registry/coverage.html"

test-benchmark: ## Run benchmark tests
	@echo "Running benchmark tests..."
	cd registry && go test -bench=. -benchmem ./...

# Code quality targets
fmt: ## Format Go code
	@echo "Formatting code..."
	cd registry && go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	cd registry && go vet ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "Running linter..."
	@which golangci-lint > /dev/null || ( \
		echo "golangci-lint not found. Installing..." && \
		GOPATH_BIN=$$(go env GOPATH)/bin && \
		mkdir -p $$GOPATH_BIN && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$GOPATH_BIN v1.54.2 && \
		echo "golangci-lint installed to $$GOPATH_BIN" && \
		echo "Please add $$GOPATH_BIN to your PATH or run: export PATH=\"$$GOPATH_BIN:\$$PATH\"" \
	)
	@cd registry && GOPATH_BIN=$$(go env GOPATH)/bin; \
	if [ -f "$$GOPATH_BIN/golangci-lint" ]; then \
		$$GOPATH_BIN/golangci-lint run --config .golangci.yml ./...; \
	elif which golangci-lint > /dev/null; then \
		golangci-lint run --config .golangci.yml ./...; \
	else \
		echo "golangci-lint not found and installation failed"; \
		exit 1; \
	fi

lint-full: ## Run full golangci-lint with all linters
	@echo "Running full linter..."
	@cd registry && GOPATH_BIN=$$(go env GOPATH)/bin; \
	if [ -f "$$GOPATH_BIN/golangci-lint" ]; then \
		$$GOPATH_BIN/golangci-lint run ./...; \
	elif which golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Please run 'make dev-setup' first"; \
		exit 1; \
	fi

# Dependency management
deps: ## Download and tidy dependencies
	@echo "Managing dependencies..."
	cd registry && go mod download
	cd registry && go mod tidy

deps-update: ## Update all dependencies
	@echo "Updating dependencies..."
	cd registry && go get -u ./...
	cd registry && go mod tidy

# Development targets
run: ## Run TerraPeak with default config
	@echo "Starting TerraPeak..."
	cd registry && ./terrapeak -c .cfg.default.yml

run-dev: build ## Build and run TerraPeak
	@echo "Building and starting TerraPeak..."
	cd registry && ./terrapeak -c .cfg.default.yml

# Docker Compose detection
DOCKER_COMPOSE := $(shell which docker-compose 2>/dev/null)
DOCKER := $(shell which docker 2>/dev/null)

# Determine which docker compose command to use
ifeq ($(DOCKER_COMPOSE),)
    ifeq ($(DOCKER),)
        COMPOSE_CMD = $(error Docker is not installed. Please install Docker first.)
    else
        # Check if 'docker compose' (plugin) is available
        COMPOSE_CHECK := $(shell docker compose version 2>/dev/null)
        ifneq ($(COMPOSE_CHECK),)
            COMPOSE_CMD = docker compose
        else
            COMPOSE_CMD = $(error Neither 'docker compose' nor 'docker-compose' is available. Please install Docker Compose.)
        endif
    endif
else
    COMPOSE_CMD = docker-compose
endif

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t ghcr.io/aliharirian/terrapeak-registry:latest registry/

docker-run: ## Run TerraPeak in Docker container
	@echo "Running TerraPeak in Docker..."
	docker run -p 8081:8081 -v $(PWD)/cfg.yml:/app/cfg.yml:ro ghcr.io/aliharirian/terrapeak-registry:latest

# Docker Compose commands (with auto-detection)
docker-up: docker-down ## Start TerraPeak registry only (filesystem storage)
	@echo "Starting TerraPeak registry (filesystem storage)..."
	@echo "Using: $(COMPOSE_CMD)"
	@echo "Config: registry/.cfg.default.yml (S3 disabled)"
	$(COMPOSE_CMD) up -d terrapeak
	@echo ""
	@echo "TerraPeak started!"
	@echo "Registry: http://localhost:8081"
	@echo "Health check: http://localhost:8081/healthz"

docker-up-minio: docker-down ## Start TerraPeak registry with MinIO (S3 storage)
	@echo "Starting TerraPeak with MinIO (S3 storage)..."
	@echo "Using: $(COMPOSE_CMD)"
	@echo "Config: Will be set to use MinIO"
	@echo ""
	@echo "Updating docker-compose to use MinIO config..."
	@cp registry/.cfg.default.yml registry/.cfg.default.yml.backup 2>/dev/null || true
	@cp registry/.cfg.minio.yml registry/.cfg.default.yml
	@echo "Config updated to enable S3/MinIO"
	@echo ""
	$(COMPOSE_CMD) --profile minio up -d
	@echo ""
	@echo "TerraPeak and MinIO started!"
	@echo "Registry: http://localhost:8081"
	@echo "MinIO Console: http://localhost:9001"
	@echo "   Username: minioadmin"
	@echo "   Password: minioadmin"
	@echo "Health check: http://localhost:8081/healthz"
	@echo ""
	@echo "Note: Config file was updated. Original backed up to .cfg.default.yml.backup"

docker-down: ## Stop all TerraPeak services
	@echo "Stopping all TerraPeak services..."
	@echo "Using: $(COMPOSE_CMD)"
	$(COMPOSE_CMD) --profile minio down
	@echo ""
	@echo "Restoring original config..."
	@if [ -f registry/.cfg.default.yml.backup ]; then \
		mv registry/.cfg.default.yml.backup registry/.cfg.default.yml; \
		echo "Config restored to filesystem storage"; \
	fi
	@echo "All services stopped"

docker-logs: ## Show docker-compose logs
	@echo "Showing docker-compose logs..."
	@echo "Using: $(COMPOSE_CMD)"
	$(COMPOSE_CMD) --profile minio logs -f

docker-restart: ## Restart all services
	@echo "Restarting all TerraPeak services..."
	@echo "Using: $(COMPOSE_CMD)"
	$(COMPOSE_CMD) --profile minio restart

docker-ps: ## Show running containers
	@echo "Running TerraPeak containers:"
	@echo "Using: $(COMPOSE_CMD)"
	$(COMPOSE_CMD) --profile minio ps

# Cleanup targets
clean: ## Clean build artifacts and test files
	@echo "Cleaning up..."
	cd registry && rm -f terrapeak terrapeak-linux
	cd registry && rm -f coverage.out coverage.html
	cd registry && rm -rf ./registry/ # Test storage directory
	@echo "Cleanup complete"

clean-all: clean ## Clean everything including dependencies
	cd registry && go clean -modcache
	docker system prune -f

# Installation targets
install: build ## Install TerraPeak binary to $GOPATH/bin
	@echo "Installing TerraPeak..."
	cd registry && go install .

# Release targets
release-check: test lint vet ## Run all checks for release
	@echo "Running release checks..."
	@echo "All release checks passed"

# CI/CD targets
ci: deps fmt vet lint test-coverage ## Run CI pipeline
	@echo "CI pipeline complete"

# Pre-commit targets
pre-commit: ## Run all checks before commit/push
	@echo "Running pre-commit checks..."
	@echo "=================================="
	@echo ""
	@echo "1. Managing dependencies..."
	@$(MAKE) deps
	@echo ""
	@echo "2. Formatting code..."
	@$(MAKE) fmt
	@echo ""
	@echo "3. Running go vet..."
	@$(MAKE) vet
	@echo ""
	@echo "4. Running unit tests..."
	@$(MAKE) test-unit
	@echo ""
	@echo "5. Building application..."
	@$(MAKE) build
	@echo ""
	@echo "All pre-commit checks passed!"
	@echo "Ready to commit and push!"

pre-commit-quick: ## Quick pre-commit checks (faster)
	@echo "Running quick pre-commit checks..."
	@echo "====================================="
	@echo ""
	@echo "1. Formatting code..."
	@$(MAKE) fmt
	@echo ""
	@echo "2. Running go vet..."
	@$(MAKE) vet
	@echo ""
	@echo "3. Running unit tests..."
	@$(MAKE) test-unit
	@echo ""
	@echo "4. Building application..."
	@$(MAKE) build
	@echo ""
	@echo "Quick pre-commit checks passed!"
	@echo "Ready to commit and push!"

pre-commit-full: ## Full pre-commit checks (comprehensive)
	@echo "Running full pre-commit checks..."
	@echo "===================================="
	@echo ""
	@echo "1. Managing dependencies..."
	@$(MAKE) deps
	@echo ""
	@echo "2. Formatting code..."
	@$(MAKE) fmt
	@echo ""
	@echo "3. Running go vet..."
	@$(MAKE) vet
	@echo ""
	@echo "4. Running unit tests..."
	@$(MAKE) test-unit
	@echo ""
	@echo "5. Running integration tests..."
	@$(MAKE) test-integration
	@echo ""
	@echo "6. Running tests with coverage..."
	@$(MAKE) test-coverage
	@echo ""
	@echo "7. Building application..."
	@$(MAKE) build
	@echo ""
	@echo "8. Testing API endpoints..."
	@$(MAKE) test-api
	@echo ""
	@echo "Full pre-commit checks passed!"
	@echo "Ready to commit and push!"

git-push: ## Run full checks and push
	@echo "Running pre-push checks..."
	@$(MAKE) pre-commit-full
	@echo ""
	@echo "Ready to push!"
	git push origin main
	@echo "Push complete!"

# Quick targets for common workflows
quick-test: fmt vet test-unit ## Quick test cycle (format, vet, unit tests)

# API Testing targets
test-api: ## Test API endpoints on localhost:8081
	@echo "Testing TerraPeak API endpoints..."
	@echo "Testing health endpoint..."
	@curl -s -f "http://localhost:8081/healthz" && echo "Health check passed" || echo "Health check failed"
	@echo ""
	@echo "Testing AWS provider versions..."
	@curl -s "http://localhost:8081/v1/providers/hashicorp/aws/versions" | head -c 200 && echo "... AWS versions endpoint working" || echo "AWS versions failed"
	@echo ""
	@echo "Testing Kubernetes provider versions..."
	@curl -s "http://localhost:8081/v1/providers/hashicorp/kubernetes/versions" | head -c 200 && echo "... Kubernetes versions endpoint working" || echo "Kubernetes versions failed"
	@echo ""
	@echo "Testing proxy info..."
	@curl -s "http://localhost:8081/proxy/info" | head -c 200 && echo "... Proxy info endpoint working" || echo "Proxy info failed"
	@echo ""
	@echo "API testing complete!"

test-api-verbose: ## Test API endpoints with verbose output
	@echo "Testing TerraPeak API endpoints (verbose)..."
	@echo "=============================================="
	@echo ""
	@echo "1. Health Check:"
	@curl -v "http://localhost:8081/healthz"
	@echo ""
	@echo "2. AWS Provider Versions:"
	@curl -v "http://localhost:8081/v1/providers/hashicorp/aws/versions"
	@echo ""
	@echo "3. Kubernetes Provider Versions:"
	@curl -v "http://localhost:8081/v1/providers/hashicorp/kubernetes/versions"
	@echo ""
	@echo "4. Proxy Info:"
	@curl -v "http://localhost:8081/proxy/info"
	@echo ""
	@echo "Verbose API testing complete!"


test-api-provider: ## Test Terraform provider discovery and download endpoints (AWS only)
	@echo "Testing Terraform Provider API Endpoints (AWS Provider)..."
	@echo "=============================================================="
	@echo ""
	@FAILED=0; \
	echo "1. Testing .well-known endpoint (service discovery)..."; \
	echo "   GET http://localhost:8081/.well-known/terraform.json"; \
	HTTP_CODE=$$(curl -s -o /tmp/test-response.json -w "%{http_code}" "http://localhost:8081/.well-known/terraform.json"); \
	if [ "$$HTTP_CODE" = "200" ]; then \
		cat /tmp/test-response.json | python3 -m json.tool; \
		echo "Service discovery working (HTTP $$HTTP_CODE)"; \
	else \
		echo "Service discovery failed (HTTP $$HTTP_CODE)"; \
		cat /tmp/test-response.json 2>/dev/null || true; \
		FAILED=$$((FAILED + 1)); \
	fi; \
	echo ""; \
	echo "2. Testing AWS provider versions endpoint..."; \
	echo "   GET http://localhost:8081/v1/providers/hashicorp/aws/versions"; \
	HTTP_CODE=$$(curl -s -o /tmp/test-response.json -w "%{http_code}" "http://localhost:8081/v1/providers/hashicorp/aws/versions"); \
	if [ "$$HTTP_CODE" = "200" ]; then \
		cat /tmp/test-response.json | python3 -m json.tool | head -20; \
		echo "... AWS versions endpoint working (HTTP $$HTTP_CODE)"; \
	else \
		echo "AWS versions failed (HTTP $$HTTP_CODE)"; \
		cat /tmp/test-response.json 2>/dev/null || true; \
		FAILED=$$((FAILED + 1)); \
	fi; \
	echo ""; \
	echo "3. Testing AWS provider download endpoint (5.100.0 darwin/arm64)..."; \
	echo "   GET http://localhost:8081/v1/providers/hashicorp/aws/5.100.0/download/darwin/arm64"; \
	HTTP_CODE=$$(curl -s -o /tmp/test-response.json -w "%{http_code}" "http://localhost:8081/v1/providers/hashicorp/aws/5.100.0/download/darwin/arm64"); \
	curl -s -I "http://localhost:8081/v1/providers/hashicorp/aws/5.100.0/download/darwin/arm64" | head -5; \
	if [ "$$HTTP_CODE" = "200" ]; then \
		cat /tmp/test-response.json | python3 -m json.tool; \
		echo "AWS download (darwin/arm64) endpoint working (HTTP $$HTTP_CODE)"; \
	else \
		echo "AWS download (darwin/arm64) failed (HTTP $$HTTP_CODE)"; \
		cat /tmp/test-response.json 2>/dev/null || true; \
		FAILED=$$((FAILED + 1)); \
	fi; \
	echo ""; \
	rm -f /tmp/test-response.json; \
	if [ $$FAILED -eq 0 ]; then \
		echo ""; \
		echo "Provider API testing complete! All tests PASSED"; \
		echo "   Tested: .well-known, AWS versions, AWS download (darwin/arm64)"; \
		exit 0; \
	else \
		echo ""; \
		echo "Provider API testing FAILED! $$FAILED test(s) failed"; \
		exit 1; \
	fi

dev-setup: deps ## Setup development environment
	@echo "Setting up development environment..."
	@echo ""
	@echo "Installing development dependencies..."
	@echo ""
	@echo "Installing golangci-lint..."
	@which golangci-lint > /dev/null || ( \
		GOPATH_BIN=$$(go env GOPATH)/bin; \
		mkdir -p $$GOPATH_BIN; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$GOPATH_BIN v1.54.2; \
		echo "Add $$GOPATH_BIN to your PATH if not already there" \
	)
	@echo "golangci-lint installed"
	@echo ""
	@echo "Installing entr (for watch mode)..."
	@which entr > /dev/null || (echo "Please install entr manually:" && echo "  macOS: brew install entr" && echo "  Ubuntu/Debian: sudo apt-get install entr" && echo "  CentOS/RHEL: sudo yum install entr")
	@echo "entr check completed"
	@echo ""
	@echo "Installing curl (for API testing)..."
	@which curl > /dev/null || (echo "Please install curl manually:" && echo "  macOS: brew install curl" && echo "  Ubuntu/Debian: sudo apt-get install curl" && echo "  CentOS/RHEL: sudo yum install curl")
	@echo "curl check completed"
	@echo ""
	@echo "Installing tree (for directory structure)..."
	@which tree > /dev/null || (echo "Please install tree manually:" && echo "  macOS: brew install tree" && echo "  Ubuntu/Debian: sudo apt-get install tree" && echo "  CentOS/RHEL: sudo yum install tree")
	@echo "tree check completed"
	@echo ""
	@echo "Development environment setup complete!"
	@echo "Installed tools:"
	@echo "  golangci-lint (Go linter)"
	@echo "  entr (file watcher)"
	@echo "  curl (API testing)"
	@echo "  tree (directory structure)"
	@echo ""
	@echo "Ready for development!"

# Watch mode (requires entr)
watch-test: ## Watch files and run tests on change (requires 'entr')
	find registry -name "*.go" | entr -c make test-unit


# =============================================================================
# Help Commands
# =============================================================================

help-registry: ## Show registry-specific commands
	@echo "Registry (Go) Commands:"
	@echo "======================="
	@echo "  make build              Build the TerraPeak binary"
	@echo "  make build-linux        Build for Linux"
	@echo "  make test               Run all tests"
	@echo "  make test-unit          Run unit tests only"
	@echo "  make test-integration   Run integration tests"
	@echo "  make test-coverage      Run tests with coverage"
	@echo "  make test-benchmark     Run benchmark tests"
	@echo "  make fmt                Format Go code"
	@echo "  make vet                Run go vet"
	@echo "  make lint               Run golangci-lint"
	@echo "  make deps               Download and tidy dependencies"
	@echo "  make run                Run TerraPeak server"
	@echo "  make docker-build       Build registry Docker image"
	@echo "  make docker-run         Run registry in Docker"


help-docker: ## Show Docker-specific commands
	@echo "Docker Commands:"
	@echo "================"
	@echo ""
	@echo "Build & Run:"
	@echo "  make docker-build         Build registry Docker image"
	@echo "  make docker-run           Run registry in Docker (standalone)"
	@echo ""
	@echo "Docker Compose (auto-detects 'docker compose' or 'docker-compose'):"
	@echo "  make docker-up            Start registry only (filesystem storage)"
	@echo "  make docker-up-minio      Start registry + MinIO (S3 storage)"
	@echo "  make docker-down          Stop all services"
	@echo "  make docker-logs          Show service logs"
	@echo "  make docker-restart       Restart all services"
	@echo "  make docker-ps            Show running containers"
	@echo ""
	@echo "Configuration:"
	@echo "  - docker-up:       Uses .cfg.default.yml (S3 disabled)"
	@echo "  - docker-up-minio: Auto-switches to MinIO config (S3 enabled)"
	@echo ""
	@echo "MinIO Access (when using docker-up-minio):"
	@echo "  - Console: http://localhost:9001"
	@echo "  - API:     http://localhost:9000"
	@echo "  - User:    minioadmin / minioadmin"

# Status check
status: ## Check project status
	@echo "TerraPeak Status"
	@echo "================================"
	@echo "Go version: $(shell go version)"
	@echo "Docker version: $(shell docker --version 2>/dev/null || echo 'Docker not installed')"
	@echo "Git branch: $(shell git branch --show-current 2>/dev/null || echo 'not a git repo')"
	@echo "Git status: $(shell git status --porcelain 2>/dev/null | wc -l | xargs) files changed"
	@echo "Registry dependencies: $(shell cd registry && go list -m all | wc -l | xargs) modules"
	@echo "Registry test files: $(shell find registry -name "*_test.go" | wc -l | xargs) files"
	@echo "Registry source files: $(shell find registry -name "*.go" -not -name "*_test.go" | wc -l | xargs) files"
	@echo ""
	@echo "Services:"
	@echo "  Backend: http://localhost:8081 ($(shell curl -s http://localhost:8081/healthz 2>/dev/null && echo 'running' || echo 'not running'))"
	@echo "  Landing: http://localhost:3000 ($(shell curl -s http://localhost:3000/api/health 2>/dev/null && echo 'running' || echo 'not running'))"

# =============================================================================
# Landing (Next.js) Targets
# =============================================================================

landing-deps: ## Install landing dependencies
	@echo "Installing landing dependencies..."
	@which pnpm > /dev/null || ( \
		echo "pnpm not found. Installing..." && \
		npm install -g pnpm \
	)
	cd landing && pnpm install --frozen-lockfile
	@echo "Landing dependencies installed!"

landing-dev: ## Run landing in development mode
	@echo "Starting landing in development mode..."
	cd landing && pnpm dev

landing-build: ## Build landing for production
	@echo "Building landing for production..."
	cd landing && pnpm build
	@echo "Landing build complete!"

landing-start: ## Start landing in production mode (requires build first)
	@echo "Starting landing in production mode..."
	cd landing && pnpm start

landing-lint: ## Lint landing code
	@echo "Linting landing code..."
	cd landing && pnpm lint

landing-clean: ## Clean landing build artifacts
	@echo "Cleaning landing build artifacts..."
	cd landing && rm -rf .next out node_modules/.cache
	@echo "Landing cleanup complete!"

landing-docker-build: ## Build landing Docker image
	@echo "Building landing Docker image..."
	docker build -t terrapeak-landing:latest \
		--build-arg NODE_VERSION=20 \
		--build-arg PNPM_VERSION=9 \
		landing/
	@echo "Landing Docker image built successfully!"
	@echo "Image: terrapeak-landing:latest"

landing-docker-run: ## Run landing Docker container standalone
	@echo "Running landing in Docker..."
	docker run -p 3000:3000 --name terrapeak-landing terrapeak-landing:latest

landing-docker-up: ## Start landing service with docker-compose
	@echo "Starting landing service..."
	@echo "Using: $(COMPOSE_CMD)"
	$(COMPOSE_CMD) up -d landing
	@echo ""
	@echo "Landing service started!"
	@echo "Landing: http://localhost:3000"
	@echo "Health check: http://localhost:3000/api/health"

landing-docker-down: ## Stop landing service
	@echo "Stopping landing service..."
	@echo "Using: $(COMPOSE_CMD)"
	$(COMPOSE_CMD) stop landing
	@echo "Landing service stopped"

landing-docker-logs: ## Show landing service logs
	@echo "Showing landing service logs..."
	@echo "Using: $(COMPOSE_CMD)"
	$(COMPOSE_CMD) logs -f landing

landing-docker-restart: ## Restart landing service
	@echo "Restarting landing service..."
	@echo "Using: $(COMPOSE_CMD)"
	$(COMPOSE_CMD) restart landing
	@echo "Landing service restarted"

# =============================================================================
# All Services Management
# =============================================================================

all-up: ## Start all services (registry + landing)
	@echo "Starting all TerraPeak services..."
	@echo "Using: $(COMPOSE_CMD)"
	$(COMPOSE_CMD) up -d terrapeak landing
	@echo ""
	@echo "All services started!"
	@echo "Registry: http://localhost:8081"
	@echo "Landing:  http://localhost:3000"
	@echo ""
	@echo "Health checks:"
	@echo "  Registry: http://localhost:8081/healthz"
	@echo "  Landing:  http://localhost:3000/api/health"

all-down: docker-down ## Stop all services (alias for docker-down)

all-logs: ## Show logs for all services
	@echo "Showing logs for all services..."
	@echo "Using: $(COMPOSE_CMD)"
	$(COMPOSE_CMD) logs -f

all-restart: ## Restart all services
	@echo "Restarting all services..."
	@echo "Using: $(COMPOSE_CMD)"
	$(COMPOSE_CMD) restart
	@echo "All services restarted"

all-clean: clean landing-clean ## Clean all build artifacts
	@echo "All build artifacts cleaned!"
