package defaults

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

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
	return schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}

func GetPathSchema(markdownDescription string) schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Computed:            true,
	}
}
