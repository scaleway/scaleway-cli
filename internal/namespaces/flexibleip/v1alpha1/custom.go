package flexibleip

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("fip").Groups = []string{"baremetal"}

	humanv2.RegisterEnumFormat(ipStatusFormat)
	humanv2.RegisterEnumFormat(macAddressStatusFormat)
	cmds.MustFind("fip", "ip", "create").Override(createIPBuilder)

	return cmds
}
