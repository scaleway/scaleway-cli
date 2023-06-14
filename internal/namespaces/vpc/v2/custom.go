package vpc

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()
	for _, cmd := range cmds.GetAll() {
		cmd.Namespace = "vpc_v2"
	}
	return cmds
}
