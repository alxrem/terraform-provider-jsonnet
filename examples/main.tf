terraform {
  required_providers {
    jsonnet = {
      source  = "alxrem/jsonnet"
      version = "~> 1.0"
    }
  }
}

provider "jsonnet" {
  jsonnet_path = "./lib/"
}

data "jsonnet_file" "template" {
  ext_str = {
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

data "jsonnet_file" "extra_paths" {
  source       = "extra_paths.jsonnet"
  jsonnet_path = "./lib-extra"
}

output "example" {
  value = data.jsonnet_file.template.rendered
}

output "tla" {
  value = data.jsonnet_file.tla.rendered
}

output "extra_paths" {
  value = data.jsonnet_file.extra_paths.rendered
}
