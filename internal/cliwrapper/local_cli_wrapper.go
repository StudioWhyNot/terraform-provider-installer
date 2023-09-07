package cliwrapper

import (
	"context"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
)

var _ CliWrapper = LocalCliWrapper{}

// The base struct for all LocalCliWrappers, that wrap the local CLI.
type LocalCliWrapper struct {
	Sudo        bool
	ProgramName string
}

// Default constructor for LocalCliWrapper.
func NewLocalCliWrapper(sudo bool, programName string) LocalCliWrapper {
	return LocalCliWrapper{
		Sudo:        sudo,
		ProgramName: programName,
	}
}

// ExecuteCommand executes a command with the given parameters, taking into consideration whether or not it should be sudo.
func (c LocalCliWrapper) ExecuteCommand(ctx context.Context, params ...string) CliOutput {
	programName, params := GetProgramAndParams(c.Sudo, c.ProgramName, params...)
	cmd := exec.CommandContext(ctx, programName, params...)
	out, err := cmd.CombinedOutput()
	strout := string(out)
	if err != nil {
		return CliOutput{strout, errors.Wrap(errors.WithDetail(err, strout), strings.Join(cmd.Args, " "))}
	}

	return CliOutput{strout, err}
}
