provider "jsonnet" {
  jsonnet_path = ["./lib/"]

  version = "~> 0.2"
}

data "jsonnet_file" "template" {
  ext_str  = {
    hello = "from external string"
  }

  ext_code = {
    calc = "2 + 2"
  }

  source = "example.jsonnet"
}

data "jsonnet_file" "tla" {
  tla_str = {
    hello = "from top level argument"
  }

  tla_code = {
    calc = "3 + 3"
  }

  source = "tla.jsonnet"
}

output "example" {
  value = data.jsonnet_file.template.rendered
}

output "tla" {
  value = data.jsonnet_file.tla.rendered
}
