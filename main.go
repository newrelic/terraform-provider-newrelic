package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-newrelic/newrelic"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: newrelic.Provider})
}
