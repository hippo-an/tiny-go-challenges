# protem-gen Makefile

.PHONY: build install test test-coverage test-e2e clean help

# Build variables
VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
MODULE := github.com/hippo-an/tiny-go-challenges/protem-gen/internal/cli
LDFLAGS := -ldflags "-X $(MODULE).Version=$(VERSION) -X $(MODULE).Commit=$(COMMIT) -X $(MODULE).BuildDate=$(BUILD_DATE)"

# Default target
.DEFAULT_GOAL := help

## build: Build the binary to ./bin/protem-gen
build:
	@echo "Building protem-gen..."
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/protem-gen ./cmd/protem-gen
	@echo "Built: bin/protem-gen"

## install: Install protem-gen to $GOPATH/bin
install:
	@echo "Installing protem-gen..."
	go install $(LDFLAGS) ./cmd/protem-gen
	@echo "Installed: $(shell go env GOPATH)/bin/protem-gen"

## test: Run all tests
test:
	go test -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	go test -cover ./...

## test-e2e: Run E2E tests (requires build first)
test-e2e: build
	@echo "Running E2E tests..."
	./scripts/e2e_test.sh

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	@echo "Done"

## help: Show this help message
help:
	@echo "protem-gen - Go Web Application Project Generator"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/ /'
