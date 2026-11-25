package core

import (
	"context"
	"errors"
	"os"
	"os/exec"
)

// OverrideExecFunc replaces subprocess execution in CLI tests.
type OverrideExecFunc func(cmd *exec.Cmd) (exitCode int, err error)

func defaultOverrideExec(cmd *exec.Cmd) (exitCode int, err error) {
	err = cmd.Run()
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		return exitErr.ExitCode(), nil
	}

	return 0, err
}

// ExecCmd runs cmd with the CLI streams and honors the configured execution hook.
//
// Deprecated: Retained for compatibility with external Go callers. New commands
// in scaleway-cli should use internal/connect.ExecCmd.
func ExecCmd(ctx context.Context, cmd *exec.Cmd) (exitCode int, err error) {
	meta := ExtractMeta(ctx)

	// We do not support override of stdin
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}

	cmd.Stdout = meta.stdout
	cmd.Stderr = meta.stderr

	// A nil OverrideExec means the child process must really be executed.
	if meta.OverrideExec != nil {
		return meta.OverrideExec(cmd)
	}

	return defaultOverrideExec(cmd)
}
