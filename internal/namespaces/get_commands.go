package namespaces

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	account "github.com/scaleway/scaleway-cli/internal/namespaces/account/v2alpha1"
	autocompleteNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/autocomplete"
	baremetal "github.com/scaleway/scaleway-cli/internal/namespaces/baremetal/v1"
	configNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/config"
	"github.com/scaleway/scaleway-cli/internal/namespaces/feedback"
	"github.com/scaleway/scaleway-cli/internal/namespaces/info"
	initNamespace "github.com/scaleway/scaleway-cli/internal/namespaces/init"
	"github.com/scaleway/scaleway-cli/internal/namespaces/instance/v1"
	k8s "github.com/scaleway/scaleway-cli/internal/namespaces/k8s/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/marketplace/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/object/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/rdb/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/registry/v1"
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
	commands.Merge(object.GetCommands())
	commands.Merge(versionNamespace.GetCommands())
	commands.Merge(registry.GetCommands())
	commands.Merge(feedback.GetCommands())
	commands.Merge(info.GetCommands())
	commands.Merge(rdb.GetCommands())
	return commands
}
