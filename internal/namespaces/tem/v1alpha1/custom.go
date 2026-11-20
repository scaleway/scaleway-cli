package tem

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("tem").Groups = []string{"Domain & WebHosting"}

	humanv2.RegisterEnumFormat(domainStatusFormat)
	humanv2.RegisterEnumFormat(emailStatusFormat)
	cmds.MustFind("tem", "domain", "get").Override(domainGetBuilder)

	return cmds
}
