package datawarehouse

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("datawarehouse").Groups = []string{"database"}

	cmds.MustFind("datawarehouse", "deployment", "create").Override(deploymentCreateBuilder)
	cmds.MustFind("datawarehouse", "deployment", "delete").Override(deploymentDeleteBuilder)
	cmds.MustFind("datawarehouse", "deployment", "get").SeeAlsos = connectSeeAlsos
	cmds.MustFind("datawarehouse", "user", "update").Override(userUpdateBuilder)

	cmds.Merge(core.NewCommands(
		deploymentConnectCommand(),
	))

	return cmds
}
