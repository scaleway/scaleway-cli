package ssh

import (
	"context"
	"os/exec"
	"strconv"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Request holds the common arguments of the server ssh commands.
type Request struct {
	Zone     scw.Zone
	ServerID string
	Username string
	Port     uint64
	Command  string
}

// Connect opens an ssh session to the given address and propagates the
// remote exit code.
func Connect(ctx context.Context, req Request, address string) (any, error) {
	sshArgs := []string{
		address,
		"-p", strconv.FormatUint(req.Port, 10),
		"-l", req.Username,
		"-t",
	}
	if req.Command != "" {
		sshArgs = append(sshArgs, req.Command)
	}

	exitCode, err := core.ExecCmd(ctx, exec.Command("ssh", sshArgs...))
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		return nil, &core.CliError{Empty: true, Code: exitCode}
	}

	return &core.SuccessResult{Empty: true}, nil
}
