package defaults

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/sources/schemastrings"
)

func getDefaultStringSchema(markdownDescription string, optional bool) schema.StringAttribute {
	schma := schema.StringAttribute{
		MarkdownDescription: markdownDescription,
		Optional:            optional,
		Required:            !optional,
	}
	return schma
}

func getDefaultBoolSchema(markdownDescription string, optional bool) schema.BoolAttribute {
	schma := schema.BoolAttribute{
		MarkdownDescription: markdownDescription,
		Optional:            optional,
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
	return getDefaultBoolSchema(schemastrings.DefaultSudoDescription, true)
}

func GetCaskSchema() schema.BoolAttribute {
	return getDefaultBoolSchema(schemastrings.BrewCaskDescription, true)
}

func GetConnectionNameSchema() schema.StringAttribute {
	schma := getDefaultStringSchema(schemastrings.DefaultConnectionNameDescription, false)
	schma.Required = false
	schma.Computed = true
	return schma
}

func GetEnvironmentSchema() schema.MapAttribute {
	return schema.MapAttribute{
		ElementType:         types.StringType,
		MarkdownDescription: schemastrings.DefaultEnvironmentDescription,
		Optional:            true,
	}
}

func GetSecretsSchema() schema.MapAttribute {
	schema := GetEnvironmentSchema()
	schema.Sensitive = true
	return schema
}
