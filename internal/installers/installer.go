package installers

import (
	"context"

	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/versionfinders"
)

type InstallerOptions interface {
	GetSudo() bool
	GetEnvironment() map[string]string
}

// The basic interface for all installers.
type Installer[T any] interface {
	GetInstallerType() enums.InstallerType
	Install(ctx context.Context, options T) error
	FindInstalled(ctx context.Context, options T) (*models.TypedInstalledProgramInfo, error)
	Uninstall(ctx context.Context, options T) (bool, error)
}

func GetInfoFromVersionFinder(installerType enums.InstallerType, versionFinder versionfinders.VersionFinder, options versionfinders.VersionFinderOptions, ctx context.Context) (*models.TypedInstalledProgramInfo, error) {
	info, err := versionFinder.FindInstalled(ctx, options)
	if info == nil {
		return nil, err
	}
	newInfo := models.NewTypedInstalledProgramInfoFromInfo(installerType, *info)
	return &newInfo, err
}
