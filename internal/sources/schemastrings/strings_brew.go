package schemastrings

const BrewSourceDescription = "`installer_brew` manages an application using [Homebrew](https://brew.sh/).\n\n" +
	"It works on systems that use Homebrew as the package management system. " +
	"Adding an `installer_brew` resource means that Terraform will ensure that " +
	"the application defined in the `name` argument is made available via brew."

const BrewNameDescription = "Name of the application that `brew` recognizes, e.g., `homebrew/cask/alfred` for a cask, `goreleaser/tap/goreleaser` for tap. " +
	"Treats a package as a formula if `cask` is not set or is false"

const BrewVersionDescription = "Optional version of the application that `brew` recognizes. e.g., `2:8.2.3995-1ubuntu2.7`"

const BrewCaskDescription = "Treat name argument as cask."

const BrewPathDescription = "The path where the application is installed by `brew` after Terraform creates this resource."
