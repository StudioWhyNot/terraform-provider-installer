package installers

import (
	"context"

	"github.com/shihanng/terraform-provider-installer/internal/sources"
)

type IInstallerBase interface {
	Install(ctx context.Context, name string) error
	FindInstalled(ctx context.Context, name string) (string, error)
	Uninstall(ctx context.Context, name string) error
}

type InstallerBase struct {
	sources.SourceType
}

//var _ IInstallerBase = &InstallerBase{}
