// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasources

import (
	"context"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
	"github.com/shihanng/terraform-provider-installer/internal/sources/datasources/defaults"
	"github.com/shihanng/terraform-provider-installer/internal/sources/schemastrings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DataSourceApt{}
var _ sources.SourceData = &DataSourceAptModel{}

// DataSourceAptModel describes the data source data model.
type DataSourceAptModel struct {
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
}

func (m *DataSourceAptModel) GetName() types.String {
	return m.Name
}

func (m *DataSourceAptModel) GetVersion() *version.Version {
	return nil
}

func (m *DataSourceAptModel) SetDataFromTypedInstalledProgramInfo(info *models.TypedInstalledProgramInfo) {
	m.Name = types.StringValue(info.Name)
	m.Path = types.StringValue(info.Path)
}

// DataSourceApt defines the data source implementation.
type DataSourceApt struct {
	*DataSource[*DataSourceAptModel]
}

func NewDataSourceApt() datasource.DataSource {
	return &DataSourceApt{
		DataSource: NewDataSource[*DataSourceAptModel](enums.InstallerApt),
	}
}

func (d *DataSourceApt) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": defaults.GetNameSchema(schemastrings.AptNameDescription),
			"path": defaults.GetPathSchema(schemastrings.AptPathDescription),
		},
	}
}
