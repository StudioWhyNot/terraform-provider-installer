package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
)

// DataSource[T] defines the generic data source implementation.
type DataSource[T sources.SourceData] struct {
	sources.SourceBase[T]
}

func NewDataSource[T sources.SourceData](installer installers.Installer[T]) *DataSource[T] {
	return &DataSource[T]{
		SourceBase: *sources.NewSourceBase[T](installer),
	}
}

func (d *DataSource[T]) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.Installer.GetInstallerType().GetSourceName(req.ProviderTypeName)
}

func (d *DataSource[T]) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	sources.DefaultConfigure[T](&d.SourceBase, req.ProviderData, &resp.Diagnostics)
}

func (d *DataSource[T]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	sources.DefaultRead[T](&d.SourceBase, &resp.State, ctx, &resp.Diagnostics)
}
