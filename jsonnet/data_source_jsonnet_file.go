package jsonnet

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJsonnetFile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJsonnetFileRead,

		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source jsonnet file",
			},
			"jsonnet_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Additional directories to prepend to the Jsonnet load path",
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

func dataSourceJsonnetFileRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*providerConfig)

	source := d.Get("source").(string)

	jpaths := []string{}
	jsonnetPath := d.Get("jsonnet_path").(string)
	if jsonnetPath != "" {
		jpaths = strings.Split(jsonnetPath, ":")
	}

	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{JPaths: append(jpaths, config.jpaths...)})

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

	output, err := vm.EvaluateFile(source)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rendered", output); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hash(output))

	return nil
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
