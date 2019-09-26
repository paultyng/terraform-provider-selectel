package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-selectel/selectel"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: selectel.Provider,
	})
}
