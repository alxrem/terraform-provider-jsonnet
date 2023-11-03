// Copyright (C) 2020-2023 Alexey Remizov <alexey@remizov.org>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

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
