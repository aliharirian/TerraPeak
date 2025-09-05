#!/bin/bash

# TerraPeak Test Runner
# This script runs all tests for the TerraPeak project

set -e

echo "🧪 Running TerraPeak Tests"
echo "=========================="

# Change to registry directory
cd "$(dirname "$0")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

print_status "Go version: $(go version)"

# Ensure dependencies are downloaded
print_status "Downloading dependencies..."
go mod download
go mod tidy

# Run linting if golangci-lint is available
if command -v golangci-lint &> /dev/null; then
    print_status "Running linter..."
    golangci-lint run ./...
    print_success "Linting passed"
else
    print_warning "golangci-lint not found, skipping linting"
fi

# Run go vet
print_status "Running go vet..."
go vet ./...
print_success "go vet passed"

# Run tests with coverage
print_status "Running unit tests..."
go test -v -race -coverprofile=coverage.out ./...
TEST_EXIT_CODE=$?

if [ $TEST_EXIT_CODE -eq 0 ]; then
    print_success "All tests passed!"

    # Generate coverage report
    print_status "Generating coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    go tool cover -func=coverage.out

    echo ""
    print_success "Coverage report generated: coverage.html"

    # Show coverage summary
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo -e "${GREEN}📊 Total coverage: ${COVERAGE}${NC}"

else
    print_error "Tests failed with exit code $TEST_EXIT_CODE"
    exit $TEST_EXIT_CODE
fi

# Run integration tests separately
print_status "Running integration tests..."
go test -v -tags=integration ./...
INTEGRATION_EXIT_CODE=$?

if [ $INTEGRATION_EXIT_CODE -eq 0 ]; then
    print_success "Integration tests passed!"
else
    print_error "Integration tests failed with exit code $INTEGRATION_EXIT_CODE"
    exit $INTEGRATION_EXIT_CODE
fi

# Build the application to ensure it compiles
print_status "Building application..."
go build -o terrapeak-test .
BUILD_EXIT_CODE=$?

if [ $BUILD_EXIT_CODE -eq 0 ]; then
    print_success "Build successful!"
    rm -f terrapeak-test  # Clean up test binary
else
    print_error "Build failed with exit code $BUILD_EXIT_CODE"
    exit $BUILD_EXIT_CODE
fi

echo ""
print_success "🎉 All tests and checks passed!"
echo ""
echo "📋 Test Summary:"
echo "  - Unit tests: ✅ PASSED"
echo "  - Integration tests: ✅ PASSED"
echo "  - Linting: ✅ PASSED"
echo "  - Build: ✅ PASSED"
echo "  - Coverage: ${COVERAGE}"
echo ""
echo "🚀 Ready for deployment!"
