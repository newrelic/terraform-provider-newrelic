package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/newrelic/terraform-provider-newrelic/v2/newrelic"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: newrelic.Provider})
}
