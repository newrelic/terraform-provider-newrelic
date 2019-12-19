#
# Makefile fragment for displaying auto-generated documentation
#

GODOC       ?= godoc
GODOC_HTTP  ?= "localhost:6060"

GO_MODULE   ?= $(shell go list -m)

docs:
	@echo "=== $(PROJECT_NAME) === [ docs             ]: Starting godoc server..."
	@echo "=== $(PROJECT_NAME) === [ docs             ]: Module Docs: http://$(GODOC_HTTP)/pkg/$(GO_MODULE)"
	@$(GODOC) -http=$(GODOC_HTTP)

.PHONY: docs
