#
# Makefile Fragment for Compiling
#

GO         ?= go
BUILD_DIR  ?= ./bin/
LDFLAGS    ?= '-X main.Version=$(PROJECT_VER)'
SRCDIR     ?= .

clean-compile:
	@echo "=== $(PROJECT_NAME) === [ clean-compile    ]: removing binaries..."
	@rm -rfv $(BUILD_DIR)/*

compile: deps compile-only

compile-only: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	@mkdir -p $(BUILD_DIR)/$(GOOS)
	@for b in $(BINS); do \
		echo "=== $(PROJECT_NAME) === [ compile          ]:     $(BUILD_DIR)$(GOOS)/$$b"; \
		BUILD_FILES=`find $(SRCDIR)/cmd/$$b -type f -name "*.go"` ; \
		GOOS=$(GOOS) $(GO) build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(GOOS)/$$b $$BUILD_FILES ; \
	done

# Override GOOS for these specific targets
compile-darwin: GOOS=darwin
compile-darwin: deps-only compile-only

compile-linux: GOOS=linux
compile-linux: deps-only compile-only

compile-windows: GOOS=windows
compile-windows: deps-only compile-only


.PHONY: clean-compile compile compile-darwin compile-linux compile-only compile-windows
