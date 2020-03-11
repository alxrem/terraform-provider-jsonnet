package jsonnet

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"os/exec"
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
			"jsonnet_bin": {
				Type:     schema.TypeString,
				Default:  "jsonnet",
				Optional: true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"jsonnet_file": dataSourceJsonnetFile(),
		},
		ConfigureFunc: configureProvider,
	}
}

type providerConfig struct {
	cmd  string
	args []string
}

func (pc *providerConfig) command(source string, extStr, extCode, tlaStr, tlaCode map[string]interface{}) *exec.Cmd {
	for argKey, vars := range map[string]map[string]interface{}{
		"--ext-str":  extStr,
		"--ext-code": extCode,
		"--tla-str":  tlaStr,
		"--tla-code": tlaCode,
	} {
		for k, v := range vars {
			pc.args = append(pc.args, argKey, k+"="+v.(string))
		}
	}

	cmd := exec.Command(pc.cmd, append(pc.args, source)...)
	cmd.Env = os.Environ()
	return cmd
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	jsonnetBin := d.Get("jsonnet_bin").(string)

	cmd, err := exec.LookPath(jsonnetBin)
	if err != nil {
		return nil, err
	}

	jsonnetPath := d.Get("jsonnet_path").([]interface{})
	args := make([]string, 0)
	for _, path := range jsonnetPath {
		args = append(args, "-J")
		args = append(args, path.(string))
	}

	return &providerConfig{
		cmd:  cmd,
		args: args,
	}, nil
}
