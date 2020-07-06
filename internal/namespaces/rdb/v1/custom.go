package rdb

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

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

	human.RegisterMarshalerFunc(rdb.Instance{}, instanceMarshalerFunc)

	cmds.Merge(core.NewCommands(
		instanceWaitCommand(),
	))
	cmds.MustFind("rdb", "instance", "create").Override(instanceCreateBuilder)
	cmds.MustFind("rdb", "instance", "clone").Override(instanceCloneBuilder)
	cmds.MustFind("rdb", "instance", "create").Override(instanceCreateBuilder)

	cmds.MustFind("rdb", "instance", "upgrade").Override(instanceUpgradeBuilder)

	cmds.MustFind("rdb", "engine", "list").Override(engineListBuilder)

	return cmds
}
