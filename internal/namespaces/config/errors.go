package config

//
// Functions in this file can only return non-nil *core.CliError.
//

import (
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func defaultCliError(err error) *core.CliError {
	return &core.CliError{
		Err: err,
	}
}

func invalidDefaultOrganizationIdError(value string) *core.CliError {
	return &core.CliError{
		Err:  fmt.Errorf("invalid default_organization_id '%v'", value),
		Hint: "default_organization_id should be a valid UUID, formatted as: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX.",
	}
}

func invalidProfileKeyError(fieldName string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("invalid profile's key identifier '%v'", fieldName),
	}
}

func invalidRegionError(value string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("invalid region '%v'", value),
	}
}

func invalidZoneError(value string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("invalid zone '%v'", value),
	}
}

func notEnoughArgsForConfigSetError() *core.CliError {
	return &core.CliError{
		Err: defaultCliError(fmt.Errorf("not enough args: enter a key and a value")),
	}
}

func missingValueForConfigSetError(key string) *core.CliError {
	return &core.CliError{
		Err: defaultCliError(fmt.Errorf("missing value for key '%v'", key)),
	}
}

func tooManyArgsForConfigSetError() *core.CliError {
	return &core.CliError{
		Err: defaultCliError(fmt.Errorf("too many args: only one value can be set at a time")),
	}
}

func notEnoughArgsForConfigGetError() *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("not enough args: enter a key"),
	}
}

func notEnoughArgsForConfigUnsetError() *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("not enough args: enter a key"),
	}
}

func tooManyArgsForConfigUnsetError() *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("too many args: only one value can be unset at a time"),
	}
}

func unknownProfileError(profileName string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("no profile named '%v'", profileName),
	}
}

func invalidProfileKeyPairError(arg string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("invalid profile/key pair identifier '%v'", arg),
	}
}

func invalidProfileKeyIdentifierError(fieldName string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("invalid profile's key identifier '%v'", fieldName),
	}
}

func nilFieldError(fieldName string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("field not set '%v'", fieldName),
	}
}

func invalidKindForKeyError(kind reflect.Kind, fieldName string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("invalid kind '%v' for key '%v'", kind, fieldName),
	}
}

func invalidProfileAttributeError(key string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("invalid profile attribute: '%v'", key),
	}
}
