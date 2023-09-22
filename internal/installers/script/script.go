package script

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
	"github.com/shihanng/terraform-provider-installer/internal/models"
)

type ScriptInstallerOptions interface {
	installers.InstallerOptions
	GetId() string
	GetPath() string
	GetInstallScript() string
	GetFindInstalledScript() string
	GetUninstallScript() string
	GetAdditionalArgs(ctx context.Context) []string
}

var _ installers.Installer[ScriptInstallerOptions] = &ScriptInstaller[ScriptInstallerOptions]{}

type ScriptInstaller[T ScriptInstallerOptions] struct {
	CliWrapper cliwrapper.CliWrapper
}

func NewScriptInstaller[T ScriptInstallerOptions]() *ScriptInstaller[T] {
	const sudo = false
	const program = "bash"
	return &ScriptInstaller[T]{
		CliWrapper: cliwrapper.NewLocalCliWrapper(sudo, program),
	}
}

func (i *ScriptInstaller[T]) GetInstallerType() enums.InstallerType {
	return enums.InstallerScript
}

func (i *ScriptInstaller[T]) Install(ctx context.Context, options T) error {
	args := append([]string{options.GetInstallScript()}, options.GetAdditionalArgs(ctx)...)
	out := i.CliWrapper.ExecuteCommand(ctx, args...)
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
			newInfo := models.NewTypedInstalledProgramInfo(i.GetInstallerType(), options.GetId(), nil, path)
			return &newInfo, nil
		}
	}
	// Otherwise, if a find installed script is specified, run the script to find the program.
	findInstalledScript := options.GetFindInstalledScript()
	if findInstalledScript == "" {
		return nil, nil
	}
	args := append([]string{findInstalledScript}, options.GetAdditionalArgs(ctx)...)
	out := i.CliWrapper.ExecuteCommand(ctx, args...)

	jsonData := out.CombinedOutput

	var info models.InstalledProgramInfo = models.InstalledProgramInfo{}
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
	args := append([]string{options.GetUninstallScript()}, options.GetAdditionalArgs(ctx)...)
	out := i.CliWrapper.ExecuteCommand(ctx, args...)
	return out.Error == nil, out.Error
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
