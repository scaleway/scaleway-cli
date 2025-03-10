package audit_trail

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("audit-trail", "event", "list").Override(eventListBuilder)

	return cmds
}
