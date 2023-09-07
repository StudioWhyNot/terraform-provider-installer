package enums

type InstallerType int

const (
	InstallerNone InstallerType = iota
	InstallerApt
	InstallerScript
)

var sourceTypeToString = map[InstallerType]string{
	InstallerNone:   "none",
	InstallerApt:    "apt",
	InstallerScript: "script",
}

func (s InstallerType) String() string {
	return sourceTypeToString[s]
}
