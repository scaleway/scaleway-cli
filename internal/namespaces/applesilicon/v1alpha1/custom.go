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

	cmds.MustFind("apple-silicon", "os", "list").Override(core.ListViewBuilder("ID", "Name", "Label", "ImageURL", "Family", "IsBeta", "Version", "XcodeVersion", "CompatibleServerTypes", "ReleaseNotesURL", "Description", "Tags", "SupportedServerTypes", "Zone"))
	cmds.MustFind("apple-silicon", "private-network", "list").Override(core.ListViewBuilder("ID", "ProjectID", "ServerID", "PrivateNetworkID", "Vlan", "Status", "CreatedAt", "UpdatedAt", "IpamIPIDs"))
	cmds.MustFind("apple-silicon", "server", "list").Override(core.ListViewBuilder("ID", "Type", "Name", "ProjectID", "OrganizationID", "IP", "VncURL", "SSHUsername", "SudoPassword", "VncPort", "Os", "FullOsVersion", "Status", "CreatedAt", "UpdatedAt", "DeletableAt", "DeletionScheduled", "Zone", "Delivered", "VpcStatus", "PublicBandwidthBps", "Tags", "AppliedRunnerConfigurationIDs", "KextEnabled"))

	cmds.MustFind("apple-silicon", "runner", "list").Override(core.ListViewBuilder("ID", "Status", "ErrorMessage", "Zone"))

	return cmds
}
