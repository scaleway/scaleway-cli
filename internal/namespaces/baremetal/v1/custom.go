package baremetal

import (
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
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

	humanv2.RegisterLayout(reflect.TypeFor[baremetal.Server](), serverLayout)
	humanv2.RegisterLayout(reflect.TypeFor[baremetal.Offer](), offerLayout)
	humanv2.RegisterEnumFormat(offerStockFormat)
	humanv2.RegisterEnumFormat(pingStatusFormat)

	cmds.MustFind("baremetal", "server", "create").Override(serverCreateBuilder)
	cmds.MustFind("baremetal", "server", "install").Override(serverInstallBuilder)
	cmds.MustFind("baremetal", "server", "delete").Override(serverDeleteBuilder)
	cmds.MustFind("baremetal", "server", "list").Override(serverListBuilder)

	// Action commands
	cmds.MustFind("baremetal", "server", "start").Override(serverStartBuilder)
	cmds.MustFind("baremetal", "server", "stop").Override(serverStopBuilder)
	cmds.MustFind("baremetal", "server", "reboot").Override(serverRebootBuilder)
	cmds.MustFind("baremetal", "offer", "list").Override(serverOfferListBuilder)

	cmds.MergeAll(baremetalV3.GetCommands())

	return cmds
}
