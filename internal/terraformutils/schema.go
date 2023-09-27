package terraformutils

import (
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator/shared"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/configs/configschema"
	"github.com/zclconf/go-cty/cty"
)

const RemoteConnectionBlockDescription = "Block used to configure a remote connection to a host. Uses the same values as a `connection` block."

func GetRemoteConnectionBlockSchema() schema.SingleNestedBlock {
	block := ConvertConfigSchemaBlockToSchemaBlock(shared.ConnectionBlockSupersetSchema)
	block.MarkdownDescription = RemoteConnectionBlockDescription
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
			Sensitive: isNameSensitive(name),
		}
	case cty.Number:
		return schema.Int64Attribute{
			Optional:  true,
			Sensitive: isNameSensitive(name),
		}
	case cty.Bool:
		return schema.BoolAttribute{
			Optional:  true,
			Sensitive: isNameSensitive(name),
		}
	default:
		panic("unsupported type")
	}
}

func isNameSensitive(name string) bool {
	switch name {
	case "password",
		"private_key",
		"proxy_user_password",
		"bastion_password",
		"bastion_private_key":
		return true
	default:
		return false
	}
}
