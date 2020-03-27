package namespaces

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/namespaces/account"
	autocompleteNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/autocomplete"
	baremetal "github.com/scaleway/scaleway-cli/internal/namespaces/baremetal/v1alpha1"
	configNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/config"
	initNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/init"
	"github.com/scaleway/scaleway-cli/internal/namespaces/instance/v1"
	k8s "github.com/scaleway/scaleway-cli/internal/namespaces/k8s/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/marketplace/v1"
	versionNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/version"
)

// GetCommands returns a list of all commands in the CLI.
// It is used by both scw and scw-qa.
// We can not put it in `core` package as it would result in a import cycle `core` -> `namespaces/autocomplete` -> `core`.
func GetCommands() *core.Commands {
	// Import all commands available in CLI from various packages.
	// NB: Merge order impacts scw usage sort.
	commands := core.NewCommands()
	commands.Merge(instance.GetCommands())
	commands.Merge(baremetal.GetCommands())
	commands.Merge(k8s.GetCommands())
	commands.Merge(marketplace.GetCommands())
	commands.Merge(initNamespace.GetCommands())
	commands.Merge(configNamespace.GetCommands())
	commands.Merge(account.GetCommands())
	commands.Merge(autocompleteNamespace.GetCommands())
	commands.Merge(versionNamespace.GetCommands())
	return commands
}
