package dpkg

import (
	"context"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/system"
	"github.com/shihanng/terraform-provider-installer/internal/versionfinders"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

var _ versionfinders.VersionFinder = &DpkgVersionFinder{}

type DpkgVersionFinder struct {
	CliWrapper cliwrapper.CliWrapper
}

func NewDpkgVersionFinder() *DpkgVersionFinder {
	const sudo = false
	const program = "dpkg"
	return &DpkgVersionFinder{
		CliWrapper: cliwrapper.NewLocalCliWrapper(sudo, program),
	}
}

func (i *DpkgVersionFinder) FindInstalled(ctx context.Context, options models.InstallerOptions) (*models.InstalledProgramInfo, error) {
	info := models.InstalledProgramInfo{}
	info.Name = options.Name
	programFound, out := i.DpkgContains(ctx, info.Name)
	if !programFound {
		return nil, out.Error
	}

	if options.Version != nil {
		statusOut := i.DpkgStatus(ctx, info.Name)
		if statusOut.Error != nil {
			return nil, statusOut.Error
		}

		installedVersion, err := versionfinders.ExtractVersion(statusOut.CombinedOutput)
		if options.Version != installedVersion {
			return nil, err
		}
		info.Version = *installedVersion
	}

	paths := strings.Split(out.CombinedOutput, versionfinders.OutputNewline)

	info.Path, out.Error = system.FindExecutablePath(paths) // nolint:wrapcheck
	if out.Error != nil {
		return nil, out.Error
	}
	return &info, out.Error
}

func (i *DpkgVersionFinder) DpkgList(ctx context.Context, name string) cliwrapper.CliOutput {
	return i.CliWrapper.ExecuteCommand(ctx, "-L", name)
}

func (i *DpkgVersionFinder) DpkgContains(ctx context.Context, name string) (bool, cliwrapper.CliOutput) {
	out := i.DpkgList(ctx, name)
	hasError := out.Error != nil
	const notInstalledString = "is not installed"
	if hasError && strings.Contains(out.CombinedOutput, notInstalledString) {
		out.Error = errors.Wrap(out.Error, xerrors.ErrNotInstalled.Error())
	}
	return !hasError, out
}

func (i *DpkgVersionFinder) DpkgStatus(ctx context.Context, name string) cliwrapper.CliOutput {
	return i.CliWrapper.ExecuteCommand(ctx, "-s", name)
}
