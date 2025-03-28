package cockpit

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("cockpit").Groups = []string{"monitoring"}

	cmds.MustFind("cockpit", "token", "get").Override(cockpitTokenGetBuilder)

	return cmds
}
