package datawarehouse

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("datawarehouse").Groups = []string{"database"}

	return cmds
}
