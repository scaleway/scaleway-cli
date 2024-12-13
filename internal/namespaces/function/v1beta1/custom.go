package function

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/human"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(function.NamespaceStatus(""), human.EnumMarshalFunc(namespaceStatusMarshalSpecs))
	human.RegisterMarshalerFunc(function.FunctionStatus(""), human.EnumMarshalFunc(functionStatusMarshalSpecs))
	human.RegisterMarshalerFunc(function.CronStatus(""), human.EnumMarshalFunc(cronStatusMarshalSpecs))

	if cmdDeploy := functionDeploy(); cmdDeploy != nil {
		cmds.Add(cmdDeploy)
	}

	return cmds
}
