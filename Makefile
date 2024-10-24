# Makefile for job-scraper project

.PHONY: all clean build test test-unit test-integration lint deps docker-up docker-down

# Default target
all: clean deps test build

# Clean up build artifacts and temporary files
clean:
	rm -f job-scraper
	rm -rf dist/
	go clean -testcache
	go clean

# Install/Update dependencies
deps:
	go mod download
	go mod tidy

# Run linter
lint:
	golangci-lint run

# Build the application
build:
	mkdir -p dist
	go build -o dist/job-scraper ./cmd/scraper

# Run only unit tests
test-unit:
	go test -v $(shell go list ./... | grep -v /tests/integration)

# Run integration tests
test-integration:
	go test -v -tags=integration ./tests/integration

# Run all tests
test: test-unit test-integration

# Start Docker services
docker-up:
	docker compose up -d

# Stop Docker services
docker-down:
	docker compose down

# Run the application
run: build
	./dist/job-scraper

# Generate test coverage report
cover:
	mkdir -p coverage
	go test -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# Check for security vulnerabilities
sec-check:
	govulncheck ./...