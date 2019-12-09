PROJECT_NAME := $(shell basename $(shell pwd))
PROJECT_VER  := $(shell git describe --tags --always --dirty | sed -e '/^v/s/^v\(.*\)$$/\1/g') # Strip leading 'v' if found
GO_PKGS      := $(shell go list ./... | grep -v -e "/vendor/" -e "/example")
NATIVEOS     := $(shell go version | awk -F '[ /]' '{print $$4}')
NATIVEARCH   := $(shell go version | awk -F '[ /]' '{print $$5}')
SRCDIR       ?= .
BUILD_DIR    := ./bin/
COVERAGE_DIR := ./coverage/
COVERMODE     = atomic

GO            = go
GODOC         = godocdown
DOC_DIR       = ./docs/
GOLINTER      = golangci-lint

# Main API Entry point
PACKAGES = ${SRCDIR}/newrelic

# Determine packages by looking into pkg/*
ifneq ("$(wildcard ${SRCDIR}/pkg/*)","")
	PACKAGES += $(wildcard ${SRCDIR}/pkg/*)
endif
ifneq ("$(wildcard ${SRCDIR}/internal/*)","")
	PACKAGES += $(wildcard ${SRCDIR}/internal/*)
endif

# Determine commands by looking into cmd/*
COMMANDS = $(wildcard ${SRCDIR}/cmd/*)

GO_FILES := $(shell find $(COMMANDS) $(PACKAGES) -type f -name "*.go")

# Determine binary names by stripping out the dir names
BINS=$(foreach cmd,${COMMANDS},$(notdir ${cmd}))

LDFLAGS='-X main.Version=$(PROJECT_VER)'

all: build

# Humans running make:
build: check-version clean lint test cover-report compile document

# Build command for CI tooling
build-ci: check-version clean lint test compile-only

clean:
	@echo "=== $(PROJECT_NAME) === [ clean            ]: removing binaries and coverage file..."
	@rm -rfv $(BUILD_DIR)/* $(COVERAGE_DIR)/*

# Import fragments
include build/compile.mk
include build/deps.mk
include build/document.mk
include build/testing.mk
include build/util.mk

.PHONY: all build build-ci package
