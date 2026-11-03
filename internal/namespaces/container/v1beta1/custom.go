package container

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("container").Groups = []string{"serverless"}

	cmds.MustFind("container", "token", "list").
		Override(core.ListViewBuilder("ID", "Token", "ContainerID", "NamespaceID", "PublicKey", "Status", "Description", "ExpiresAt"))

	return cmds
}
