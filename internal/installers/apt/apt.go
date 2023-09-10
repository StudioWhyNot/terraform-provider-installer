package apt

import (
	"context"

	"github.com/hashicorp/go-version"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/versionfinders"
	"github.com/shihanng/terraform-provider-installer/internal/versionfinders/factory"
)

var _ installers.Installer = &AptInstaller{}

type AptInstaller struct {
	CliWrapper    cliwrapper.CliWrapper
	VersionFinder versionfinders.VersionFinder
}

func NewAptInstaller() *AptInstaller {
	const sudo = true
	const program = "apt-get"
	return &AptInstaller{
		CliWrapper:    cliwrapper.NewLocalCliWrapper(sudo, program),
		VersionFinder: factory.VersionFinderFactory(enums.VersionFinderDpkg),
	}
}

func (i *AptInstaller) GetInstallerType() enums.InstallerType {
	return enums.InstallerApt
}

func (i *AptInstaller) Install(ctx context.Context, options models.InstallerOptions) error {
	out := i.AptInstall(ctx, options.Name, options.Version)
	return out.Error
}

func (i *AptInstaller) FindInstalled(ctx context.Context, options models.InstallerOptions) (*models.TypedInstalledProgramInfo, error) {
	return installers.GetInfoFromVersionFinder(i.GetInstallerType(), i.VersionFinder, options, ctx)
}

func (i *AptInstaller) Uninstall(ctx context.Context, options models.InstallerOptions) (bool, error) {
	info, _ := i.FindInstalled(ctx, options)
	if info == nil {
		// Not installed, no error.
		return false, nil
	}
	out := i.AptRemove(ctx, options.Name, options.Version)
	return out.Error == nil, out.Error
}

func (i *AptInstaller) AptInstall(ctx context.Context, name string, version *version.Version) cliwrapper.CliOutput {
	return i.CliWrapper.ExecuteCommand(ctx, "-y", "install", models.GetVersionedName(name, version))
}

func (i *AptInstaller) AptRemove(ctx context.Context, name string, version *version.Version) cliwrapper.CliOutput {
	return i.CliWrapper.ExecuteCommand(ctx, "-y", "remove", models.GetVersionedName(name, version))
}
