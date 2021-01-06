#
# Makefile fragment for displaying auto-generated documentation
#

GODOC       ?= godoc
GODOC_HTTP  ?= "localhost:6060"

CHANGELOG_CMD      ?= git-chglog
CHANGELOG_FILE     ?= CHANGELOG.md
RELEASE_NOTES_FILE ?= relnotes.md

changelog: tools
	@echo "=== $(PROJECT_NAME) === [ changelog        ]: Generating changelog..."
	@$(CHANGELOG_CMD) --silent -o $(CHANGELOG_FILE)

docs: tools
	@echo "=== $(PROJECT_NAME) === [ docs             ]: Starting godoc server..."
	@echo "=== $(PROJECT_NAME) === [ docs             ]:"
	@echo "=== $(PROJECT_NAME) === [ docs             ]: NOTE: This only works if this codebase is in your GOPATH!"
	@echo "=== $(PROJECT_NAME) === [ docs             ]:    godoc issue: https://github.com/golang/go/issues/26827"
	@echo "=== $(PROJECT_NAME) === [ docs             ]:"
	@echo "=== $(PROJECT_NAME) === [ docs             ]: Module Docs: http://$(GODOC_HTTP)/pkg/$(PROJECT_MODULE)"
	@$(GODOC) -http=$(GODOC_HTTP)

release-notes: tools
	@echo "=== $(PROJECT_NAME) === [ release-notes    ]: Generating release notes..."
	@mkdir -p $(SRCDIR)/tmp
	@$(CHANGELOG_CMD) --silent -o $(SRCDIR)/tmp/$(RELEASE_NOTES_FILE) v$(PROJECT_VER_TAGGED)

.PHONY: docs changelog release-notes
