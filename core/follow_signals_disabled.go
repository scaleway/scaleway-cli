//go:build wasm

package core

import "os"

var followStopSignals = []os.Signal{os.Interrupt}
