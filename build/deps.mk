#
# Makefile fragment for installing deps
#

# These should be mirrored in /tools.go to keep versions consistent
GOTOOLS       = github.com/axw/gocov/gocov \
                github.com/AlekSi/gocov-xml \
                github.com/stretchr/testify/assert \
                github.com/robertkrimen/godocdown/godocdown \
                github.com/golangci/golangci-lint/cmd/golangci-lint

VENDOR_CMD = ${GO} mod tidy

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

