#
# Makefile fragment for installing deps
#

GO           ?= go

deps: tools deps-only

deps-only:
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Installing package dependencies required by the project..."
	#@$(GO) mod tidy
	@$(GO) mod download

tools: check-version
	@echo "=== $(PROJECT_NAME) === [ tools            ]: Installing tools required by the project..."
	@$(GO) install $(GOTOOLS)

tools-update: check-version
	@echo "=== $(PROJECT_NAME) === [ tools-update     ]: Updating tools required by the project..."
	@$(GO) get -u $(GOTOOLS)


.PHONY: deps deps-only tools tools-update
