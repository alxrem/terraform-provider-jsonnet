// Copyright (C) 2020-2023 Alexey Remizov <alexey@remizov.org>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"os"
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

func TestDataSourceJsonnetFile_RenderContentWithDefaultJPath(t *testing.T) {
	expected := `{
   "data": "from global library"
}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck: func() {
			_ = os.Setenv("JSONNET_PATH", "tests/global_libs")
		},
		Steps: []resource.TestStep{
			{
				Config: `
					data "jsonnet_file" "template" {
						content = "import \"common.libsonnet\""
					}
				`,
				Check: resource.TestCheckResourceAttr("data.jsonnet_file.template", "rendered", expected),
			},
		},
	})
	_ = os.Setenv("JSONNET_PATH", "global_libs")
}

func TestDataSourceJsonnetFile_RenderFileWithLocalJPath(t *testing.T) {
	expected := `{
   "data": "from local library"
}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck: func() {
			_ = os.Setenv("JSONNET_PATH", "tests/global_libs")
		},
		Steps: []resource.TestStep{
			{
				Config: `
					data "jsonnet_file" "template" {
						jsonnet_path = "tests/local_libs"
						source = "` + testsDir + `/import.jsonnet"
					}
				`,
				Check: resource.TestCheckResourceAttr("data.jsonnet_file.template", "rendered", expected),
			},
		},
	})
}

func TestDataSourceJsonnetFile_RenderContentWithTrace(t *testing.T) {
	expectedResult := `{
   "a": "rest"
}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "jsonnet_file" "template" {
						content = <<EOF
						{
							"a": std.trace("str", "rest")
						}
						EOF
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.jsonnet_file.template", "rendered", expectedResult),
					resource.TestCheckResourceAttr("data.jsonnet_file.template", "trace", "TRACE: data:2 str\n"),
				),
			},
			{
				Config: `
					data "jsonnet_file" "template" {
						content = <<EOF
						{
							"a": "rest"
						}
						EOF
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.jsonnet_file.template", "rendered", expectedResult),
					resource.TestCheckResourceAttr("data.jsonnet_file.template", "trace", ""),
				),
			},
		},
	})
}

func TestDataSourceJsonnetFile_AddTraceToDiagnostics(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "jsonnet_file" "template" {
						content = <<EOF
						{
							"a": std.trace("str", "rest"),
							"b": 1/0,
						}
						EOF
					}
				`,
				ExpectError: regexp.MustCompile(`(?s)RUNTIME ERROR: Division by zero\..*TRACE: data:2 str`),
			},
		},
	})
}
