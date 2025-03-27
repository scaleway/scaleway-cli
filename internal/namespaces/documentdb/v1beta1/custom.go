package documentdb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	documentdb "github.com/scaleway/scaleway-sdk-go/api/documentdb/v1beta1"
)

var instanceStatusMarshalSpecs = human.EnumMarshalSpecs{
	documentdb.InstanceStatusUnknown: &human.EnumMarshalSpec{
		Attribute: color.Faint,
		Value:     "unknown",
	},
	documentdb.InstanceStatusReady: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
		Value:     "ready",
	},
	documentdb.InstanceStatusProvisioning: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "provisioning",
	},
	documentdb.InstanceStatusConfiguring: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "configuring",
	},
	documentdb.InstanceStatusDeleting: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "deleting",
	},
	documentdb.InstanceStatusError: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "error",
	},
	documentdb.InstanceStatusAutohealing: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "auto-healing",
	},
	documentdb.InstanceStatusLocked: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "locked",
	},
	documentdb.InstanceStatusInitializing: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "initialized",
	},
	documentdb.InstanceStatusDiskFull: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "disk_full",
	},
	documentdb.InstanceStatusBackuping: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "backuping",
	},
	documentdb.InstanceStatusSnapshotting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
	documentdb.InstanceStatusRestarting:   &human.EnumMarshalSpec{Attribute: color.FgBlue},
}

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("document-db").Groups = []string{"database"}

	human.RegisterMarshalerFunc(
		documentdb.InstanceStatus(""),
		human.EnumMarshalFunc(instanceStatusMarshalSpecs),
	)

	cmds.MustFind("document-db", "engine", "list").Override(engineListBuilder)

	return cmds
}
