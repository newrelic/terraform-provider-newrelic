#
# Makefile fragment for Testing
#

GO           ?= go
GOLINTER     ?= golangci-lint
MISSPELL     ?= misspell
GOFMT        ?= gofmt

SRCDIR       ?= .
GO_PKGS      ?= $(shell ${GO} list ./... | grep -v -e "/vendor/" -e "/example")
FILES        ?= $(shell find ${SRCDIR} -type f | grep -v -e '.git/' -e '/vendor/')

GOTOOLS += github.com/client9/misspell/cmd/misspell \
           github.com/fzipp/gocyclo \
           github.com/golangci/golangci-lint/cmd/golangci-lint \
           github.com/gordonklaus/ineffassign \
           github.com/remyoudompheng/go-misc/deadcode \
           github.com/timakin/bodyclose \
           golang.org/x/lint/golint


#lint: deps spell-check
#	@echo "=== $(PROJECT_NAME) === [ lint             ]: Validating source code running $(GOLINTER)..."
#	@$(GOLINTER) -v run ./...

lint: deps spell-check gofmt govet golint ineffassign gocyclo deadcode bodyclose
lint-fix: deps spell-check-fix gofmt-fix

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
	@$(GOFMT) -e -l -s -d ${SRCDIR}

gofmt-fix: deps
	@echo "=== $(PROJECT_NAME) === [ gofmt-fix        ]: Fixing file format with $(GOFMT)..."
	@$(GOFMT) -e -l -s -w ${SRCDIR}

govet: deps
	@echo "=== $(PROJECT_NAME) === [ govet            ]: Checking file format with $(GO) vet..."
	@$(GO) vet $(GO_PKGS)

golint: deps
	@echo "=== $(PROJECT_NAME) === [ golint           ]: Checking source files with golint..."
	@golint -set_exit_status $(SRCDIR)/...

ineffassign: deps
	@echo "=== $(PROJECT_NAME) === [ ineffassign      ]: Checking for ineffectual assignments..."
	@ineffassign $(SRCDIR)

gocyclo: deps
	@echo "=== $(PROJECT_NAME) === [ gocyclo          ]: Calculating cyclomatic complexities of functions (gocyclo)..."
	@gocyclo -over 20 $(SRCDIR)

deadcode: deps
	@echo "=== $(PROJECT_NAME) === [ deadcode         ]: Checking for unused code (deadcode)..."
	@deadcode $(GO_PKGS)

bodyclose: deps
	@echo "=== $(PROJECT_NAME) === [ bodyclose        ]: Checking that http response bodies are closed (bodyclose)..."
	@$(GO) vet -vettool=$(shell which bodyclose) $(SRCDIR)/...

.PHONY: lint spell-check spell-check-fix gofmt gofmt-fix lint-fix
