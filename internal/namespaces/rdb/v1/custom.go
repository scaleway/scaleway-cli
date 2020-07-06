package rdb

import "github.com/scaleway/scaleway-cli/internal/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(core.NewCommands(
		instanceWaitCommand(),
	))
	cmds.MustFind("rdb", "instance", "create").Override(instanceCreateBuilder)
	cmds.MustFind("rdb", "instance", "clone").Override(instanceCloneBuilder)
	cmds.MustFind("rdb", "instance", "upgrade").Override(instanceUpgradeBuilder)

	cmds.MustFind("rdb", "engine", "list").Override(engineListBuilder)

	return cmds
}
