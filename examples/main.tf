provider "jsonnet" {
  jsonnet_path = ["./lib/"]
}

data "jsonnet_file" "template" {
  source = "example.jsonnet"
}

output "result" {
  value = data.jsonnet_file.template.rendered
}
