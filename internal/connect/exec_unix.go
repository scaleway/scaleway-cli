//go:build unix

package connect

import (
	"os"
	"os/exec"
	"os/signal"
)

func runChild(cmd *exec.Cmd) (int, error) {
	session := listenForSignals()
	defer signal.Stop(session.signals)

	return session.runChild(cmd)
}

func forwardSignal(process *os.Process, sig os.Signal) {
	// The child may already have exited after receiving the terminal signal.
	_ = process.Signal(sig)
}
