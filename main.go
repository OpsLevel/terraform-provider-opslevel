package main

import (
	opslevel "terraform-provider-opslevel/pkg/opslevel"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"

	// goreleaser can also pass the specific commit if you want
	// commit  string = ""
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opslevel.Provider})
}
