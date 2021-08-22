package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"gitlab.com/alxrem/terraform-provider-jsonnet/jsonnet"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: jsonnet.Provider})
}
