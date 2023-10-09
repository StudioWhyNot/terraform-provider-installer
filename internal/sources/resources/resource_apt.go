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
	"github.com/shihanng/terraform-provider-installer/internal/installers/apt"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
	"github.com/shihanng/terraform-provider-installer/internal/sources/resources/defaults"
	"github.com/shihanng/terraform-provider-installer/internal/sources/schemastrings"
	"github.com/shihanng/terraform-provider-installer/internal/terraformutils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ResourceApt{}
var _ resource.ResourceWithImportState = &ResourceApt{}
var _ sources.SourceData = &ResourceAptModel{}

// ResourceAptModel describes the resource data model.
type ResourceAptModel struct {
	Id                                   types.String `tfsdk:"id"`
	Name                                 types.String `tfsdk:"name"`
	Version                              types.String `tfsdk:"version"`
	Path                                 types.String `tfsdk:"path"`
	Sudo                                 types.Bool   `tfsdk:"sudo"`
	*terraformutils.RemoteConnectionInfo `tfsdk:"remote_connection"`
}

func (m *ResourceAptModel) GetSudo() bool {
	return m.Sudo.ValueBool()
}

func (m *ResourceAptModel) GetEnvironment() map[string]string {
	return nil
}

func (m *ResourceAptModel) GetNamedVersion() models.NamedVersion {
	return models.NewNamedVersionFromStrings(apt.VersionSeperator, m.Name.ValueString(), m.Version.ValueString())
}

func (m *ResourceAptModel) GetName() string {
	return m.GetNamedVersion().Name
}

func (m *ResourceAptModel) GetVersion() *version.Version {
	return m.GetNamedVersion().Version
}

func (m *ResourceAptModel) Initialize() bool {
	m.Id = sources.GetIDFromNameAndVersion(apt.VersionSeperator, m.Name, m.Version, enums.InstallerApt)
	return !m.Name.IsNull()
}

func (m *ResourceAptModel) GetRemoteConnectionInfo() *terraformutils.RemoteConnectionInfo {
	return m.RemoteConnectionInfo
}

func (m *ResourceAptModel) CopyFromTypedInstalledProgramInfo(installedInfo *models.TypedInstalledProgramInfo) {
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

// ResourceApt defines the resource implementation.
type ResourceApt struct {
	*Resource[*ResourceAptModel]
}

func NewResourceApt() resource.Resource {
	resource := &ResourceApt{}
	resource.Resource = NewResource[*ResourceAptModel](apt.NewAptInstaller[*ResourceAptModel](resource))
	return resource
}

func (r *ResourceApt) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: schemastrings.AptSourceDescription,
		Attributes: map[string]schema.Attribute{
			"id":      defaults.GetIdSchema(),
			"name":    defaults.GetNameSchema(schemastrings.AptNameDescription),
			"version": defaults.GetVersionSchema(schemastrings.AptVersionDescription),
			"path":    defaults.GetPathSchema(schemastrings.AptPathDescription),
			"sudo":    defaults.GetSudoSchema(apt.DefaultSudo),
		},
		Blocks: map[string]schema.Block{
			"remote_connection": defaults.GetRemoteConnectionBlockSchema(),
		},
	}
}
