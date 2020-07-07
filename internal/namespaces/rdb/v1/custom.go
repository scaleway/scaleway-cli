package rdb

import "github.com/scaleway/scaleway-cli/internal/core"

var nodeTypes = []string{
	"DB-DEV-S",
	"DB-DEV-M",
	"DB-DEV-L",
	"DB-DEV-XL",
	"DB-GP-XS",
	"DB-GP-S",
	"DB-GP-M",
	"DB-GP-L",
	"DB-GP-XL",
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
