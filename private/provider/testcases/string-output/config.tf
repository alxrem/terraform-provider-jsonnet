data "jsonnet_file" "template" {
  tla_code = {
    vars = jsonencode({
      foo = "bar"
      bar = "baz"
    })
  }
  string_output = true
  source = "%s"
}
