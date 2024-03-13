//go:build wasm

package interactive

import (
	"io"
)

type ValidateFunc func(string) error

var (
	// IsInteractive must be set to print anything with Printer functions (Print, Printf,...).
	IsInteractive = false

	// OutputWriter is the writer used by Printer functions (Print, Printf,...).
	outputWriter io.Writer
)

// SetOutputWriter set the output writer that will be used by both Printer functions (Print, Printf,...) and
// readline prompter. This should be called once from the bootstrap function.
func SetOutputWriter(w io.Writer) {
	outputWriter = w
}
