package cliwrapper

import (
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator"
)

type CliWrapperConfig interface {
	GetCommunicator() communicator.Communicator
}
