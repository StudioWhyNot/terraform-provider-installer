package script

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper/clioutput"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/system"
)

type ScriptInstallerOptions interface {
	installers.InstallerOptions
	GetId() string
	GetPath() string
	GetShell() string
	GetInstallScript() string
	GetFindInstalledScript() string
	GetUninstallScript() string
	GetAdditionalArgs(ctx context.Context) []string
	GetDefaultArgs(ctx context.Context) []string
}

var _ installers.Installer[ScriptInstallerOptions] = &ScriptInstaller[ScriptInstallerOptions]{}

type ScriptInstaller[T ScriptInstallerOptions] struct {
	installers.InstallerConfig
}

const DefaultSudo = false
const DefaultProgram = "sh"
const DefaultArg = "-c"
const VersionSeperator = "="

func NewScriptInstaller[T ScriptInstallerOptions](config installers.InstallerConfig) *ScriptInstaller[T] {
	return &ScriptInstaller[T]{
		InstallerConfig: config,
	}
}

func (i *ScriptInstaller[T]) GetInstallerType() enums.InstallerType {
	return enums.InstallerScript
}

func (i *ScriptInstaller[T]) Install(ctx context.Context, options T) error {
	out := i.executeScript(ctx, options, options.GetInstallScript())
	return out.Error
}

func (i *ScriptInstaller[T]) FindInstalled(ctx context.Context, options T) (*models.TypedInstalledProgramInfo, error) {
	// If a path is specified, check if the path has a program installed.
	path := options.GetPath()
	if path != "" {
		installed, err := IsInstalled(path)
		if err != nil {
			return nil, err
		}
		if installed {
			newInfo := models.NewTypedInstalledProgramInfo(i.GetInstallerType(), VersionSeperator, options.GetId(), nil, path)
			return &newInfo, nil
		}
	}
	// Otherwise, if a find installed script is specified, run the script to find the program.
	findInstalledScript := options.GetFindInstalledScript()
	if findInstalledScript == "" {
		return nil, nil
	}
	out := i.executeScript(ctx, options, findInstalledScript)

	jsonData := out.CombinedOutput

	var info models.InstalledProgramInfo = models.InstalledProgramInfo{}
	info.Seperator = VersionSeperator
	if out.CombinedOutput != "" {
		err := json.Unmarshal([]byte(jsonData), &info)
		if err != nil {
			out.Error = errors.Wrap(err, "Failed to parse JSON output of `find_installed_script`: "+findInstalledScript)
		}
	}
	typedInfo := models.NewTypedInstalledProgramInfoFromInfo(i.GetInstallerType(), info)
	return &typedInfo, out.Error
}

func (i *ScriptInstaller[T]) Uninstall(ctx context.Context, options T) (bool, error) {
	info, _ := i.FindInstalled(ctx, options)
	if info == nil {
		// Not installed, no error.
		return false, nil
	}
	out := i.executeScript(ctx, options, options.GetUninstallScript())
	return out.Error == nil, out.Error
}

func (i *ScriptInstaller[T]) executeScript(ctx context.Context, options T, script string) clioutput.CliOutput {
	// Use single quote to wrap the script to avoid shell expansion.
	const wrapperCharacter = "'"
	wrapper := i.GetCliWrapper(options)
	args := append(options.GetDefaultArgs(ctx), system.WrapString(script, wrapperCharacter))
	args = append(args, options.GetAdditionalArgs(ctx)...)
	return wrapper.ExecuteCommand(ctx, args...)
}

func (i *ScriptInstaller[T]) GetCliWrapper(options T) cliwrapper.CliWrapper {
	return cliwrapper.New(i, options.GetSudo(), options.GetEnvironment(), options.GetShell())
}

func IsInstalled(path string) (bool, error) {
	if _, err := exec.LookPath(path); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return false, nil
		}
		return false, errors.Wrap(err, "check if path is installed")
	}
	return true, nil
}
