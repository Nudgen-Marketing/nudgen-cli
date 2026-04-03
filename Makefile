SHELL := /bin/bash

# Binary name
BINARY := nudgen

# Build directory
BUILD_DIR := ./bin

# Version info (can be overridden: make build VERSION=1.0.0)
VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOVET := $(GOCMD) vet
GOFMT := gofmt
GOMOD := $(GOCMD) mod

# Version package path
VERSION_PKG := github.com/nudgen/nudgen-cli/internal/version

# Build flags
LDFLAGS := -s -w -X $(VERSION_PKG).Version=$(VERSION) -X $(VERSION_PKG).Commit=$(COMMIT) -X $(VERSION_PKG).Date=$(DATE)
BUILD_FLAGS := -trimpath -ldflags "$(LDFLAGS)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY) .

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run go vet
.PHONY: vet
vet:
	$(GOVET) ./...

# Format code
.PHONY: fmt
fmt:
	$(GOFMT) -s -w .

# Check formatting (for CI)
.PHONY: fmt-check
fmt-check:
	@test -z "$$($(GOFMT) -s -l . | tee /dev/stderr)" || (echo "Code is not formatted. Run 'make fmt'" && exit 1)

# Tidy dependencies
.PHONY: tidy
tidy:
	$(GOMOD) tidy

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install to GOPATH/bin
.PHONY: install
install:
	$(GOCMD) install $(BUILD_FLAGS) .

# Show help
.PHONY: help
help:
	@echo "nudgen Makefile targets:"
	@echo ""
	@echo "Build:"
	@echo "  build          Build the binary"
	@echo "  install        Install to GOPATH/bin"
	@echo "  clean          Remove build artifacts"
	@echo ""
	@echo "Quality:"
	@echo "  test           Run Go unit tests"
	@echo "  vet            Run go vet"
	@echo "  fmt            Format code"
	@echo "  fmt-check      Check code formatting"
	@echo "  tidy           Tidy go.mod dependencies"
	@echo ""
	@echo "  help           Show this help"
