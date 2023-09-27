package defaults

import (
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator/shared"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/configs/configschema"
	"github.com/shihanng/terraform-provider-installer/internal/terraformutils"
	"github.com/zclconf/go-cty/cty"
)

func GetRemoteConnectionBlockSchema() schema.SingleNestedBlock {
	block := ConvertConfigSchemaBlockToSchemaBlock(shared.ConnectionBlockSupersetSchema)
	block.MarkdownDescription = terraformutils.RemoteConnectionBlockDescription
	return block
}

func ConvertConfigSchemaBlockToSchemaBlock(config *configschema.Block) schema.SingleNestedBlock {
	block := schema.SingleNestedBlock{
		Attributes: map[string]schema.Attribute{},
	}
	for name, attr := range config.Attributes {
		block.Attributes[name] = ConvertConfigSchemaAttrToSchemaAttr(attr, name)
	}
	return block
}

func ConvertConfigSchemaAttrToSchemaAttr(config *configschema.Attribute, name string) schema.Attribute {
	switch config.Type {
	case cty.String:
		return schema.StringAttribute{
			Optional:  true,
			Sensitive: terraformutils.IsNameSensitive(name),
		}
	case cty.Number:
		return schema.Int64Attribute{
			Optional:  true,
			Sensitive: terraformutils.IsNameSensitive(name),
		}
	case cty.Bool:
		return schema.BoolAttribute{
			Optional:  true,
			Sensitive: terraformutils.IsNameSensitive(name),
		}
	default:
		panic("unsupported type")
	}
}
