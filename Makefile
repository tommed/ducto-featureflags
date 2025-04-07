# ----------------------
# Configuration
# ----------------------

COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html
GO=go
LINTER=golangci-lint
LINTER_REMOTE=github.com/golangci/golangci-lint/cmd/golangci-lint@latest
LINTER_OPTS=--timeout=2m

# ----------------------
# General Targets
# ----------------------

.PHONY: ci test-full lint lint-install coverage clean build-all ducto-flags-macos ducto-flags-windows

check: lint test-full coverage

ci: check build-all

build-all: ducto-flags-macos ducto-flags-windows

clean:
	@rm -f $(COVERAGE_OUT) $(COVERAGE_HTML) ducto-flags*

# ----------------------
# Linting
# ----------------------

lint:
	@echo "==> Running linter"
	$(LINTER) run $(LINTER_OPTS)

lint-install:
	go install $(LINTER_REMOTE)

# ----------------------
# Testing
# ----------------------

test-full:
	@echo "==> Running all tests"
	$(GO) test -coverpkg=./... -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v ./...
	$(GO) tool cover -func=$(COVERAGE_OUT)

coverage:
	@echo "==> Generating coverage HTML report"
	$(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

# ----------------------
# Builds
# ----------------------

ducto-flags-macos:
	@echo "==> Building macOS CLI"
	GOOS=darwin GOARCH=arm64 $(GO) build -o ducto-flags ./cmd/ducto-flags

ducto-flags-windows:
	@echo "==> Building Windows CLI"
	GOOS=windows GOARCH=amd64 $(GO) build -o ducto-flags.exe ./cmd/ducto-flags
