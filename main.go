package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"gitlab.com/alxrem/terraform-provider-jsonnet/private/provider"
	"log"
)

func main() {
	if err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/alxrem/jsonnet",
	}); err != nil {
		log.Fatal(err)
	}
}
