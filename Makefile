#############################
# Global vars
#############################
PROJECT_NAME := $(shell basename $(shell pwd))
PROJECT_VER  := $(shell git describe --tags --always --dirty | sed -e '/^v/s/^v\(.*\)$$/\1/g') # Strip leading 'v' if found
SRCDIR       ?= .
GO            = go

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



#############################
# Targets
#############################
all: build

# Humans running make:
build: check-version clean lint test cover-report compile

# Build command for CI tooling
build-ci: check-version clean lint test compile-only

# All clean commands
clean: clean-cover clean-compile

# Import fragments
include build/compile.mk
include build/deps.mk
include build/document.mk
include build/testing.mk
include build/util.mk

.PHONY: all build build-ci clean
