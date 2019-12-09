#
# Makefile fragment for auto-documenting Golang projects
#

GODOC        ?= godocdown

documentation: document

document:
	@echo "=== $(PROJECT_NAME) === [ documentation    ]: Generating Godoc in Markdown..."
	@for p in $(PACKAGES); do \
		echo "=== $(PROJECT_NAME) === [ documentation    ]:     $$p"; \
		mkdir -p $(DOC_DIR)/$$p ; \
		$(GODOC) $$p > $(DOC_DIR)/$$p/README.md ; \
	done
	@for c in $(COMMANDS); do \
		echo "=== $(PROJECT_NAME) === [ documentation    ]:     $$c"; \
		mkdir -p $(DOC_DIR)/$$c ; \
		$(GODOC) $$c > $(DOC_DIR)/$$c/README.md ; \
	done
