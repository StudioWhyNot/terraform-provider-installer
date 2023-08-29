// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shihanng/terraform-provider-installer/internal/installers/apt"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DataSourceApt{}

func NewDataSourceApt() datasource.DataSource {
	source := DataSourceApt{}
	source.SourceType = sources.SourceTypeApt
	return &source
}

// DataSourceApt defines the data source implementation.
type DataSourceApt struct {
	sources.SourceBase
	client *http.Client
}

// DataSourceAptModel describes the data source data model.
type DataSourceAptModel struct {
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
}

func (d *DataSourceApt) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.GetSourceName(req.ProviderTypeName)
}

func (d *DataSourceApt) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"path": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *DataSourceApt) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *DataSourceApt) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, success := sources.TryGetData[DataSourceAptModel](ctx, req.Config, &resp.Diagnostics)
	if !success {
		return
	}

	var diags diag.Diagnostics

	path, err := apt.FindInstalled(ctx, data.Name.ValueString())
	if err != nil {
		diags = xerrors.ToDiags(err)
		resp.Diagnostics.Append(diags...)
	}

	data.Path = types.StringValue(path)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "Read a data source")

	// Save data into Terraform state
	sources.SetStateData(ctx, &resp.State, &resp.Diagnostics, &data)
}
