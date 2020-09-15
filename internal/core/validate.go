package core

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

// CommandValidateFunc validates en entire command.
// Used in core.cobraRun().
type CommandValidateFunc func(ctx context.Context, cmd *Command, cmdArgs interface{}, rawArgs args.RawArgs) error

// ArgSpecValidateFunc validates one argument of a command.
type ArgSpecValidateFunc func(argSpec *ArgSpec, value interface{}) error

// DefaultCommandValidateFunc is the default validation function for commands.
func DefaultCommandValidateFunc() CommandValidateFunc {
	return func(ctx context.Context, cmd *Command, cmdArgs interface{}, rawArgs args.RawArgs) error {
		err := validateArgValues(cmd, cmdArgs)
		if err != nil {
			return err
		}
		err = validateRequiredArgs(cmd, cmdArgs, rawArgs)
		if err != nil {
			return err
		}
		err = validateNoConflict(cmd, rawArgs)
		if err != nil {
			return err
		}

		validateDeprecated(ctx, cmd)
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
// TODO refactor this method which uses a mix of reflect and string arrays
func validateRequiredArgs(cmd *Command, cmdArgs interface{}, rawArgs args.RawArgs) error {
	for _, arg := range cmd.ArgSpecs {
		if !arg.Required {
			continue
		}

		fieldName := strcase.ToPublicGoName(arg.Name)
		fieldValues, err := getValuesForFieldByName(reflect.ValueOf(cmdArgs), strings.Split(fieldName, "."))
		if err != nil {
			validationErr := fmt.Errorf("could not validate arg value for '%v': invalid field name '%v': %v", arg.Name, fieldName, err.Error())
			if !arg.Required {
				logger.Infof(validationErr.Error())
				continue
			}
			panic(validationErr)
		}

		// Either fieldsValues have a length for 1 and we check for existence in the rawArgs
		// or it has multiple values and we loop through each one to get the right element in
		// the corresponding rawArgs array and replace {index} by the element's index.
		// TODO handle required maps
		for i := range fieldValues {
			if !rawArgs.ExistsArgByName(strings.Replace(arg.Name, "{index}", strconv.Itoa(i), 1)) {
				return MissingRequiredArgumentError(strings.Replace(arg.Name, "{index}", strconv.Itoa(i), 1))
			}
		}
	}
	return nil
}

func validateNoConflict(cmd *Command, rawArgs args.RawArgs) error {
	for _, arg1 := range cmd.ArgSpecs {
		for _, arg2 := range cmd.ArgSpecs {
			if !arg1.ConflictWith(arg2) || arg1 == arg2 {
				continue
			}
			if rawArgs.Has(arg1.Name) && rawArgs.Has(arg2.Name) {
				return ArgumentConflictError(arg1.Name, arg2.Name)
			}
		}
	}
	return nil
}

// validateDeprecated print a warning message if a deprecated argument is used
func validateDeprecated(ctx context.Context, cmd *Command) {
	deprecatedArgs := cmd.ArgSpecs.GetDeprecated(true)
	for _, argSpec := range deprecatedArgs {
		helpCmd := cmd.GetCommandLine(extractMeta(ctx).BinaryName) + " --help"
		ExtractLogger(ctx).Warningf("The argument '%s' is deprecated, more info with: %s\n", argSpec.Name, helpCmd)
	}
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
		value, isStr := valueI.(string)
		valuePtr, isPtr := valueI.(*string)
		if !isStr && isPtr && valuePtr != nil {
			value = *valuePtr
		}

		if value == "" && !argSpec.Required {
			return nil
		}
		if !validation.IsOrganizationID(value) {
			return InvalidOrganizationIDError(value)
		}
		return nil
	}
}

// ValidateProjectID validates a non-required project ID.
// By default, for most command, the project ID is not required.
// In that case, we allow the empty-string value "".
func ValidateProjectID() ArgSpecValidateFunc {
	return func(argSpec *ArgSpec, valueI interface{}) error {
		value, isStr := valueI.(string)
		valuePtr, isPtr := valueI.(*string)
		if !isStr && isPtr && valuePtr != nil {
			value = *valuePtr
		}

		if value == "" && !argSpec.Required {
			return nil
		}
		if !validation.IsProjectID(value) {
			return InvalidProjectIDError(value)
		}
		return nil
	}
}
