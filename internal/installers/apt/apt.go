package apt

import (
	"context"

	"github.com/hashicorp/go-version"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper/clioutput"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/system"
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
	installers.InstallerConfig
	VersionFinder versionfinders.VersionFinder
}

const DefaultSudo = true
const DefaultProgram = "apt-get"
const VersionSeperator = "="

var DefaultEnvironment = map[string]string{
	"DEBIAN_FRONTEND": "noninteractive",
}

func NewAptInstaller[T AptInstallerOptions](config installers.InstallerConfig) *AptInstaller[T] {
	return &AptInstaller[T]{
		InstallerConfig: config,
		VersionFinder:   factory.VersionFinderFactory(enums.VersionFinderDpkg, config),
	}
}

func (i *AptInstaller[T]) GetInstallerType() enums.InstallerType {
	return enums.InstallerApt
}

func (i *AptInstaller[T]) Install(ctx context.Context, options T) error {
	wrapper := i.GetCliWrapper(ctx, options)
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
	wrapper := i.GetCliWrapper(ctx, options)
	out := aptRemove(ctx, wrapper, options.GetName(), options.GetVersion())
	return out.Error == nil, out.Error
}

func (i *AptInstaller[T]) GetCliWrapper(ctx context.Context, options T) cliwrapper.CliWrapper {
	environment := system.MergeMaps(DefaultEnvironment, options.GetEnvironmentAndSecrets(ctx))
	return cliwrapper.New(i, options.GetSudo(), environment, DefaultProgram)
}

func aptInstall(ctx context.Context, wrapper cliwrapper.CliWrapper, name string, version *version.Version) clioutput.CliOutput {
	return wrapper.ExecuteCommand(ctx, "-y", "-o", "DPkg::Lock::Timeout=-1", "install", models.GetVersionedName(VersionSeperator, name, version))
}

func aptRemove(ctx context.Context, wrapper cliwrapper.CliWrapper, name string, version *version.Version) clioutput.CliOutput {
	return wrapper.ExecuteCommand(ctx, "-y", "-o", "DPkg::Lock::Timeout=-1", "remove", models.GetVersionedName(VersionSeperator, name, version))
}
