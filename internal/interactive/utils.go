//go:build !wasm

package interactive

import (
	"errors"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/scaleway/scaleway-sdk-go/validation"
	"golang.org/x/term"
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
	// x/term.IsTerminal uses GetConsoleMode, so Cygwin terminals on
	// Windows (pseudo-ttys over pipes) are no longer detected as interactive.
	// If Cygwin users report broken prompts, restore the isatty.IsCygwinTerminal check.
	return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stderr.Fd()))
}

func ValidateOrganizationID() ValidateFunc {
	return func(s string) error {
		if !validation.IsOrganizationID(s) {
			return errors.New("invalid organization-id")
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
