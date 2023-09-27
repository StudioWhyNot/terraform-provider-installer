package sources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shihanng/terraform-provider-installer/internal/terraformutils"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

func TryGetData[T any](ctx context.Context, provider TerraformDataProvider, diagnostics *diag.Diagnostics) (T, bool) {
	var data T
	diags := provider.Get(ctx, &data)
	diagnostics.Append(diags...)
	return data, !diagnostics.HasError()
}

func TryGetInitializedData[T SourceData](ctx context.Context, provider TerraformDataProvider, diagnostics *diag.Diagnostics) (T, bool) {
	data, success := TryGetData[T](ctx, provider, diagnostics)
	return data, success && data.Initialize()
}

func SetStateData(ctx context.Context, state *tfsdk.State, diagnostics *diag.Diagnostics, val interface{}) {
	diags := state.Set(ctx, val)
	diagnostics.Append(diags...)
}

func FillAndSetStateData[T SourceData](source *SourceBase[T], ctx context.Context, state *tfsdk.State, diagnostics *diag.Diagnostics, data T) {
	SetCommunicator(source, data.GetRemoteConnectionInfo(), diagnostics)
	err := source.TryConnect()
	if err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
		return
	}
	info, err := source.Installer.FindInstalled(ctx, data)
	if err != nil {
		state.RemoveResource(ctx)
	}
	if info != nil {
		data.CopyFromTypedInstalledProgramInfo(info)
	}
	err = source.TryDisconnect()
	if err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
	}
	SetStateData(ctx, state, diagnostics, &data)
}

func DefaultCreate[T SourceData](source *SourceBase[T], plan tfsdk.Plan, state *tfsdk.State, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	data, success := TryGetInitializedData[T](ctx, plan, diagnostics)
	if !success {
		return false
	}

	SetCommunicator(source, data.GetRemoteConnectionInfo(), diagnostics)
	err := source.TryConnect()
	if err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
		return false
	}

	err = source.Installer.Install(ctx, data)
	if err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
		state.RemoveResource(ctx)
		return false
	}

	err = source.TryDisconnect()
	if err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "Created resource of type: "+source.Installer.GetInstallerType().String())
	FillAndSetStateData(source, ctx, state, diagnostics, data)
	return true
}

func DefaultConfigure[T SourceData](source *SourceBase[T], providerData any, diagnostics *diag.Diagnostics) {
	if providerData == nil {
		return
	}
	connInfo := providerData.(*terraformutils.RemoteConnectionInfo)
	communicator, err := terraformutils.MakeCommunicator(connInfo)
	if err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
	}
	source.Communicator = communicator
}

func DefaultRead[T SourceData](source *SourceBase[T], state *tfsdk.State, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	data, success := TryGetInitializedData[T](ctx, state, diagnostics)
	if !success {
		state.RemoveResource(ctx)
		return false
	}

	FillAndSetStateData(source, ctx, state, diagnostics, data)
	return true
}

func DefaultUpdate[T SourceData](source *SourceBase[T], plan tfsdk.Plan, state *tfsdk.State, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	data, success := TryGetInitializedData[T](ctx, plan, diagnostics)
	if !success {
		return false
	}

	FillAndSetStateData(source, ctx, state, diagnostics, data)
	// Save updated data into Terraform state
	diags := state.Set(ctx, &data)
	diagnostics.Append(diags...)
	return true
}

func DefaultDelete[T SourceData](source *SourceBase[T], state *tfsdk.State, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	data, success := TryGetInitializedData[T](ctx, state, diagnostics)
	if !success {
		return false
	}

	SetCommunicator(source, data.GetRemoteConnectionInfo(), diagnostics)
	err := source.TryConnect()
	if err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
		return false
	}

	if _, err := source.Installer.Uninstall(ctx, data); err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
	}
	state.RemoveResource(ctx)

	err = source.TryDisconnect()
	if err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
	}
	return true
}

func SetCommunicator[T SourceData](source *SourceBase[T], connInfo *terraformutils.RemoteConnectionInfo, diagnostics *diag.Diagnostics) {
	if connInfo == nil {
		return
	}
	communicator, err := terraformutils.MakeCommunicator(connInfo)
	if err != nil {
		xerrors.AppendToDiagnostics(diagnostics, err)
	}
	source.Communicator = communicator
}
