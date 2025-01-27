package qa

import (
	"encoding/json"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
)

type CommandInvalidJSONExampleError struct {
	Command *core.Command
}

func (err CommandInvalidJSONExampleError) Error() string {
	return fmt.Sprintf("command has invalid json examples '%s'",
		err.Command.GetCommandLine("scw"),
	)
}

// testArgSpecInvalidError tests that all argspecs have a corresponding in their command's argstype.
func testCommandInvalidJSONExampleError(commands *core.Commands) []error {
	errors := []error(nil)

	for _, command := range commands.GetAll() {
		for _, example := range command.Examples {
			if example.ArgsJSON != "" {
				out := map[string]any{}
				err := json.Unmarshal([]byte(example.ArgsJSON), &out)
				if err != nil {
					errors = append(errors, &CommandInvalidJSONExampleError{
						Command: command,
					})
				}
			}
		}
	}

	return errors
}
