package main

import (
	opslevel "terraform-provider-opslevel/pkg/opslevel"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opslevel.Provider})
}
