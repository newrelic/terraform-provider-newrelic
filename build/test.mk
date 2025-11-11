#
# Makefile fragment for Testing
#

GO           ?= go
GOLINTER     ?= golangci-lint
MISSPELL     ?= misspell
GOFMT        ?= gofmt
TEST_RUNNER  ?= gotestsum

COVERAGE_DIR ?= ./coverage/
COVERMODE    ?= atomic
SRCDIR       ?= .
GO_PKGS      ?= $(shell $(GO) list ./... | grep -v -e "/example")
FILES        ?= $(shell find $(SRCDIR) -type f | grep -v -e '.git/')

PROJECT_MODULE ?= $(shell $(GO) list -m)

# Set a few vars and run the test suite
LDFLAGS_TEST ?= "-X=$(PROJECT_MODULE)/version.ProviderVersion=acc"

TEST_MAPPER_SCRIPT ?= ./integration_test_utilities/integration_test_utilities_tag_mapper.sh
TEST_MAPPER_SCRIPT_OUTPUT ?= $(shell bash $(TEST_MAPPER_SCRIPT) | tail -n 1)

test: test-only
test-only: test-unit test-integration-all

test-unit: tools
	@echo "=== $(PROJECT_NAME) === [ test-unit        ]: running unit tests..."
	@mkdir -p $(COVERAGE_DIR)
	@$(TEST_RUNNER) -f testname --junitfile $(COVERAGE_DIR)/unit.xml --packages "$(GO_PKGS)" --jsonfile $(COVERAGE_DIR)/unit.report \
		-- -v -parallel 10 -tags=unit $(TEST_ARGS) -covermode=$(COVERMODE) -coverprofile $(COVERAGE_DIR)/unit.tmp

test-integration: tools
	@echo "=== $(PROJECT_NAME) === [ test-integration ]: running integration tests..."
	@mkdir -p $(COVERAGE_DIR)
	@if [ "$(TEST_MAPPER_SCRIPT_OUTPUT)" = "NONE" ]; then \
		echo "=== $(PROJECT_NAME) === [ test-integration ]: no changes found in code affecting tests, skipping integration tests..."; \
	else \
		echo "=== $(PROJECT_NAME) === [ test-integration ]: running integration tests with tags $(TEST_MAPPER_SCRIPT_OUTPUT)"; \
		TF_ACC=1 $(TEST_RUNNER) -f testname --junitfile $(COVERAGE_DIR)/integration.xml --rerun-fails=4 --packages "$(GO_PKGS)" --jsonfile $(COVERAGE_DIR)/integration.report \
			-- -v -parallel 14 $(TEST_MAPPER_SCRIPT_OUTPUT) $(TEST_ARGS) -covermode=$(COVERMODE) -coverprofile $(COVERAGE_DIR)/integration.tmp \
			-timeout 120m -ldflags=$(LDFLAGS_TEST); \
	fi

test-integration-all: tools
	@echo "=== $(PROJECT_NAME) === [ test-integration-all ]: running all integration tests..."
	@mkdir -p $(COVERAGE_DIR)
	@if [ "$(TEST_MAPPER_SCRIPT_OUTPUT)" = "-tags=integration" ]; then \
		echo "=== $(PROJECT_NAME) === [ test-integration-all ]: skipping the test-integration-all suite, as the test-integration suite is currently running all tests..."; \
	else \
		TF_ACC=1 $(TEST_RUNNER) -f testname --junitfile $(COVERAGE_DIR)/integration.xml --rerun-fails=4 --packages "$(GO_PKGS)" --jsonfile $(COVERAGE_DIR)/integration.report \
			-- -v -parallel 14 -tags=integration $(TEST_ARGS) -covermode=$(COVERMODE) -coverprofile $(COVERAGE_DIR)/integration.tmp \
			   -timeout 120m -ldflags=$(LDFLAGS_TEST); \
	fi

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

.PHONY: test test-only test-unit test-integration test-integration-all cover-report cover-view
