package main

import (
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/newrelic/terraform-provider-newrelic/v2/newrelic"
	"log"
)

var (
	// ProviderVersion is set during the release process to the release version of the binary.
	// See .goreleaser.yml for more details.
	ProviderVersion = "dev"

	// UserAgentServiceName can be set via -ldflags and used to customize
	// the provider's user agent string in request headers to facilitate
	// a better understanding of what additional services may "wrap" our
	// provider, such as Pulumi.
	UserAgentServiceName = ""
)

//func main() {
//	// We need to set the following package variables in the newrelic package
//	// to ensure they are properly set as part of the provider's user agent header.
//	newrelic.ProviderVersion = ProviderVersion
//	newrelic.UserAgentServiceName = UserAgentServiceName
//
//	var debugMode bool
//
//	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
//	flag.Parse()
//
//	if debugMode {
//		err := plugin.Debug(context.Background(), "registry.terraform.io/newrelic/newrelic",
//			&plugin.ServeOpts{
//				ProviderFunc: newrelic.Provider,
//			})
//		if err != nil {
//			log.Println(err.Error())
//		}
//	} else {
//		plugin.Serve(&plugin.ServeOpts{
//			ProviderFunc: newrelic.Provider})
//	}
//}

func main() {
	ctx := context.Background()

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	upgradedSdkServer, err := tf5to6server.UpgradeServer(
		ctx,
		newrelic.Provider().GRPCProvider, // Example terraform-plugin-sdk provider
	)

	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(newrelic.New()), // Example terraform-plugin-framework provider
		func() tfprotov6.ProviderServer {
			return upgradedSdkServer
		},
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/newrelic/newrelic",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
