package applesilicon

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(core.NewCommands(
		serverSSHCommand(),
		serverWaitCommand(),
	))

	cmds.MustFind("apple-silicon").Groups = []string{"baremetal"}

	human.RegisterMarshalerFunc(applesilicon.ServerTypeCPU{}, cpuMarshalerFunc)
	human.RegisterMarshalerFunc(applesilicon.ServerTypeDisk{}, diskMarshalerFunc)
	human.RegisterMarshalerFunc(applesilicon.ServerTypeMemory{}, memoryMarshalerFunc)
	human.RegisterMarshalerFunc(applesilicon.OS{}, OSMarshalerFunc)

	human.RegisterMarshalerFunc(
		applesilicon.ServerStatus(""),
		human.EnumMarshalFunc(serverStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		applesilicon.ServerTypeStock(""),
		human.EnumMarshalFunc(serverTypeStockMarshalSpecs),
	)

	cmds.MustFind("apple-silicon", "server", "create").Override(serverCreateBuilder)
	cmds.MustFind("apple-silicon", "server", "reboot").Override(serverRebootBuilder)
	cmds.MustFind("apple-silicon", "server", "delete").Override(serverDeleteBuilder)

	cmds.MustFind("apple-silicon", "server-type", "list").Override(serverTypeBuilder)

	return cmds
}
