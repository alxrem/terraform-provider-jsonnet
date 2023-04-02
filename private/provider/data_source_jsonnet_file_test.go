package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"path"
	"regexp"
	"runtime"
	"testing"
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"jsonnet": providerserver.NewProtocol6WithError(New()),
	}
	testsDir string
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	testsDir = path.Join(path.Dir(filename), "tests")
}

func TestDataSourceJsonnetFile_FailWithoutContentNorSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `data "jsonnet_file" "template" {}`,
				ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured: \[content,source]`),
			},
		},
	})
}

func TestDataSourceJsonnetFile_FailWithBothContentAndSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "jsonnet_file" "template" {
					source  = "/path/to/file.jsonnet"
					content = "{}"
                }`,
				ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured: \[content,source]`),
			},
		},
	})
}

func TestDataSourceJsonnetFile_RenderFile(t *testing.T) {
	expected := `{
   "say": "hello world",
   "who": "world"
}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "jsonnet_file" "template" {
						source = "` + testsDir + `/no-vars.jsonnet"
					}
				`,
				Check: resource.TestCheckResourceAttr("data.jsonnet_file.template", "rendered", expected),
			},
		},
	})
}

func TestDataSourceJsonnetFile_RenderContentWithExtVars(t *testing.T) {
	expected := `{
   "say": "a",
   "sayAgain": "4"
}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "jsonnet_file" "template" {
						ext_str = {
							a = "a"
						}
						ext_code = {
							b = "2 + 2"
						}
						
						content = <<-EOF
						{
						  say: '%s' % [std.extVar('a')],
						  sayAgain: '%d' % [std.extVar('b')],
						}
						EOF
					}
				`,
				Check: resource.TestCheckResourceAttr("data.jsonnet_file.template", "rendered", expected),
			},
		},
	})
}

func TestDataSourceJsonnetFile_RenderContentWithTLAVars(t *testing.T) {
	expected := `{
   "say": "b",
   "sayAgain": "6"
}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "jsonnet_file" "template" {
						tla_str = {
							a = "b"
						}
						
						tla_code = {
							b = "3 + 3"
						}
						
						content = <<-EOF
						function(a, b){
						  say: '%s' % [a],
						  sayAgain: '%d' % [b],
						}
						EOF
					}
				`,
				Check: resource.TestCheckResourceAttr("data.jsonnet_file.template", "rendered", expected),
			},
		},
	})
}

func TestDataSourceJsonnetFile_RenderContentToString(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "jsonnet_file" "template" {
						tla_code = {
							vars = jsonencode({
								foo = "bar"
								bar = "baz"
							})
						}

						content = <<-EOF
						function(vars)
						  std.lines(["%s=%s" % [k, std.escapeStringBash(vars[k])] for k in std.objectFields(vars)])
						EOF

						string_output = true
					}
				`,
				Check: resource.TestCheckResourceAttr("data.jsonnet_file.template", "rendered", "bar='baz'\nfoo='bar'\n\n"),
			},
		},
	})
}
