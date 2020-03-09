package main

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/qa"
)

func main() {
	commands := getCommands()
	errors := qa.LintCommands(commands)
	fmt.Println(errors)
}

func getCommands() *core.Commands {
	// Import all commands available in CLI from various packages.
	// NB: Merge order impacts scw usage sort.
	commands := core.NewCommands()
	//commands.Merge(instance.GetCommands())
	//commands.Merge(k8s.GetCommands())
	//commands.Merge(marketplace.GetCommands())
	//commands.Merge(initNamespace.GetCommands())
	//commands.Merge(configNamespace.GetCommands())
	//commands.Merge(autocompleteNamespace.GetCommands())
	//commands.Merge(versionNamespace.GetCommands())
	return commands
}
