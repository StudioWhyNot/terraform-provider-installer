package versionfinders

import (
	"context"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/shihanng/terraform-provider-installer/internal/models"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

type VersionFinderOptions interface {
	GetName() string
	GetVersion() *version.Version
}

// The basic interface for all installers.
type VersionFinder interface {
	FindInstalled(ctx context.Context, options VersionFinderOptions) (*models.InstalledProgramInfo, error)
}

const OutputNewline = "\n"

// ExtractVersion extracts version value from the output of dpkg -s <package>.
func ExtractVersion(input string) (*version.Version, error) {
	const aptVersionPrefix string = "Version: "

	for _, line := range strings.Split(input, OutputNewline) {
		if strings.HasPrefix(line, aptVersionPrefix) {
			versionString := strings.TrimSpace(strings.TrimPrefix(line, aptVersionPrefix))
			return version.NewVersion(versionString)
		}
	}

	return nil, xerrors.ErrVersionNotFound
}
