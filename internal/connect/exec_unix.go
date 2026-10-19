//go:build unix

package connect

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// runChild runs cmd and forwards SIGINT/SIGTERM to the child process so an
// interrupted session terminates the child and the caller can return through
// its cleanup instead of being killed by the runtime's default signal handler.
func runChild(cmd *exec.Cmd) (exitCode int, err error) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigs)

	if err := cmd.Start(); err != nil {
		return 0, err
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case sig := <-sigs:
				// Forward to the child. It may have already exited (the terminal
				// sends the signal to the whole foreground group too), ignore errors.
				// ponytail: never escalate to SIGKILL; a child that ignores the
				// signal keeps the session alive. Upgrade path: cmd.WaitDelay.
				if s, ok := sig.(syscall.Signal); ok {
					_ = cmd.Process.Signal(s)
				}
			case <-done:
				return
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			return 0, err
		}
		// The child was killed by a signal (e.g. the forwarded SIGINT):
		// report the shell convention of 128+signal.
		if exitErr.ExitCode() == -1 {
			if ws, ok := exitErr.ProcessState.Sys().(syscall.WaitStatus); ok {
				if sig := ws.Signal(); sig != 0 {
					return 128 + int(sig), nil
				}
			}
		}

		return exitErr.ExitCode(), nil
	}

	return 0, nil
}
