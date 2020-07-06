package rdb

import "github.com/scaleway/scaleway-cli/internal/core"

var nodeTypes = []string{
	"db-dev-s",
	"db-dev-m",
	"db-dev-l",
	"db-dev-xl",
	"db-gp-xs",
	"db-gp-s",
	"db-gp-m",
	"db-gp-l",
	"db-gp-xl",
}

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(core.NewCommands(
		instanceWaitCommand(),
	))
	cmds.MustFind("rdb", "instance", "create").Override(instanceCreateBuilder)
	cmds.MustFind("rdb", "instance", "clone").Override(instanceCloneBuilder)
	cmds.MustFind("rdb", "instance", "create").Override(instanceCreateBuilder)

	cmds.MustFind("rdb", "instance", "upgrade").Override(instanceUpgradeBuilder)

	return cmds
}
