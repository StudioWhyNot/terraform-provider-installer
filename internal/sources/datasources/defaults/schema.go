package defaults

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func getDefaultStringSchema(markdownDescription string, optional bool) schema.StringAttribute {
	schma := schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Optional:            optional,
		Required:            !optional,
	}
	return schma
}

func GetNameSchema(markdownDescription string) schema.StringAttribute {
	return getDefaultStringSchema(markdownDescription, false)
}

func GetPathSchema(markdownDescription string) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Computed:            true,
	}
}

func GetVersionSchema(markdownDescription string) schema.StringAttribute {
	return getDefaultStringSchema(markdownDescription, true)
}

func GetSudoSchema() schema.BoolAttribute {
	return schema.BoolAttribute{
		MarkdownDescription: "Whether or not to run the installer as a sudo user.",
		Optional:            true,
		Computed:            true,
	}
}
