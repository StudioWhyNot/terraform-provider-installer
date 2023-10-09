package cliwrapper

import (
	"context"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper/clibuilder"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper/clioutput"
)

var _ CliWrapper = LocalCliWrapper{}

// The base struct for all LocalCliWrappers, that wrap the local CLI.
type LocalCliWrapper struct {
	clibuilder.CliBuilder
}

// Default constructor for LocalCliWrapper.
func NewLocalCliWrapper(sudo bool, environment map[string]string, programName string) LocalCliWrapper {
	return LocalCliWrapper{
		CliBuilder: clibuilder.NewCliBuilder(sudo, environment, programName),
	}
}

// ExecuteCommand executes a command with the given parameters, taking into consideration whether or not it should be sudo.
func (c LocalCliWrapper) ExecuteCommand(ctx context.Context, params ...string) clioutput.CliOutput {
	programName, params := c.GetProgramAndParams(params...)
	cmd := exec.CommandContext(ctx, programName, params...)
	cmd.Env = c.EnvironmentList()

	out, err := cmd.CombinedOutput()
	strout := string(out)
	if err != nil {
		return clioutput.CliOutput{CombinedOutput: strout, Error: errors.Wrap(errors.WithDetail(err, strout), strings.Join(cmd.Args, clioutput.CliParamSeperator))}
	}

	return clioutput.CliOutput{CombinedOutput: strout, Error: err}
}
