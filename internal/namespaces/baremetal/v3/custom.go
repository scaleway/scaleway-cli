package baremetal

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()
	cmds.MustFind("baremetal").Groups = []string{"baremetal"}

	return cmds
}
