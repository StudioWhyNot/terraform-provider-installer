package sources

import (
	"context"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

func TryGetData[T any](ctx context.Context, provider TerraformDataProvider, diagnostics *diag.Diagnostics) (T, bool) {
	var data T
	diags := provider.Get(ctx, &data)
	diagnostics.Append(diags...)
	return data, !diagnostics.HasError()
}

func SetStateData(ctx context.Context, state *tfsdk.State, diagnostics *diag.Diagnostics, val interface{}) {
	diags := state.Set(ctx, val)
	diagnostics.Append(diags...)
}

func getLoggedOptions(data SourceData, diagnostics *diag.Diagnostics) *models.InstallerOptions {
	options, err := models.NewInstallerOptions(data.GetName().ValueString(), nil)
	if err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
		return nil
	}
	return &options
}

func getInstallationInfo(fullName string, version *version.Version, installer installers.Installer, ctx context.Context) (*models.TypedInstalledProgramInfo, error) {
	options, err := models.NewInstallerOptions(fullName, version)
	if err != nil {
		return nil, err
	}
	return installer.FindInstalled(ctx, options)
}

func GetInstallationInfo(data SourceData, installer installers.Installer, ctx context.Context, diagnostics *diag.Diagnostics) *models.TypedInstalledProgramInfo {
	info, err := getInstallationInfo(data.GetName().ValueString(), data.GetVersion(), installer, ctx)
	if err != nil {
		diags := xerrors.ToDiags(err)
		diagnostics.Append(diags...)
		return nil
	}
	return info
}

func DefaultCreate[T SourceData](source *SourceBase, plan tfsdk.Plan, state *tfsdk.State, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	data, success := TryGetData[T](ctx, plan, diagnostics)
	if !success {
		//TODO: See if this needs an error in all cases
		return false
	}

	options := getLoggedOptions(data, diagnostics)
	if options == nil {
		return false
	}
	err := source.Installer.Install(ctx, *options)
	if err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
	}
	if !source.UpdateFromInstallation(data, ctx, diagnostics) {
		state.RemoveResource(ctx)
		return false
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "Created an "+source.Installer.GetInstallerType().String()+" resource")
	SetStateData(ctx, state, diagnostics, &data)
	return true
}

func DefaultRead[T SourceData](source *SourceBase, state *tfsdk.State, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	data, success := TryGetData[T](ctx, state, diagnostics)
	if !success {
		return false
	}

	if !source.UpdateFromInstallation(data, ctx, diagnostics) {
		state.RemoveResource(ctx)
		return false
	}

	SetStateData(ctx, state, diagnostics, &data)
	return true
}

func DefaultUpdate[T SourceData](source *SourceBase, state *tfsdk.State, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	data, success := TryGetData[T](ctx, state, diagnostics)
	if !success {
		return false
	}

	// Save updated data into Terraform state
	diags := state.Set(ctx, &data)
	diagnostics.Append(diags...)
	return true
}

func DefaultDelete[T SourceData](source *SourceBase, state *tfsdk.State, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	data, success := TryGetData[T](ctx, state, diagnostics)
	if !success {
		return false
	}
	options := getLoggedOptions(data, diagnostics)
	if options == nil {
		return false
	}
	if _, err := source.Installer.Uninstall(ctx, *options); err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
	}
	state.RemoveResource(ctx)
	return true
}
