package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"gitlab.com/peikk0/terraform-provider-jsonnet/jsonnet"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: jsonnet.Provider})
}
