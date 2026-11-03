package function

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("function").Groups = []string{"serverless"}

	human.RegisterMarshalerFunc(
		function.NamespaceStatus(""),
		human.EnumMarshalFunc(namespaceStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		function.FunctionStatus(""),
		human.EnumMarshalFunc(functionStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		function.CronStatus(""),
		human.EnumMarshalFunc(cronStatusMarshalSpecs),
	)

	if cmdDeploy := functionDeploy(); cmdDeploy != nil {
		cmds.Add(cmdDeploy)
	}

	cmds.MustFind("function", "cron", "list").
		Override(core.ListViewBuilder("ID", "FunctionID", "Schedule", "Args", "Status", "Name", "Region"))
	cmds.MustFind("function", "domain", "list").
		Override(core.ListViewBuilder("ID", "Hostname", "FunctionID", "URL", "Status", "ErrorMessage", "Region"))
	cmds.MustFind("function", "function", "list").
		Override(core.ListViewBuilder("ID", "Name", "NamespaceID", "Status", "MinScale", "MaxScale", "Runtime", "MemoryLimit", "CPULimit", "Timeout", "Handler", "ErrorMessage", "BuildMessage", "Privacy", "Description", "DomainName", "SecretEnvironmentVariables", "Region", "HTTPOption", "RuntimeMessage", "Sandbox", "CreatedAt", "UpdatedAt", "ReadyAt", "Tags", "PrivateNetworkID"))
	cmds.MustFind("function", "token", "list").
		Override(core.ListViewBuilder("ID", "Token", "FunctionID", "NamespaceID", "PublicKey", "Status", "Description", "ExpiresAt"))
	cmds.MustFind("function", "trigger", "list").
		Override(core.ListViewBuilder("ID", "Name", "Description", "FunctionID", "InputType", "Status", "ErrorMessage", "Region"))

	return cmds
}
