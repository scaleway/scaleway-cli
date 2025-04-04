package autocomplete

import (
	"errors"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func unsupportedShellError(shell string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("unsupported shell '%v'", shell),
	}
}

func unsupportedOsError(os string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("unsupported OS '%v'", os),
	}
}

func installationCancelledError(shellName string, script string) *core.CliError {
	return &core.CliError{
		Err:  errors.New("installation cancelled"),
		Hint: fmt.Sprintf("To manually enable autocomplete for %v, run: %v", shellName, script),
	}
}

func installationNotFound(shellName string, location string, script string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("cannot find where to install autocomplete script (tried %s)", location),
		Hint: fmt.Sprintf(
			"You can add this line: `%s` in your %s configuration file",
			script,
			shellName,
		),
	}
}
