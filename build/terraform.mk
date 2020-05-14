PKG_NAME      ?= newrelic
GO_PKGS       ?= $(shell ${GO} list ./... | grep -v -e "/vendor/" -e "/example")
WEBSITE_REPO  ?= github.com/hashicorp/terraform-website
TF_LINTER     ?= tfproviderlint
GOTOOLS       += github.com/bflad/tfproviderlint/cmd/tfproviderlint

lint-terraform: deps
	@echo "=== $(PROJECT_NAME) === [ lint-terraform   ]: running terraform linter $(TF_LINTER) ..."
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

# Set a few vars and run the test suite
testacc: TF_ACC = 1
testacc: TEST_ARGS = "-timeout 120m"
testacc: LDFLAGS_TEST = "-X=github.com/terraform-providers/terraform-provider-newrelic/version.ProviderVersion=acc"
testacc: test-only


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

.PHONY: lint-terraform testacc website website-lint website-test
