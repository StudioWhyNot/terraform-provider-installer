package models

import (
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

// Name and version of a program.
type NamedVersion struct {
	Name      string           `json:"name"`
	Version   *version.Version `json:"version"`
	Seperator string
}

func NewNamedVersion(seperator string, name string, version *version.Version) NamedVersion {
	return NamedVersion{
		Name:      name,
		Version:   version,
		Seperator: seperator,
	}
}

func NewNamedVersionFromString(seperator string, name string) NamedVersion {
	name, version, _ := GetNameAndVersion(seperator, name)
	return NewNamedVersion(seperator, name, version)
}

func NewNamedVersionFromStrings(seperator string, name string, ver string) NamedVersion {
	if ver == "" {
		return NewNamedVersionFromString(seperator, name)
	}
	newVer, _ := version.NewVersion(ver)
	return NewNamedVersion(seperator, name, newVer)
}

func (n NamedVersion) String() string {
	return GetVersionedName(n.Seperator, n.Name, n.Version)
}

func (n NamedVersion) Equals(other NamedVersion) bool {
	return n.Name == other.Name && n.Version.Equal(other.Version)
}

// Information about the installed program.
type InstalledProgramInfo struct {
	NamedVersion `json:",inline"`
	Path         string `json:"path"`
}

func NewInstalledProgramInfo(seperator string, name string, version *version.Version, path string) InstalledProgramInfo {
	return InstalledProgramInfo{
		NamedVersion: NewNamedVersion(seperator, name, version),
		Path:         path,
	}
}

func (i *InstalledProgramInfo) IsNamedVersion(other NamedVersion) bool {
	if i == nil {
		return false
	}
	return i.NamedVersion.Equals(other)
}

// Information about the installed program, with the installer type.
type TypedInstalledProgramInfo struct {
	InstalledProgramInfo
	InstallerType enums.InstallerType
}

func NewTypedInstalledProgramInfo(installerType enums.InstallerType, seperator string, name string, version *version.Version, path string) TypedInstalledProgramInfo {
	return NewTypedInstalledProgramInfoFromInfo(installerType, NewInstalledProgramInfo(seperator, name, version, path))
}

func NewTypedInstalledProgramInfoFromInfo(installerType enums.InstallerType, info InstalledProgramInfo) TypedInstalledProgramInfo {
	return TypedInstalledProgramInfo{
		InstalledProgramInfo: info,
		InstallerType:        installerType,
	}
}

// Information about the program to install.
type InstallerOptions struct {
	Name      string
	Version   *version.Version
	Seperator string
}

// HasVersion returns true if the version is not nil.
func (o InstallerOptions) HasVersion() bool {
	return o.Version == nil
}

func GetNameAndVersion(seperator string, nameVersionString string) (string, *version.Version, error) {
	const expectedParts = 2
	var name string
	var ver *version.Version
	var err error

	split := strings.SplitN(nameVersionString, seperator, expectedParts)
	name = split[0]

	if len(split) == expectedParts {
		ver, err = version.NewVersion(split[1])
	}

	return name, ver, err
}

// GetOptions splits the name and version from string "name=version" and puts the
// values into InstallerOptions.
func getOptions(seperator string, nameVersionString string) (InstallerOptions, error) {
	var info InstallerOptions
	var err error
	info.Name, info.Version, err = GetNameAndVersion(seperator, nameVersionString)

	return info, err
}

// NewInstallerOptions splits the name and version from string "name=version" and puts the
// values into InstallerOptions. If a version is provided, it will use that version.
func NewInstallerOptions(seperator string, nameVersionString string, version *version.Version) (InstallerOptions, error) {
	opt, err := getOptions(seperator, nameVersionString)
	if err != nil {
		return opt, err
	}
	if (opt.Version != nil) && (version != nil) {
		return opt, xerrors.ErrDoubleVersions
	}
	opt.Version = version
	return opt, nil
}

func (o *InstallerOptions) GetVersionedName() string {
	return GetVersionedName(o.Seperator, o.Name, o.Version)
}

// GetVersionedName returns the name and version as a combined string.
func GetVersionedName(seperator string, program string, version *version.Version) string {
	if version == nil {
		return program
	}
	return program + seperator + version.String()
}

func GetCombinedNameVersionStrings(seperator string, name string, version string) string {
	if version == "" {
		return name
	}
	return name + seperator + version
}

func GetIDFromNameAndVersion(seperator string, name string, version string, installerType enums.InstallerType) string {
	if version != "" {
		name = GetCombinedNameVersionStrings(seperator, name, version)
	}
	return installerType.GetIDFromName(name)
}
