# terraform-provider-jsonnet

Terraform provider for generating JSON documents from [Jsonnet](https://jsonnet.org/) templates. It initially aimed to
rendering [Grafana](https://grafana.com) dashboards using [grafonnet library](https://github.com/grafana/grafonnet-lib).

## Migration 1.x -> 2.0.0

In versions 1.x parameter `jsonnet_path` of provider was of type **list**.
Starting from version 2.0.0 path to jsonnet libraries must be prepresented
as strings with the paths divided by colon as in shell `PATH` variable.

The easiest way to migrate provider definition to 2.x is to use `join` function i.e

```
 provider "jsonnet" {
-  jsonnet_path = ["${path.module}/jsonnet", "${path.module}/jsonnet/grafonnet-lib"]
+  jsonnet_path = join(":", ["${path.module}/jsonnet", "${path.module}/jsonnet/grafonnet-lib"])
 }
```

## Installation

### terraform 0.13+

Add into your Terraform configuration this code:

```hcl-terraform
terraform {
  required_providers {
    jsonnet = {
      source = "alxrem/jsonnet"
    }
  }
}
```

and run `terraform init`

### terraform 0.12 and earlier

1. Download archive with the latest version of provider for your operating system from
   [Github releases page](https://github.com/alxrem/terraform-provider-jsonnet/releases).
2. Unpack provider to `$HOME/.terraform.d/plugins`, i.e.
   ```
   unzip terraform-provider-jsonnet_vX.Y.Z_linux_amd64.zip terraform-provider-jsonnet_* -d $HOME/.terraform.d/plugins/
   ```
3. Init your terraform project
   ```
   terraform init
   ```

## Usage

Read the [documentation on Terraform Registry site](https://registry.terraform.io/providers/alxrem/jsonnet/latest/docs).
