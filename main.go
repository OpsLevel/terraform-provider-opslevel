package main

import (
	opslevel "github.com/opslevel/terraform-provider-opslevel/pkg/opslevel"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

var (
	version string = "dev"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: opslevel.Provider})
}
