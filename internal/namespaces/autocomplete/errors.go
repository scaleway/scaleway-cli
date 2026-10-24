package autocomplete

import (
	"errors"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
)

var (
	ErrUnsupportedShell             = errors.New("unsupported shell")
	ErrUnsupportedOS                = errors.New("unsupported OS")
	ErrAutocompleteLocationNotFound = errors.New(
		"cannot find where to install autocomplete script",
	)
)

func unsupportedShellError(shell string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("%w '%v'", ErrUnsupportedShell, shell),
	}
}

func unsupportedOsError(os string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("%w '%v'", ErrUnsupportedOS, os),
	}
}

func installationCancelledError(shellName string, script string) *core.CliError {
	return &core.CliError{
		Err:  core.ErrInstallationCancelled,
		Hint: fmt.Sprintf("To manually enable autocomplete for %v, run: %v", shellName, script),
	}
}

func installationNotFound(shellName string, location string, script string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("%w (tried %s)", ErrAutocompleteLocationNotFound, location),
		Hint: fmt.Sprintf(
			"You can add this line: `%s` in your %s configuration file",
			script,
			shellName,
		),
	}
}
