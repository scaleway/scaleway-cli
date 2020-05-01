// +build windows

package gotty

import "os"

func subscribeToResize(resizeChan chan os.Signal) func() {
	// Platform not supported
	return func() {}
}
