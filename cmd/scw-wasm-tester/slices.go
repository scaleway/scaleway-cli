//go:build wasm && js

package main

import (
	"syscall/js"

	"github.com/scaleway/scaleway-cli/v2/internal/jshelpers"
)

func wasmTestFromSlice(_ js.Value, _ []js.Value) any {
	return jshelpers.FromSlice([]string{"1", "2", "3"})
}
