package defaults

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getDefaultStringSchema(markdownDescription string, optional bool, requiresReplace bool) schema.Attribute {
	schma := schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Optional:            optional,
	}
	if requiresReplace {
		schma.PlanModifiers = []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		}
	}
	return schma
}

func GetIdSchema() schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: "Internal ID of the resource.",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func GetNameSchema(markdownDescription string) schema.Attribute {
	return getDefaultStringSchema(markdownDescription, true, true)
}

func GetPathSchema(markdownDescription string) schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Computed:            true,
	}
}

func GetScriptPathSchema(markdownDescription string) schema.Attribute {
	return getDefaultStringSchema(markdownDescription, false, true)
}

func GetVersionSchema(markdownDescription string) schema.Attribute {
	return getDefaultStringSchema(markdownDescription, true, true)
}

func GetInstallScriptSchema(markdownDescription string) schema.Attribute {
	return getDefaultStringSchema(markdownDescription, false, true)
}

func GetFindInstalledScriptSchema(markdownDescription string) schema.Attribute {
	return getDefaultStringSchema(markdownDescription, true, true)
}

func GetUninstallScriptSchema(markdownDescription string) schema.Attribute {
	return getDefaultStringSchema(markdownDescription, false, true)
}

func GetAdditionalArgsSchema(markdownDescription string) schema.ListAttribute {
	return schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.RequiresReplace(),
		},
	}
}
