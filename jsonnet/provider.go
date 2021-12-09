package jsonnet

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"jsonnet_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JSONNET_PATH", nil),
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
	jsonnetPath := strings.Split(d.Get("jsonnet_path").(string), ":")
	config := &providerConfig{
		jpaths: make([]string, len(jsonnetPath)),
	}
	for i, path := range jsonnetPath {
		config.jpaths[i] = strings.TrimSpace(path)
	}

	return config, nil
}
