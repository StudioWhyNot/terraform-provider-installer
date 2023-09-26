package brew

import (
	"context"

	"github.com/hashicorp/go-version"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper/clioutput"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/versionfinders"
	"github.com/shihanng/terraform-provider-installer/internal/versionfinders/factory"
)

type BrewInstallerOptions interface {
	installers.InstallerOptions
	GetName() string
	GetVersion() *version.Version
	GetCask() bool
}

var _ installers.Installer[BrewInstallerOptions] = &BrewInstaller[BrewInstallerOptions]{}

type BrewInstaller[T BrewInstallerOptions] struct {
	installers.InstallerConfig
	VersionFinder versionfinders.VersionFinder
}

const DefaultSudo = false
const DefaultProgram = "brew"
const DefaultCask = false
const VersionSeperator = "@"

func NewBrewInstaller[T BrewInstallerOptions](config installers.InstallerConfig) *BrewInstaller[T] {
	return &BrewInstaller[T]{
		InstallerConfig: config,
		VersionFinder:   factory.VersionFinderFactory(enums.VersionFinderDpkg, config),
	}
}

func (i *BrewInstaller[T]) GetInstallerType() enums.InstallerType {
	return enums.InstallerBrew
}

func (i *BrewInstaller[T]) Install(ctx context.Context, options T) error {
	// wrapper := installers.GetCliWrapper(options.GetSudo(), DefaultProgram)
	// out := brewInstall(ctx, wrapper, options.GetName(), options.GetVersion())
	// return out.Error
	return nil
}

func (i *BrewInstaller[T]) FindInstalled(ctx context.Context, options T) (*models.TypedInstalledProgramInfo, error) {
	return installers.GetInfoFromVersionFinder(i.GetInstallerType(), i.VersionFinder, options, ctx)
}

func (i *BrewInstaller[T]) Uninstall(ctx context.Context, options T) (bool, error) {
	info, _ := i.FindInstalled(ctx, options)
	if info == nil {
		// Not installed, no error.
		return false, nil
	}
	wrapper := i.GetCliWrapper(options.GetSudo())
	out := brewUninstall(ctx, wrapper, options.GetName(), options.GetVersion())
	return out.Error == nil, out.Error
}

func (i *BrewInstaller[T]) getBrewCommand(options T, command string, withJSONV2 bool) []string {
	commandArray := []string{command}
	brewOptions := newBrewOptions(options.GetCask(), withJSONV2)
	args := getBrewOptionsWithVersionedName(brewOptions, options.GetName(), options.GetVersion())
	args = append(commandArray, args...)
	return args
}

func (i *BrewInstaller[T]) GetCliWrapper(sudo bool) cliwrapper.CliWrapper {
	return cliwrapper.New(i, sudo, DefaultProgram)
}

type brewOptions int

const (
	WithCask brewOptions = 1 << iota
	WithJSONV2
)

func newBrewOptions(cask bool, jsonV2 bool) brewOptions {
	var options brewOptions
	if cask {
		options |= WithCask
	}
	if jsonV2 {
		options |= WithJSONV2
	}
	return options
}

func getBrewOptions(options brewOptions) []string {
	var brewOptions []string
	if options&WithCask != 0 {
		brewOptions = append(brewOptions, "--cask")
	}
	if options&WithJSONV2 != 0 {
		brewOptions = append(brewOptions, "--json=v2")
	}
	return brewOptions
}

func getBrewOptionsWithVersionedName(options brewOptions, name string, version *version.Version) []string {
	args := getBrewOptions(options)
	args = append(args, models.GetVersionedName(VersionSeperator, name, version))
	return args
}

// func brewInstall(ctx context.Context, wrapper cliwrapper.CliWrapper, options brewOptions, name string, version *version.Version) cliwrapper.CliOutput {
// 	return wrapper.ExecuteCommand(ctx, "install", getBrewOptionsWithVersionedName(options, name, version))
// }

func brewUninstall(ctx context.Context, wrapper cliwrapper.CliWrapper, name string, version *version.Version) clioutput.CliOutput {
	return wrapper.ExecuteCommand(ctx, "uninstall", models.GetVersionedName(VersionSeperator, name, version))
}
