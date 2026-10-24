package core

import (
	"errors"
	"fmt"
)

var (
	ErrMissingRequiredArgument        = errors.New("missing required argument")
	ErrInvalidValue                   = errors.New("invalid value")
	ErrInvalidSecretKey               = errors.New("invalid secret_key")
	ErrInvalidAccessKey               = errors.New("invalid access_key")
	ErrInvalidOrganizationID          = errors.New("invalid organization-id")
	ErrInvalidProjectID               = errors.New("invalid project-id")
	ErrInvalidRegion                  = errors.New("invalid region")
	ErrInvalidZone                    = errors.New("invalid zone")
	ErrInvalidAPIURL                  = errors.New("invalid api_url")
	ErrOnlyOneOfThoseTwoArguments     = errors.New("only one of those two arguments")
	ErrWindowsIsNotCurrentlySupported = errors.New("windows is not currently supported")
)

func MissingRequiredArgumentError(argumentName string) *CliError {
	return &CliError{
		Err: fmt.Errorf("%w '%v'", ErrMissingRequiredArgument, argumentName),
	}
}

func InvalidValueForEnumError(
	argSpecName string,
	argSpecEnumValues []string,
	value string,
) *CliError {
	return &CliError{
		Err:  fmt.Errorf("%w '%v' for arg '%v'", ErrInvalidValue, value, argSpecName),
		Hint: fmt.Sprintf("Accepted values for '%v' are %q", argSpecName, argSpecEnumValues),
	}
}

func InvalidSecretKeyError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("%w '%v'", ErrInvalidSecretKey, value),
		Hint: "secret_key should be a valid UUID, formatted as: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX.",
	}
}

func InvalidAccessKeyError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("%w '%v'", ErrInvalidAccessKey, value),
		Hint: "access_key should look like: SCWXXXXXXXXXXXXXXXXX.",
	}
}

func InvalidOrganizationIDError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("%w '%v'", ErrInvalidOrganizationID, value),
		Hint: "organization-id should be a valid UUID, formatted as: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX.",
	}
}

func InvalidProjectIDError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("%w '%v'", ErrInvalidProjectID, value),
		Hint: "project-id should be a valid UUID, formatted as: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX.",
	}
}

func InvalidRegionError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("%w '%v'", ErrInvalidRegion, value),
		Hint: "region format should look like: XX-XXX (e.g. fr-par).",
	}
}

func InvalidZoneError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("%w '%v'", ErrInvalidZone, value),
		Hint: "zone format should look like XX-XXX-X: (e.g. fr-par-1).",
	}
}

func InvalidAPIURLError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("%w '%v'", ErrInvalidAPIURL, value),
		Hint: "api_url should look like: https://www.example.com (e.g. https://api.scaleway.com).",
	}
}

func ArgumentConflictError(arg1 string, arg2 string) *CliError {
	return &CliError{
		Err: fmt.Errorf(
			"%w '%s' and '%s' can be specified in the same time",
			ErrOnlyOneOfThoseTwoArguments,
			arg1,
			arg2,
		),
	}
}

func WindowIsNotSupportedError() *CliError {
	return &CliError{
		Err: ErrWindowsIsNotCurrentlySupported,
	}
}
