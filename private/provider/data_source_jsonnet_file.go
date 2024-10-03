// Copyright (C) 2020-2023 Alexey Remizov <alexey@remizov.org>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/go-jsonnet"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var _ datasource.DataSource = &JsonnetFileDataSource{}
var _ datasource.DataSourceWithConfigure = &JsonnetFileDataSource{}
var _ datasource.DataSourceWithConfigValidators = &JsonnetFileDataSource{}

type JsonnetFileDataSource struct {
	defaultJsonnetPaths []string
}

func (d *JsonnetFileDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("content"),
			path.MatchRoot("source"),
		),
	}
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
	JsonnetPath  types.String `tfsdk:"jsonnet_path"`
	Source       types.String `tfsdk:"source"`
	Content      types.String `tfsdk:"content"`
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
			"jsonnet_path": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Paths used to search additional Jsonnet libraries. " +
					"Overrides paths from provider config.",
			},
			"source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Path to the Jsonnet template file.",
			},
			"content": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Jsonnet template.",
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

	rendered, err := func() (string, error) {
		var jPaths []string

		if state.JsonnetPath.IsNull() {
			jPaths = d.defaultJsonnetPaths
		} else {
			jPaths = strings.Split(state.JsonnetPath.ValueString(), ":")
		}

		vm.Importer(&jsonnet.FileImporter{JPaths: jPaths})

		if state.Source.IsNull() {
			return vm.EvaluateAnonymousSnippet("data", state.Content.ValueString())
		} else {
			return vm.EvaluateFile(state.Source.ValueString())
		}
	}()

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
