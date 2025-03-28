package ipam

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("ipam").Groups = []string{"network"}

	return cmds
}
