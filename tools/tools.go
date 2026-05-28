//go:build tools
// +build tools

package tools

import (
	_ "github.com/bflad/tfproviderlint/cmd/tfproviderlint"
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/git-chglog/git-chglog/cmd/git-chglog"
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/masahiro331/go-commitlinter"
	_ "github.com/psampaz/go-mod-outdated"
	_ "github.com/stretchr/testify/assert"
	_ "golang.org/x/tools/cmd/godoc"
	_ "golang.org/x/tools/cmd/goimports"
	_ "gotest.tools/gotestsum"
)
