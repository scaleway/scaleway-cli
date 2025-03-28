package account

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	commands := GetGeneratedCommands()

	commands.MustFind("account").Groups = []string{"security"}

	return commands
}
