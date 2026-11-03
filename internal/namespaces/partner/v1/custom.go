package partner

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("partner").Groups = []string{"security"}

	cmds.MustFind("partner", "organization", "list").Override(core.ListViewBuilder("ID", "Name", "Email", "Status", "OwnerFirstname", "OwnerLastname", "CreatedAt", "PhoneNumber", "CustomerID", "LockReasonMessage", "LockedBy", "LockedAt", "PictureLink", "Comment"))

	return cmds
}
