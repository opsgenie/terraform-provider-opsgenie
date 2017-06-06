package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-opsgenie/opsgenie"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opsgenie.Provider})
}
