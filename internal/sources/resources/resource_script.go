// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers/script"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
	"github.com/shihanng/terraform-provider-installer/internal/sources/resources/defaults"
	"github.com/shihanng/terraform-provider-installer/internal/sources/schemastrings"
	"github.com/shihanng/terraform-provider-installer/internal/terraformutils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ResourceScript{}
var _ resource.ResourceWithImportState = &ResourceScript{}
var _ sources.SourceData = &ResourceScriptModel{}

// ResourceScriptModel describes the resource data model.
type ResourceScriptModel struct {
	Id                                   types.String `tfsdk:"id"`
	Path                                 types.String `tfsdk:"path"`
	Script                               types.String `tfsdk:"script"`
	InstallScript                        types.String `tfsdk:"install_script"`
	FindInstalledScript                  types.String `tfsdk:"find_installed_script"`
	UninstallScript                      types.String `tfsdk:"uninstall_script"`
	DefaultArgs                          types.List   `tfsdk:"default_args"`
	AdditionalArgs                       types.List   `tfsdk:"additional_args"`
	Sudo                                 types.Bool   `tfsdk:"sudo"`
	Environment                          types.Map    `tfsdk:"environment"`
	Shell                                types.String `tfsdk:"shell"`
	*terraformutils.RemoteConnectionInfo `tfsdk:"remote_connection"`
}

func (m *ResourceScriptModel) GetId() string {
	return m.Id.ValueString()
}

func (m *ResourceScriptModel) GetPath() string {
	return m.Path.ValueString()
}

func (m *ResourceScriptModel) GetScript() string {
	return m.Script.ValueString()
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
	return sources.ListValueToList[string](ctx, &m.AdditionalArgs)
}

func (m *ResourceScriptModel) GetDefaultArgs(ctx context.Context) []string {
	return sources.ListValueToList[string](ctx, &m.DefaultArgs)
}

func (m *ResourceScriptModel) GetSudo() bool {
	return m.Sudo.ValueBool()
}

func (m *ResourceScriptModel) GetEnvironment(ctx context.Context) map[string]string {
	return sources.MapValueToMap(ctx, &m.Environment)
}

func (m *ResourceScriptModel) GetShell() string {
	return m.Shell.ValueString()
}

func (m *ResourceScriptModel) Initialize() bool {
	scriptString := m.GetPath() + m.GetInstallScript() + m.GetFindInstalledScript() + m.GetUninstallScript()
	hash := sha256.Sum256([]byte(scriptString))
	hashString := hex.EncodeToString(hash[:])

	m.Id = sources.GetIDFromName(hashString, enums.InstallerScript)
	return !m.Id.IsNull()
}

func (m *ResourceScriptModel) GetRemoteConnectionInfo() *terraformutils.RemoteConnectionInfo {
	return m.RemoteConnectionInfo
}

func (m *ResourceScriptModel) CopyFromTypedInstalledProgramInfo(installedInfo *models.TypedInstalledProgramInfo) {
	if installedInfo == nil {
		m.Path = types.StringNull()
		return
	}
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
			"script":                defaults.GetScriptSchema(schemastrings.ScriptScriptDescription),
			"install_script":        defaults.GetInstallScriptSchema(schemastrings.ScriptInstallScriptDescription),
			"find_installed_script": defaults.GetFindInstalledScriptSchema(schemastrings.ScriptFindInstalledScriptDescription),
			"uninstall_script":      defaults.GetUninstallScriptSchema(schemastrings.ScriptUninstallScriptDescription),
			"additional_args":       defaults.GetAdditionalArgsSchema(schemastrings.ScriptAdditionalArgsDescription),
			"default_args":          defaults.GetDefaultArgsSchema(schemastrings.ScriptDefaultArgsDescription, script.DefaultArg),
			"sudo":                  defaults.GetSudoSchema(script.DefaultSudo),
			"environment":           defaults.GetEnvironmentSchema(),
			"shell":                 defaults.GetShellSchema(schemastrings.ScriptShellDescription, script.DefaultProgram),
		},
		Blocks: map[string]schema.Block{
			"remote_connection": defaults.GetRemoteConnectionBlockSchema(),
		},
	}
}
