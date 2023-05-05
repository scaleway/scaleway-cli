package mnq

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("mnq", "credential", "create").Override(credentialCreateBuilder)
	cmds.MustFind("mnq", "credential", "get").Override(credentialGetBuilder)

	return cmds
}
