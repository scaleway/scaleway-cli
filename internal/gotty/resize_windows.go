//go:build windows

package gotty

func subscribeToResize(resizeChan chan bool) func() {
	// Platform not supported
	return func() {}
}
