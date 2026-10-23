package messageq

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(core.NewCommands(deploymentCockpitCommand()))

	return cmds
}
