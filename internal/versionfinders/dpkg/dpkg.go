package dpkg

import (
	"context"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper/clioutput"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/system"
	"github.com/shihanng/terraform-provider-installer/internal/versionfinders"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

var _ versionfinders.VersionFinder = &DpkgVersionFinder{}

type DpkgVersionFinder struct {
	versionfinders.VersionFinderConfig
}

const DefaultSudo = false
const DefaultProgram = "dpkg"

func NewDpkgVersionFinder(config versionfinders.VersionFinderConfig) *DpkgVersionFinder {
	return &DpkgVersionFinder{
		VersionFinderConfig: config,
	}
}

func (i *DpkgVersionFinder) FindInstalled(ctx context.Context, options versionfinders.VersionFinderOptions) (*models.InstalledProgramInfo, error) {
	info := models.InstalledProgramInfo{}
	info.Name = options.GetName()
	programFound, out := i.DpkgContains(ctx, info.Name)
	if !programFound {
		return nil, out.Error
	}
	version := options.GetVersion()
	if version != nil {
		statusOut := i.DpkgStatus(ctx, info.Name)
		if statusOut.Error != nil {
			return nil, statusOut.Error
		}

		installedVersion, err := versionfinders.ExtractVersion(statusOut.CombinedOutput)
		if version != installedVersion {
			return nil, err
		}
		info.Version = installedVersion
	}

	paths := strings.Split(out.CombinedOutput, versionfinders.OutputNewline)

	info.Path, out.Error = system.FindExecutablePath(paths, info.Name)
	if out.Error != nil {
		return nil, out.Error
	}
	return &info, out.Error
}

func (i *DpkgVersionFinder) DpkgList(ctx context.Context, name string) clioutput.CliOutput {
	return i.getCliWrapper().ExecuteCommand(ctx, "-L", name)
}

func (i *DpkgVersionFinder) DpkgContains(ctx context.Context, name string) (bool, clioutput.CliOutput) {
	out := i.DpkgList(ctx, name)
	hasError := out.Error != nil
	const notInstalledString = "is not installed"
	if hasError && strings.Contains(out.CombinedOutput, notInstalledString) {
		out.Error = errors.Wrap(out.Error, xerrors.ErrNotInstalled.Error())
	}
	return !hasError, out
}

func (i *DpkgVersionFinder) DpkgStatus(ctx context.Context, name string) clioutput.CliOutput {
	return i.getCliWrapper().ExecuteCommand(ctx, "-s", name)
}

func (i *DpkgVersionFinder) getCliWrapper() cliwrapper.CliWrapper {
	return cliwrapper.New(i, DefaultSudo, nil, DefaultProgram)
}
