package sources

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type ISourceBase interface {
	GetAptIDFromName(name string) string
}

var _ ISourceBase = &SourceBase{}

type SourceBase struct {
	SourceType
}

const IDSeparator = ":"

func (b *SourceBase) GetAptIDFromName(name string) string {
	return strings.Join([]string{b.SourceType.String(), name}, IDSeparator)
}

const NameSeparator = "_"

func (b *SourceBase) GetSourceName(prefix string) string {
	return strings.Join([]string{prefix, b.SourceType.String()}, NameSeparator)
}

func TryGetStateData[T any](ctx context.Context, state tfsdk.State, diagnostics diag.Diagnostics) (T, bool) {
	var data T
	diags := state.Get(ctx, &data)
	diagnostics.Append(diags...)

	if diagnostics.HasError() {
		return data, false
	}
	return data, true
}

func SetStateData(ctx context.Context, state tfsdk.State, diagnostics diag.Diagnostics, val interface{}) {
	diags := state.Set(ctx, val)
	diagnostics.Append(diags...)

}

func TryGetConfigData[T any](ctx context.Context, config tfsdk.Config, diagnostics diag.Diagnostics) (T, bool) {
	var data T
	diags := config.Get(ctx, &data)
	diagnostics.Append(diags...)

	if diagnostics.HasError() {
		return data, false
	}
	return data, true
}
