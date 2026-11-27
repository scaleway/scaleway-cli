//go:build !unix && !windows

package connect

import (
	"os"
	"os/exec"
)

func runChild(cmd *exec.Cmd) (int, error) {
	return childResult(cmd.Run())
}

func forwardSignal(process *os.Process, sig os.Signal) {
	_ = process.Signal(sig)
}
