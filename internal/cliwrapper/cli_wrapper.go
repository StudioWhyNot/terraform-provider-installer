package cliwrapper

import "context"

type CliOutput struct {
	CombinedOutput string
	Error          error
}

type CliWrapper interface {
	ExecuteCommand(ctx context.Context, params ...string) CliOutput
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
