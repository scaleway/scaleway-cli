package connect

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
)

// SetEnv replaces any existing key in the command's environment with
// key=value, so the child process only sees this session's value.
func SetEnv(cmd *exec.Cmd, key, value string) {
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}
	for i, e := range cmd.Env {

		if !strings.HasPrefix(e, key+"=") {
			env = append(env, e)
		}
	}
	env = append(env, key+"="+value)
	cmd.Env = env
}

// ExecCmd runs the child process of a connect session.
// On unix, SIGINT/SIGTERM received while the session is running are forwarded
// to the child process, so an interrupted session (e.g. Ctrl-C during
// `kubectl get pods -w`) terminates the child and returns through the
// command's cleanup (e.g. temporary kubeconfig removal).
// A child killed by a signal exits with 128+signal, the shell convention.
func ExecCmd(ctx context.Context, cmd *exec.Cmd) (exitCode int, err error) {
	meta := core.ExtractMeta(ctx)

	// We do not support override of stdin
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = core.ExtractStdout(ctx)
	cmd.Stderr = core.ExtractStderr(ctx)

	if meta != nil && meta.OverrideExec != nil {
		return meta.OverrideExec(cmd)
	}

	return runChild(cmd)
}
