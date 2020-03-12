package jsonnet

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"os/exec"
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
			// TODO: validate variable names
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
	extStr := d.Get("ext_str").(map[string]interface{})
	extCode := d.Get("ext_code").(map[string]interface{})
	tlaStr := d.Get("tla_str").(map[string]interface{})
	tlaCode := d.Get("tla_code").(map[string]interface{})
	command := config.command(source, extStr, extCode, tlaStr, tlaCode)

	stdout, err := command.Output()
	if err != nil {
		exitError := err.(*exec.ExitError)
		return fmt.Errorf(string(exitError.Stderr))
	}

	rendered := string(stdout)
	if err := d.Set("rendered", rendered); err != nil {
		return err
	}
	d.SetId(hash(rendered))

	return nil
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
