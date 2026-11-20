package searchdb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	searchdb "github.com/scaleway/scaleway-sdk-go/api/searchdb/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("searchdb").Groups = []string{"database"}

	deploymentStatusFormat := humanv2.EnumFormats[searchdb.DeploymentStatus]{
		searchdb.DeploymentStatusUnknownStatus: humanv2.EnumFormat{
			Color: color.Faint,
			Value: "unknown",
		},
		searchdb.DeploymentStatusReady: {
			Color: color.FgGreen,
			Value: "ready",
		},
		searchdb.DeploymentStatusCreating: {
			Color: color.FgBlue,
			Value: "creating",
		},
		searchdb.DeploymentStatusInitializing: {
			Color: color.FgBlue,
			Value: "initializing",
		},
		searchdb.DeploymentStatusUpgrading: {
			Color: color.FgBlue,
			Value: "upgrading",
		},
		searchdb.DeploymentStatusDeleting: {
			Color: color.FgBlue,
			Value: "deleting",
		},
		searchdb.DeploymentStatusError: {
			Color: color.FgRed,
			Value: "error",
		},
		searchdb.DeploymentStatusLocked: {
			Color: color.FgRed,
			Value: "locked",
		},
		searchdb.DeploymentStatusLocking: {
			Color: color.FgBlue,
			Value: "locking",
		},
		searchdb.DeploymentStatusUnlocking: {
			Color: color.FgBlue,
			Value: "unlocking",
		},
	}

	nodeStatusFormat := humanv2.EnumFormats[searchdb.NodeTypeStockStatus]{
		searchdb.NodeTypeStockStatusUnknownStock: humanv2.EnumFormat{
			Color: color.Faint,
			Value: "unknown",
		},
		searchdb.NodeTypeStockStatusLowStock: humanv2.EnumFormat{
			Color: color.FgYellow,
			Value: "low_stock",
		},
		searchdb.NodeTypeStockStatusOutOfStock: humanv2.EnumFormat{
			Color: color.FgRed,
			Value: "out_of_stock",
		},
		searchdb.NodeTypeStockStatusAvailable: humanv2.EnumFormat{
			Color: color.FgGreen,
			Value: "available",
		},
	}

	humanv2.RegisterEnumFormat(deploymentStatusFormat)
	humanv2.RegisterEnumFormat(nodeStatusFormat)

	return cmds
}
