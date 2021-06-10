package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

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

	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/newrelic/newrelic",
			&plugin.ServeOpts{
				ProviderFunc: newrelic.Provider,
			})
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: newrelic.Provider})
	}
}
