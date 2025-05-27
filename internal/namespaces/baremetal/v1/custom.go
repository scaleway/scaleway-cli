package baremetal

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	baremetalV3 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v3"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("baremetal").Groups = []string{"baremetal"}

	cmds.Merge(core.NewCommands(
		serverWaitCommand(),
		serverAddFlexibleIP(),
	))

	human.RegisterMarshalerFunc(
		baremetal.ServerPingStatus(""),
		human.EnumMarshalFunc(serverPingStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		baremetal.OfferStock(""),
		human.EnumMarshalFunc(offerAvailabilityMarshalSpecs),
	)
	human.RegisterMarshalerFunc(baremetal.Server{}, serverMarshalerFunc)

	human.RegisterMarshalerFunc(baremetal.Offer{}, listOfferMarshalerFunc)

	cmds.MustFind("baremetal", "server", "create").Override(serverCreateBuilder)
	cmds.MustFind("baremetal", "server", "install").Override(serverInstallBuilder)
	cmds.MustFind("baremetal", "server", "delete").Override(serverDeleteBuilder)
	cmds.MustFind("baremetal", "server", "list").Override(serverListBuilder)

	// Action commands
	cmds.MustFind("baremetal", "server", "start").Override(serverStartBuilder)
	cmds.MustFind("baremetal", "server", "stop").Override(serverStopBuilder)
	cmds.MustFind("baremetal", "server", "reboot").Override(serverRebootBuilder)

	cmds.MergeAll(baremetalV3.GetCommands())

	return cmds
}
