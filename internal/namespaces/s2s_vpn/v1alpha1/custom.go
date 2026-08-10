package s2s_vpn

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("s2s-vpn").Groups = []string{"network"}

	return cmds
}
