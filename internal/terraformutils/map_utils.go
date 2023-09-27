package terraformutils

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/iancoleman/strcase"
	"github.com/zclconf/go-cty/cty"
)

func StructToCtyValueMap[T any](obj *T) *cty.Value {
	if obj == nil {
		return nil
	}
	const metadataTag = "tfsdk"
	objValue := reflect.ValueOf(*obj)
	objType := objValue.Type()

	result := make(map[string]cty.Value)

	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		fieldName := objType.Field(i).Name
		// Get the tfsdk metadata tag value for the field name
		tagValue := objType.Field(i).Tag.Get(metadataTag)
		if tagValue != "" {
			fieldName = tagValue
		} else {
			// Convert the field name to snake case if no tfsdk tag is present
			fieldName = strcase.ToSnake(fieldName)
		}

		fieldValue := ConvertToCtyValue(field)

		result[fieldName] = fieldValue
	}
	newVal := cty.ObjectVal(result)
	return &newVal
}

// ConvertToCtyValue converts a reflect.Value to a cty.Value
func ConvertToCtyValue(value reflect.Value) cty.Value {
	valType := value.Type()
	switch valType {
	case reflect.TypeOf(basetypes.StringValue{}):
		return cty.StringVal(value.Interface().(basetypes.StringValue).ValueString())
	case reflect.TypeOf(basetypes.BoolValue{}):
		return cty.BoolVal(value.Interface().(basetypes.BoolValue).ValueBool())
	case reflect.TypeOf(basetypes.Int64Value{}):
		return cty.NumberIntVal(value.Interface().(basetypes.Int64Value).ValueInt64())
	default:
		panic("unsupported type: " + valType.Name())
	}
}
