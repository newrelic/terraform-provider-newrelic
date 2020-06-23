package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/newrelic/terraform-provider-newrelic/newrelic"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: newrelic.Provider})
}
