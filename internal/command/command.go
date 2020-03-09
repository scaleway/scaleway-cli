package command

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	autocompleteNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/autocomplete"
	configNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/config"
	initNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/init"
	"github.com/scaleway/scaleway-cli/internal/namespaces/instance/v1"
	k8s "github.com/scaleway/scaleway-cli/internal/namespaces/k8s/v1beta4"
	"github.com/scaleway/scaleway-cli/internal/namespaces/marketplace/v1"
	versionNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/version"
)

func GetCommands() *core.Commands {
	// Import all commands available in CLI from various packages.
	// NB: Merge order impacts scw usage sort.
	commands := core.NewCommands()
	commands.Merge(instance.GetCommands())
	commands.Merge(k8s.GetCommands())
	commands.Merge(marketplace.GetCommands())
	commands.Merge(initNamespace.GetCommands())
	commands.Merge(configNamespace.GetCommands())
	commands.Merge(autocompleteNamespace.GetCommands())
	commands.Merge(versionNamespace.GetCommands())
	return commands
}
