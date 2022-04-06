package redis

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	return cmds
}
