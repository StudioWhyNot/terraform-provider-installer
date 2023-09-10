package sources

import (
	"context"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
	"github.com/shihanng/terraform-provider-installer/internal/installers/factory"
	"github.com/shihanng/terraform-provider-installer/internal/models"
)

type SourceData interface {
	GetName() types.String
	GetVersion() *version.Version
	SetDataFromTypedInstalledProgramInfo(info *models.TypedInstalledProgramInfo)
}

type SourceBase struct {
	Installer installers.Installer
}

func NewSourceBase(installerType enums.InstallerType) *SourceBase {
	return &SourceBase{
		Installer: factory.InstallerFinderFactory(installerType),
	}
}

func (s *SourceBase) GetDefaultTypeName(providerTypeName string) string {
	return s.Installer.GetInstallerType().GetSourceName(providerTypeName)
}

func (s *SourceBase) UpdateFromInstallation(data SourceData, ctx context.Context, diagnostics *diag.Diagnostics) bool {
	info := GetInstallationInfo(data, s.Installer, ctx, diagnostics)
	success := info != nil
	if success {
		data.SetDataFromTypedInstalledProgramInfo(info)
	}
	return success
}

type TerraformDataProvider interface {
	Get(ctx context.Context, target interface{}) diag.Diagnostics
}
