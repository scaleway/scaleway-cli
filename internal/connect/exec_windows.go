//go:build windows

package connect

import (
	"os"
	"os/exec"
)

func runChild(cmd *exec.Cmd) (int, error) {
	return childResult(cmd.Run())
}

func forwardSignal(process *os.Process, _ os.Signal) {
	// Windows cannot send os.Interrupt to an individual process. End only
	// this session’s direct child; other console processes may be unrelated.
	_ = process.Kill()
}
