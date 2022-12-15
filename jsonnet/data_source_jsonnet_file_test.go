package jsonnet

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"io/ioutil"
	"os"
	"testing"
)

var testProviderFactories = map[string]func() (*schema.Provider, error){
	"jsonnet": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

func TestJsonnetRendering(t *testing.T) {
	var cases = []struct {
		name     string
		manifest string
		template string
		want     string
	}{
		{
			name: "no vars",
			manifest: `
data "jsonnet_file" "template" {
  source = "%s"
}

output "rendered" {
  value = data.jsonnet_file.template.rendered
}
`,
			template: `{
	who: 'world',
	say: 'hello %(who)s' % (self)
}`,
			want: `{
   "say": "hello world",
   "who": "world"
}
`},
		{
			name: "ext vars",
			manifest: `
data "jsonnet_file" "template" {
  ext_str = {
    a = "a"
  }
  ext_code = {
    b = "2 + 2"
  }

  source = "%s"
}

output "rendered" {
  value = data.jsonnet_file.template.rendered
}
`,
			template: `{
	say: '%s' % [std.extVar('a')],
    sayAgain: '%d' % [std.extVar('b')],
}`,
			want: `{
   "say": "a",
   "sayAgain": "4"
}
`},
		{
			name: "tla vars",
			manifest: `
data "jsonnet_file" "template" {
  tla_str = {
    a = "b"
  }
  tla_code = {
    b = "3 + 3"
  }

  source = "%s"
}

output "rendered" {
  value = data.jsonnet_file.template.rendered
}
`,
			template: `function(a, b){
	say: '%s' % [a],
    sayAgain: '%d' % [b],
}`,
			want: `{
   "say": "b",
   "sayAgain": "6"
}
`},
		{
			name: "string_output",
			manifest: `
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

output "rendered" {
  value = data.jsonnet_file.template.rendered
}
`,
			template: `function(vars)
  std.lines(["%s=%s" % [k, std.escapeStringBash(vars[k])] for k in std.objectFields(vars)])
`,
			want: "bar='baz'\nfoo='bar'\n\n",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			templateFile, err := testJsonnetTemplate(tt.template)
			if templateFile != nil {
				//noinspection GoUnhandledErrorResult
				defer os.Remove(templateFile.Name())
			}

			if err != nil {
				t.Fatalf("error: %s", err.Error())
			}

			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(tt.manifest, templateFile.Name()),
						Check: func(s *terraform.State) error {
							got := s.RootModule().Outputs["rendered"]
							if tt.want != got.Value {
								return fmt.Errorf("template:\n%s\ngot:\n%s\nwant:\n%s\n", tt.template, got.Value, tt.want)
							}
							return nil
						},
					},
				},
			})
		})
	}
}

func testJsonnetTemplate(template string) (*os.File, error) {
	templateFile, err := ioutil.TempFile("", "*.jsonnet")
	if err != nil {
		return nil, err
	}
	//noinspection GoUnhandledErrorResult
	defer templateFile.Close()

	if _, err := templateFile.Write([]byte(template)); err != nil {
		return templateFile, err
	}

	return templateFile, nil
}
