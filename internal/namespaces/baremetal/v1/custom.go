package baremetal

import (
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	baremetalV3 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v3"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("baremetal").Groups = []string{"baremetal"}

	cmds.Merge(core.NewCommands(
		serverWaitCommand(),
		serverAddFlexibleIP(),
	))

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[baremetal.OfferStock]{
		baremetal.OfferStockEmpty:     humanv2.EnumFormat{Color: color.FgRed, Value: "empty"},
		baremetal.OfferStockLow:       humanv2.EnumFormat{Color: color.FgYellow, Value: "low"},
		baremetal.OfferStockAvailable: humanv2.EnumFormat{Color: color.FgGreen, Value: "available"},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[baremetal.ServerPingStatus]{
		baremetal.ServerPingStatusPingStatusDown:    humanv2.EnumFormat{Color: color.FgRed, Value: "down"},
		baremetal.ServerPingStatusPingStatusUp:      humanv2.EnumFormat{Color: color.FgGreen, Value: "up"},
		baremetal.ServerPingStatusPingStatusUnknown: humanv2.EnumFormat{Color: color.Faint, Value: "unknown"},
	})

	humanv2.RegisterLayout(reflect.TypeFor[baremetal.Offer](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "Disks", Title: "Disks"},
			{Path: "CPUs", Title: "CPUs"},
			{Path: "Memories", Title: "Memories"},
			{Path: "Options", Title: "Options", Fields: []humanv2.Field{
				{Path: "ID"},
				{Path: "Name"},
				{Path: "Enabled"},
				{Path: "SubscriptionPeriod"},
				{Path: "Price"},
				{Path: "Manageable"},
				{Path: "OsID"},
				{Path: "Type"},
			}},
			{Path: "Bandwidth", Title: "Bandwidth(bit/s)"},
			{Path: "PrivateBandwidth", Title: "PrivateBandwidth(bit/s)"},
		},
	})

	humanv2.RegisterLayout(reflect.TypeFor[baremetal.Server](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "IPs", Title: "IPs"},
			{Path: "Options", Title: "Options", Fields: []humanv2.Field{
				{Path: "ID"},
				{Path: "Name"},
				{Path: "Status"},
				{Path: "Manageable"},
				{Path: "ExpiresAt"},
				{Path: "Type"},
			}},
			{Path: "Install", Title: "Install", Omittable: true},
		},
	})

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
