package models

import (
	"strings"

	"github.com/hashicorp/go-version"
)

// Information about the isntalled program.
type InstalledProgramInfo struct {
	Path    string
	Version version.Version
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
func GetOptions(nameVersionString string) (InstallerOptions, error) {
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

func GetVersionedName(program string, version *version.Version) string {
	if version == nil {
		return program
	}
	return program + VersionSeperator + version.String()
}
