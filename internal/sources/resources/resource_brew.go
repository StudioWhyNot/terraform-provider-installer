// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"context"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers/brew"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
	"github.com/shihanng/terraform-provider-installer/internal/sources/resources/defaults"
	"github.com/shihanng/terraform-provider-installer/internal/sources/schemastrings"
	"github.com/shihanng/terraform-provider-installer/internal/system"
	"github.com/shihanng/terraform-provider-installer/internal/terraformutils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ResourceBrew{}
var _ resource.ResourceWithImportState = &ResourceBrew{}
var _ sources.SourceData = &ResourceBrewModel{}

// ResourceBrewModel describes the resource data model.
type ResourceBrewModel struct {
	Id                                   types.String `tfsdk:"id"`
	Name                                 types.String `tfsdk:"name"`
	Version                              types.String `tfsdk:"version"`
	Path                                 types.String `tfsdk:"path"`
	Sudo                                 types.Bool   `tfsdk:"sudo"`
	Environment                          types.Map    `tfsdk:"environment"`
	Cask                                 types.Bool   `tfsdk:"sudo"`
	Secrets                              types.Map    `tfsdk:"secrets"`
	*terraformutils.RemoteConnectionInfo `tfsdk:"remote_connection"`
}

func (m *ResourceBrewModel) GetSudo() bool {
	return m.Sudo.ValueBool()
}

func (m *ResourceBrewModel) GetEnvironmentAndSecrets(ctx context.Context) map[string]string {
	return system.MergeMaps(sources.MapValueToMap(ctx, &m.Environment), sources.MapValueToMap(ctx, &m.Secrets))
}

func (m *ResourceBrewModel) GetNamedVersion() models.NamedVersion {
	return models.NewNamedVersionFromStrings(brew.VersionSeperator, m.Name.ValueString(), m.Version.ValueString())
}

func (m *ResourceBrewModel) GetName() string {
	return m.GetNamedVersion().Name
}

func (m *ResourceBrewModel) GetVersion() *version.Version {
	return m.GetNamedVersion().Version
}

func (m *ResourceBrewModel) GetCask() bool {
	return m.Cask.ValueBool()
}

func (m *ResourceBrewModel) Initialize() bool {
	m.Id = sources.GetIDFromNameAndVersion(brew.VersionSeperator, m.Name, m.Version, enums.InstallerBrew)
	return !m.Name.IsNull()
}

func (m *ResourceBrewModel) GetRemoteConnectionInfo() *terraformutils.RemoteConnectionInfo {
	return m.RemoteConnectionInfo
}

func (m *ResourceBrewModel) CopyFromTypedInstalledProgramInfo(installedInfo *models.TypedInstalledProgramInfo) {
	if installedInfo == nil {
		m.Name = types.StringNull()
		m.Path = types.StringNull()
		m.Version = types.StringNull()
		return
	}
	m.Name = types.StringValue(installedInfo.Name)
	m.Path = types.StringValue(installedInfo.Path)
	if installedInfo.Version != nil {
		m.Version = types.StringValue(installedInfo.Version.String())
	}
}

// ResourceBrew defines the resource implementation.
type ResourceBrew struct {
	*Resource[*ResourceBrewModel]
}

func NewResourceBrew() resource.Resource {
	resource := &ResourceBrew{}
	resource.Resource = NewResource[*ResourceBrewModel](brew.NewBrewInstaller[*ResourceBrewModel](resource))
	return resource
}

func (r *ResourceBrew) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: schemastrings.BrewSourceDescription,
		Attributes: map[string]schema.Attribute{
			"id":          defaults.GetIdSchema(),
			"name":        defaults.GetNameSchema(schemastrings.BrewNameDescription),
			"version":     defaults.GetVersionSchema(schemastrings.BrewVersionDescription),
			"path":        defaults.GetPathSchema(schemastrings.BrewPathDescription),
			"sudo":        defaults.GetSudoSchema(brew.DefaultSudo),
			"environment": defaults.GetEnvironmentSchema(),
			"secrets":     defaults.GetSecretsSchema(),
			"cask":        defaults.GetCaskSchema(schemastrings.BrewCaskDescription, brew.DefaultCask),
		},
		Blocks: map[string]schema.Block{
			"remote_connection": defaults.GetRemoteConnectionBlockSchema(),
		},
	}
}
