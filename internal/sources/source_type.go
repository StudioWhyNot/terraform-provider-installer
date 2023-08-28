package sources

type SourceType int

const (
	SourceTypeNone SourceType = iota
	SourceTypeApt
	SourceTypeScript
)

var sourceTypeToString = map[SourceType]string{
	SourceTypeNone:   "none",
	SourceTypeApt:    "apt",
	SourceTypeScript: "script",
}

func (s SourceType) String() string {
	return sourceTypeToString[s]
}
