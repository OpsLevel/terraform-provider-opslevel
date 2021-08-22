package main

import (
	"github.com/opslevel/terraform-provider-opslevel/opslevel"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	version string = "dev"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: opslevel.Provider})
}
