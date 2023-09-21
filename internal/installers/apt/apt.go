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

type AptInstallerOptions interface {
	installers.InstallerOptions
	GetName() string
	GetVersion() *version.Version
}

var _ installers.Installer[AptInstallerOptions] = &AptInstaller[AptInstallerOptions]{}

type AptInstaller[T AptInstallerOptions] struct {
	CliWrapper    cliwrapper.CliWrapper
	VersionFinder versionfinders.VersionFinder
}

func NewAptInstaller[T AptInstallerOptions]() *AptInstaller[T] {
	const sudo = true
	const program = "apt-get"
	return &AptInstaller[T]{
		CliWrapper:    cliwrapper.NewLocalCliWrapper(sudo, program),
		VersionFinder: factory.VersionFinderFactory(enums.VersionFinderDpkg),
	}
}

func (i *AptInstaller[T]) GetInstallerType() enums.InstallerType {
	return enums.InstallerApt
}

func (i *AptInstaller[T]) Install(ctx context.Context, options T) error {
	out := i.aptInstall(ctx, options.GetName(), options.GetVersion())
	return out.Error
}

func (i *AptInstaller[T]) FindInstalled(ctx context.Context, options T) (*models.TypedInstalledProgramInfo, error) {
	return installers.GetInfoFromVersionFinder(i.GetInstallerType(), i.VersionFinder, options, ctx)
}

func (i *AptInstaller[T]) Uninstall(ctx context.Context, options T) (bool, error) {
	info, _ := i.FindInstalled(ctx, options)
	if info == nil {
		// Not installed, no error.
		return false, nil
	}
	out := i.aptRemove(ctx, options.GetName(), options.GetVersion())
	return out.Error == nil, out.Error
}

func (i *AptInstaller[T]) aptInstall(ctx context.Context, name string, version *version.Version) cliwrapper.CliOutput {
	return i.CliWrapper.ExecuteCommand(ctx, "-y", "install", models.GetVersionedName(name, version))
}

func (i *AptInstaller[T]) aptRemove(ctx context.Context, name string, version *version.Version) cliwrapper.CliOutput {
	return i.CliWrapper.ExecuteCommand(ctx, "-y", "remove", models.GetVersionedName(name, version))
}
