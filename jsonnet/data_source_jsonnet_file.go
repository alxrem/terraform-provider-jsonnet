package jsonnet

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/go-jsonnet"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
)

func dataSourceJsonnetFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceJsonnetFileRead,

		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source jsonnet file",
			},
			"ext_str": {
				Type:        schema.TypeMap,
				Optional:    true,
				Default:     make(map[string]interface{}),
				Description: "External variables providing value as a string",
			},
			"ext_code": {
				Type:        schema.TypeMap,
				Optional:    true,
				Default:     make(map[string]interface{}),
				Description: "External variables providing value as a code",
			},
			"tla_str": {
				Type:        schema.TypeMap,
				Optional:    true,
				Default:     make(map[string]interface{}),
				Description: "Top-level arguments providing value as a string",
			},
			"tla_code": {
				Type:        schema.TypeMap,
				Optional:    true,
				Default:     make(map[string]interface{}),
				Description: "Top-level arguments providing value as a code",
			},
			"rendered": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered template",
			},
		},
	}
}

func dataSourceJsonnetFileRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*providerConfig)

	source := d.Get("source").(string)

	vm := jsonnet.MakeVM()
	vm.Importer(config.importer)

	for name, value := range d.Get("ext_str").(map[string]interface{}) {
		vm.ExtVar(name, value.(string))
	}
	for name, value := range d.Get("ext_code").(map[string]interface{}) {
		vm.ExtCode(name, value.(string))
	}
	for name, value := range d.Get("tla_str").(map[string]interface{}) {
		vm.TLAVar(name, value.(string))
	}
	for name, value := range d.Get("tla_code").(map[string]interface{}) {
		vm.TLACode(name, value.(string))
	}

	snippet, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	output, err := vm.EvaluateSnippet(source, string(snippet))
	if err != nil {
		return err
	}

	if err := d.Set("rendered", output); err != nil {
		return err
	}
	d.SetId(hash(output))

	return nil
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
