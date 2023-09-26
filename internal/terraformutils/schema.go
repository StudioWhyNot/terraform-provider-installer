package terraformutils

import (
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/configs/configschema"
	"github.com/zclconf/go-cty/cty"
)

func ConvertConfigSchemaBlockToSchemaBlock(config *configschema.Block) schema.SingleNestedBlock {
	block := schema.SingleNestedBlock{
		Attributes: map[string]schema.Attribute{},
	}
	for name, attr := range config.Attributes {
		block.Attributes[name] = ConvertConfigSchemaAttrToSchemaAttr(attr)
	}
	return block
}

func ConvertConfigSchemaAttrToSchemaAttr(config *configschema.Attribute) schema.Attribute {
	switch config.Type {
	case cty.String:
		return schema.StringAttribute{
			Optional: true,
		}
	case cty.Number:
		return schema.Int64Attribute{
			Optional: true,
		}
	case cty.Bool:
		return schema.BoolAttribute{
			Optional: true,
		}
	default:
		panic("unsupported type")
	}
}
