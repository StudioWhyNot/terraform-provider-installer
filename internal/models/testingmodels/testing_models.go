package testingmodels

import (
	"github.com/shihanng/terraform-provider-installer/internal/installers/apt"
	"github.com/shihanng/terraform-provider-installer/internal/models"
)

type TestInfo[T any] struct {
	Input    T
	Expected models.NamedVersion
}

const errInvalidVersion = "invalid version in test"

func NewAptTestInfo(name string, version string) TestInfo[apt.AptInstallerOptions] {
	seperator := apt.VersionSeperator
	name = models.GetCombinedNameVersionStrings(seperator, name, version)
	options, err := models.NewInstallerOptions(seperator, name, nil)

	if err != nil {
		panic(errInvalidVersion)
	}

	return TestInfo[apt.AptInstallerOptions]{
		//TODO: Fix this
		//Input:    options,
		Expected: models.NewNamedVersion(seperator, options.Name, options.Version),
	}
}

func (i TestInfo[T]) String() string {
	return i.Expected.String()
}
