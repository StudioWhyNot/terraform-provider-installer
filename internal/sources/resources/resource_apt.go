// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"context"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
	"github.com/shihanng/terraform-provider-installer/internal/sources/resources/defaults"
	"github.com/shihanng/terraform-provider-installer/internal/sources/schemastrings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ResourceApt{}
var _ resource.ResourceWithImportState = &ResourceApt{}
var _ sources.SourceData = &ResourceAptModel{}

// ResourceAptModel describes the resource data model.
type ResourceAptModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
}

func (m *ResourceAptModel) GetName() types.String {
	return m.Name
}

func (m *ResourceAptModel) GetVersion() *version.Version {
	return nil
}

func (m *ResourceAptModel) SetDataFromTypedInstalledProgramInfo(info *models.TypedInstalledProgramInfo) {
	m.Id = types.StringValue(info.InstallerType.GetIDFromName(info.Name))
	m.Name = types.StringValue(info.Name)
	m.Path = types.StringValue(info.Path)
}

// ResourceApt defines the resource implementation.
type ResourceApt struct {
	*Resource[*ResourceAptModel]
}

func NewResourceApt() resource.Resource {
	return &ResourceApt{
		Resource: NewResource[*ResourceAptModel](enums.InstallerApt),
	}
}

func (r *ResourceApt) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.GetDefaultTypeName(req.ProviderTypeName)
}

func (r *ResourceApt) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: schemastrings.AptSourceDescription,
		Attributes: map[string]schema.Attribute{
			"id":   defaults.GetIdSchema(),
			"name": defaults.GetNameSchema(schemastrings.AptNameDescription),
			"path": defaults.GetPathSchema(schemastrings.AptPathDescription),
		},
	}
}
