package models

import (
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

// Information about the installed program.
type InstalledProgramInfo struct {
	Name    string
	Version version.Version
	Path    string
}

// Information about the installed program, with the installer type.
type TypedInstalledProgramInfo struct {
	InstalledProgramInfo
	InstallerType enums.InstallerType
}

// Information about the program to install.
type InstallerOptions struct {
	Name    string
	Version *version.Version
}

// HasVersion returns true if the version is not nil.
func (o InstallerOptions) HasVersion() bool {
	return o.Version == nil
}

// The separator between name and version.
const VersionSeperator = "="

// GetOptions splits the name and version from string "name=version" and puts the
// values into InstallerOptions.
func getOptions(nameVersionString string) (InstallerOptions, error) {
	const expectedParts = 2

	var info InstallerOptions

	split := strings.SplitN(nameVersionString, VersionSeperator, expectedParts)

	info.Name = split[0]

	if len(split) == expectedParts {
		version, error := version.NewVersion(split[1])
		if error != nil {
			return info, error
		}
		info.Version = version
	}

	return info, nil
}

// GetOptions splits the name and version from string "name=version" and puts the
// values into InstallerOptions. If a version is provided, it will use that version.
func GetOptions(nameVersionString string, version *version.Version) (InstallerOptions, error) {
	opt, err := getOptions(nameVersionString)
	if err != nil {
		return opt, err
	}
	if (opt.Version != nil) && (version != nil) {
		return opt, xerrors.ErrDoubleVersions
	}
	opt.Version = version
	return opt, nil
}

func GetVersionedName(program string, version *version.Version) string {
	if version == nil {
		return program
	}
	return program + VersionSeperator + version.String()
}
