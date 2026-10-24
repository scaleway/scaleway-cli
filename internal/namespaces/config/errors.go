package config

//
// Functions in this file can only return non-nil *core.CliError.
//

import (
	"errors"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
)

var (
	ErrInvalidProfilesKeyIdentifier = errors.New("invalid profile's key identifier")
	ErrNoProfileNamed               = errors.New("no profile named")
)

func invalidProfileKeyError(fieldName string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("%w %s", ErrInvalidProfilesKeyIdentifier, fieldName),
	}
}

func unknownProfileError(profileName string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("%w %s", ErrNoProfileNamed, profileName),
	}
}
