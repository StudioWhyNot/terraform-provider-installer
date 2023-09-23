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
	VersionFinder versionfinders.VersionFinder
}

const DefaultSudo = true
const DefaultProgram = "apt-get"

func NewAptInstaller[T AptInstallerOptions]() *AptInstaller[T] {
	return &AptInstaller[T]{
		VersionFinder: factory.VersionFinderFactory(enums.VersionFinderDpkg),
	}
}

func (i *AptInstaller[T]) GetInstallerType() enums.InstallerType {
	return enums.InstallerApt
}

func (i *AptInstaller[T]) Install(ctx context.Context, options T) error {
	wrapper := installers.GetCliWrapper(options.GetSudo(), DefaultProgram)
	out := aptInstall(ctx, wrapper, options.GetName(), options.GetVersion())
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
	wrapper := installers.GetCliWrapper(options.GetSudo(), DefaultProgram)
	out := aptRemove(ctx, wrapper, options.GetName(), options.GetVersion())
	return out.Error == nil, out.Error
}

func aptInstall(ctx context.Context, wrapper cliwrapper.CliWrapper, name string, version *version.Version) cliwrapper.CliOutput {
	return wrapper.ExecuteCommand(ctx, "-y", "install", models.GetVersionedName(name, version))
}

func aptRemove(ctx context.Context, wrapper cliwrapper.CliWrapper, name string, version *version.Version) cliwrapper.CliOutput {
	return wrapper.ExecuteCommand(ctx, "-y", "remove", models.GetVersionedName(name, version))
}
