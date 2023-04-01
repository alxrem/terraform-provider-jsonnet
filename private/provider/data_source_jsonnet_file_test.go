package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"os"
	"path"
	"runtime"
	"testing"
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"jsonnet": providerserver.NewProtocol6WithError(New()),
	}
	testCasesDir string
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	testCasesDir = path.Join(path.Dir(filename), "testcases")
}

func LoadTestCase(name string) resource.TestCase {
	testCaseDir := path.Join(testCasesDir, name)
	source := path.Join(testCaseDir, "template.jsonnet")
	configTpl, err := os.ReadFile(path.Join(testCaseDir, "config.tf"))
	if err != nil {
		panic(err)
	}
	config := fmt.Sprintf(string(configTpl), source)
	expected, err := os.ReadFile(path.Join(testCaseDir, "expected"))
	if err != nil {
		panic(err)
	}

	return resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.jsonnet_file.template", "rendered", string(expected)),
				),
			},
		},
	}
}

func TestJsonnetRendering(t *testing.T) {
	resource.UnitTest(t, LoadTestCase("no-vars"))
	resource.UnitTest(t, LoadTestCase("ext-vars"))
	resource.UnitTest(t, LoadTestCase("tla-vars"))
	resource.UnitTest(t, LoadTestCase("string-output"))
}
