package terraformutils

const RemoteConnectionBlockDescription = "Block used to configure a remote connection to a host. Uses the same values as a `connection` block."

func IsNameSensitive(name string) bool {
	switch name {
	case "password",
		"private_key",
		"proxy_user_password",
		"bastion_password",
		"bastion_private_key":
		return true
	default:
		return false
	}
}
