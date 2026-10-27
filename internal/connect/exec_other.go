//go:build !unix

package connect

import (
	"os/exec"
)

// runChild runs cmd. On platforms without unix signals there is nothing to
// forward, so the child is run as-is.
func runChild(cmd *exec.Cmd) (exitCode int, err error) {
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), nil
		}

		return 0, err
	}

	return 0, nil
}
