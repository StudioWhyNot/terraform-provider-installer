// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/shihanng/terraform-provider-installer/internal/sources/datasources"
	"github.com/shihanng/terraform-provider-installer/internal/sources/resources"
)

// Ensure InstallerProvider satisfies various provider interfaces.
var _ provider.Provider = &InstallerProvider{}

// InstallerProvider defines the provider implementation.
type InstallerProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// The prefix name of the provider.
const ProviderName = "installer"

// InstallerProviderModel describes the provider data model.
type InstallerProviderModel struct {
}

func (p *InstallerProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = ProviderName
	resp.Version = p.version
}

func (p *InstallerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	// No provider-level schema.
}

func (p *InstallerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// No provider-level configuration.
}

func (p *InstallerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewResourceApt,
	}
}

func (p *InstallerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewDataSourceApt,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &InstallerProvider{
			version: version,
		}
	}
}
