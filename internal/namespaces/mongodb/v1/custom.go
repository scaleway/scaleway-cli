package mongodb

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("mongodb").Groups = []string{"database"}

	human.RegisterMarshalerFunc(
		mongodb.SnapshotStatus(""),
		human.EnumMarshalFunc(snapshotStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		mongodb.InstanceStatus(""),
		human.EnumMarshalFunc(instanceStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		mongodb.NodeTypeStock(""),
		human.EnumMarshalFunc(nodeTypeStockMarshalSpecs),
	)

	cmds.MustFind("mongodb", "instance", "create").Override(instanceCreateBuilder)

	cmds.Merge(core.NewCommands(
		instanceWaitCommand(),
	))

	cmds.MustFind("mongodb", "instance", "list").
		Override(core.ListViewBuilder("ID", "Name", "ProjectID", "OrganizationID", "Status", "Version", "Tags", "NodeAmount", "NodeType", "Endpoints", "CreatedAt", "Region", "Settings", "Maintenances", "UpgradableVersions"))
	cmds.MustFind("mongodb", "maintenance", "list").
		Override(core.ListViewBuilder("ID", "InstanceID", "CreatedAt", "StartsAt", "StopsAt", "Status", "ForcedAt", "AppliedAt", "AppliedBy", "Reason", "Region"))
	cmds.MustFind("mongodb", "node-type", "list").
		Override(core.ListViewBuilder("Name", "StockStatus", "Description", "Vcpus", "MemoryBytes", "AvailableVolumeTypes", "Disabled", "Beta", "InstanceRange"))
	cmds.MustFind("mongodb", "snapshot", "list").
		Override(core.ListViewBuilder("ID", "InstanceID", "Name", "Status", "SizeBytes", "ExpiresAt", "CreatedAt", "UpdatedAt", "InstanceName", "NodeType", "VolumeType", "Region"))
	cmds.MustFind("mongodb", "version", "list").
		Override(core.ListViewBuilder("Version", "EndOfLifeAt", "ReleasedAt"))

	cmds.MustFind("mongodb", "database", "list").Override(core.ListViewBuilder("Name"))
	cmds.MustFind("mongodb", "user", "list").Override(core.ListViewBuilder("Name", "Roles"))

	return cmds
}
