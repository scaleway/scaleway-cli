//go:build !wasm

package interactive

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	isatty "github.com/mattn/go-isatty"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

type ValidateFunc func(string) error

var (
	// IsInteractive must be set to print anything with Printer functions (Print, Printf,...).
	IsInteractive = isInteractive()

	// TerminalOutput define if the
	TerminalOutput = IsInteractive

	// outputWriter is the writer used by Printer functions (Print, Printf,...).
	outputWriter io.Writer

	// defaultValidateFunc is used by readline to validate user input.
	defaultValidateFunc ValidateFunc = func(string) error { return nil }
)

// SetOutputWriter set the output writer that will be used by both Printer functions (Print, Printf,...) and
// readline prompter. This should be called once from the bootstrap function.
func SetOutputWriter(w io.Writer) {
	outputWriter = w
	readline.Stdout = newWriteCloser(w)
}

// we should expect both Stdin and Stderr to enable interactive mode
func isInteractive() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) && isatty.IsTerminal(os.Stderr.Fd()) ||
		isatty.IsCygwinTerminal(os.Stdin.Fd()) && isatty.IsCygwinTerminal(os.Stderr.Fd()) // windows cygwin terminal
}

func ValidateOrganizationID() ValidateFunc {
	return func(s string) error {
		if !validation.IsOrganizationID(s) {
			return fmt.Errorf("invalid organization-id")
		}
		return nil
	}
}

func newWriteCloser(w io.Writer) io.WriteCloser {
	return &writeCloser{w}
}

type writeCloser struct {
	io.Writer
}

func (wc *writeCloser) Close() error {
	if closer, ok := wc.Writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
