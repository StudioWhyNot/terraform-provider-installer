package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
)

// Resource[T] defines the generic data source implementation.
type Resource[T sources.SourceData] struct {
	sources.SourceBase[T]
}

func NewResource[T sources.SourceData](installer installers.Installer[T]) *Resource[T] {
	return &Resource[T]{
		SourceBase: *sources.NewSourceBase[T](installer),
	}
}

func (d *Resource[T]) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = d.Installer.GetInstallerType().GetIDFromName(req.ProviderTypeName)
}

func (r *Resource[T]) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Nothing to configure
}

func (r *Resource[T]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	sources.DefaultCreate[T](&r.SourceBase, req.Plan, &resp.State, ctx, &resp.Diagnostics)
}

func (r *Resource[T]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	sources.DefaultRead[T](&r.SourceBase, &resp.State, ctx, &resp.Diagnostics)
}

func (r *Resource[T]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	sources.DefaultUpdate[T](&r.SourceBase, &resp.State, ctx, &resp.Diagnostics)
}

func (r *Resource[T]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	sources.DefaultDelete[T](&r.SourceBase, &resp.State, ctx, &resp.Diagnostics)
}

func (r *Resource[T]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
