package qa

import (
	"strconv"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func LintCommands(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	errors = append(errors, testShortEndWithDotError(commands)...)
	errors = append(errors, testShortIsNotPresentError(commands)...)
	errors = append(errors, testArgMustUseDashError(commands)...)
	errors = append(errors, testExampleCanHaveOnlyOneTypeOfExampleError(commands)...)
	return errors
}

type ShortMustNotEndWithDotError struct {
	Command *core.Command
}

func (err ShortMustNotEndWithDotError) Error() string {
	return "short must not end with '.' for command '" + err.Command.GetCommandLine() + "'"
}

func testShortEndWithDotError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		if strings.HasSuffix(command.Short, ".") {
			errors = append(errors, &ShortMustNotEndWithDotError{
				Command: command,
			})
		}
	}
	return errors
}

type ShortMustBePresentError struct {
	Command *core.Command
}

func (err ShortMustBePresentError) Error() string {
	return "short must be present for command '" + err.Command.GetCommandLine() + "'"
}

func testShortIsNotPresentError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		if command.Short == "" {
			errors = append(errors, &ShortMustBePresentError{
				Command: command,
			})
		}
	}
	return errors
}

type ArgMustUseDashError struct {
	Command *core.Command
	Argspec *core.ArgSpec
}

func (err ArgMustUseDashError) Error() string {
	return "arg must use dash for command '" + err.Command.GetCommandLine() + "', arg '" + err.Argspec.Name + "'"
}

func testArgMustUseDashError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		for _, argspec := range command.ArgSpecs {
			if strings.Contains(argspec.Name, "_") {
				errors = append(errors, &ArgMustUseDashError{
					Command: command,
					Argspec: argspec,
				})
			}
		}
	}
	return errors
}

type ExampleCanHaveOnlyOneTypeOfExampleError struct {
	Command      *core.Command
	ExampleIndex int
}

func (err ExampleCanHaveOnlyOneTypeOfExampleError) Error() string {
	return "arg must use dash for command '" + err.Command.GetCommandLine() + "', example #" + strconv.Itoa(err.ExampleIndex)
}

func testExampleCanHaveOnlyOneTypeOfExampleError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		for i, example := range command.Examples {
			if example.Request != "" && example.Raw != "" {
				errors = append(errors, &ExampleCanHaveOnlyOneTypeOfExampleError{
					Command:      command,
					ExampleIndex: i,
				})
			}
		}
	}
	return errors
}
