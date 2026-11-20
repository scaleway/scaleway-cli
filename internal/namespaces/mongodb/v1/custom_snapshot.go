package mongodb

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
)

var snapshotStatusFormat = humanv2.EnumFormats[mongodb.SnapshotStatus]{
	mongodb.SnapshotStatusCreating: humanv2.EnumFormat{
		Color: color.FgBlue,
		Value: "creating",
	},
	mongodb.SnapshotStatusDeleting: humanv2.EnumFormat{
		Color: color.FgBlue,
		Value: "deleting",
	},
	mongodb.SnapshotStatusError: humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
	mongodb.SnapshotStatusLocked: humanv2.EnumFormat{
		Color: color.FgRed,
		Value: "locked",
	},
	mongodb.SnapshotStatusReady: humanv2.EnumFormat{
		Color: color.FgGreen,
		Value: "ready",
	},
	mongodb.SnapshotStatusRestoring: humanv2.EnumFormat{
		Color: color.FgBlue,
		Value: "restoring",
	},
}
