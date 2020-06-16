package autocomplete

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func unsupportedShellError(shell string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("unsupported shell '%v'", shell),
	}
}

func unsupportedOsError(OS string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("unsupported OS '%v'", OS),
	}
}

func installationCancelledError(shellName string, script string) *core.CliError {
	return &core.CliError{
		Err:  fmt.Errorf("installation cancelled"),
		Hint: fmt.Sprintf("To manually enable autocomplete for %v, run: %v", shellName, script),
	}
}

func installationNotFound(shellName string, location string, script string) *core.CliError {
	return &core.CliError{
		Err:  fmt.Errorf("not found where to install autocomplete script (tried %s)", location),
		Hint: fmt.Sprintf("You can add this line: `%s` in your %s configuration file", script, shellName),
	}
}
