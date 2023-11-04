data "jsonnet_file" "dashboard" {
  jsonnet_path = "${root.module}/jsonnet/lib"
  source       = "${root.module}/jsonnet/dashboard.jsonnet"

  ext_str = {
    service = "my_service"
  }

  ext_code = {
    installations = jsonencode(var.installations)
  }

  tla_str = {
    description = "My service"
  }
}
