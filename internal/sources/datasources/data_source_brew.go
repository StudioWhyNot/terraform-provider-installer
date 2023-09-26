// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasources

import (
	"context"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/installers/brew"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
	"github.com/shihanng/terraform-provider-installer/internal/sources/datasources/defaults"
	"github.com/shihanng/terraform-provider-installer/internal/sources/schemastrings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DataSourceBrew{}
var _ sources.SourceData = &DataSourceBrewModel{}

// DataSourceBrewModel describes the data source data model.
type DataSourceBrewModel struct {
	Name    types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`
	Path    types.String `tfsdk:"path"`
	Sudo    types.Bool   `tfsdk:"sudo"`
	Cask    types.Bool   `tfsdk:"cask"`
}

func (m *DataSourceBrewModel) GetSudo() bool {
	return m.Sudo.ValueBool()
}

func (m *DataSourceBrewModel) GetNamedVersion() models.NamedVersion {
	return models.NewNamedVersionFromStrings(brew.VersionSeperator, m.Name.ValueString(), m.Version.ValueString())
}

func (m *DataSourceBrewModel) GetName() string {
	return m.GetNamedVersion().Name
}

func (m *DataSourceBrewModel) GetVersion() *version.Version {
	return m.GetNamedVersion().Version
}

func (m *DataSourceBrewModel) GetCask() bool {
	return m.Cask.ValueBool()
}

func (m *DataSourceBrewModel) Initialize() bool {
	return !m.Name.IsNull()
}

func (m *DataSourceBrewModel) CopyFromTypedInstalledProgramInfo(installedInfo *models.TypedInstalledProgramInfo) {
	m.Name = types.StringValue(installedInfo.Name)
	m.Path = types.StringValue(installedInfo.Path)
	if installedInfo.Version != nil {
		m.Version = types.StringValue(installedInfo.Version.String())
	}
}

// DataSourceBrew defines the data source implementation.
type DataSourceBrew struct {
	*DataSource[*DataSourceBrewModel]
}

func NewDataSourceBrew() datasource.DataSource {
	resource := &DataSourceBrew{}
	resource.DataSource = NewDataSource[*DataSourceBrewModel](brew.NewBrewInstaller[*DataSourceBrewModel](resource))
	return resource
}

func (d *DataSourceBrew) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":    defaults.GetNameSchema(schemastrings.BrewNameDescription),
			"version": defaults.GetVersionSchema(schemastrings.BrewVersionDescription),
			"path":    defaults.GetPathSchema(schemastrings.BrewPathDescription),
			"sudo":    defaults.GetSudoSchema(),
			"cask":    defaults.GetCaskSchema(),
		},
	}
}
