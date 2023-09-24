package enums

import (
	"strings"
)

type InstallerType int

const (
	InstallerNone InstallerType = iota
	InstallerApt
	InstallerScript
	InstallerBrew
)

var sourceTypeToString = map[InstallerType]string{
	InstallerNone:   "none",
	InstallerApt:    "apt",
	InstallerScript: "script",
	InstallerBrew:   "brew",
}

func (s InstallerType) String() string {
	return sourceTypeToString[s]
}

const IDSeparator = ":"

func (s InstallerType) GetIDFromName(name string) string {
	return strings.Join([]string{s.String(), name}, IDSeparator)
}

const NameSeparator = "_"

func (s InstallerType) GetSourceName(prefix string) string {
	return strings.Join([]string{prefix, s.String()}, NameSeparator)
}
