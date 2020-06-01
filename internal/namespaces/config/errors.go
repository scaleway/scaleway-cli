package config

//
// Functions in this file can only return non-nil *core.CliError.
//

import (
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func invalidDefaultOrganizationIDError(value string) *core.CliError {
	return &core.CliError{
		Message: fmt.Sprintf("invalid default_organization_id '%v'", value),
		Hint:    "default_organization_id should be a valid UUID, formatted as: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX.",
	}
}

func invalidProfileKeyError(fieldName string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("invalid profile's key identifier %s", fieldName),
	}
}

func invalidRegionError(value string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("invalid region %s", value),
	}
}

func invalidZoneError(value string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("invalid zone %s", value),
	}
}

func notEnoughArgsForConfigSetError() *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("not enough args: enter a key and a value"),
	}
}

func missingValueForConfigSetError(key string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("missing value for key %s", key),
	}
}

func tooManyArgsForConfigSetError() *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("too many args: only one value can be set at a time"),
	}
}

func unknownProfileError(profileName string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("no profile named %s", profileName),
	}
}

func invalidKindForKeyError(kind reflect.Kind, fieldName string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("invalid kind %v for key %s", kind, fieldName),
	}
}
