package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
	"strings"
)

type JsonnetProvider struct{}

type JsonnetProviderModel struct {
	JsonnetPath types.String `tfsdk:"jsonnet_path"`
}

func (p *JsonnetProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jsonnet"
}

func (p *JsonnetProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"jsonnet_path": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Paths used to search additional Jsonnet libraries. " +
					"Can be specified by `JSONNET_PATH` environment variable " +
					"with multiple paths separated by colons.",
			},
		},
	}
}

func (p *JsonnetProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config JsonnetProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.JsonnetPath.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("jsonnet_path"),
			"Unknown Jsonnet path",
			"The provider cannot configure paths to search jsonnet libraries as there is an unknown configuration value. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the JSONNET_PATH environment variable.",
		)
	}

	var jsonnetPath string
	if config.JsonnetPath.IsNull() {
		jsonnetPath = os.Getenv("JSONNET_PATH")
	} else {
		jsonnetPath = config.JsonnetPath.String()
	}

	jsonnetPaths := strings.Split(jsonnetPath, ":")
	for i, jpath := range jsonnetPaths {
		jsonnetPaths[i] = strings.TrimSpace(jpath)
	}

	resp.DataSourceData = jsonnetPaths
}

func (p *JsonnetProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewJsonnetFileDataSource,
	}
}

func (p *JsonnetProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func New() provider.Provider {
	return &JsonnetProvider{}
}
