// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/installers/script"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	providerdefaults "github.com/shihanng/terraform-provider-installer/internal/provider/defaults"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
	"github.com/shihanng/terraform-provider-installer/internal/sources/datasources/defaults"
	"github.com/shihanng/terraform-provider-installer/internal/sources/schemastrings"
	"github.com/shihanng/terraform-provider-installer/internal/system"
	"github.com/shihanng/terraform-provider-installer/internal/terraformutils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DataSourceScript{}
var _ sources.SourceData = &DataSourceScriptModel{}

// DataSourceScriptModel describes the resource data model.
type DataSourceScriptModel struct {
	Script                               types.String `tfsdk:"script"`
	FindInstalledScript                  types.String `tfsdk:"find_installed_script"`
	DefaultArgs                          types.List   `tfsdk:"default_args"`
	AdditionalArgs                       types.List   `tfsdk:"additional_args"`
	Sudo                                 types.Bool   `tfsdk:"sudo"`
	Environment                          types.Map    `tfsdk:"environment"`
	Shell                                types.String `tfsdk:"shell"`
	Secrets                              types.Map    `tfsdk:"secrets"`
	Output                               types.String `tfsdk:"output"`
	*terraformutils.RemoteConnectionInfo `tfsdk:"remote_connection"`
}

func (m *DataSourceScriptModel) GetId() string {
	return ""
}

func (m *DataSourceScriptModel) GetPath() string {
	return ""
}

func (m *DataSourceScriptModel) GetScript() string {
	return m.Script.ValueString()
}

func (m *DataSourceScriptModel) GetInstallScript() string {
	return ""
}

func (m *DataSourceScriptModel) GetFindInstalledScript() string {
	return m.FindInstalledScript.ValueString()
}

func (m *DataSourceScriptModel) GetUninstallScript() string {
	return ""
}

func (m *DataSourceScriptModel) GetAdditionalArgs(ctx context.Context) []string {
	return sources.ListValueToList[string](ctx, &m.AdditionalArgs)
}

func (m *DataSourceScriptModel) GetDefaultArgs(ctx context.Context) []string {
	defaultArgs := sources.ListValueToList[string](ctx, &m.DefaultArgs)
	if len(defaultArgs) == 0 {
		defaultArgs = append(defaultArgs, script.DefaultArg)
	}
	return defaultArgs
}

func (m *DataSourceScriptModel) GetSudo() bool {
	return m.Sudo.ValueBool()
}

func (m *DataSourceScriptModel) GetEnvironmentAndSecrets(ctx context.Context) map[string]string {
	return system.MergeMaps(sources.MapValueToMap(ctx, &m.Environment), sources.MapValueToMap(ctx, &m.Secrets))
}

func (m *DataSourceScriptModel) GetShell() string {
	shell := m.Shell.ValueString()
	if shell == "" {
		shell = script.DefaultProgram
	}
	return shell
}

func (m *DataSourceScriptModel) SetOutput(output string) {
	m.Output = types.StringValue(output)
}

func (m *DataSourceScriptModel) Initialize() bool {
	return !m.Script.IsNull()
}

func (m *DataSourceScriptModel) GetRemoteConnectionInfo() *terraformutils.RemoteConnectionInfo {
	return m.RemoteConnectionInfo
}

func (m *DataSourceScriptModel) CopyFromTypedInstalledProgramInfo(installedInfo *models.TypedInstalledProgramInfo) {
}

// DataSourceScript defines the resource implementation.
type DataSourceScript struct {
	*DataSource[*DataSourceScriptModel]
}

func NewDataSourceScript() datasource.DataSource {
	resource := &DataSourceScript{}
	resource.DataSource = NewDataSource[*DataSourceScriptModel](script.NewScriptInstaller[*DataSourceScriptModel](resource))
	return resource
}

func (r *DataSourceScript) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: schemastrings.ScriptSourceDescription,
		Attributes: map[string]schema.Attribute{
			"script":                defaults.GetScriptSchema(schemastrings.ScriptScriptDescription),
			"find_installed_script": defaults.GetFindInstalledScriptSchema(schemastrings.ScriptFindInstalledScriptDescription),
			"additional_args":       defaults.GetAdditionalArgsSchema(schemastrings.ScriptAdditionalArgsDescription),
			"default_args":          defaults.GetDefaultArgsSchema(schemastrings.ScriptDefaultArgsDescription),
			"sudo":                  defaults.GetSudoSchema(),
			"environment":           defaults.GetEnvironmentSchema(),
			"secrets":               defaults.GetSecretsSchema(),
			"shell":                 defaults.GetShellSchema(schemastrings.ScriptShellDescription),
			"output":                defaults.GetOutputSchema(schemastrings.ScriptOutputDescription),
		},
		Blocks: map[string]schema.Block{
			"remote_connection": providerdefaults.GetRemoteConnectionBlockSchema(),
		},
	}
}
