package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/opslevel/terraform-provider-opslevel/opslevel"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var version string = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address:         "registry.terraform.io/OpsLevel/OpsLevel",
		Debug:           debug,
		ProtocolVersion: 6,
	}
	err := providerserver.Serve(context.Background(), opslevel.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
