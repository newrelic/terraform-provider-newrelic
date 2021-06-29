#
# Makefile Fragment for Compiling
#

GO              ?= go
BUILD_DIR       ?= ./bin/
PROJECT_MODULE  ?= $(shell $(GO) list -m)
# $b replaced by the binary name in the compile loop, -s/w remove debug symbols
LDFLAGS         ?= "-s -w -X main.version=$(PROJECT_VER) -X main.appName=newrelic -X $(PROJECT_MODULE)/internal/client.version=$(PROJECT_VER)"
SRCDIR          ?= .
COMPILE_OS      ?= darwin linux windows
BINARY ?= terraform-provider-newrelic

compile-clean:
	@echo "=== $(PROJECT_NAME) === [ compile-clean    ]: removing binaries..."
	@rm -rfv $(BUILD_DIR)/*

compile: deps compile-only

compile-all: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	@mkdir -p $(BUILD_DIR)/$(GOOS)
	for os in $(COMPILE_OS); do \
		echo "=== $(PROJECT_NAME) === [ compile          ]:     $(BUILD_DIR)$$os/$(BINARY)"; \
		BUILD_FILES=`find $(SRCDIR)/newrelic -type f -name "*.go"` ; \
		GOOS=$$os $(GO) build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$$os/$(BINARY) $$BUILD_FILES ; \
	done

compile-only: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	@mkdir -p $(BUILD_DIR)/$(GOOS)
	echo "=== $(PROJECT_NAME) === [ compile          ]:     $(BUILD_DIR)$(GOOS)/$(BINARY)"; \
	BUILD_FILES=`find $(SRCDIR)/newrelic -type f -name "*.go"` ; \
	GOOS=$(GOOS) $(GO) build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(GOOS)/$(BINARY) $$BUILD_FILES ; \

# Override GOOS for these specific targets
compile-darwin: GOOS=darwin
compile-darwin: deps-only compile-only

compile-linux: GOOS=linux
compile-linux: deps-only compile-only

compile-windows: GOOS=windows
compile-windows: deps-only compile-only

.PHONY: clean-compile compile compile-darwin compile-linux compile-only compile-windows
