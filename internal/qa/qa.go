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

type ShortEndWithDotError struct {
	Command *core.Command
}

func (err ShortEndWithDotError) Error() string {
	return "short ends with '.' for command '" + err.Command.GetCommandLine() + "'"
}

func testShortEndWithDotError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		if strings.HasSuffix(command.Short, ".") {
			errors = append(errors, &ShortEndWithDotError{
				Command: command,
			})
		}
	}
	return errors
}

type ShortIsNotPresent struct {
	Command *core.Command
}

func (err ShortIsNotPresent) Error() string {
	return "short is not present for command '" + err.Command.GetCommandLine() + "'"
}

func testShortIsNotPresentError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		if command.Short == "" {
			errors = append(errors, &ShortIsNotPresent{
				Command: command,
			})
		}
	}
	return errors
}
