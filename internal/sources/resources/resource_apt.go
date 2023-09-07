// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
	"github.com/shihanng/terraform-provider-installer/internal/installers/apt"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/sources"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ResourceApt{}
var _ resource.ResourceWithImportState = &ResourceApt{}

// ResourceApt defines the resource implementation.
type ResourceApt struct {
	sources.SourceBase
	client    *http.Client
	installer installers.Installer
}

// ResourceAptModel describes the resource data model.
type ResourceAptModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
}

func NewResourceApt() resource.Resource {
	source := ResourceApt{}
	source.InstallerType = enums.InstallerApt
	source.installer = apt.NewAptInstaller()
	return &source
}

func (r *ResourceApt) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.GetSourceName(req.ProviderTypeName)
}

func (r *ResourceApt) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "`installer_apt` manages an application using [APT](https://en.wikipedia.org/wiki/APT_(software)).\n\n" +
			"It works on systems that use APT as the package management system. " +
			"Adding an `installer_apt` resource means that Terraform will ensure that " +
			"the application defined in the `name` argument is made available via APT.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Internal ID of the resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the application that `apt-get` recognizes. Specify a version of a package by following the package name with an equal sign and the version, e.g., `vim=2:8.2.3995-1ubuntu2.7`.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The path where the application is installed by `apt-get` after Terraform creates this resource.",
				Computed:            true,
			},
		},
	}
}

func (r *ResourceApt) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ResourceApt) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data, success := sources.TryGetData[ResourceAptModel](ctx, req.Plan, &resp.Diagnostics)
	if !success {
		return
	}

	var diags diag.Diagnostics
	options, err := models.GetOptions(data.Name.ValueString())
	if err != nil {
		diags = xerrors.ToDiags(err)
		resp.Diagnostics.Append(diags...)
		return
	}
	if err := r.installer.Install(ctx, options); err != nil {
		diags = xerrors.ToDiags(err)
		resp.Diagnostics.Append(diags...)
	}

	data.Id = types.StringValue(r.GetIDFromName(options.Name))
	if !r.UpdateFromInstallation(&data, ctx, &resp.Diagnostics) {
		resp.State.RemoveResource(ctx)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "Created an "+r.InstallerType.String()+" resource")

	sources.SetStateData(ctx, &resp.State, &resp.Diagnostics, &data)
}

func (r *ResourceApt) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data, success := sources.TryGetData[ResourceAptModel](ctx, req.State, &resp.Diagnostics)
	if !success {
		return
	}

	if !r.UpdateFromInstallation(&data, ctx, &resp.Diagnostics) {
		resp.State.RemoveResource(ctx)
		return
	}

	sources.SetStateData(ctx, &resp.State, &resp.Diagnostics, &data)
}

func (r *ResourceApt) UpdateFromInstallation(data *ResourceAptModel, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	options, err := models.GetOptions(data.Name.ValueString())
	if err != nil {
		diags := xerrors.ToDiags(err)
		diagnostics.Append(diags...)
		return false
	}
	info, err := r.installer.FindInstalled(ctx, options)
	if err != nil {
		if errors.Is(err, xerrors.ErrNotInstalled) {
			return false
		}
		diags := xerrors.ToDiags(err)
		diagnostics.Append(diags...)
	}

	data.Name = types.StringValue(options.Name)
	data.Path = types.StringValue(info.Path)
	return true
}

func (r *ResourceApt) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ResourceAptModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	diags = req.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *ResourceApt) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ResourceAptModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	options, err := models.GetOptions(data.Name.ValueString())
	if err != nil {
		diags := xerrors.ToDiags(err)
		diags.Append(diags...)
		return
	}
	if _, err := r.installer.Uninstall(ctx, options); err != nil {
		diags = xerrors.ToDiags(err)
		resp.Diagnostics.Append(diags...)
	}
	resp.State.RemoveResource(ctx)
}

func (r *ResourceApt) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
