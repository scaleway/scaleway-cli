package vpc

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("vpc", "private-network", "get").Override(privateNetworkGetBuilder)

	return cmds
}
