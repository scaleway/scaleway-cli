//go:build !unix

package connect

import (
	"os"
	"os/exec"
)

func runChild(cmd *exec.Cmd) (int, error) {
	return childResult(cmd.Run())
}

func forwardSignal(process *os.Process, _ os.Signal) {
	// Windows cannot send an interrupt to an individual process.
	_ = process.Kill()
}
