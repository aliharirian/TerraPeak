# Contributing to TerraPeak

Thank you for your interest in contributing to TerraPeak! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

We are committed to providing a welcoming and inclusive environment. Please be respectful and constructive in all interactions.

### Our Standards

- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Docker and Docker Compose (for integration testing)
- Make (optional but recommended)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/TerraPeak.git
   cd TerraPeak
   ```

3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/aliharirian/TerraPeak.git
   ```

### Setup Development Environment

```bash
# Navigate to registry directory
cd registry

# Install dependencies
go mod download

# Verify setup
go build -o terrapeak .
./terrapeak --version
```

## Development Workflow

### Create a Feature Branch

Always create a new branch for your work:

```bash
# Update main branch
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feat/your-feature-name
# or
git checkout -b fix/bug-description
```

### Branch Naming Convention

- `feat/feature-name` - New features
- `fix/bug-description` - Bug fixes
- `docs/description` - Documentation changes
- `refactor/description` - Code refactoring
- `test/description` - Test improvements
- `chore/description` - Maintenance tasks

### Making Changes

1. Make your changes in the feature branch
2. Write or update tests as needed
3. Ensure all tests pass
4. Follow coding standards
5. Commit your changes with meaningful messages

### Commit Message Format

Follow conventional commits specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(cache): add TTL support for cache entries

Implement configurable TTL for cached provider binaries.
Add expiration checking on cache hit.

Closes #123

fix(storage): handle connection timeout to S3

Add retry logic and better error messages for S3 connection failures.

docs: update installation guide for Kubernetes

Add Helm chart deployment instructions.
```

## Coding Standards

### Go Style Guide

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run vet
go vet ./...
```

### Project-Specific Guidelines

1. **Interfaces**: Use interfaces for abstraction (e.g., `Storage` interface)
2. **Error Handling**: Always check and handle errors appropriately
3. **Logging**: Use structured logging with zerolog
4. **Context**: Pass `context.Context` for cancellation support
5. **Testing**: Write tests for all new code
6. **Documentation**: Document exported functions and types

### Example Code Style

```go
// Package documentation
package storage

import (
    "context"
    "io"
)

// Storage defines the interface for backend storage operations.
// All storage implementations must satisfy this interface.
type Storage interface {
    // Save stores data with the given key.
    // Returns an error if the operation fails.
    Save(ctx context.Context, key string, data io.Reader) error
    
    // Get retrieves data for the given key.
    // Returns an error if the key doesn't exist or retrieval fails.
    Get(ctx context.Context, key string) (io.ReadCloser, error)
}

// s3Storage implements Storage interface using S3-compatible storage.
type s3Storage struct {
    client *s3.Client
    bucket string
}

// NewS3Storage creates a new S3 storage backend.
func NewS3Storage(cfg *config.Config) (Storage, error) {
    if cfg.Storage.S3.Bucket == "" {
        return nil, errors.New("S3 bucket name is required")
    }
    
    // Implementation...
    return &s3Storage{}, nil
}
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run with coverage
make test-coverage

# Run specific package tests
cd registry
go test ./api/...

# Run specific test
go test -run TestProviderVersions ./api/
```

### Writing Tests

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **Table-Driven Tests**: Use for multiple test cases

**Example:**

```go
func TestCacheKeyGeneration(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "provider versions",
            input:    "/v1/providers/hashicorp/aws/versions",
            expected: "registry/v1/versions/hashicorp/aws",
        },
        {
            name:     "provider download",
            input:    "/v1/providers/hashicorp/aws/5.0.0/download/linux/amd64",
            expected: "registry/v1/download/hashicorp/aws/5.0.0/linux/amd64",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := generateCacheKey(tt.input)
            if result != tt.expected {
                t.Errorf("expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

### Test Coverage

- Maintain minimum 60% overall coverage
- Critical paths should have 90% coverage
- Test all error conditions

## Submitting Changes

### Pre-Commit Checklist

Before committing, ensure:

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Or use pre-commit command
make pre-commit
```

### Push Your Changes

```bash
# Commit changes
git add .
git commit -m "feat: add new feature"

# Push to your fork
git push origin feat/your-feature-name
```

### Create Pull Request

1. Go to your fork on GitHub
2. Click "New Pull Request"
3. Select your feature branch
4. Fill in the PR template:
   - Description of changes
   - Related issues
   - Testing performed
   - Screenshots (if applicable)

### PR Guidelines

- Keep PRs focused on a single feature or fix
- Update documentation as needed
- Add tests for new functionality
- Ensure CI passes
- Respond to review comments

### PR Review Process

1. Automated checks (tests, linting) must pass
2. At least one maintainer approval required
3. Address review comments
4. Squash commits if requested
5. Maintainer will merge when approved

## Reporting Issues

### Before Reporting

1. Check existing issues to avoid duplicates
2. Try to reproduce with latest version
3. Gather relevant information

### Issue Template

**Bug Report:**
```markdown
**Description:**
Clear description of the bug

**Steps to Reproduce:**
1. Step 1
2. Step 2
3. Expected behavior vs actual behavior

**Environment:**
- TerraPeak version:
- Go version:
- OS:
- Configuration:

**Logs:**
```log
Paste relevant logs here
```

**Feature Request:**
```markdown
**Problem:**
Describe the problem you're trying to solve

**Proposed Solution:**
Your proposed solution

**Alternatives:**
Other solutions you've considered

**Additional Context:**
Any other relevant information
```

## Development Tips

### Debugging

```bash
# Run with debug logging
./terrapeak -c cfg.yml --log-level=debug

# Use delve debugger
dlv debug -- -c cfg.yml
```

### Local Testing

```bash
# Start local MinIO
docker-compose up -d minio

# Run TerraPeak locally
cd registry
go run . -c ../cfg.yml

# Test with curl
curl http://localhost:8081/healthz
```

### IDE Setup

**VSCode:**
- Install Go extension
- Configure format on save
- Enable linting

**GoLand:**
- Enable Go modules
- Configure file watchers for gofmt
- Enable golangci-lint integration

## Getting Help

- **Documentation**: Read [README](README.md), [API docs](docs/registry/API.md), and [Architecture](docs/registry/ARCHITECTURE.md)
- **GitHub Issues**: Ask questions by creating an issue
- **Discussions**: Use GitHub Discussions for general questions

## License

By contributing to TerraPeak, you agree that your contributions will be licensed under the Apache License 2.0.

---

**Thank you for contributing to TerraPeak!**

