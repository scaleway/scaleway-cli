package interlink

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("interlink").Groups = []string{"network"}

	humanv2.RegisterEnumFormat(bgpStatusFormat)
	humanv2.RegisterEnumFormat(linkStatusFormat)

	return cmds
}
