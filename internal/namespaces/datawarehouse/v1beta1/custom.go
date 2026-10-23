package datawarehouse

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("datawarehouse").Groups = []string{"database"}

	cmds.Merge(core.NewCommands(deploymentCockpitCommand()))

	return cmds
}
