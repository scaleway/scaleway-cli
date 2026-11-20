//go:build !wasm

package core

import (
	"os"
	"syscall"
)

// followStopSignals end a stream=true run cleanly, so StreamFunc can release what it created.
// SIGHUP is sent when the terminal is closed.
var followStopSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGHUP}
