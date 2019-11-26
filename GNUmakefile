PKG_NAME      = newrelic
TEST         ?= $$(go list ./... |grep -v 'vendor')
GO_PKGS      := $(shell go list ./... | grep -v -e "/vendor/" -e "/example")
GOFMT_FILES  ?= $$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO  = github.com/hashicorp/terraform-website

# Commands
GO        := go
GO_LINTER := golangci-lint
TF_LINTER := tfproviderlint

BUILD_DIR    := ./bin/
COVERAGE_DIR := ./coverage/
COVERMODE     = atomic

default: clean build lint tidy vet test cover-report

build: fmtcheck
	@$(GO) install

clean:
	@echo "=== $(PKG_NAME) === [ clean            ]: removing binaries and coverage file..."
	@$(GO) clean $(GO_PKGS)
	@rm -rfv $(BUILD_DIR)/* $(COVERAGE_DIR)/*

lint: tools
	@echo "=== $(PKG_NAME) === [ lint             ]: running linters..."
	@echo "=== $(PKG_NAME) === [ lint             ]:     $(GO_LINTER) ..."
	@$(GO_LINTER) run ./$(PKG_NAME)
	@echo "=== $(PKG_NAME) === [ lint             ]:     $(TF_LINTER) ..."
	@$(TF_LINTER) \
		-c 1 \
		-AT001 \
		-AT002 \
		-S001 \
		-S002 \
		-S003 \
		-S004 \
		-S005 \
		-S007 \
		-S008 \
		-S009 \
		-S010 \
		-S011 \
		-S012 \
		-S013 \
		-S014 \
		-S015 \
		-S016 \
		-S017 \
		-S019 \
		./$(PKG_NAME)

test: clean fmtcheck
	@echo "=== $(PKG_NAME) === [ test             ]: running test suite..."
	@mkdir -p $(COVERAGE_DIR)
	@$(GO) test -i $(TEST) || exit 1
	@for d in $(GO_PKGS); do \
		pkg=`basename $$d` ;\
		$(GO) test $(TESTARGS) -timeout=30s -parallel=4 -covermode=$(COVERMODE) -coverprofile $(COVERAGE_DIR)/$$pkg.tmp $$d || exit 1; \
	done

testacc: fmtcheck
	TF_ACC=1 $(GO) test $(TEST) \
				 -v $(TESTARGS) \
				 -timeout 120m \
				 -ldflags="-X=github.com/terraform-providers/terraform-provider-newrelic/version.ProviderVersion=acc"

vet:
	@echo "=== $(PKG_NAME) === [ vet              ]: running go vet..."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	@echo "=== $(PKG_NAME) === [ format           ]: Running gofmt..."
	gofmt -w -s $(GOFMT_FILES)

tools:
	@echo "=== $(PKG_NAME) === [ tools            ]: installing required tooling..."
	@GO111MODULE=on $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint
	@GO111MODULE=on $(GO) install github.com/bflad/tfproviderlint/cmd/tfproviderlint
	@GO111MODULE=on $(GO) install github.com/client9/misspell/cmd/misspell

tidy:
	@echo "=== $(PKG_NAME) === [ tidy             ]: tidying modules..."
	@$(GO) mod tidy

fmtcheck:
	@echo "=== $(PKG_NAME) === [ fmtcheck         ]: Checking that code complies with gofmt requirements..."
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

cover-report: test
	@echo "=== $(PKG_NAME) === [ cover-report     ]: generating test coverage..."
	@echo 'mode: $(COVERMODE)' > $(COVERAGE_DIR)/coverage.out
	@cat $(COVERAGE_DIR)/*.tmp | grep -v 'mode: $(COVERMODE)' >> $(COVERAGE_DIR)/coverage.out || true
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "=== $(PKG_NAME) === [ cover-report     ]:     $(COVERAGE_DIR)coverage.html"

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	@echo "=== $(PKG_NAME) === [ website-lint     ]: linting website..."
	@misspell -error -source=text website/

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website website-test tools lint cover-report

