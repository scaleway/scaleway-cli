package qa

import (
	"strings"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func LintCommands(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	errors = append(errors, testShortEndWithDotError(commands)...)
	errors = append(errors, testShortIsNotPresentError(commands)...)
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
