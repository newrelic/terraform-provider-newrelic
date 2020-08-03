// +build tools

package main

import (
	_ "github.com/bflad/tfproviderlint/cmd/tfproviderlint"
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/newrelic-forks/git-chglog/cmd/git-chglog"
	_ "gotest.tools/gotestsum"
)
