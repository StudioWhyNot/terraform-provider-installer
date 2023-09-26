package enums

type VersionFinderType int

const (
	VersionFinderDefault VersionFinderType = iota
	VersionFinderDpkg
)

var versionFinderTypeToString = map[VersionFinderType]string{
	VersionFinderDefault: "default",
	VersionFinderDpkg:    "dpkg",
}

func (s VersionFinderType) String() string {
	return versionFinderTypeToString[s]
}
