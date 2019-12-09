// +build tools

package main

import (
	_ "github.com/AlekSi/gocov-xml"
	_ "github.com/axw/gocov/gocov"
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/robertkrimen/godocdown/godocdown"
	_ "github.com/stretchr/testify/assert"
)
