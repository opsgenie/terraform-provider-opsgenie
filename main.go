package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/opsgenie/terraform-provider-opsgenie/opsgenie"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opsgenie.Provider})
}
