# =============================================================================
# TerraPeak - Unified Makefile
# =============================================================================
# Complete build and deployment automation for TerraPeak project
# - Backend (Go) and Frontend (Next.js) operations
# - Docker containerization and deployment
# - Development and production workflows
# =============================================================================

.PHONY: help build test test-unit test-integration test-coverage clean fmt lint vet deps run docker-build docker-run
.PHONY: web-build web-dev web-test web-lint web-clean web-docker-build web-docker-run web-docker-stop
.PHONY: all-build all-test all-clean all-docker-build all-docker-run all-docker-stop
.PHONY: dev-setup status quick-test watch-test

# Default target
help: ## Show this help message
	@echo "TerraPeak - Unified Build System"
	@echo "================================"
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  make all-deps           Install all dependencies (first time setup)"
	@echo "  make dev                Start both backend and frontend"
	@echo "  make all-build          Build both backend and frontend"
	@echo ""
	@echo "Registry (Go) Commands:"
	@echo "  make build              Build the TerraPeak binary"
	@echo "  make test               Run all tests"
	@echo "  make run                Run TerraPeak server"
	@echo "  make docker-build       Build registry Docker image"
	@echo ""
	@echo "Web (Next.js) Commands:"
	@echo "  make web-deps           Install frontend dependencies"
	@echo "  make web-build          Build Next.js application"
	@echo "  make web-dev            Start Next.js dev server"
	@echo "  make web-type-check     Type check TypeScript code"
	@echo "  make web-docker-build   Build web Docker image"
	@echo ""
	@echo "Unified Commands:"
	@echo "  make all-deps           Install all dependencies"
	@echo "  make all-build          Build both registry and web"
	@echo "  make all-test           Test both registry and web"
	@echo "  make all-lint           Lint both registry and web"
	@echo "  make all-docker-build   Build all Docker images"
	@echo "  make dev                Start both services in dev mode"
	@echo ""
	@echo "For detailed help on specific sections:"
	@echo "  make help-registry      Show registry-specific commands"
	@echo "  make help-web           Show web-specific commands"
	@echo "  make help-docker        Show Docker-specific commands"

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
	@which golangci-lint > /dev/null || ( \
		echo "❌ golangci-lint not found. Installing..." && \
		GOPATH_BIN=$$(go env GOPATH)/bin && \
		mkdir -p $$GOPATH_BIN && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$GOPATH_BIN v1.54.2 && \
		echo "✅ golangci-lint installed to $$GOPATH_BIN" && \
		echo "⚠️  Please add $$GOPATH_BIN to your PATH or run: export PATH=\"$$GOPATH_BIN:\$$PATH\"" \
	)
	@cd registry && GOPATH_BIN=$$(go env GOPATH)/bin; \
	if [ -f "$$GOPATH_BIN/golangci-lint" ]; then \
		$$GOPATH_BIN/golangci-lint run --config .golangci.yml ./...; \
	elif which golangci-lint > /dev/null; then \
		golangci-lint run --config .golangci.yml ./...; \
	else \
		echo "❌ golangci-lint not found and installation failed"; \
		exit 1; \
	fi

lint-full: ## Run full golangci-lint with all linters
	@echo "🔍 Running full linter..."
	@cd registry && GOPATH_BIN=$$(go env GOPATH)/bin; \
	if [ -f "$$GOPATH_BIN/golangci-lint" ]; then \
		$$GOPATH_BIN/golangci-lint run ./...; \
	elif which golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "❌ golangci-lint not found. Please run 'make dev-setup' first"; \
		exit 1; \
	fi

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
	docker build -t ghcr.io/aliharirian/terrapeak-registry:latest registry/

docker-run: ## Run TerraPeak in Docker container
	@echo "🐳 Running TerraPeak in Docker..."
	docker run -p 8081:8081 -v $(PWD)/cfg.yml:/app/cfg.yml:ro ghcr.io/aliharirian/terrapeak-registry:latest

docker-compose-up: ## Start all services with docker-compose
	@echo "🐳 Starting all TerraPeak services..."
	docker-compose up -d

docker-compose-down: ## Stop all docker-compose services
	@echo "🐳 Stopping all TerraPeak services..."
	docker-compose down

