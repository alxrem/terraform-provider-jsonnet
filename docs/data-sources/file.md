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

* `source` &mdash; (Required) Path to the Jsonnet template file.
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
