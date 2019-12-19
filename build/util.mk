#
# Makefile fragment for utility items
#

NATIVEOS    ?= $(shell go version | awk -F '[ /]' '{print $$4}')
NATIVEARCH  ?= $(shell go version | awk -F '[ /]' '{print $$5}')


check-version:
ifdef GOOS
ifneq "$(GOOS)" "$(NATIVEOS)"
	$(error GOOS is not $(NATIVEOS). Cross-compiling is only allowed for 'clean', 'deps-only' and 'compile-only' targets)
endif
else
GOOS = ${NATIVEOS}
endif
ifdef GOARCH
ifneq "$(GOARCH)" "$(NATIVEARCH)"
	$(error GOARCH variable is not $(NATIVEARCH). Cross-compiling is only allowed for 'clean', 'deps-only' and 'compile-only' targets)
endif
else
GOARCH = ${NATIVEARCH}
endif

.PHONY: check-version
