# terraform-provider-jsonnet

Terraform provider for generating JSON documents from [Jsonnet](https://jsonnet.org/) templates. It initially aimed to
rendering [Grafana](https://grafana.com) dashboards using [grafonnet library](https://github.com/grafana/grafonnet-lib).

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
