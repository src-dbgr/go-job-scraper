# Makefile for job-scraper project

.PHONY: build run test lint docker-up docker-down

# Build the application
build:
	go build -o job-scraper ./cmd/scraper

# Run the application
run: build
	./job-scraper

# Run tests
test:
	go test ./...

# Run linter
lint:
	golangci-lint run

# Start Docker services
docker-up:
	docker compose up -d

# Stop Docker services
docker-down:
	docker compose down

# Clean up build artifacts
clean:
	rm -f job-scraper
	go clean

# Update dependencies
deps-update:
	go get -u ./...
	go mod tidy