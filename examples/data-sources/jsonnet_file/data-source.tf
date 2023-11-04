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
