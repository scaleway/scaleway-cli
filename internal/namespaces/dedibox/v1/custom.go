package dedibox

import "github.com/scaleway/scaleway-cli/v2/internal/core"

func GetCommand() *core.Commands {
	cmds := GetGeneratedCommands()

	return cmds
}
