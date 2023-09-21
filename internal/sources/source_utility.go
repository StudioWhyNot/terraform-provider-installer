package sources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func DefaultCreate[T SourceData](source *SourceBase[T], plan tfsdk.Plan, state *tfsdk.State, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	data, success := TryGetData[T](ctx, plan, diagnostics)
	if !success || !data.Initialize() {
		//TODO: See if this needs an error in all cases
		return false
	}

	err := source.Installer.Install(ctx, data)
	if err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
		state.RemoveResource(ctx)
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "Created resource of type: "+source.Installer.GetInstallerType().String())
	SetStateData(ctx, state, diagnostics, &data)
	return true
}

func DefaultRead[T SourceData](source *SourceBase[T], state *tfsdk.State, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	data, success := TryGetData[T](ctx, state, diagnostics)
	if !success || !data.Initialize() {
		state.RemoveResource(ctx)
		return false
	}

	SetStateData(ctx, state, diagnostics, &data)
	return true
}

func DefaultUpdate[T SourceData](source *SourceBase[T], state *tfsdk.State, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	data, success := TryGetData[T](ctx, state, diagnostics)
	if !success {
		return false
	}

	// Save updated data into Terraform state
	diags := state.Set(ctx, &data)
	diagnostics.Append(diags...)
	return true
}

func DefaultDelete[T SourceData](source *SourceBase[T], state *tfsdk.State, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	data, success := TryGetData[T](ctx, state, diagnostics)
	if !success {
		return false
	}

	if _, err := source.Installer.Uninstall(ctx, data); err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
	}
	state.RemoveResource(ctx)
	return true
}
