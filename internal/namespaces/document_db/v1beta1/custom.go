package document_db

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	document_db "github.com/scaleway/scaleway-sdk-go/api/document_db/v1beta1"
)

var (
	instanceStatusMarshalSpecs = human.EnumMarshalSpecs{
		document_db.InstanceStatusUnknown:      &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
		document_db.InstanceStatusReady:        &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
		document_db.InstanceStatusProvisioning: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "provisioning"},
		document_db.InstanceStatusConfiguring:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "configuring"},
		document_db.InstanceStatusDeleting:     &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "deleting"},
		document_db.InstanceStatusError:        &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
		document_db.InstanceStatusAutohealing:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "auto-healing"},
		document_db.InstanceStatusLocked:       &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "locked"},
		document_db.InstanceStatusInitializing: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "initialized"},
		document_db.InstanceStatusDiskFull:     &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "disk_full"},
		document_db.InstanceStatusBackuping:    &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "backuping"},
		document_db.InstanceStatusSnapshotting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		document_db.InstanceStatusRestarting:   &human.EnumMarshalSpec{Attribute: color.FgBlue},
	}
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(document_db.InstanceStatus(""), human.EnumMarshalFunc(instanceStatusMarshalSpecs))

	cmds.MustFind("document-db", "engine", "list").Override(engineListBuilder)

	return cmds
}
