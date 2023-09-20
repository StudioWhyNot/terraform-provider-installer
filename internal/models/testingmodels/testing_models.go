package testingmodels

import (
	"github.com/shihanng/terraform-provider-installer/internal/models"
)

type TestInfo struct {
	Input    models.InstallerOptions
	Expected models.NamedVersion
}

const errInvalidVersion = "invalid version in test"

func NewTestInfo(name string, version string) TestInfo {
	name = models.GetCombinedNameVersionStrings(name, version)
	options, err := models.NewInstallerOptions(name, nil)

	if err != nil {
		panic(errInvalidVersion)
	}

	return TestInfo{
		Input:    options,
		Expected: models.NewNamedVersion(options.Name, options.Version),
	}
}

func (i TestInfo) String() string {
	return i.Expected.String()
}
