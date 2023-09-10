package schemastrings

const AptSourceDescription = "`installer_apt` manages an application using [APT](https://en.wikipedia.org/wiki/APT_(software)).\n\n" +
	"It works on systems that use APT as the package management system. " +
	"Adding an `installer_apt` resource means that Terraform will ensure that " +
	"the application defined in the `name` argument is made available via APT."

const AptNameDescription = "Name of the application that `apt-get` recognizes." +
	" Specify a version of a package by following the package name with an equal sign and the version, e.g., `vim=2:8.2.3995-1ubuntu2.7`."

const AptPathDescription = "The path where the application is installed by `apt-get` after Terraform creates this resource."
