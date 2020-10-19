package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/newrelic/terraform-provider-newrelic/v2/newrelic"
)

var (
	// ProviderVersion is set during the release process to the release version of the binary.
	// See .goreleaser.yml for more details.
	ProviderVersion = "dev"
)

func main() {
	// We need to set the ProviderVersion variable in the newrelic package
	// to ensure it gets properly set as part of the user agent header.
	newrelic.ProviderVersion = ProviderVersion

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: newrelic.Provider})
}
