// +build tools

package main

import (
	// build/test.mk
	_ "github.com/stretchr/testify/assert"

	// build/lint.mk
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/fzipp/gocyclo"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/gordonklaus/ineffassign"
	_ "github.com/llorllale/go-gitlint/cmd/go-gitlint"
	_ "github.com/remyoudompheng/go-misc/deadcode"
	_ "github.com/timakin/bodyclose"
	_ "golang.org/x/lint/golint"

	// build/document.mk
	_ "github.com/git-chglog/git-chglog/cmd/git-chglog"
	_ "golang.org/x/tools/cmd/godoc"
)
