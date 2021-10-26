package function

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(function.NamespaceStatus(""), human.EnumMarshalFunc(namespaceStatusMarshalSpecs))
	human.RegisterMarshalerFunc(function.FunctionStatus(""), human.EnumMarshalFunc(functionStatusMarshalSpecs))
	human.RegisterMarshalerFunc(function.CronStatus(""), human.EnumMarshalFunc(cronStatusMarshalSpecs))

	return cmds
}
