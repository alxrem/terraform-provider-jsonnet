# jsonnet_file

The `jsonnet_file` data source renders a JSON document from a Jsonnet template file.

## Example Usage

```hcl-terraform
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
```

## Argument Reference

The following arguments are supported:

* `source` (string) &mdash; Path to the Jsonnet template file.
* `content` (string) &mdash; Content of Jsonnet template.

Exactly one of `source` or `content` is required.

* `jsonnet_path` &mdash; (Optional) Paths used to search additional Jsonnet libraries. Multiple paths are separated
  by colons. Overrides value of `jsonnet_path` configured in provider.
* `ext_str` &mdash; (Optional) Map of string for passing to the interpreter as external variables.
* `ext_code` &mdash; (Optional) Map of string representing a Jsonnet code for passing to the interpreter
  as external variables.
* `tla_str` &mdash; (Optional) Map of string for passing to the interpreter as top level argument.
* `tla_code` &mdash; (Optional) Map of string representing a Jsonnet code for passing to the interpreter
  as top-level argument.
* `string_output` &mdash; (Optional) When rendering a textual manifest, does not convert to a json string;
  "false" by default.

## Attributes Reference

The following attribute is exported:

* `rendered` &mdash; Rendered JSON document.
* `trace` &mdash; Output of std.trace() function.
