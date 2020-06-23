#############################
# Global vars
#############################
PROJECT_NAME := $(shell basename $(shell pwd))
PROJECT_VER  ?= $(shell git describe --tags --always --dirty | sed -e '/^v/s/^v\(.*\)$$/\1/g') # Strip leading 'v' if found
# Last released version (not dirty)
PROJECT_VER_TAGGED  ?= $(shell git describe --tags --always --abbrev=0 | sed -e '/^v/s/^v\(.*\)$$/\1/g') # Strip leading 'v' if found

SRCDIR       ?= .
GO            = go

# The root module (from go.mod)
PROJECT_MODULE  ?= $(shell $(GO) list -m)
GO_PKGS         ?= $(shell $(GO) list ./... | grep -v -e "/vendor/" -e "/example")

#############################
# Targets
#############################
all: build

# Humans running make:
build: git-hooks check-version clean lint lint-terraform test cover-report compile
	@echo "=== $(PROJECT_NAME) === [ build            ]: installing package..."
	@$(GO) install

# Build command for CI tooling
build-ci: check-version clean lint lint-terraform test compile-only

# All clean commands
clean: cover-clean compile-clean release-clean
	@echo "=== $(PROJECT_NAME) === [ clean            ]: removing go object files..."
	@$(GO) clean $(GO_PKGS)

# Import fragments
include build/compile.mk
include build/deps.mk
include build/docker.mk
include build/document.mk
include build/lint.mk
include build/release.mk
include build/terraform.mk
include build/test.mk
include build/util.mk

.PHONY: all build build-ci clean