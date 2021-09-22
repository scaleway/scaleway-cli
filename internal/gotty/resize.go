//nolint
// +build !windows

package gotty

import (
	"os"
	"os/signal"
	"syscall"
)

func subscribeToResize(resizeChan chan bool) func() {
	sigChan := make(chan os.Signal, 1)

	go func() {
		for {
			sig := <-sigChan
			if sig == nil {
				return
			}
			resizeChan <- true
		}
	}()

	signal.Notify(sigChan, syscall.SIGWINCH)
	return func() {
		signal.Stop(sigChan)
		close(sigChan)
	}
}
