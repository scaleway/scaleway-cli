package searchdb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	searchdb "github.com/scaleway/scaleway-sdk-go/api/searchdb/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("searchdb").Groups = []string{"database"}

	deploymentStatusMarshalSpecs := human.EnumMarshalSpecs{
		searchdb.DeploymentStatusUnknownStatus: &human.EnumMarshalSpec{
			Attribute: color.Faint,
			Value:     "unknown",
		},
		searchdb.DeploymentStatusReady: {
			Attribute: color.FgGreen,
			Value:     "ready",
		},
		searchdb.DeploymentStatusCreating: {
			Attribute: color.FgBlue,
			Value:     "creating",
		},
		searchdb.DeploymentStatusInitializing: {
			Attribute: color.FgBlue,
			Value:     "initializing",
		},
		searchdb.DeploymentStatusUpgrading: {
			Attribute: color.FgBlue,
			Value:     "upgrading",
		},
		searchdb.DeploymentStatusDeleting: {
			Attribute: color.FgBlue,
			Value:     "deleting",
		},
		searchdb.DeploymentStatusError: {
			Attribute: color.FgRed,
			Value:     "error",
		},
		searchdb.DeploymentStatusLocked: {
			Attribute: color.FgRed,
			Value:     "locked",
		},
		searchdb.DeploymentStatusLocking: {
			Attribute: color.FgBlue,
			Value:     "locking",
		},
		searchdb.DeploymentStatusUnlocking: {
			Attribute: color.FgBlue,
			Value:     "unlocking",
		},
	}

	nodeStatusMarshalSpecs := human.EnumMarshalSpecs{
		searchdb.NodeTypeStockStatusUnknownStock: &human.EnumMarshalSpec{
			Attribute: color.Faint,
			Value:     "unknown",
		},
		searchdb.NodeTypeStockStatusLowStock: &human.EnumMarshalSpec{
			Attribute: color.FgYellow,
			Value:     "low_stock",
		},
		searchdb.NodeTypeStockStatusOutOfStock: &human.EnumMarshalSpec{
			Attribute: color.FgRed,
			Value:     "out_of_stock",
		},
		searchdb.NodeTypeStockStatusAvailable: &human.EnumMarshalSpec{
			Attribute: color.FgGreen,
			Value:     "available",
		},
	}

	human.RegisterMarshalerFunc(
		searchdb.DeploymentStatus(""),
		human.EnumMarshalFunc(deploymentStatusMarshalSpecs),
	)

	human.RegisterMarshalerFunc(
		searchdb.NodeTypeStockStatus(""),
		human.EnumMarshalFunc(nodeStatusMarshalSpecs),
	)

	cmds.MustFind("searchdb", "deployment", "list").
		Override(core.ListViewBuilder("ID", "Name", "OrganizationID", "ProjectID", "Status", "Tags", "NodeAmount", "NodeCount", "NodeType", "Endpoints", "CreatedAt", "UpdatedAt", "Version", "Region"))
	cmds.MustFind("searchdb", "node-types", "list").
		Override(core.ListViewBuilder("StockStatus", "Name", "Description", "Vcpus", "MemoryBytes", "Disabled", "Beta", "InstanceRange", "AvailableVolumeTypes", "Region"))
	cmds.MustFind("searchdb", "versions", "list").
		Override(core.ListViewBuilder("Version", "EndOfLife", "Disabled", "Beta", "Region"))

	cmds.MustFind("searchdb", "user", "list").
		Override(core.ListViewBuilder("Username", "DeploymentID", "Region"))

	return cmds
}