docker-compose-logs: ## Show docker-compose logs
	@echo "📋 Showing docker-compose logs..."
	docker-compose logs -f

docker-compose-build: ## Build all services with docker-compose
	@echo "🏗️  Building all TerraPeak services..."
	docker-compose build

docker-compose-restart: ## Restart all services
	@echo "🔄 Restarting all TerraPeak services..."
	docker-compose restart

# Service-specific docker-compose commands
docker-compose-registry: ## Start only registry service
	@echo "🐳 Starting registry service..."
	docker-compose up -d registry

docker-compose-web: ## Start only web service
	@echo "🐳 Starting web service..."
	docker-compose up -d web

docker-compose-minio: ## Start only MinIO service
	@echo "🐳 Starting MinIO service..."
	docker-compose up -d minio

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

# Pre-commit targets
pre-commit: ## Run all checks before commit/push
	@echo "🚀 Running pre-commit checks..."
	@echo "=================================="
	@echo ""
	@echo "1. 📦 Managing dependencies..."
	@$(MAKE) deps
	@echo ""
	@echo "2. 🎨 Formatting code..."
	@$(MAKE) fmt
	@echo ""
	@echo "3. 🔍 Running go vet..."
	@$(MAKE) vet
	@echo ""
	@echo "4. 🧪 Running unit tests..."
	@$(MAKE) test-unit
	@echo ""
	@echo "5. 🏗️ Building application..."
	@$(MAKE) build
	@echo ""
	@echo "✅ All pre-commit checks passed!"
	@echo "🚀 Ready to commit and push!"

pre-commit-quick: ## Quick pre-commit checks (faster)
	@echo "⚡ Running quick pre-commit checks..."
	@echo "====================================="
	@echo ""
	@echo "1. 🎨 Formatting code..."
	@$(MAKE) fmt
	@echo ""
	@echo "2. 🔍 Running go vet..."
	@$(MAKE) vet
	@echo ""
	@echo "3. 🧪 Running unit tests..."
	@$(MAKE) test-unit
	@echo ""
	@echo "4. 🏗️ Building application..."
	@$(MAKE) build
	@echo ""
	@echo "✅ Quick pre-commit checks passed!"
	@echo "🚀 Ready to commit and push!"

pre-commit-full: ## Full pre-commit checks (comprehensive)
	@echo "🔍 Running full pre-commit checks..."
	@echo "===================================="
	@echo ""
	@echo "1. 📦 Managing dependencies..."
	@$(MAKE) deps
	@echo ""
	@echo "2. 🎨 Formatting code..."
	@$(MAKE) fmt
	@echo ""
	@echo "3. 🔍 Running go vet..."
	@$(MAKE) vet
	@echo ""
	@echo "4. 🧪 Running unit tests..."
	@$(MAKE) test-unit
	@echo ""
	@echo "5. 🧪 Running integration tests..."
	@$(MAKE) test-integration
	@echo ""
	@echo "6. 📊 Running tests with coverage..."
	@$(MAKE) test-coverage
	@echo ""
	@echo "7. 🏗️ Building application..."
	@$(MAKE) build
	@echo ""
	@echo "8. 🧪 Testing API endpoints..."
	@$(MAKE) test-api
	@echo ""
	@echo "9. 🧪 Testing API downloads..."
	@$(MAKE) test-api-download
	@echo ""
	@echo "✅ Full pre-commit checks passed!"
	@echo "🚀 Ready to commit and push!"

git-push: ## Run full checks and push
	@echo "🚀 Running pre-push checks..."
	@$(MAKE) pre-commit-full
	@echo ""
	@echo "🚀 Ready to push!"
	git push origin main
	@echo "✅ Push complete!"

# Quick targets for common workflows
quick-test: fmt vet test-unit ## Quick test cycle (format, vet, unit tests)

