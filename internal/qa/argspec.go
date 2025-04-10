package qa

import (
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/args"
)

type ArgSpecInvalidError struct {
	Command    *core.Command
	argSpec    *core.ArgSpec
	innerError error
}

func (err ArgSpecInvalidError) Error() string {
	return fmt.Sprintf("command has invalid argspecs '%s' '%s' '%s'",
		err.Command.GetCommandLine("scw"),
		err.argSpec.Name,
		err.innerError,
	)
}

// testArgSpecInvalidError tests that all argspecs have a corresponding in their command's argstype.
func testArgSpecInvalidError(commands *core.Commands) []error {
	errors := []error(nil)

	for _, command := range commands.GetAll() {
		for _, arg := range command.ArgSpecs {
			_, err := args.GetArgType(command.ArgsType, arg.Name)
			if err != nil {
				errors = append(
					errors,
					&ArgSpecInvalidError{Command: command, argSpec: arg, innerError: err},
				)

				continue
			}
		}
	}

	return errors
}

type ArgSpecMissingError struct {
	Command *core.Command
	argName string
}

func (err ArgSpecMissingError) Error() string {
	return fmt.Sprintf("command has a missing argspec '%s' '%s'",
		err.Command.GetCommandLine("scw"),
		err.argName,
	)
}

// testArgSpecInvalidError tests that all argstype fields have a corresponding argspec.
func testArgSpecMissingError(commands *core.Commands) []error {
	errors := []error(nil)

	// Check all commands
	for _, command := range commands.GetAll() {
		if command.ArgsType == nil || command.ArgsType == reflect.TypeOf(args.RawArgs{}) {
			continue
		}

		supposedArgSpecs := args.ListArgTypeFields(command.ArgsType)

		for _, argSpecName := range supposedArgSpecs {
			if command.ArgSpecs.GetByName(argSpecName) == nil {
				errors = append(
					errors,
					&ArgSpecMissingError{Command: command, argName: argSpecName},
				)
			}
		}
	}

	return errors
}
