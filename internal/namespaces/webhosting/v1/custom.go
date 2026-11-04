package webhosting

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("webhosting").Groups = []string{"Domain & WebHosting"}

	humanv2.RegisterEnumFormat(hostingStatusFormat)
	humanv2.RegisterEnumFormat(hostingDNSFormat)
	humanv2.RegisterEnumFormat(nameserverFormat)
	cmds.MustFind("webhosting", "offer", "list").Override(webhostingOfferListBuilder)

	return cmds
}
