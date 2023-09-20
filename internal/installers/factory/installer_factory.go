package factory

import (
	"github.com/shihanng/terraform-provider-installer/internal/enums"
	"github.com/shihanng/terraform-provider-installer/internal/installers"
	"github.com/shihanng/terraform-provider-installer/internal/installers/apt"
)

func InstallerFactory(installerType enums.InstallerType) installers.Installer {
	switch installerType {
	default:
		fallthrough
	case enums.InstallerApt:
		return apt.NewAptInstaller()
	}
}
