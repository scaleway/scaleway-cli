package block

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	block "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
)

var (
	volumeStatusFormat = humanv2.EnumFormats[block.VolumeStatus]{
		block.VolumeStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},
		block.VolumeStatusCreating:      humanv2.EnumFormat{Color: color.FgBlue},
		block.VolumeStatusAvailable:     humanv2.EnumFormat{Color: color.FgGreen},
		block.VolumeStatusInUse:         humanv2.EnumFormat{Color: color.FgHiGreen},
		block.VolumeStatusDeleting:      humanv2.EnumFormat{Color: color.FgBlue},
		block.VolumeStatusResizing:      humanv2.EnumFormat{Color: color.FgBlue},
		block.VolumeStatusError:         humanv2.EnumFormat{Color: color.FgRed},
		block.VolumeStatusSnapshotting:  humanv2.EnumFormat{Color: color.FgBlue},
		block.VolumeStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
	}

	snapshotStatusFormat = humanv2.EnumFormats[block.SnapshotStatus]{
		block.SnapshotStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},
		block.SnapshotStatusCreating:      humanv2.EnumFormat{Color: color.FgBlue},
		block.SnapshotStatusAvailable:     humanv2.EnumFormat{Color: color.FgGreen},
		block.SnapshotStatusError:         humanv2.EnumFormat{Color: color.FgRed},
		block.SnapshotStatusDeleting:      humanv2.EnumFormat{Color: color.FgBlue},
		block.SnapshotStatusDeleted:       humanv2.EnumFormat{Color: color.Faint},
		block.SnapshotStatusInUse:         humanv2.EnumFormat{Color: color.FgHiGreen},
		block.SnapshotStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
	}
	referenceStatusFormat = humanv2.EnumFormats[block.ReferenceStatus]{
		block.ReferenceStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},

		block.ReferenceStatusAttaching: humanv2.EnumFormat{Color: color.FgBlue},
		block.ReferenceStatusAttached:  humanv2.EnumFormat{Color: color.FgHiGreen},
		block.ReferenceStatusDetaching: humanv2.EnumFormat{Color: color.FgBlue},
		block.ReferenceStatusDetached:  humanv2.EnumFormat{Color: color.FgGreen},
		block.ReferenceStatusCreating:  humanv2.EnumFormat{Color: color.FgBlue},
		block.ReferenceStatusError:     humanv2.EnumFormat{Color: color.FgRed},
	}
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("block").Groups = []string{"storage"}

	cmds.Add(volumeWaitCommand())
	cmds.Add(snapshotWaitCommand())

	cmds.MustFind("block", "snapshot", "create").Override(blockSnapshotCreateBuilder)
	cmds.MustFind("block", "volume", "create").Override(blockVolumeCreateBuilder)

	humanv2.RegisterEnumFormat(volumeStatusFormat)
	humanv2.RegisterEnumFormat(snapshotStatusFormat)
	humanv2.RegisterEnumFormat(referenceStatusFormat)

	return cmds
}
