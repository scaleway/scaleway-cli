package applesilicon

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(
		core.NewCommands(
			serverWaitCommand(),
		),
	)

	human.RegisterMarshalerFunc(applesilicon.ServerStatus(""), human.EnumMarshalFunc(serverStatusMarshalSpecs))

	cmds.MustFind("apple-silicon", "server", "create").Override(serverCreateBuilder)
	cmds.MustFind("apple-silicon", "server", "reboot").Override(serverRebootBuilder)
	cmds.MustFind("apple-silicon", "server", "delete").Override(serverDeleteBuilder)

	return cmds
}
