package defaults

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/provider/defaults"
	"github.com/shihanng/terraform-provider-installer/internal/sources/schemastrings"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator/shared"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/configs/configschema"
	"github.com/shihanng/terraform-provider-installer/internal/terraformutils"
)

func getDefaultStringListSchema(markdownDescription string, optional bool) schema.ListAttribute {
	return schema.ListAttribute{
		ElementType:         types.StringType,
		MarkdownDescription: markdownDescription,
		Optional:            optional,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.RequiresReplace(),
		},
	}
}

func getDefaultStringSchema(markdownDescription string, optional bool, requiresReplace bool) schema.StringAttribute {
	schma := schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Optional:            optional,
		Required:            !optional,
	}
	if requiresReplace {
		schma.PlanModifiers = []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		}
	}
	return schma
}

func getDefaultBoolSchema(markdownDescription string, defaultVal bool, requiresReplace bool) schema.BoolAttribute {
	schma := schema.BoolAttribute{
		MarkdownDescription: markdownDescription,
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(defaultVal),
	}
	if requiresReplace {
		schma.PlanModifiers = []planmodifier.Bool{
			boolplanmodifier.RequiresReplace(),
		}
	}
	return schma
}

func GetIdSchema() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: schemastrings.DefaultIdDescription,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func GetNameSchema(markdownDescription string) schema.StringAttribute {
	return getDefaultStringSchema(markdownDescription, true, true)
}

func GetPathSchema(markdownDescription string) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Computed:            true,
	}
}

func GetScriptPathSchema(markdownDescription string) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Computed:            true,
		Optional:            true,
	}
}

func GetVersionSchema(markdownDescription string) schema.StringAttribute {
	return getDefaultStringSchema(markdownDescription, true, true)
}

func GetSudoSchema(defaultVal bool) schema.BoolAttribute {
	return getDefaultBoolSchema(schemastrings.DefaultSudoDescription, defaultVal, true)
}

func GetCaskSchema(markdownDescription string, defaultVal bool) schema.BoolAttribute {
	return getDefaultBoolSchema(markdownDescription, defaultVal, true)
}

func GetInstallScriptSchema(markdownDescription string) schema.StringAttribute {
	return getDefaultStringSchema(markdownDescription, false, true)
}

func GetFindInstalledScriptSchema(markdownDescription string) schema.StringAttribute {
	return getDefaultStringSchema(markdownDescription, true, true)
}

func GetUninstallScriptSchema(markdownDescription string) schema.StringAttribute {
	return getDefaultStringSchema(markdownDescription, false, true)
}

func GetScriptSchema(markdownDescription string) schema.StringAttribute {
	return getDefaultStringSchema(markdownDescription, false, true)
}

func GetAdditionalArgsSchema(markdownDescription string) schema.ListAttribute {
	return getDefaultStringListSchema(markdownDescription, true)
}

func GetEnvironmentSchema() schema.MapAttribute {
	return schema.MapAttribute{
		ElementType:         types.StringType,
		MarkdownDescription: schemastrings.DefaultEnvironmentDescription,
		Optional:            true,
		PlanModifiers: []planmodifier.Map{
			mapplanmodifier.RequiresReplace(),
		},
	}
}

func GetDefaultArgsSchema(markdownDescription string, defaultArg string) schema.ListAttribute {
	args, _ := types.ListValue(types.StringType, []attr.Value{types.StringValue(defaultArg)})
	return schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.RequiresReplace(),
		},
		Default: listdefault.StaticValue(args),
	}
}

func GetShellSchema(markdownDescription string, defaultVal string) schema.StringAttribute {
	schma := getDefaultStringSchema(markdownDescription, true, false)
	schma.Computed = true
	schma.Default = stringdefault.StaticString(defaultVal)
	return schma
}

func GetConnectionNameSchema() schema.StringAttribute {
	schma := getDefaultStringSchema(schemastrings.DefaultConnectionNameDescription, false, true)
	schma.Required = false
	schma.Computed = true
	return schma
}

func GetRemoteConnectionBlockSchema() schema.SingleNestedBlock {
	block := convertConfigSchemaBlockToSchemaBlock(shared.ConnectionBlockSupersetSchema)
	block.MarkdownDescription = terraformutils.RemoteConnectionBlockDescription
	return block
}

func convertConfigSchemaBlockToSchemaBlock(config *configschema.Block) schema.SingleNestedBlock {
	block := schema.SingleNestedBlock{
		Attributes: map[string]schema.Attribute{},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
	}
	for name, attr := range config.Attributes {
		block.Attributes[name] = defaults.ConvertConfigSchemaAttrToSchemaAttr(attr, name)
	}
	return block
}
