package cliwrapper

import (
	"context"

	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper/clioutput"
)

type CliWrapper interface {
	ExecuteCommand(ctx context.Context, params ...string) clioutput.CliOutput
	EscapeScript(script string) string
}

func New(config CliWrapperConfig, sudo bool, environment map[string]string, programName string) CliWrapper {
	comm := config.GetCommunicator()
	if comm == nil {
		return NewLocalCliWrapper(sudo, environment, programName)
	}
	return NewRemoteCliWrapper(comm, sudo, environment, programName)
}
