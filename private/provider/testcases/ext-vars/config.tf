data "jsonnet_file" "template" {
  ext_str = {
    a = "a"
  }
  ext_code = {
    b = "2 + 2"
  }

  source = "%s"
}