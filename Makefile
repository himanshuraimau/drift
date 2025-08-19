# Distributed File System (Drift) Makefile

# Go parameters
GOCMD=/usr/local/go/bin/go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=drift
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) -v .

# Build for production (with optimizations)
build-prod:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -ldflags="-w -s" -o $(BINARY_NAME) .

# Test all packages
test:
	$(GOTEST) -v ./...

# Test with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Test with race detection
test-race:
	$(GOTEST) -v -race ./...

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v .
	./$(BINARY_NAME)

# Run tests and then the application
run-test:
	$(GOTEST) -v ./...
	$(GOBUILD) -o $(BINARY_NAME) -v .
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out coverage.html
	rm -rf *_network

# Clean test data
clean-test:
	rm -rf test_store*
	rm -rf *_network

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
update-deps:
	$(GOMOD) tidy
	$(GOGET) -u ./...

# Cross compilation for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v .

# Cross compilation for Windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME).exe -v .

# Cross compilation for macOS
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_darwin -v .

# Build for all platforms
build-all: build-linux build-windows build-darwin

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Vet code
vet:
	$(GOCMD) vet ./...

# Check code quality
check: fmt vet lint test

# Run benchmark tests
bench:
	$(GOTEST) -bench=. -benchmem ./...

# Install the binary
install:
	$(GOBUILD) -o $(BINARY_NAME) -v .
	sudo mv $(BINARY_NAME) /usr/local/bin/

# Show help
help:
	@echo "Available commands:"
	@echo "  build         - Build the binary"
	@echo "  build-prod    - Build optimized binary for production"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-race     - Run tests with race detection"
	@echo "  run           - Build and run the application"
	@echo "  run-test      - Run tests then the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  clean-test    - Clean test data"
	@echo "  deps          - Install dependencies"
	@echo "  update-deps   - Update dependencies"
	@echo "  build-linux   - Cross compile for Linux"
	@echo "  build-windows - Cross compile for Windows"
	@echo "  build-darwin  - Cross compile for macOS"
	@echo "  build-all     - Build for all platforms"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  vet           - Vet code"
	@echo "  check         - Run fmt, vet, lint, and test"
	@echo "  bench         - Run benchmark tests"
	@echo "  install       - Install binary to /usr/local/bin"
	@echo "  help          - Show this help message"

.PHONY: build build-prod test test-coverage test-race run run-test clean clean-test deps update-deps build-linux build-windows build-darwin build-all fmt lint vet check bench install help
