package models

import (
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

// Name and version of a program.
type NamedVersion struct {
	Name    string
	Version *version.Version
}

func NewNamedVersion(name string, version *version.Version) NamedVersion {
	return NamedVersion{
		Name:    name,
		Version: version,
	}
}

func NewNamedVersionFromString(name string) NamedVersion {
	name, version, _ := GetNameAndVersion(name)
	return NewNamedVersion(name, version)
}

func NewNamedVersionFromStrings(name string, ver string) NamedVersion {
	if ver == "" {
		return NewNamedVersionFromString(name)
	}
	newVer, _ := version.NewVersion(ver)
	return NewNamedVersion(name, newVer)
}

func (n NamedVersion) String() string {
	return GetVersionedName(n.Name, n.Version)
}

func (n NamedVersion) Equals(other NamedVersion) bool {
	return n.Name == other.Name && n.Version.Equal(other.Version)
}

// Information about the installed program.
type InstalledProgramInfo struct {
	NamedVersion
	Path string
}

func NewInstalledProgramInfo(name string, version *version.Version, path string) InstalledProgramInfo {
	return InstalledProgramInfo{
		NamedVersion: NewNamedVersion(name, version),
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

func NewTypedInstalledProgramInfo(installerType enums.InstallerType, name string, version *version.Version, path string) TypedInstalledProgramInfo {
	return NewTypedInstalledProgramInfoFromInfo(installerType, NewInstalledProgramInfo(name, version, path))
}

func NewTypedInstalledProgramInfoFromInfo(installerType enums.InstallerType, info InstalledProgramInfo) TypedInstalledProgramInfo {
	return TypedInstalledProgramInfo{
		InstalledProgramInfo: info,
		InstallerType:        installerType,
	}
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

func GetNameAndVersion(nameVersionString string) (string, *version.Version, error) {
	const expectedParts = 2
	var name string
	var ver *version.Version
	var err error

	split := strings.SplitN(nameVersionString, VersionSeperator, expectedParts)
	name = split[0]

	if len(split) == expectedParts {
		ver, err = version.NewVersion(split[1])
	}

	return name, ver, err
}

// GetOptions splits the name and version from string "name=version" and puts the
// values into InstallerOptions.
func getOptions(nameVersionString string) (InstallerOptions, error) {
	var info InstallerOptions
	var err error
	info.Name, info.Version, err = GetNameAndVersion(nameVersionString)

	return info, err
}

// NewInstallerOptions splits the name and version from string "name=version" and puts the
// values into InstallerOptions. If a version is provided, it will use that version.
func NewInstallerOptions(nameVersionString string, version *version.Version) (InstallerOptions, error) {
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

func (o *InstallerOptions) GetVersionedName() string {
	return GetVersionedName(o.Name, o.Version)
}

// GetVersionedName returns the name and version as a combined string.
func GetVersionedName(program string, version *version.Version) string {
	if version == nil {
		return program
	}
	return program + VersionSeperator + version.String()
}

func GetCombinedNameVersionStrings(name string, version string) string {
	if version == "" {
		return name
	}
	return name + VersionSeperator + version
}

func GetIDFromNameAndVersion(name string, version string, installerType enums.InstallerType) string {
	if version != "" {
		name = GetCombinedNameVersionStrings(name, version)
	}
	return installerType.GetIDFromName(name)
}
