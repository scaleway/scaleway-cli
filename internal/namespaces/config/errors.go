package config

//
// Functions in this file can only return non-nil *core.CliError.
//

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func invalidProfileKeyError(fieldName string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("invalid profile's key identifier %s", fieldName),
	}
}

func unknownProfileError(profileName string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("no profile named %s", profileName),
	}
}
