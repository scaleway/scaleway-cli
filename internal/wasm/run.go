//go:build wasm

package wasm

import (
	"bytes"
	"io"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces"
	"github.com/scaleway/scaleway-cli/v2/internal/platform/web"
)

var commands *core.Commands

func getCommands() *core.Commands {
	if commands == nil {
		commands = namespaces.GetCommands()
	}
	return commands
}

type RunConfig struct {
	JWT                   string `js:"jwt"`
	DefaultProjectID      string `js:"defaultProjectID"`
	DefaultOrganizationID string `js:"defaultOrganizationID"`
	APIUrl                string `js:"apiUrl"`
}

type RunResponse struct {
	Stdout   string `js:"stdout"`
	Stderr   string `js:"stderr"`
	ExitCode int    `js:"exitCode"`
}

func runCommand(cfg *RunConfig, args []string, stdout io.Writer, stderr io.Writer) int {
	exitCode, _, _ := core.Bootstrap(&core.BootstrapConfig{
		Args:      args,
		Commands:  getCommands(),
		BuildInfo: &core.BuildInfo{},
		Stdout:    stdout,
		Stderr:    stderr,
		Stdin:     nil,
		Platform: &web.Platform{
			JWT:                   cfg.JWT,
			DefaultProjectID:      cfg.DefaultProjectID,
			DefaultOrganizationID: cfg.DefaultOrganizationID,
			APIUrl:                cfg.APIUrl,
		},
	})

	return exitCode
}

func Run(cfg *RunConfig, args []string) (*RunResponse, error) {
	cliArgs := []string{"scw"}
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	cliArgs = append(cliArgs, args...)

	exitCode := runCommand(cfg, cliArgs, stdout, stderr)

	return &RunResponse{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, nil
}
