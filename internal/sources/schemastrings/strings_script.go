package schemastrings

const ScriptSourceDescription = "`installer_apt` manages an application using [APT](https://en.wikipedia.org/wiki/APT_(software)).\n\n" +
	"It works on systems that use APT as the package management system. " +
	"Adding an `installer_apt` resource means that Terraform will ensure that " +
	"the application defined in the `name` argument is made available via APT."

const ScriptPathDescription = "The path where the application is installed by `apt-get` after Terraform creates this resource."

const ScriptAdditionalArgsDescription = "Optional version of the application that `apt-get` recognizes. e.g., `2:8.2.3995-1ubuntu2.7`"

const ScriptInstallScriptDescription = "Optional version of the application that `apt-get` recognizes. e.g., `2:8.2.3995-1ubuntu2.7`"

const ScriptFindInstalledScriptDescription = "Optional version of the application that `apt-get` recognizes. e.g., `2:8.2.3995-1ubuntu2.7`"

const ScriptUninstallScriptDescription = "Optional version of the application that `apt-get` recognizes. e.g., `2:8.2.3995-1ubuntu2.7`"
