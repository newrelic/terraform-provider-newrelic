#
# Makefile fragment for Linting
#

GO               ?= go
GOFMT            ?= gofmt
GOIMPORTS        ?= goimports
GOLINTER         ?= golangci-lint
GO_MOD_OUTDATED  ?= go-mod-outdated
MISSPELL         ?= misspell

COMMIT_LINT_CMD   ?= go-gitlint
COMMIT_LINT_REGEX ?= "(chore|docs|feat|fix|refactor|tests?)(\([^\)]+\))?: .*"
COMMIT_LINT_START ?= "2020-06-17"


EXCLUDEDIR      ?= .git
SRCDIR          ?= .
GO_PKGS         ?= $(shell ${GO} list ./... | grep -v -e "/vendor/" -e "/example")
FILES           ?= $(shell find ${SRCDIR} -type f | grep -v -e '.git/' -e '/vendor/')
GO_FILES        ?= $(shell find $(SRCDIR) -type f -name "*.go" | grep -v -e ".git/" -e '/vendor/' -e '/example/')
PROJECT_MODULE  ?= $(shell $(GO) list -m)


GOTOOLS += github.com/client9/misspell/cmd/misspell \
           github.com/llorllale/go-gitlint/cmd/go-gitlint \
           github.com/psampaz/go-mod-outdated \
           github.com/golangci/golangci-lint/cmd/golangci-lint \
           golang.org/x/tools/cmd/goimports


lint: deps spell-check gofmt lint-commit golangci goimports outdated
lint-fix: deps spell-check-fix gofmt-fix goimports

#
# Check spelling on all the files, not just source code
#
spell-check: deps
	@echo "=== $(PROJECT_NAME) === [ spell-check      ]: Checking for spelling mistakes with $(MISSPELL)..."
	@$(MISSPELL) -source text $(FILES)

spell-check-fix: deps
	@echo "=== $(PROJECT_NAME) === [ spell-check-fix  ]: Fixing spelling mistakes with $(MISSPELL)..."
	@$(MISSPELL) -source text -w $(FILES)

gofmt: deps
	@echo "=== $(PROJECT_NAME) === [ gofmt            ]: Checking file format with $(GOFMT)..."
	@find . -path "$(EXCLUDEDIR)" -prune -print0 | xargs -0 $(GOFMT) -e -l -s -d ${SRCDIR}

gofmt-fix: deps
	@echo "=== $(PROJECT_NAME) === [ gofmt-fix        ]: Fixing file format with $(GOFMT)..."
	@find . -path "$(EXCLUDEDIR)" -prune -print0 | xargs -0 $(GOFMT) -e -l -s -w ${SRCDIR}

goimports: deps
	@echo "=== $(PROJECT_NAME) === [ goimports        ]: Checking imports with $(GOIMPORTS)..."
	@$(GOIMPORTS) -l -w -local $(PROJECT_MODULE) $(GO_FILES)

lint-commit: deps
	@echo "=== $(PROJECT_NAME) === [ lint-commit      ]: Checking that commit messages are properly formatted ($(COMMIT_LINT_CMD))..."
	@$(COMMIT_LINT_CMD) --since=$(COMMIT_LINT_START) --subject-minlen=10 --subject-maxlen=120 --subject-regex=$(COMMIT_LINT_REGEX)

golangci: deps
	@echo "=== $(PROJECT_NAME) === [ golangci-lint    ]: Linting using $(GOLINTER) ($(COMMIT_LINT_CMD))..."
	@$(GOLINTER) run

outdated: deps
	@echo "=== $(PROJECT_NAME) === [ outdated         ]: Finding outdated deps with $(GO_MOD_OUTDATED)..."
	@$(GO) list -u -m -json all | $(GO_MOD_OUTDATED) -direct -update

.PHONY: lint spell-check spell-check-fix gofmt gofmt-fix lint-fix lint-commit outdated goimports
