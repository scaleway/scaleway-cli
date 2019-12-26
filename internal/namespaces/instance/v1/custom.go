package instance

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

// custom_commands: contains handwritten commands
// custom_runs: contains handwritten Run functions
// custom_marshal: contains all the registered MarshalFunc and the human.Marshaler implementations

// GetCommands gets the generated commands, apply custom runs on it,
// merge it with custom commands and return the result.
func GetCommands() *core.Commands {
	// get generated commands
	instanceCommands := GetGeneratedCommands()

	// updates commands with custom fields
	updateCommands(instanceCommands)

	// apply custom runs on generated commands only
	applyCustomRuns(instanceCommands)

	// apply custom error handling on generated commands only
	applyCustomErrors(instanceCommands)

	// merge custom commands
	instanceCommands.Merge(getCustomCommands())

	return instanceCommands
}
