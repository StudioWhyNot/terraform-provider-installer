package factory

import (
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/versionfinders"
	"github.com/shihanng/terraform-provider-installer/internal/versionfinders/dpkg"
)

func VersionFinderFactory(vfType enums.VersionFinderType, config versionfinders.VersionFinderConfig) versionfinders.VersionFinder {
	switch vfType {
	default:
		fallthrough
	case enums.VersionFinderDpkg:
		return dpkg.NewDpkgVersionFinder(config)
	}
}
