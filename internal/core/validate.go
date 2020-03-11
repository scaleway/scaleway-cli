package core

import (
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

// CommandValidateFunc validates en entire command.
// Used in core.cobraRun().
type CommandValidateFunc func(cmd *Command, cmdArgs interface{}) error

// ArgSpecValidateFunc validates one argument of a command.
type ArgSpecValidateFunc func(argSpec *ArgSpec, value interface{}) error

// DefaultCommandValidateFunc is the default validation function for commands.
func DefaultCommandValidateFunc() CommandValidateFunc {
	return func(cmd *Command, cmdArgs interface{}) error {
		err := validateArgValues(cmd, cmdArgs)
		if err != nil {
			return err
		}
		err = validateRequiredArgs(cmd, cmdArgs)
		if err != nil {
			return err
		}
		return nil
	}
}

// validateArgValues validates values passed to the different args of a Command.
func validateArgValues(cmd *Command, cmdArgs interface{}) error {
	for _, argSpec := range cmd.ArgSpecs {
		fieldName := strcase.ToPublicGoName(argSpec.Name)
		fieldValues, err := getValuesForFieldByName(reflect.ValueOf(cmdArgs), strings.Split(fieldName, "."))
		if err != nil {
			logger.Infof("could not validate arg value for '%v': invalid fieldName: %v: %v", argSpec.Name, fieldName, err.Error())
			continue
		}
		validateFunc := DefaultArgSpecValidateFunc()
		if argSpec.ValidateFunc != nil {
			validateFunc = argSpec.ValidateFunc
		}
		for _, fieldValue := range fieldValues {
			err := validateFunc(argSpec, fieldValue.Interface())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// validateRequiredArgs checks for missing required args with no default value.
// Returns an error for the first missing required arg.
// Returns nil otherwise.
func validateRequiredArgs(cmd *Command, cmdArgs interface{}) error {
	for _, arg := range cmd.ArgSpecs {
		fieldName := strcase.ToPublicGoName(arg.Name)
		fieldValues, err := getValuesForFieldByName(reflect.ValueOf(cmdArgs), strings.Split(fieldName, "."))
		if err != nil {
			logger.Infof("could not validate arg value for '%v': invalid fieldName: %v: %v", arg.Name, fieldName, err.Error())
			if arg.Required {
				return err
			}
			continue
		}

		for _, fieldValue := range fieldValues {
			if arg.Required && (fieldValue.IsZero() || !fieldValue.IsValid()) {
				return MissingRequiredArgumentError(arg.Name)
			}
		}
	}
	return nil
}

// DefaultArgSpecValidateFunc validates a value passed for an ArgSpec
// Uses ArgSpec.EnumValues
func DefaultArgSpecValidateFunc() ArgSpecValidateFunc {
	return func(argSpec *ArgSpec, value interface{}) error {
		if len(argSpec.EnumValues) < 1 {
			return nil
		}

		strValue, err := args.MarshalValue(value)
		if err != nil {
			return err
		}

		// When an enum is not provided as an argument args.MarshalValue will in most cases return "" (go default value)
		// In those cases we ignore validation. This is not ideal but covers most of the use cases.
		// The only caveat would be that `my-enum=""` would not trigger an error, which is acceptable.
		if strValue == "" {
			return nil
		}

		if !stringExists(argSpec.EnumValues, strValue) {
			return InvalidValueForEnumError(argSpec.Name, argSpec.EnumValues, strValue)
		}
		return nil
	}
}

func stringExists(strs []string, s string) bool {
	for _, s2 := range strs {
		if s == s2 {
			return true
		}
	}
	return false
}

func ValidateSecretKey() ArgSpecValidateFunc {
	return func(argSpec *ArgSpec, valueI interface{}) error {
		value := valueI.(string)
		err := DefaultArgSpecValidateFunc()(argSpec, value)
		if err != nil {
			return err
		}
		if !validation.IsSecretKey(value) {
			return InvalidSecretKeyError(value)
		}
		return nil
	}
}

// ValidateOrganizationID validates a non-required organization ID.
// By default, for most command, the organization ID is not required.
// In that case, we allow the empty-string value "".
func ValidateOrganizationID() ArgSpecValidateFunc {
	return func(argSpec *ArgSpec, valueI interface{}) error {
		value := valueI.(string)
		if value == "" && !argSpec.Required {
			return nil
		}
		return ValidateOrganizationIDRequired()(argSpec, valueI)
	}
}

// ValidateOrganizationIDRequired validates a required organization ID.
// We do not allow empty-string value "".
func ValidateOrganizationIDRequired() ArgSpecValidateFunc {
	return func(argSpec *ArgSpec, valueI interface{}) error {
		value := valueI.(string)
		err := DefaultArgSpecValidateFunc()(argSpec, value)
		if err != nil {
			return err
		}
		if !validation.IsOrganizationID(value) {
			return InvalidOrganizationIDError(value)
		}
		return nil
	}
}
