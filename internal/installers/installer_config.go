package installers

import "github.com/shihanng/terraform-provider-installer/internal/terraform/communicator"

type InstallerConfig interface {
	GetCommunicator() communicator.Communicator
}
