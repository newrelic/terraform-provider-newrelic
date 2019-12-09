#
# Makefile Fragment for Compiling
#

compile: deps compile-only
compile-all: compile-linux compile-darwin compile-windows
build-all: compile-linux compile-darwin compile-windows

compile-only: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	@mkdir -p $(BUILD_DIR)/$(GOOS)
	@for b in $(BINS); do \
		echo "=== $(PROJECT_NAME) === [ compile          ]:     $(BUILD_DIR)$(GOOS)/$$b"; \
		BUILD_FILES=`find $(SRCDIR)/cmd/$$b -type f -name "*.go"` ; \
		$(GO) build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(GOOS)/$$b $$BUILD_FILES ; \
	done

build-linux: compile-linux
compile-linux: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile-linux    ]: building commands:"
	@mkdir -p $(BUILD_DIR)/linux
	@for b in $(BINS); do \
		OUTPUT_FILE="$(BUILD_DIR)linux/$$b" ; \
		echo "=== $(PROJECT_NAME) === [ compile-linux    ]:     $$OUTPUT_FILE"; \
		BUILD_FILES=`find $(SRCDIR)/cmd/$$b -type f -name "*.go"` ; \
		GOOS=linux $(GO) build -ldflags=$(LDFLAGS) -o $$OUTPUT_FILE $$BUILD_FILES ; \
	done

build-darwin: compile-darwin
compile-darwin: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile-darwin   ]: building commands:"
	@mkdir -p $(BUILD_DIR)/darwin
	@for b in $(BINS); do \
		OUTPUT_FILE="$(BUILD_DIR)darwin/$$b" ; \
		echo "=== $(PROJECT_NAME) === [ compile-darwin   ]:     $$OUTPUT_FILE"; \
		BUILD_FILES=`find $(SRCDIR)/cmd/$$b -type f -name "*.go"` ; \
		GOOS=darwin $(GO) build -ldflags=$(LDFLAGS) -o $$OUTPUT_FILE $$BUILD_FILES ; \
	done

build-windows: compile-windows
compile-windows: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile-windows  ]: building commands:"
	@mkdir -p $(BUILD_DIR)/windows
	@for b in $(BINS); do \
		OUTPUT_FILE="$(BUILD_DIR)windows/$$b.exe" ; \
		echo "=== $(PROJECT_NAME) === [ compile-windows  ]:     $$OUTPUT_FILE"; \
		BUILD_FILES=`find $(SRCDIR)/cmd/$$b -type f -name "*.go"` ; \
		GOOS=windows $(GO) build -ldflags=$(LDFLAGS) -o $$OUTPUT_FILE $$BUILD_FILES ; \
	done

