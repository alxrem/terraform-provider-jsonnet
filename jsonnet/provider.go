package jsonnet

import (
	"github.com/google/go-jsonnet"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
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
		ConfigureFunc: configureProvider,
	}
}

type providerConfig struct {
	importer jsonnet.Importer
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	jsonnetPath := d.Get("jsonnet_path").([]interface{})
	jpaths := make([]string, len(jsonnetPath))
	for i, path := range jsonnetPath {
		jpaths[i] = path.(string)
	}

	return &providerConfig{
		importer: &jsonnet.FileImporter{JPaths: jpaths},
	}, nil
}
