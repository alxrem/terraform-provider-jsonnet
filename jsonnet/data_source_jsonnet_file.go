package jsonnet

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	command := config.command(source)

	stdout, err := command.Output()
	if err != nil {
		return err
	}

	if err := d.Set("rendered", string(stdout)); err != nil {
		return err
	}
	d.SetId(hash(string(stdout)))

	return nil
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
