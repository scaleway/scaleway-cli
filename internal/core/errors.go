package core

import (
	"fmt"
)

func MissingRequiredArgumentError(argumentName string) *CliError {
	return &CliError{
		Err: fmt.Errorf("missing required argument '%v'", argumentName),
	}
}

func InvalidValueForEnumError(argSpecName string, argSpecEnumValues []string, value string) *CliError {
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

func InvalidOrganizationIdError(value string) *CliError {
	return &CliError{
		Err:  fmt.Errorf("invalid organization-id '%v'", value),
		Hint: "organization-id should be a valid UUID, formatted as: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX.",
	}
}
