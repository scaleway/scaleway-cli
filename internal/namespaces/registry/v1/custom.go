package registry

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

// GetCommands returns registry commands.
//
// This function:
// - Gets the generated commands
// - Register handwritten marshalers
// - Apply handwritten overrides (of Command.Run)
func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("registry").Groups = []string{"container", "storage"}

	cmds.Merge(core.NewCommands(
		registryLoginCommand(),
		registryLogoutCommand(),
		registryDockerHelperGetCommand(),
		registryDockerHelperEraseCommand(),
		registryDockerHelperListCommand(),
		registryDockerHelperStoreCommand(),
		registryInstallDockerHelperCommand(),
	))

	cmds.MustFind("registry", "tag", "delete").Override(tagDeleteBuilder)
	cmds.MustFind("registry", "tag", "get").Override(tagGetBuilder)
	cmds.MustFind("registry", "tag", "list").Override(tagListBuilder)

	cmds.MustFind("registry", "image", "get").Override(imageGetBuilder)
	cmds.MustFind("registry", "image", "list").Override(imageListBuilder)

	// TODO
	// human.RegisterMarshalerFunc(
	// 	registry.NamespaceStatus(""),
	// 	human.EnumMarshalFunc(namespaceStatusMarshalSpecs),
	// )
	// human.RegisterMarshalerFunc(
	// 	registry.ImageStatus(""),
	// 	human.EnumMarshalFunc(imageStatusMarshalSpecs),
	// )
	// human.RegisterMarshalerFunc(
	// 	registry.TagStatus(""),
	// 	human.EnumMarshalFunc(tagStatusMarshalSpecs),
	// )

	return cmds
}
