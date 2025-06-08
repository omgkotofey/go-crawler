LINTER_TIMEOUT=10m
GOLANGCI_VERSION=latest
GOLANGCI_BIN=$(HOME)/go/bin/golangci-lint

lint: install-linter
	@echo "Running linter..."
	@$(GOLANGCI_BIN) run --timeout $(LINTER_TIMEOUT)

install-linter:
	@if [ ! -x "$(GOLANGCI_BIN)" ]; then \
		echo "Installing golangci-lint $(LINTER_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(HOME)/go/bin/ $(LINTER_VERSION); \
	else \
		echo "golangci-lint already installed."; \
	fi