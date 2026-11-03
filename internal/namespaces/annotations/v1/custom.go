package annotations

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("annotations", "binding", "list").Override(core.ListViewBuilder("ID", "TargetSrn"))
	cmds.MustFind("annotations", "key", "list").Override(core.ListViewBuilder("ID", "Name", "Description"))
	cmds.MustFind("annotations", "value", "list").Override(core.ListViewBuilder("ID", "KeyID", "Name", "Description"))

	return cmds
}