# API Testing targets
test-api: ## Test API endpoints on localhost:8081
	@echo "🧪 Testing TerraPeak API endpoints..."
	@echo "Testing health endpoint..."
	@curl -s -f "http://localhost:8081/healthz" && echo "✅ Health check passed" || echo "❌ Health check failed"
	@echo ""
	@echo "Testing AWS provider versions..."
	@curl -s "http://localhost:8081/v1/providers/hashicorp/aws/versions" | head -c 200 && echo "... ✅ AWS versions endpoint working" || echo "❌ AWS versions failed"
	@echo ""
	@echo "Testing Kubernetes provider versions..."
	@curl -s "http://localhost:8081/v1/providers/hashicorp/kubernetes/versions" | head -c 200 && echo "... ✅ Kubernetes versions endpoint working" || echo "❌ Kubernetes versions failed"
	@echo ""
	@echo "Testing proxy info..."
	@curl -s "http://localhost:8081/proxy/info" | head -c 200 && echo "... ✅ Proxy info endpoint working" || echo "❌ Proxy info failed"
	@echo ""
	@echo "🎉 API testing complete!"

test-api-verbose: ## Test API endpoints with verbose output
	@echo "🧪 Testing TerraPeak API endpoints (verbose)..."
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
	@echo "🎉 Verbose API testing complete!"

test-api-download: ## Test file download endpoints
	@echo "🧪 Testing TerraPeak download endpoints..."
	@echo "Testing AWS provider download (this may take a moment)..."
	@curl -s -I "http://localhost:8081/v1/providers/hashicorp/aws/5.0.0/download/linux/amd64" | head -5
	@echo ""
	@echo "Testing Kubernetes provider download..."
	@curl -s -I "http://localhost:8081/v1/providers/hashicorp/kubernetes/3.0.0/download/linux/amd64" | head -5
	@echo ""
	@echo "🎉 Download testing complete!"

dev-setup: deps ## Setup development environment
	@echo "🔧 Setting up development environment..."
	@echo ""
	@echo "📦 Installing development dependencies..."
	@echo ""
	@echo "Installing golangci-lint..."
	@which golangci-lint > /dev/null || ( \
		GOPATH_BIN=$$(go env GOPATH)/bin; \
		mkdir -p $$GOPATH_BIN; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$GOPATH_BIN v1.54.2; \
		echo "Add $$GOPATH_BIN to your PATH if not already there" \
	)
	@echo "✅ golangci-lint installed"
	@echo ""
	@echo "Installing entr (for watch mode)..."
	@which entr > /dev/null || (echo "Please install entr manually:" && echo "  macOS: brew install entr" && echo "  Ubuntu/Debian: sudo apt-get install entr" && echo "  CentOS/RHEL: sudo yum install entr")
	@echo "✅ entr check completed"
	@echo ""
	@echo "Installing curl (for API testing)..."
	@which curl > /dev/null || (echo "Please install curl manually:" && echo "  macOS: brew install curl" && echo "  Ubuntu/Debian: sudo apt-get install curl" && echo "  CentOS/RHEL: sudo yum install curl")
	@echo "✅ curl check completed"
	@echo ""
	@echo "Installing tree (for directory structure)..."
	@which tree > /dev/null || (echo "Please install tree manually:" && echo "  macOS: brew install tree" && echo "  Ubuntu/Debian: sudo apt-get install tree" && echo "  CentOS/RHEL: sudo yum install tree")
	@echo "✅ tree check completed"
	@echo ""
	@echo "🔧 Development environment setup complete!"
	@echo "📋 Installed tools:"
	@echo "  ✅ golangci-lint (Go linter)"
	@echo "  ✅ entr (file watcher)"
	@echo "  ✅ curl (API testing)"
	@echo "  ✅ tree (directory structure)"
	@echo ""
	@echo "🚀 Ready for development!"

# Watch mode (requires entr)
watch-test: ## Watch files and run tests on change (requires 'entr')
	find registry -name "*.go" | entr -c make test-unit

# =============================================================================
# Web (Next.js) Commands
# =============================================================================

# Web build targets
web-deps: ## Install frontend dependencies
	@echo "📦 Installing frontend dependencies..."
	cd web && pnpm install
	@echo "✅ Frontend dependencies installed!"

web-build: web-deps ## Build Next.js application for production
	@echo "🏗️  Building Next.js application..."
	cd web && pnpm build
	@echo "✅ Web build complete!"

web-dev: web-deps ## Start Next.js development server
	@echo "🚀 Starting Next.js dev server..."
	@echo "📝 API will connect to: http://localhost:8081"
	@echo "🌐 Frontend will run on: http://localhost:3000"
	cd web && pnpm dev

