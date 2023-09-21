package defaults

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func GetNameSchema(markdownDescription string) schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Required:            true,
	}
}

func GetPathSchema(markdownDescription string) schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Computed:            true,
	}
}

func GetVersionSchema(markdownDescription string) schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Optional:            true,
	}
}
