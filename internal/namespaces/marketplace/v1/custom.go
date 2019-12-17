package marketplace

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

// custom_marshal: contains all the registered MarshalFunc and the human.Marshaler implementations
// updated_commands: contains updated commands

// GetCommands gets the generated commands, apply custom runs on it,
// merge it with custom commands and return the result.
func GetCommands() *core.Commands {
	// get generated commands
	marketplaceCommands := GetGeneratedCommands()

	// Update commands with custom implementation of fields
	updateCommands(marketplaceCommands)

	return marketplaceCommands
}
