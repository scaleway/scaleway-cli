package registry

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
)

// GetCommands returns registry commands.
//
// This function:
// - Gets the generated commands
// - Register handwritten marshalers
// - Apply handwritten overrides (of Command.Run)
func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()
	cmds.Merge(core.NewCommands())

	human.RegisterMarshalerFunc(registry.NamespaceStatus(0), human.EnumMarshalFunc(namespaceStatusMarshalSpecs))
	human.RegisterMarshalerFunc(registry.ImageStatus(0), human.EnumMarshalFunc(imageStatusMarshalSpecs))
	human.RegisterMarshalerFunc(registry.TagStatus(0), human.EnumMarshalFunc(tagStatusMarshalSpecs))

	return cmds
}
