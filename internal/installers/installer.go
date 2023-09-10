package installers

import (
	"context"

	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/versionfinders"
)

// The basic interface for all installers.
type Installer interface {
	GetInstallerType() enums.InstallerType
	Install(ctx context.Context, options models.InstallerOptions) error
	FindInstalled(ctx context.Context, options models.InstallerOptions) (*models.TypedInstalledProgramInfo, error)
	Uninstall(ctx context.Context, options models.InstallerOptions) (bool, error)
}

func GetInfoFromVersionFinder(installerType enums.InstallerType, versionFinder versionfinders.VersionFinder, options models.InstallerOptions, ctx context.Context) (*models.TypedInstalledProgramInfo, error) {
	info, err := versionFinder.FindInstalled(ctx, options)
	if info == nil {
		return nil, err
	}
	return &models.TypedInstalledProgramInfo{
		InstalledProgramInfo: *info,
		InstallerType:        installerType,
	}, err
}
