#
# Makefile fragment for installing deps
#

GO               ?= go
VENDOR_CMD       ?= ${GO} mod tidy

deps: tools deps-only

deps-only:
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Installing package dependencies required by the project..."
	@$(VENDOR_CMD)

.PHONY: deps deps-only tools tools-update tools-outdated
