# Jsonnet Provider

Terraform provider for generating JSON documents from [Jsonnet](https://jsonnet.org/) templates. It initially aimed to
rendering [Grafana](https://grafana.com) dashboards using [grafonnet library](https://github.com/grafana/grafonnet-lib).

## Example Usage

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

## Argument Reference

The following argument is supported in the provider block:

* `jsonnet_path` &mdash; (Optional) Array of paths used to search additional Jsonnet libraries.
