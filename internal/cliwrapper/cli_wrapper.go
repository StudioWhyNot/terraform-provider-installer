package cliwrapper

import (
	"context"

	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper/clioutput"
)

type CliWrapper interface {
	ExecuteCommand(ctx context.Context, params ...string) clioutput.CliOutput
}

// Adds sudo to the command if sudo is true, shifting params.
func GetProgramAndParams(sudo bool, programName string, params ...string) (string, []string) {
	const sudoProgramName = "sudo"
	if sudo {
		params = append([]string{programName}, params...)
		programName = sudoProgramName
	}
	return programName, params
}

func New(config CliWrapperConfig, sudo bool, programName string) CliWrapper {
	comm := config.GetCommunicator()
	if comm == nil {
		return NewLocalCliWrapper(sudo, programName)
	}
	return NewRemoteCliWrapper(comm, sudo, programName)
}
