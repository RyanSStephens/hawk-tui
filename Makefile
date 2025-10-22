# Hawk TUI Universal Makefile
# Supports building for multiple platforms and languages

.DEFAULT_GOAL := help
.PHONY: help build test clean install package release dev examples docs

# Configuration
APP_NAME := hawk
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v1.0.0")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)

# Directories
BUILD_DIR := build
DIST_DIR := dist
SCRIPTS_DIR := scripts
EXAMPLES_DIR := examples
DOCS_DIR := docs

# Go build settings
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GO_BUILD_FLAGS := -ldflags "$(LDFLAGS)" -trimpath

# Node.js settings
NODE_VERSION := $(shell node --version 2>/dev/null || echo "not-installed")
NPM_VERSION := $(shell npm --version 2>/dev/null || echo "not-installed")

# Python settings
PYTHON_VERSION := $(shell python3 --version 2>/dev/null || echo "not-installed")

# Colors for output
BOLD := \033[1m
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

define print_success
	@echo "$(GREEN)✓$(RESET) $(1)"
endef

define print_info
	@echo "$(BLUE)ℹ$(RESET) $(1)"
endef

define print_warning
	@echo "$(YELLOW)⚠$(RESET) $(1)"
endef

define print_error
	@echo "$(RED)✗$(RESET) $(1)"
endef

## help: Show this help message
help:
	@echo "$(BOLD)🦅 Hawk TUI - Universal TUI Framework$(RESET)"
	@echo ""
	@echo "$(BOLD)Available targets:$(RESET)"
	@awk '/^##/ { \
		split($$0, a, ":"); \
		split(a[1], b, " "); \
		printf "  $(BLUE)%-15s$(RESET) %s\n", b[2], substr($$0, index($$0, ":")+2) \
	}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BOLD)Current configuration:$(RESET)"
	@echo "  Version:     $(VERSION)"
	@echo "  Commit:      $(COMMIT)"
	@echo "  Platform:    $(GOOS)/$(GOARCH)"
	@echo "  Go:          $(shell go version 2>/dev/null || echo 'not installed')"
	@echo "  Node.js:     $(NODE_VERSION)"
	@echo "  Python:      $(PYTHON_VERSION)"

## build: Build the main binary for current platform
build:
	$(call print_info,"Building Hawk TUI for $(GOOS)/$(GOARCH)")
	@mkdir -p $(BUILD_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/hawk
	$(call print_success,"Binary built: $(BUILD_DIR)/$(APP_NAME)")

## build-all: Build binaries for all supported platforms
build-all:
	$(call print_info,"Building for all supported platforms")
	@mkdir -p $(BUILD_DIR)
	
	# Linux builds
	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./cmd/hawk
	GOOS=linux GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 ./cmd/hawk
	GOOS=linux GOARCH=arm go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-arm ./cmd/hawk
	
	# macOS builds
	GOOS=darwin GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./cmd/hawk
	GOOS=darwin GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ./cmd/hawk
	
	# Windows builds
	GOOS=windows GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ./cmd/hawk
	GOOS=windows GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe ./cmd/hawk
	
	$(call print_success,"All platform binaries built in $(BUILD_DIR)/")

## test: Run all tests (Go, Node.js, Python)
test: test-go test-nodejs test-python

## test-go: Run Go tests
test-go:
	$(call print_info,"Running Go tests")
	go test -v -race -coverprofile=coverage.out ./pkg/... ./internal/... ./cmd/...
	go tool cover -html=coverage.out -o coverage.html
	$(call print_success,"Go tests completed")

## test-nodejs: Run Node.js client tests
test-nodejs:
	$(call print_info,"Running Node.js client tests")
	@if [ "$(NODE_VERSION)" = "not-installed" ]; then \
		$(call print_warning,"Node.js not installed, skipping tests"); \
	else \
		cd $(EXAMPLES_DIR)/nodejs && npm test; \
		$(call print_success,"Node.js tests completed"); \
	fi

## test-python: Run Python client tests
test-python:
	$(call print_info,"Running Python client tests")
	@if [ "$(PYTHON_VERSION)" = "not-installed" ]; then \
		$(call print_warning,"Python not installed, skipping tests"); \
	else \
		cd $(EXAMPLES_DIR)/python && python3 -m pytest test_hawk.py -v; \
		$(call print_success,"Python tests completed"); \
	fi

## lint: Run linters for all languages
lint: lint-go lint-nodejs lint-python

## lint-go: Run Go linting
lint-go:
	$(call print_info,"Running Go linter")
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		$(call print_success,"Go linting completed"); \
	else \
		$(call print_warning,"golangci-lint not installed, using go vet"); \
		go vet ./...; \
	fi

## lint-nodejs: Run Node.js linting
lint-nodejs:
	$(call print_info,"Running Node.js linter")
	@if [ "$(NODE_VERSION)" = "not-installed" ]; then \
		$(call print_warning,"Node.js not installed, skipping linting"); \
	else \
		cd $(EXAMPLES_DIR)/nodejs && npm run lint; \
		$(call print_success,"Node.js linting completed"); \
	fi

## lint-python: Run Python linting
lint-python:
	$(call print_info,"Running Python linter")
	@if [ "$(PYTHON_VERSION)" = "not-installed" ]; then \
		$(call print_warning,"Python not installed, skipping linting"); \
	else \
		cd $(EXAMPLES_DIR)/python && python3 -m flake8 hawk.py; \
		$(call print_success,"Python linting completed"); \
	fi

## clean: Clean build artifacts and temporary files
clean:
	$(call print_info,"Cleaning build artifacts")
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html
	rm -rf $(EXAMPLES_DIR)/nodejs/node_modules
	find . -name "*.pyc" -delete
	find . -name "__pycache__" -delete
	$(call print_success,"Cleanup completed")

## install: Install the binary to system PATH
install: build
	$(call print_info,"Installing Hawk TUI")
	@if [ "$(GOOS)" = "windows" ]; then \
		$(call print_error,"Use install.ps1 script for Windows installation"); \
		exit 1; \
	fi
	
	@if [ "$$(id -u)" = "0" ]; then \
		cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/; \
		chmod +x /usr/local/bin/$(APP_NAME); \
		$(call print_success,"Installed to /usr/local/bin/$(APP_NAME)"); \
	else \
		mkdir -p $$HOME/.local/bin; \
		cp $(BUILD_DIR)/$(APP_NAME) $$HOME/.local/bin/; \
		chmod +x $$HOME/.local/bin/$(APP_NAME); \
		$(call print_success,"Installed to $$HOME/.local/bin/$(APP_NAME)"); \
		$(call print_warning,"Make sure $$HOME/.local/bin is in your PATH"); \
	fi

## package: Create distribution packages
package: build-all
	$(call print_info,"Creating distribution packages")
	@mkdir -p $(DIST_DIR)
	
	# Linux packages
	@for arch in amd64 arm64 arm; do \
		mkdir -p $(DIST_DIR)/hawk-tui-$(VERSION)-linux-$$arch; \
		cp $(BUILD_DIR)/$(APP_NAME)-linux-$$arch $(DIST_DIR)/hawk-tui-$(VERSION)-linux-$$arch/hawk; \
		cp README.md LICENSE $(DIST_DIR)/hawk-tui-$(VERSION)-linux-$$arch/; \
		tar -czf $(DIST_DIR)/hawk-tui-$(VERSION)-linux-$$arch.tar.gz -C $(DIST_DIR) hawk-tui-$(VERSION)-linux-$$arch; \
		rm -rf $(DIST_DIR)/hawk-tui-$(VERSION)-linux-$$arch; \
	done
	
	# macOS packages
	@for arch in amd64 arm64; do \
		mkdir -p $(DIST_DIR)/hawk-tui-$(VERSION)-darwin-$$arch; \
		cp $(BUILD_DIR)/$(APP_NAME)-darwin-$$arch $(DIST_DIR)/hawk-tui-$(VERSION)-darwin-$$arch/hawk; \
		cp README.md LICENSE $(DIST_DIR)/hawk-tui-$(VERSION)-darwin-$$arch/; \
		tar -czf $(DIST_DIR)/hawk-tui-$(VERSION)-darwin-$$arch.tar.gz -C $(DIST_DIR) hawk-tui-$(VERSION)-darwin-$$arch; \
		rm -rf $(DIST_DIR)/hawk-tui-$(VERSION)-darwin-$$arch; \
	done
	
	# Windows packages
	@for arch in amd64 arm64; do \
		mkdir -p $(DIST_DIR)/hawk-tui-$(VERSION)-windows-$$arch; \
		cp $(BUILD_DIR)/$(APP_NAME)-windows-$$arch.exe $(DIST_DIR)/hawk-tui-$(VERSION)-windows-$$arch/hawk.exe; \
		cp README.md LICENSE $(DIST_DIR)/hawk-tui-$(VERSION)-windows-$$arch/; \
		cd $(DIST_DIR) && zip -r hawk-tui-$(VERSION)-windows-$$arch.zip hawk-tui-$(VERSION)-windows-$$arch; \
		rm -rf $(DIST_DIR)/hawk-tui-$(VERSION)-windows-$$arch; \
	done
	
	$(call print_success,"Distribution packages created in $(DIST_DIR)/")

## release: Build, test, package, and prepare for release
release: clean test lint build-all package
	$(call print_info,"Preparing release $(VERSION)")
	@echo "$(BOLD)Release artifacts:$(RESET)"
	@ls -la $(DIST_DIR)/
	$(call print_success,"Release $(VERSION) ready!")

## dev: Set up development environment
dev:
	$(call print_info,"Setting up development environment")
	
	# Install Go dependencies
	go mod download
	go mod verify
	
	# Install Node.js dependencies
	@if [ "$(NODE_VERSION)" != "not-installed" ]; then \
		cd $(EXAMPLES_DIR)/nodejs && npm install; \
		$(call print_success,"Node.js dependencies installed"); \
	else \
		$(call print_warning,"Node.js not installed, skipping npm install"); \
	fi
	
	# Install Python dependencies
	@if [ "$(PYTHON_VERSION)" != "not-installed" ]; then \
		cd $(EXAMPLES_DIR)/python && pip3 install -r requirements.txt 2>/dev/null || true; \
		$(call print_success,"Python dependencies installed"); \
	else \
		$(call print_warning,"Python not installed, skipping pip install"); \
	fi
	
	# Install development tools
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		$(call print_info,"Installing golangci-lint"); \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	
	$(call print_success,"Development environment ready!")

## examples: Run example applications
examples:
	$(call print_info,"Running example applications")
	
	@echo "$(BOLD)Available examples:$(RESET)"
	@echo "  Node.js: cd $(EXAMPLES_DIR)/nodejs && node demo.js | $(BUILD_DIR)/$(APP_NAME)"
	@echo "  Python:  cd $(EXAMPLES_DIR)/python && python3 demo.py | $(BUILD_DIR)/$(APP_NAME)"
	@echo ""
	@echo "Run these commands to see Hawk TUI in action!"

## docs: Generate documentation
docs:
	$(call print_info,"Generating documentation")
	@mkdir -p $(DOCS_DIR)
	
	# Generate Go documentation
	go doc -all ./pkg/types > $(DOCS_DIR)/go-api.md
	
	# TODO: Add more documentation generation
	$(call print_success,"Documentation generated in $(DOCS_DIR)/")

## version: Show version information
version:
	@echo "$(BOLD)Hawk TUI Version Information$(RESET)"
	@echo "Version:    $(VERSION)"
	@echo "Commit:     $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo "Platform:   $(GOOS)/$(GOARCH)"

## docker: Build Docker image
docker:
	$(call print_info,"Building Docker image")
	docker build -t hawk-tui:$(VERSION) .
	docker tag hawk-tui:$(VERSION) hawk-tui:latest
	$(call print_success,"Docker image built: hawk-tui:$(VERSION)")

## docker-run: Run Hawk TUI in Docker
docker-run:
	$(call print_info,"Running Hawk TUI in Docker")
	docker run -it --rm hawk-tui:latest

# Development helpers
watch:
	$(call print_info,"Watching for changes (requires entr)")
	find . -name "*.go" | entr -r make build

fmt:
	$(call print_info,"Formatting Go code")
	go fmt ./...
	$(call print_success,"Code formatted")

tidy:
	$(call print_info,"Tidying Go modules")
	go mod tidy
	$(call print_success,"Modules tidied")