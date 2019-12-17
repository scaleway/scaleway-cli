package interactive

import (
	"fmt"
	"io"
	"os"

	"github.com/scaleway/scaleway-sdk-go/validation"

	"github.com/chzyer/readline"
	isatty "github.com/mattn/go-isatty"
)

var IsInteractive = false
var outputWriter io.Writer = os.Stderr

type ValidateFunc func(string) error

var defaultValidateFunc = func(string) error { return nil }

func init() {
	IsInteractive = isInteractive()
	readline.Stdout = os.Stderr
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
