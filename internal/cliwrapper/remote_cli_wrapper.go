package cliwrapper

import (
	"bytes"
	"context"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/terraform-provider-installer/internal/cliwrapper/clioutput"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator"
	"github.com/shihanng/terraform-provider-installer/internal/terraform/communicator/remote"
)

var _ CliWrapper = RemoteCliWrapper{}

// The base struct for all RemoteCliWrappers, that wrap the local CLI.
type RemoteCliWrapper struct {
	Sudo         bool
	ProgramName  string
	Communicator communicator.Communicator
}

// Default constructor for RemoteCliWrapper.
func NewRemoteCliWrapper(communicator communicator.Communicator, sudo bool, programName string) RemoteCliWrapper {
	return RemoteCliWrapper{
		Sudo:         sudo,
		ProgramName:  programName,
		Communicator: communicator,
	}
}

// ExecuteCommand executes a command with the given parameters, taking into consideration whether or not it should be sudo.
func (c RemoteCliWrapper) ExecuteCommand(ctx context.Context, params ...string) clioutput.CliOutput {
	programName, params := GetProgramAndParams(c.Sudo, c.ProgramName, params...)
	buffer := bytes.Buffer{}
	errBuff := bytes.Buffer{}
	cmd := getCommand(&buffer, &errBuff, programName, params...)
	err := c.Communicator.Start(cmd)
	cmd.Wait()
	// Print the combined output
	strout := buffer.String()
	if err != nil {
		return clioutput.CliOutput{CombinedOutput: strout, Error: errors.Wrap(errors.WithDetail(err, strout), cmd.Command)}
	}

	return clioutput.CliOutput{CombinedOutput: strout, Error: err}
}

func getCommand(output *bytes.Buffer, err *bytes.Buffer, programName string, params ...string) *remote.Cmd {
	params = append([]string{programName}, params...)
	joined := strings.Join(params, clioutput.CliParamSeperator)
	return &remote.Cmd{
		Command: joined,
		Stdout:  output,
		Stderr:  err,
	}
}
