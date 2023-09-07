package installers

import (
	"context"

	"github.com/shihanng/terraform-provider-installer/internal/models"
)

// The basic interface for all installers.
type Installer interface {
	Install(ctx context.Context, options models.InstallerOptions) error
	FindInstalled(ctx context.Context, options models.InstallerOptions) (*models.InstalledProgramInfo, error)
	Uninstall(ctx context.Context, options models.InstallerOptions) (bool, error)
}
