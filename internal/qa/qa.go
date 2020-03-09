package qa

import (
	"strings"

	"github.com/scaleway/scaleway-cli/internal/core"
)

type shortEndWithDotError struct {
	Command *core.Command
}

func (err shortEndWithDotError) Error() string {
	return "short ends with '.' for command '" + err.Command.GetCommandLine() + "'"
}

func testShortEndWithDotError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		if strings.HasSuffix(command.Short, ".") {
			errors = append(errors, &shortEndWithDotError{
				Command: command,
			})
		}
	}
	return errors
}

type shortIsNotPresent struct {
	Command *core.Command
}

func (err shortIsNotPresent) Error() string {
	return "short is not present for command '" + err.Command.GetCommandLine() + "'"
}

func testShortIsNotPresentError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		if strings.HasSuffix(command.Short, ".") {
			errors = append(errors, &shortIsNotPresent{
				Command: command,
			})
		}
	}
	return errors
}

func LintCommands(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	errors = append(errors, testShortEndWithDotError(commands)...)
	errors = append(errors, testShortIsNotPresentError(commands)...)
	return errors
}
