package function

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("function").Groups = []string{"serverless"}

	humanv2.RegisterEnumFormat(namespaceStatusFormat)
	humanv2.RegisterEnumFormat(functionStatusFormat)
	humanv2.RegisterEnumFormat(cronStatusFormat)
	if cmdDeploy := functionDeploy(); cmdDeploy != nil {
		cmds.Add(cmdDeploy)
	}

	cmds.Add(functionCockpitCommand())

	return cmds
}
