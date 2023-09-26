package versionfinders

import (
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator"
)

type VersionFinderConfig interface {
	GetCommunicator() communicator.Communicator
}
