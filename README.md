# terraform-provider-jsonnet

Terraform provider for generating JSON documents from [Jsonnet](https://jsonnet.org/) templates. It initially aimed to
rendering [Grafana](https://grafana.com) dashboards using [grafonnet library](https://github.com/grafana/grafonnet-lib).

## Installation

1. Download archive with the latest version of provider for your operating system from
   [Gitlab releases page](https://gitlab.com/alxrem/terraform-provider-jsonnet/-/releases).
2. Unpack provider to `$HOME/.terraform.d/plugins`, i.e.
   ```
   unzip terraform-provider-jsonnet_vX.Y.Z-linux-amd64.zip -d $HOME/.terraform.d/plugins/
   ```
3. Init your terraform project
   ```
   terraform init
   ```
4. Download from the [releases page](https://gitlab.com/alxrem/terraform-provider-jsonnet/-/releases) the file
   containing checksum of downloaded version of plugin and compare it with checksum contained in the file
   `.terraform/plugins/<os_arch>/lock.json` in the root of your terraform project.

## Jsonnet Provider

### Example Usage

```hcl-terraform
# Configure Jsonnet provider
provider "jsonnet" {
    jsonnet_path = ["${root.module}/jsonnet/grafonnet-lib"]
}

# Template of Grafana dashboard
data "jsonnet_file" "dashboard" {
    source = "${root.module}/jsonnet/dashboard.jsonnet"
}

# Install dashboard using Grafana API
resource "grafana_dashboard" "service_overview" {
  config_json = data.jsonnet_file.dashboard.rendered
}
```

### Argument Reference

The following argument is supported in the provider block:

* `jsonnet_path` &mdash; (Optional) Array of paths used to search additional Jsonnet libraries.

## jsonnet_file data source

The `jsonnet_file` data source renders a JSON document from a Jsonnet template file.

## Example Usage

```hcl-terraform
data "jsonnet_file" "dashboard" {
    ext_str = {
        service = "my_service"
    }

    ext_code = {
        installations = jsonencode(var.installations)
    }

    tla_str = {
        description = "My service"
    }

    source = "${root.module}/jsonnet/dashboard.jsonnet"
}
```

### Argument Reference

The following arguments are supported:

* `source` &mdash; (Required) Path to the Jsonnet template file.
* `ext_str` &mdash; (Optional) Map of string for passing to the interpreter as external variables.
* `ext_code` &mdash; (Optional) Map of string representing a Jsonnet code for passing to the interpreter
                                as external variables.
* `tla_str` &mdash; (Optional) Map of string for passing to the interpreter as top level argument.
* `tla_code` &mdash; (Optional) Map of string representing a Jsonnet code for passing to the interpreter
                                as top-level argument.

### Attributes Reference

The following attribute is exported:

* `rendered` &mdash; Rendered JSON document.