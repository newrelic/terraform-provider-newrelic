#
# Makefile fragment for Testing
#

GO           ?= go
GOLINTER     ?= golangci-lint
MISSPELL     ?= misspell
GOFMT        ?= gofmt

COVERAGE_DIR ?= ./coverage/
COVERMODE    ?= atomic
SRCDIR       ?= .
GO_PKGS      ?= $(shell $(GO) list ./... | grep -v -e "/vendor/" -e "/example")
FILES        ?= $(shell find $(SRCDIR) -type f | grep -v -e '.git/' -e '/vendor/')

PROJECT_MODULE ?= $(shell $(GO) list -m)

LDFLAGS_TEST ?= '-X $(PROJECT_MODULE)/internal/version.GitTag=$(PROJECT_VER_TAGGED)'
GOTOOLS += github.com/stretchr/testify/assert

test: test-only
test-only: test-unit test-integration

test-unit:
	@echo "=== $(PROJECT_NAME) === [ test-unit        ]: running unit tests..."
	@mkdir -p $(COVERAGE_DIR)
	@$(GO) test -v -ldflags=$(LDFLAGS_TEST) -parallel 4 $(TEST_ARGS) -tags unit -covermode=$(COVERMODE) -coverprofile $(COVERAGE_DIR)/unit.tmp $(GO_PKGS)

test-integration:
	@echo "=== $(PROJECT_NAME) === [ test-integration ]: running integration tests..."
	@mkdir -p $(COVERAGE_DIR)
	@$(GO) test -v -ldflags=$(LDFLAGS_TEST) -parallel 4 $(TEST_ARGS) -tags integration -covermode=$(COVERMODE) -coverprofile $(COVERAGE_DIR)/integration.tmp $(GO_PKGS)

test-one:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-one TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	@echo "=== $(PROJECT_NAME) === [ test-one         ]: running test: $(TEST) ..."
	@$(GO) test -v -ldflags=$(LDFLAGS_TEST) -c $(TEST) $(TEST_ARGS)


#
# Coverage
#
cover-clean:
	@echo "=== $(PROJECT_NAME) === [ cover-clean      ]: removing coverage files..."
	@rm -rfv $(COVERAGE_DIR)/*

cover-report:
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]: generating coverage results..."
	@mkdir -p $(COVERAGE_DIR)
	@echo 'mode: $(COVERMODE)' > $(COVERAGE_DIR)/coverage.out
	@cat $(COVERAGE_DIR)/*.tmp | grep -v 'mode: $(COVERMODE)' >> $(COVERAGE_DIR)/coverage.out || true
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]:     $(COVERAGE_DIR)coverage.html"

cover-view: cover-report
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out

.PHONY: test test-only test-unit test-integration cover-report cover-view
