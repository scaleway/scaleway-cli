// +build !windows
package gotty

import (
	"os"
	"os/signal"
	"syscall"
)

func subscribeToResize(resizeChan chan os.Signal) func() {
	signal.Notify(resizeChan, syscall.SIGWINCH)
	return func() {
		signal.Stop(resizeChan)
	}
}
