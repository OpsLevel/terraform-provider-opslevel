package main

import (
	"context"
	"flag"
	"log"

	"github.com/opslevel/terraform-provider-opslevel/internal/provider"
	// _ "github.com/opslevel/terraform-provider-opslevel/opslevel"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	// _ "github.com/hashicorp/terraform-plugin-sdk/plugin"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

//	func main() {
//		plugin.Serve(&plugin.ServeOpts{ProviderFunc: opslevel.Provider})
//	}
var version string = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like dlv")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address:         "registry.terraform.io/OpsLevel/opslevel",
		Debug:           debug,
		ProtocolVersion: 6,
	}
	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err)
	}
}
