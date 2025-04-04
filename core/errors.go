package core

import (
	"errors"
	"fmt"
)

func MissingRequiredArgumentError(argumentName string) *CliError {
	return &CliError{
		Err: fmt.Errorf("missing required argument '%v'", argumentName),
	}
}

func InvalidValueForEnumError(
	argSpecName string,
	argSpecEnumValues []string,
	value string,
) *CliError {
	return &CliError{
		Err:  fmt.Errorf("invalid value '%v' for arg '%v'", value, argSpecName),
		Hint: fmt.Sprintf("Accepted values for '%v' are %v", argSpecName, argSpecEnumValues),
	}
}

func InvalidSecretKeyError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("invalid secret_key '%v'", value),
		Hint: "secret_key should be a valid UUID, formatted as: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX.",
	}
}

func InvalidAccessKeyError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("invalid access_key '%v'", value),
		Hint: "access_key should look like: SCWXXXXXXXXXXXXXXXXX.",
	}
}

func InvalidOrganizationIDError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("invalid organization-id '%v'", value),
		Hint: "organization-id should be a valid UUID, formatted as: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX.",
	}
}

func InvalidProjectIDError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("invalid project-id '%v'", value),
		Hint: "project-id should be a valid UUID, formatted as: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX.",
	}
}

func InvalidRegionError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("invalid region '%v'", value),
		Hint: "region format should look like: XX-XXX (e.g. fr-par).",
	}
}

func InvalidZoneError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("invalid zone '%v'", value),
		Hint: "zone format should look like XX-XXX-X: (e.g. fr-par-1).",
	}
}

func InvalidAPIURLError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("invalid api_url '%v'", value),
		Hint: "api_url should look like: https://www.example.com (e.g. https://api.scaleway.com).",
	}
}

func ArgumentConflictError(arg1 string, arg2 string) *CliError {
	return &CliError{
		Err: fmt.Errorf(
			"only one of those two arguments '%s' and '%s' can be specified in the same time",
			arg1,
			arg2,
		),
	}
}

func WindowIsNotSupportedError() *CliError {
	return &CliError{
		Err: errors.New("windows is not currently supported"),
	}
}
