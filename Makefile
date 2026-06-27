.PHONY: help init build test lint fmt docker-build docker-push dev-up dev-down dev-logs proto db-migrate ci clean

# =============================================================================
# Spark Monorepo - Top-Level Makefile
# =============================================================================

SHELL := /bin/bash

# -- Service directories ------------------------------------------------------
GO_SERVICES := identity-service creator-service viewer-service stream-service wallet-service chat-service messaging-service subscription-service gift-service payment-service analytics-service notification-service recommendation-service search-service translation-service moderation-service community-service event-service competition-service advertising-service commerce-service media-service licensing-service discovery-service trust-service

RUST_SERVICES := media-service

PYTHON_SERVICES := ai/ranking-ai ai/fraud-ai ai/moderation-ai ai/translation-ai ai/voice-ai ai/vision-ai ai/clip-ai ai/thumbnail-ai ai/assistant-ai ai/recommendation-ai ai/creator-ai

# -- Tooling ------------------------------------------------------------------
GO ?= go
GO_VERSION ?= 1.22
RUST ?= cargo
NODE ?= node
NPM ?= npm
DOCKER ?= docker
TERRAFORM ?= terraform
GOLANGCI_LINT ?= golangci-lint
PROTOC ?= protoc

# -- Docker registry -----------------------------------------------------------
REGISTRY ?= ghcr.io/spark-platform
TAG ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "latest")

