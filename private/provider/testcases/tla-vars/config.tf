data "jsonnet_file" "template" {
  tla_str = {
    a = "b"
  }

  tla_code = {
    b = "3 + 3"
  }

  source = "%s"
}