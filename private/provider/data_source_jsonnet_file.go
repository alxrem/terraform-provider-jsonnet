package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/go-jsonnet"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &JsonnetFileDataSource{}
var _ datasource.DataSourceWithConfigure = &JsonnetFileDataSource{}

type JsonnetFileDataSource struct {
	defaultJsonnetPaths []string
}

func (d *JsonnetFileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.defaultJsonnetPaths = req.ProviderData.([]string)
}

func NewJsonnetFileDataSource() datasource.DataSource {
	return &JsonnetFileDataSource{}
}

type JsonnetFileDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Source       types.String `tfsdk:"source"`
	ExtStr       types.Map    `tfsdk:"ext_str"`
	ExtCode      types.Map    `tfsdk:"ext_code"`
	TlaStr       types.Map    `tfsdk:"tla_str"`
	TlaCode      types.Map    `tfsdk:"tla_code"`
	StringOutput types.Bool   `tfsdk:"string_output"`
	Rendered     types.String `tfsdk:"rendered"`
}

func (d *JsonnetFileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (d *JsonnetFileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"source": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Path to the Jsonnet template file.",
			},
			"ext_str": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Map of string for passing to the interpreter as external variables.",
			},
			"ext_code": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Map of string representing a Jsonnet code for passing to the interpreter as external variables.",
			},
			"tla_str": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Map of string for passing to the interpreter as top level argument.",
			},
			"tla_code": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Map of string representing a Jsonnet code for passing to the interpreter as top-level argument.",
			},
			"string_output": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: `When rendering a textual manifest, does not convert to a json string; "false" by default.`,
			},
			"rendered": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Rendered text.",
			},
		},
	}
}

func (d *JsonnetFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state JsonnetFileDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{JPaths: d.defaultJsonnetPaths})

	for _, attr := range []struct {
		m types.Map
		f func(string, string)
	}{
		{state.ExtStr, vm.ExtVar},
		{state.ExtCode, vm.ExtCode},
		{state.TlaStr, vm.TLAVar},
		{state.TlaCode, vm.TLACode}} {
		var args map[string]string
		resp.Diagnostics.Append(attr.m.ElementsAs(ctx, &args, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for name, value := range args {
			attr.f(name, value)
		}
	}

	vm.StringOutput = state.StringOutput.ValueBool()

	rendered, err := vm.EvaluateFile(state.Source.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to render test from jsonnet template", err.Error())
		return
	}

	state.Rendered = types.StringValue(rendered)
	state.Id = types.StringValue(hash(rendered))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
