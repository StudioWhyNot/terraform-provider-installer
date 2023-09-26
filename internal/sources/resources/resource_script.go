// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers/script"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
	"github.com/shihanng/terraform-provider-installer/internal/sources/resources/defaults"
	"github.com/shihanng/terraform-provider-installer/internal/sources/schemastrings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ResourceScript{}
var _ resource.ResourceWithImportState = &ResourceScript{}
var _ sources.SourceData = &ResourceScriptModel{}

// ResourceScriptModel describes the resource data model.
type ResourceScriptModel struct {
	Id                  types.String `tfsdk:"id"`
	Path                types.String `tfsdk:"path"`
	InstallScript       types.String `tfsdk:"install_script"`
	FindInstalledScript types.String `tfsdk:"find_installed_script"`
	UninstallScript     types.String `tfsdk:"uninstall_script"`
	DefaultArgs         types.List   `tfsdk:"default_args"`
	AdditionalArgs      types.List   `tfsdk:"additional_args"`
	Sudo                types.Bool   `tfsdk:"sudo"`
	Shell               types.String `tfsdk:"shell"`
}

func (m *ResourceScriptModel) GetId() string {
	return m.Id.ValueString()
}

func (m *ResourceScriptModel) GetPath() string {
	return m.Path.ValueString()
}

func (m *ResourceScriptModel) GetInstallScript() string {
	return m.InstallScript.ValueString()
}

func (m *ResourceScriptModel) GetFindInstalledScript() string {
	return m.FindInstalledScript.ValueString()
}

func (m *ResourceScriptModel) GetUninstallScript() string {
	return m.UninstallScript.ValueString()
}

func (m *ResourceScriptModel) GetAdditionalArgs(ctx context.Context) []string {
	var args []string
	m.AdditionalArgs.ElementsAs(ctx, args, false)
	return args
}

func (m *ResourceScriptModel) GetDefaultArgs(ctx context.Context) []string {
	var args []string
	m.DefaultArgs.ElementsAs(ctx, args, false)
	return args
}

func (m *ResourceScriptModel) GetSudo() bool {
	return m.Sudo.ValueBool()
}

func (m *ResourceScriptModel) GetShell() string {
	return m.Shell.ValueString()
}

func (m *ResourceScriptModel) Initialize() bool {
	idString := m.Path
	if m.Path.IsNull() {
		idString = m.FindInstalledScript
	}
	m.Id = sources.GetIDFromName(idString, enums.InstallerScript)
	return !idString.IsNull()
}

func (m *ResourceScriptModel) CopyFromTypedInstalledProgramInfo(installedInfo *models.TypedInstalledProgramInfo) {
	m.Path = types.StringValue(installedInfo.Path)
}

// ResourceScript defines the resource implementation.
type ResourceScript struct {
	*Resource[*ResourceScriptModel]
}

func NewResourceScript() resource.Resource {
	resource := &ResourceScript{}
	resource.Resource = NewResource[*ResourceScriptModel](script.NewScriptInstaller[*ResourceScriptModel](resource))
	return resource
}

func (r *ResourceScript) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: schemastrings.ScriptSourceDescription,
		Attributes: map[string]schema.Attribute{
			"id":                    defaults.GetIdSchema(),
			"path":                  defaults.GetScriptPathSchema(schemastrings.ScriptPathDescription),
			"install_script":        defaults.GetInstallScriptSchema(schemastrings.ScriptInstallScriptDescription),
			"find_installed_script": defaults.GetFindInstalledScriptSchema(schemastrings.ScriptFindInstalledScriptDescription),
			"uninstall_script":      defaults.GetUninstallScriptSchema(schemastrings.ScriptUninstallScriptDescription),
			"additional_args":       defaults.GetAdditionalArgsSchema(schemastrings.ScriptAdditionalArgsDescription),
			"default_args":          defaults.GetDefaultArgsSchema(schemastrings.ScriptDefaultArgsDescription, script.DefaultArg),
			"sudo":                  defaults.GetSudoSchema(script.DefaultSudo),
			"shell":                 defaults.GetShellSchema(schemastrings.ScriptShellDescription, script.DefaultProgram),
		},
	}
}