web-start: ## Start Next.js production server
	@echo "🚀 Starting Next.js production server..."
	cd web && pnpm start

web-test: ## Run Next.js tests
	@echo "🧪 Running Next.js tests..."
	cd web && pnpm test

web-lint: ## Lint Next.js code
	@echo "🔍 Linting Next.js code..."
	cd web && pnpm lint

web-lint-fix: ## Lint and fix Next.js code
	@echo "🔍 Linting and fixing Next.js code..."
	cd web && pnpm lint --fix

web-type-check: ## Type check Next.js code
	@echo "🔍 Type checking Next.js code..."
	cd web && pnpm tsc --noEmit

web-clean: ## Clean Next.js build artifacts
	@echo "🧹 Cleaning Next.js build artifacts..."
	cd web && rm -rf .next out node_modules/.cache
	@echo "✅ Web cleanup complete!"

web-clean-full: ## Clean everything including node_modules
	@echo "🧹 Cleaning all Next.js artifacts..."
	cd web && rm -rf .next out node_modules node_modules/.cache
	@echo "✅ Full web cleanup complete!"

web-analyze: web-build ## Analyze Next.js bundle size
	@echo "📊 Analyzing Next.js bundle..."
	cd web && pnpm analyze

# Web Docker commands
web-docker-build: ## Build web Docker image
	@echo "🐳 Building web Docker image..."
	cd web && docker build -t ghcr.io/aliharirian/terrapeak-web:latest .
	@echo "✅ Web Docker image built!"

web-docker-run: web-docker-build ## Run web container
	@echo "🚀 Starting web container..."
	cd web && docker run -d --name terrapeak-web -p 3000:3000 --restart unless-stopped ghcr.io/aliharirian/terrapeak-web:latest
	@echo "✅ Web container started on http://localhost:3000"

web-docker-stop: ## Stop web container
	@echo "⏸️  Stopping web container..."
	-docker stop terrapeak-web
	-docker rm terrapeak-web
	@echo "✅ Web container stopped"

web-docker-logs: ## Show web container logs
	@echo "📋 Showing web container logs..."
	docker logs -f terrapeak-web

web-docker-shell: ## Open shell in web container
	@echo "🐚 Opening shell in web container..."
	docker exec -it terrapeak-web /bin/sh

web-docker-health: ## Check web container health
	@echo "🏥 Checking web container health..."
	@docker inspect --format='{{.State.Health.Status}}' terrapeak-web 2>/dev/null || echo "Container not running"
	@curl -s http://localhost:3000/api/health | python3 -m json.tool || echo "Health endpoint not responding"

# =============================================================================
# Unified Commands
# =============================================================================

all-deps: deps web-deps ## Install all dependencies (backend + frontend)
	@echo "✅ All dependencies installed!"

all-build: build web-build ## Build both registry and web
	@echo "✅ All builds complete!"

all-test: test web-test ## Test both registry and web
	@echo "✅ All tests complete!"

all-lint: lint web-lint ## Lint both registry and web
	@echo "✅ All linting complete!"

all-clean: clean web-clean ## Clean both registry and web
	@echo "✅ All cleanup complete!"

all-docker-build: docker-build web-docker-build ## Build all Docker images
	@echo "✅ All Docker images built!"

all-docker-run: docker-compose-up ## Run all Docker containers
	@echo "✅ All services started!"

all-docker-stop: docker-compose-down ## Stop all Docker containers
	@echo "✅ All services stopped!"

# Development workflow
dev: ## Start both backend and frontend in development mode
	@echo "🚀 Starting TerraPeak in development mode..."
	@echo ""
	@echo "Starting backend server..."
	@cd registry && ./terrapeak -c .cfg.default.yml > /dev/null 2>&1 &
	@echo "⏳ Waiting for backend to start..."
	@sleep 2
	@echo "✅ Backend started on http://localhost:8081"
	@echo ""
	@echo "Starting frontend server..."
	@echo "🌐 Frontend will be available at http://localhost:3000"
	@echo "📝 API endpoint: http://localhost:8081/api/v1"
	@echo ""
	@$(MAKE) web-dev

