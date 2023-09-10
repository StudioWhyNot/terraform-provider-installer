package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
)

// DataSource[T] defines the generic data source implementation.
type DataSource[T sources.SourceData] struct {
	sources.SourceBase
}

func NewDataSource[T sources.SourceData](installerType enums.InstallerType) *DataSource[T] {
	return &DataSource[T]{
		SourceBase: *sources.NewSourceBase(installerType),
	}
}

func (d *DataSource[T]) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.GetDefaultTypeName(req.ProviderTypeName)
}

func (d *DataSource[T]) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	//Nothing to configure.
}

func (d *DataSource[T]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	sources.DefaultRead[T](&d.SourceBase, &resp.State, ctx, &resp.Diagnostics)
	tflog.Trace(ctx, "Read a data source")
}
