//go:build unix

package connect

import (
	"os"
	"os/exec"
)

func runChild(cmd *exec.Cmd) (int, error) {
	session := ListenForSignals()
	defer session.Close()

	return session.runChild(cmd)
}

func forwardSignal(process *os.Process, sig os.Signal) {
	// The child may already have exited after receiving the terminal signal.
	_ = process.Signal(sig)
}
