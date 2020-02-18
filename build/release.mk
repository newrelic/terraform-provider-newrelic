RELEASE_SCRIPT_PATH ?= ./tools/release.sh

# Example usage: make release version=0.11.0
release:
	@echo "=== $(PROJECT_NAME) === [ release          ]: Generating release."
	./tools/release.sh $(version)