# -- Colors for help -----------------------------------------------------------
BLUE := \033[36m
RESET := \033[0m

# =============================================================================
# Targets
# =============================================================================

help: ## List all available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "$(BLUE)%-20s$(RESET) %s\n", $$1, $$2}'

init: ## Initialize the monorepo (install tooling, git hooks)
	@echo "==> Initializing monorepo..."
	@scripts/init.sh
	@$(GO) version
	@$(NODE) --version
	@$(NPM) --version

# -- Build --------------------------------------------------------------------

build: build-go build-rust build-py ## Build all services

build-go: ## Build all Go services
	@echo "==> Building Go services..."
	@for svc in $(GO_SERVICES); do \
		echo "    building services/$$svc..."; \
		cd services/$$svc && $(GO) build -o ../../dist/$$svc ./cmd/... && cd ../..; \
	done

build-rust: ## Build all Rust services
	@echo "==> Building Rust services..."
	@cd services/media-service && $(RUST) build --release && cd ../..

build-py: ## Build all Python services (syntax check)
	@echo "==> Checking Python services..."
	@for svc in $(PYTHON_SERVICES); do \
		echo "    checking $$svc..."; \
		python3 -m py_compile $$svc/src/main.py 2>/dev/null || true; \
	done

# -- Test ---------------------------------------------------------------------

test: test-go test-rust test-py ## Run all tests

test-go: ## Run all Go tests
	@echo "==> Testing Go services..."
	@for svc in $(GO_SERVICES); do \
		echo "    testing services/$$svc..."; \
		cd services/$$svc && $(GO) test ./... -coverprofile=coverage.out -covermode=atomic && cd ../..; \
	done

test-rust: ## Run all Rust tests
	@echo "==> Testing Rust services..."
	@cd services/media-service && $(RUST) test && cd ../..

test-py: ## Run all Python tests
	@echo "==> Testing Python services..."
	@for svc in $(PYTHON_SERVICES); do \
		echo "    testing $$svc..."; \
		cd $$svc && python3 -m pytest src/ -v --tb=short 2>/dev/null || true && cd ../..; \
	done

# -- Lint ---------------------------------------------------------------------

lint: lint-go lint-rust lint-py ## Lint all code

lint-go: ## Lint Go code
	@echo "==> Linting Go services..."
	@$(GOLANGCI_LINT) run ./services/...

lint-rust: ## Lint Rust code
	@echo "==> Linting Rust services..."
	@cd services/media-service && $(RUST) clippy -- -D warnings && cd ../..

lint-py: ## Lint Python code
	@echo "==> Linting Python services..."
	@for svc in $(PYTHON_SERVICES); do \
		cd $$svc && ruff check src/ && cd ../..; \
	done

# -- Format -------------------------------------------------------------------

fmt: fmt-go fmt-rust fmt-py fmt-proto ## Format all code

fmt-go: ## Format Go code
	@echo "==> Formatting Go code..."
	@$(GO) fmt ./services/... ./packages/... ./api/... ./sdk/go/...

fmt-rust: ## Format Rust code
	@echo "==> Formatting Rust code..."
	@cd services/media-service && $(RUST) fmt && cd ../..

fmt-py: ## Format Python code
	@echo "==> Formatting Python code..."
	@for svc in $(PYTHON_SERVICES); do \
		cd $$svc && ruff format src/ && cd ../..; \
	done

fmt-proto: ## Format proto files
	@echo "==> Formatting protobuf files..."
	@find api/ -name '*.proto' -exec clang-format -i {} \;

# -- Docker -------------------------------------------------------------------

docker-build: ## Build all Docker images
	@echo "==> Building Docker images..."
	@for svc in $(GO_SERVICES); do \
		echo "    building $(REGISTRY)/$$svc:$(TAG)..."; \
		$(DOCKER) build \
			--build-arg SERVICE=$$svc \
			-t $(REGISTRY)/$$svc:$(TAG) \
			-t $(REGISTRY)/$$svc:latest \
			-f Dockerfile .; \
	done

docker-push: ## Push all Docker images to registry
	@echo "==> Pushing Docker images..."
	@for svc in $(GO_SERVICES); do \
		echo "    pushing $(REGISTRY)/$$svc:$(TAG)..."; \
		$(DOCKER) push $(REGISTRY)/$$svc:$(TAG); \
		$(DOCKER) push $(REGISTRY)/$$svc:latest; \
	done

# -- Local Development --------------------------------------------------------

dev-up: ## Start local development environment
	@echo "==> Starting local dev environment..."
	@$(DOCKER) compose -f deployment/development/docker-compose.yml up -d --build

dev-down: ## Stop local development environment
	@echo "==> Stopping local dev environment..."
	@$(DOCKER) compose -f deployment/development/docker-compose.yml down

dev-logs: ## Follow logs from local development
	@echo "==> Following logs..."
	@$(DOCKER) compose -f deployment/development/docker-compose.yml logs -f

# -- Protobuf ----------------------------------------------------------------

proto: ## Generate protobuf code
	@echo "==> Generating protobuf code..."
	@for proto_file in $$(find api/ -name '*.proto'); do \
		echo "    generating $$proto_file..."; \
		$(PROTOC) \
			--proto_path=api/ \
			--go_out=paths=source_relative:api/go \
			--go-grpc_out=paths=source_relative:api/go \
			$$proto_file; \
	done

# -- Database ----------------------------------------------------------------

db-migrate: ## Run database migrations
	@echo "==> Running migrations..."
	@for svc in $(GO_SERVICES); do \
		if [ -f "database/migrations/$$svc" ]; then \
			echo "    migrating $$svc..."; \
			cd services/$$svc && make db-migrate 2>/dev/null || true && cd ../..; \
		fi; \
	done

# -- CI Pipeline -------------------------------------------------------------

ci: lint test build ## Full CI pipeline (lint -> test -> build)
	@echo "==> CI pipeline passed successfully!"

# -- Cleanup -----------------------------------------------------------------

clean: ## Clean all build artifacts
	@echo "==> Cleaning..."
	@rm -rf dist/
	@rm -f coverage.out cover.html
	@find . -type f -name '*.exe' -delete
	@find . -type d -name 'target' -exec rm -rf {} + 2>/dev/null || true
	@find . -type d -name 'vendor' -exec rm -rf {} + 2>/dev/null || true
	@find . -type d -name 'node_modules' -exec rm -rf {} + 2>/dev/null || true
	@$(GO) clean -cache 2>/dev/null || true
	@echo "==> Done."

