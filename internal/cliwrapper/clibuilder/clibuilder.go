package clibuilder

import (
	"strings"

	"github.com/shihanng/terraform-provider-installer/internal/system"
)

const SudoProgramName = "sudo"

type CliBuilder struct {
	Sudo        bool
	Environment map[string]string
	ProgramName string
}

func NewCliBuilder(sudo bool, environment map[string]string, programName string) CliBuilder {
	return CliBuilder{
		Sudo:        sudo,
		Environment: environment,
		ProgramName: programName,
	}
}

func (c *CliBuilder) GetProgramAndParams(params ...string) (string, []string) {
	programName := c.ProgramName
	if c.Sudo {
		params = append([]string{c.ProgramName}, params...)
		programName = SudoProgramName
	}
	return programName, params
}

// Adds sudo to the command if sudo is true, shifting params.
func (c *CliBuilder) GetProgramAndParamsWithEnvironment(params ...string) []string {
	program, params := c.GetProgramAndParams(params...)
	envList := c.EnvironmentList()
	if c.Sudo {
		envList = append([]string{program}, envList...)
	} else {
		envList = append(envList, program)
	}

	return append(envList, params...)
}

func (c *CliBuilder) EnvironmentList() []string {
	return EnvMapToEnvList(c.Environment)
}

func EnvMapToString(env map[string]string) string {
	const envSeperator = " "
	envList := EnvMapToEnvList(env)
	return strings.Join(envList, envSeperator)
}

const EnvSeperator = "="

func EnvMapToEnvList(env map[string]string) []string {
	if env == nil {
		return []string{}
	}
	// Use quotes to allow for spaces and shell expansion.
	const envWrapperCharacter = "\""
	envList := make([]string, 0, len(env))
	for k, v := range env {
		envList = append(envList, k+EnvSeperator+system.WrapString(v, envWrapperCharacter))
	}
	return envList
}
