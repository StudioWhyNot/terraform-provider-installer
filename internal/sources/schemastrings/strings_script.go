package schemastrings

const ScriptSourceDescription = "`installer_script` manages an application using a custom script.\n\n" +
	"Adding an `installer_script` resource means that Terraform will install " +
	"application in the `path` by running the `install_script` when creating the resource."

const ScriptPathDescription = "is an optional location of the application installed by the install script. " +
	"If the application does not exist at path, then the resource is considered not exist by Terraform.\n\n" +
	"If not specified, the value will be computed from the `find_installed` script."

const ScriptInstallScriptDescription = "is the script that will be called by Terraform when executing `terraform plan/apply`."

const ScriptFindInstalledScriptDescription = "is an optional script that will be used by terraform to find the path of the installed application."

const ScriptUninstallScriptDescription = "is the script that will be called by Terraform when executing `terraform destroy`."

const ScriptAdditionalArgsDescription = "Additional arguments to be passed to the install, uninstall, and find_installed scripts."

const ScriptShellDescription = "Which shell program to use to run the install, uninstall, and find_installed scripts. This shellis followed by the `-c` flag."
