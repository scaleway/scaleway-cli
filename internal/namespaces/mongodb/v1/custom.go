package mongodb

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("mongodb").Groups = []string{"database"}

	humanv2.RegisterEnumFormat(snapshotStatusFormat)
	humanv2.RegisterEnumFormat(instanceStatusFormat)
	humanv2.RegisterEnumFormat(nodeTypeStockFormat)
	cmds.MustFind("mongodb", "instance", "create").Override(instanceCreateBuilder)

	cmds.Merge(core.NewCommands(
		instanceWaitCommand(),
	))

	cmds.Merge(core.NewCommands(instanceCockpitCommand()))

	return cmds
}
