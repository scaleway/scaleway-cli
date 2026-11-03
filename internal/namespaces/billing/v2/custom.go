package billing

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("billing", "budget", "list").Override(core.ListViewBuilder("ID", "CreatedAt", "UpdatedAt", "OrganizationID", "Enabled", "ConsumptionLimit", "Alerts"))

	return cmds
}
