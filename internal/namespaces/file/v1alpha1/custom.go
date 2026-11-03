package file

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("file").Groups = []string{"storage"}

	cmds.Add(filesystemChargeCommand())

	return cmds
}
