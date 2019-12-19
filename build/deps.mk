#
# Makefile fragment for installing deps
#

GO           ?= go
VENDOR_CMD   ?= ${GO} mod tidy

# These should be mirrored in /tools.go to keep versions consistent
GOTOOLS      += github.com/client9/misspell/cmd/misspell


tools: check-version
	@echo "=== $(PROJECT_NAME) === [ tools            ]: Installing tools required by the project..."
	@$(GO) install $(GOTOOLS)

tools-update: check-version
	@echo "=== $(PROJECT_NAME) === [ tools-update     ]: Updating tools required by the project..."
	@$(GO) get -u $(GOTOOLS)

deps: tools deps-only

deps-only:
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Installing package dependencies required by the project..."
	@$(VENDOR_CMD)

.PHONY: deps deps-only tools tools-update
