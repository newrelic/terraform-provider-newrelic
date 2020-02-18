RELEASE_SCRIPT ?= ./tools/release.sh

# Example usage: make release version=0.11.0
release:
	@echo "=== $(PROJECT_NAME) === [ release          ]: Generating release."
	$(RELEASE_SCRIPT) $(version)
