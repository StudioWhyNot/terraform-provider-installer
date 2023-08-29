package sources

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type ISourceBase interface {
	GetIDFromName(name string) string
}

type SourceBase struct {
	SourceType
}

var _ ISourceBase = &SourceBase{}

const IDSeparator = ":"

func (b *SourceBase) GetIDFromName(name string) string {
	return strings.Join([]string{b.SourceType.String(), name}, IDSeparator)
}

const NameSeparator = "_"

func (b *SourceBase) GetSourceName(prefix string) string {
	return strings.Join([]string{prefix, b.SourceType.String()}, NameSeparator)
}

type ITerraformDataProvider interface {
	Get(ctx context.Context, target interface{}) diag.Diagnostics
}

func TryGetData[T any](ctx context.Context, provider ITerraformDataProvider, diagnostics *diag.Diagnostics) (T, bool) {
	var data T
	diags := provider.Get(ctx, &data)
	diagnostics.Append(diags...)
	return data, !diagnostics.HasError()
}

func SetStateData(ctx context.Context, state *tfsdk.State, diagnostics *diag.Diagnostics, val interface{}) {
	diags := state.Set(ctx, val)
	diagnostics.Append(diags...)
}
