// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasources

import (
	"context"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/installers/apt"
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
	Name    types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`
	Path    types.String `tfsdk:"path"`
	Sudo    types.Bool   `tfsdk:"sudo"`
}

func (m *DataSourceAptModel) GetSudo() bool {
	return m.Sudo.ValueBool()
}

func (m *DataSourceAptModel) GetNamedVersion() models.NamedVersion {
	return models.NewNamedVersionFromStrings(apt.VersionSeperator, m.Name.ValueString(), m.Version.ValueString())
}

func (m *DataSourceAptModel) GetName() string {
	return m.GetNamedVersion().Name
}

func (m *DataSourceAptModel) GetVersion() *version.Version {
	return m.GetNamedVersion().Version
}

func (m *DataSourceAptModel) Initialize() bool {
	return !m.Name.IsNull()
}

func (m *DataSourceAptModel) CopyFromTypedInstalledProgramInfo(installedInfo *models.TypedInstalledProgramInfo) {
	m.Name = types.StringValue(installedInfo.Name)
	m.Path = types.StringValue(installedInfo.Path)
	if installedInfo.Version != nil {
		m.Version = types.StringValue(installedInfo.Version.String())
	}
}

// DataSourceApt defines the data source implementation.
type DataSourceApt struct {
	*DataSource[*DataSourceAptModel]
}

func NewDataSourceApt() datasource.DataSource {
	resource := &DataSourceApt{}
	resource.DataSource = NewDataSource[*DataSourceAptModel](apt.NewAptInstaller[*DataSourceAptModel](resource))
	return resource
}

func (d *DataSourceApt) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":    defaults.GetNameSchema(schemastrings.AptNameDescription),
			"version": defaults.GetVersionSchema(schemastrings.AptVersionDescription),
			"path":    defaults.GetPathSchema(schemastrings.AptPathDescription),
			"sudo":    defaults.GetSudoSchema(),
		},
	}
}
