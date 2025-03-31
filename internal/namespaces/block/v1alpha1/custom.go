package block

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	block "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
)

var (
	volumeStatusMarshalSpecs = human.EnumMarshalSpecs{
		block.VolumeStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
		block.VolumeStatusCreating:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
		block.VolumeStatusAvailable:     &human.EnumMarshalSpec{Attribute: color.FgGreen},
		block.VolumeStatusInUse:         &human.EnumMarshalSpec{Attribute: color.FgHiGreen},
		block.VolumeStatusDeleting:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
		block.VolumeStatusResizing:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
		block.VolumeStatusError:         &human.EnumMarshalSpec{Attribute: color.FgRed},
		block.VolumeStatusSnapshotting:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		block.VolumeStatusLocked:        &human.EnumMarshalSpec{Attribute: color.FgRed},
	}

	snapshotStatusMarshalSpecs = human.EnumMarshalSpecs{
		block.SnapshotStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
		block.SnapshotStatusCreating:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
		block.SnapshotStatusAvailable:     &human.EnumMarshalSpec{Attribute: color.FgGreen},
		block.SnapshotStatusError:         &human.EnumMarshalSpec{Attribute: color.FgRed},
		block.SnapshotStatusDeleting:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
		block.SnapshotStatusDeleted:       &human.EnumMarshalSpec{Attribute: color.Faint},
		block.SnapshotStatusInUse:         &human.EnumMarshalSpec{Attribute: color.FgHiGreen},
		block.SnapshotStatusLocked:        &human.EnumMarshalSpec{Attribute: color.FgRed},
	}
	referenceStatusMarshalSpecs = human.EnumMarshalSpecs{
		block.ReferenceStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},

		block.ReferenceStatusAttaching: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		block.ReferenceStatusAttached:  &human.EnumMarshalSpec{Attribute: color.FgHiGreen},
		block.ReferenceStatusDetaching: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		block.ReferenceStatusDetached:  &human.EnumMarshalSpec{Attribute: color.FgGreen},
		block.ReferenceStatusCreating:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		block.ReferenceStatusError:     &human.EnumMarshalSpec{Attribute: color.FgRed},
	}
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("block").Groups = []string{"storage"}

	cmds.Add(volumeWaitCommand())
	cmds.Add(snapshotWaitCommand())

	cmds.MustFind("block", "snapshot", "create").Override(blockSnapshotCreateBuilder)
	cmds.MustFind("block", "volume", "create").Override(blockVolumeCreateBuilder)

	human.RegisterMarshalerFunc(
		block.VolumeStatus(""),
		human.EnumMarshalFunc(volumeStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		block.SnapshotStatus(""),
		human.EnumMarshalFunc(snapshotStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		block.ReferenceStatus(""),
		human.EnumMarshalFunc(referenceStatusMarshalSpecs),
	)

	return cmds
}