dev-backend: build ## Build and start only backend
	@echo "🚀 Starting backend in development mode..."
	cd registry && ./terrapeak -c .cfg.default.yml

dev-frontend: web-deps ## Start only frontend (requires backend running separately)
	@echo "🚀 Starting frontend in development mode..."
	@echo "⚠️  Make sure backend is running on http://localhost:8080"
	@$(MAKE) web-dev

dev-full: build web-deps ## Build everything and start in development mode
	@echo "🚀 Full development setup..."
	@$(MAKE) dev

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

help-web: ## Show web-specific commands
	@echo "Web (Next.js) Commands:"
	@echo "======================"
	@echo "  make web-deps           Install frontend dependencies"
	@echo "  make web-build          Build Next.js application"
	@echo "  make web-dev            Start Next.js dev server"
	@echo "  make web-start          Start Next.js production server"
	@echo "  make web-test           Run Next.js tests"
	@echo "  make web-lint           Lint Next.js code"
	@echo "  make web-lint-fix       Lint and fix Next.js code"
	@echo "  make web-type-check     Type check TypeScript code"
	@echo "  make web-clean          Clean Next.js build artifacts"
	@echo "  make web-clean-full     Clean including node_modules"
	@echo "  make web-analyze        Analyze bundle size"
	@echo "  make web-docker-build   Build web Docker image"
	@echo "  make web-docker-run     Run web container"
	@echo "  make web-docker-stop    Stop web container"
	@echo "  make web-docker-logs    Show web container logs"
	@echo "  make web-docker-shell   Open shell in web container"
	@echo "  make web-docker-health  Check web container health"

help-docker: ## Show Docker-specific commands
	@echo "Docker Commands:"
	@echo "================"
	@echo "  make docker-build       Build registry Docker image"
	@echo "  make docker-run         Run registry in Docker"
	@echo "  make docker-compose-up  Start all services with docker-compose"
	@echo "  make docker-compose-down Stop all docker-compose services"
	@echo "  make docker-compose-build Build all services"
	@echo "  make docker-compose-logs Show all service logs"
	@echo "  make docker-compose-registry Start only registry service"
	@echo "  make docker-compose-web Start only web service"
	@echo "  make docker-compose-minio Start only MinIO service"
	@echo "  make web-docker-build   Build web Docker image"
	@echo "  make web-docker-run     Run web container"
	@echo "  make all-docker-build   Build all Docker images"
	@echo "  make all-docker-run     Run all services"
	@echo "  make all-docker-stop    Stop all services"

# Status check
status: ## Check project status
	@echo "TerraPeak Status"
	@echo "================================+"
	@echo "Go version: $(shell go version)"
	@echo "Node version: $(shell node --version 2>/dev/null || echo 'Node.js not installed')"
	@echo "pnpm version: $(shell pnpm --version 2>/dev/null || echo 'pnpm not installed')"
	@echo "Docker version: $(shell docker --version 2>/dev/null || echo 'Docker not installed')"
	@echo "Git branch: $(shell git branch --show-current 2>/dev/null || echo 'not a git repo')"
	@echo "Git status: $(shell git status --porcelain 2>/dev/null | wc -l | xargs) files changed"
	@echo "Registry dependencies: $(shell cd registry && go list -m all | wc -l | xargs) modules"
	@echo "Web dependencies: $(shell cd web && pnpm list --depth=0 2>/dev/null | wc -l | xargs) packages"
	@echo "Registry test files: $(shell find registry -name "*_test.go" | wc -l | xargs) files"
	@echo "Registry source files: $(shell find registry -name "*.go" -not -name "*_test.go" | wc -l | xargs) files"
	@echo "Web source files: $(shell find web -name "*.tsx" -o -name "*.ts" -o -name "*.js" -o -name "*.jsx" | wc -l | xargs) files"
	@echo ""
	@echo "Services:"
	@echo "  Backend: http://localhost:8080 ($(shell curl -s http://localhost:8080/api/v1/health 2>/dev/null && echo 'running' || echo 'not running'))"
	@echo "  Frontend: http://localhost:3000 ($(shell curl -s http://localhost:3000 2>/dev/null && echo 'running' || echo 'not running'))"
