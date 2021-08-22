package jsonnet

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"jsonnet_path": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"jsonnet_file": dataSourceJsonnetFile(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

type providerConfig struct {
	jpaths []string
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	jsonnetPath := d.Get("jsonnet_path").([]interface{})
	config := &providerConfig{
		jpaths: make([]string, len(jsonnetPath)),
	}
	for i, path := range jsonnetPath {
		config.jpaths[i] = path.(string)
	}

	return config, nil
}
