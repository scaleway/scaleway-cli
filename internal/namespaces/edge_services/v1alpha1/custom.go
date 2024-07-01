package edge_services

import "github.com/scaleway/scaleway-cli/v2/internal/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	return cmds
}
