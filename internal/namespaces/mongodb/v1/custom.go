package mongodb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("mongodb").Groups = []string{"database"}

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[mongodb.InstanceStatus]{
		mongodb.InstanceStatusConfiguring:  humanv2.EnumFormat{Color: color.FgBlue, Value: "configuring"},
		mongodb.InstanceStatusDeleting:     humanv2.EnumFormat{Color: color.FgBlue, Value: "deleting"},
		mongodb.InstanceStatusError:        humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
		mongodb.InstanceStatusInitializing: humanv2.EnumFormat{Color: color.FgBlue, Value: "initializing"},
		mongodb.InstanceStatusLocked:       humanv2.EnumFormat{Color: color.FgRed, Value: "locked"},
		mongodb.InstanceStatusProvisioning: humanv2.EnumFormat{Color: color.FgBlue, Value: "provisioning"},
		mongodb.InstanceStatusReady:        humanv2.EnumFormat{Color: color.FgGreen, Value: "ready"},
		mongodb.InstanceStatusSnapshotting: humanv2.EnumFormat{Color: color.FgBlue, Value: "snapshotting"},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[mongodb.NodeTypeStock]{
		mongodb.NodeTypeStockAvailable:  humanv2.EnumFormat{Color: color.FgGreen, Value: "available"},
		mongodb.NodeTypeStockLowStock:   humanv2.EnumFormat{Color: color.FgYellow, Value: "low stock"},
		mongodb.NodeTypeStockOutOfStock: humanv2.EnumFormat{Color: color.FgRed, Value: "out of stock"},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[mongodb.SnapshotStatus]{
		mongodb.SnapshotStatusCreating:  humanv2.EnumFormat{Color: color.FgBlue, Value: "creating"},
		mongodb.SnapshotStatusDeleting:  humanv2.EnumFormat{Color: color.FgBlue, Value: "deleting"},
		mongodb.SnapshotStatusError:     humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
		mongodb.SnapshotStatusLocked:    humanv2.EnumFormat{Color: color.FgRed, Value: "locked"},
		mongodb.SnapshotStatusReady:     humanv2.EnumFormat{Color: color.FgGreen, Value: "ready"},
		mongodb.SnapshotStatusRestoring: humanv2.EnumFormat{Color: color.FgBlue, Value: "restoring"},
	})

	cmds.MustFind("mongodb", "instance", "create").Override(instanceCreateBuilder)

	cmds.Merge(core.NewCommands(
		instanceWaitCommand(),
	))

	cmds.Merge(core.NewCommands(instanceCockpitCommand()))

	return cmds
}
