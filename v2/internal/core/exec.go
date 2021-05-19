package core

import (
	"context"
	"os"
	"os/exec"
)

type OverrideExecFunc func(cmd *exec.Cmd) (exitCode int, err error)

func defaultOverrideExec(cmd *exec.Cmd) (exitCode int, err error) {
	err = cmd.Run()
	if exitErr, isExitErr := err.(*exec.ExitError); isExitErr {
		return exitErr.ExitCode(), nil
	}
	return 0, err
}

func ExecCmd(ctx context.Context, cmd *exec.Cmd) (exitCode int, err error) {
	meta := extractMeta(ctx)

	// We do not support override of stdin
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}

	cmd.Stdout = meta.stdout
	cmd.Stderr = meta.stderr
	return meta.OverrideExec(cmd)
}
