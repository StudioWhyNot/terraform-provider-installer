package sources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
)

type SourceData interface {
	Initialize() bool
}

type SourceBase[T any] struct {
	Installer installers.Installer[T]
}

func NewSourceBase[T any](installer installers.Installer[T]) *SourceBase[T] {
	return &SourceBase[T]{
		Installer: installer,
	}
}

func GetIDFromName(name types.String, installerType enums.InstallerType) types.String {
	return types.StringValue(installerType.GetIDFromName(name.ValueString()))
}

type TerraformDataProvider interface {
	Get(ctx context.Context, target interface{}) diag.Diagnostics
}
