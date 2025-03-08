# Configure Jsonnet provider
provider "jsonnet" {
    jsonnet_path = "${root.module}/lib:/usr/local/share/jsonnet/grafonnet-lib"
}

# Template of Grafana dashboard
data "jsonnet_file" "dashboard" {
    source = "${root.module}/jsonnet/dashboard.jsonnet"
}

# Install dashboard using Grafana API
resource "grafana_dashboard" "service_overview" {
  config_json = data.jsonnet_file.dashboard.rendered
}
