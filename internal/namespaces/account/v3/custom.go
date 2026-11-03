package account

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	commands := GetGeneratedCommands()

	commands.MustFind("account").Groups = []string{"security"}

	commands.MergeAll(core.NewCommands(
		listResources(),
	))

	commands.MustFind("account", "project", "list").
		Override(core.ListViewBuilder("ID", "Name", "OrganizationID", "CreatedAt", "UpdatedAt", "Description", "Status"))

	return commands
}
