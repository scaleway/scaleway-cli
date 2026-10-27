package mongodb

import (
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/color"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
)

var snapshotStatusMarshalSpecs = human.EnumMarshalSpecs{
	mongodb.SnapshotStatusCreating: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "creating",
	},
	mongodb.SnapshotStatusDeleting: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "deleting",
	},
	mongodb.SnapshotStatusError: &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
	mongodb.SnapshotStatusLocked: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "locked",
	},
	mongodb.SnapshotStatusReady: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
		Value:     "ready",
	},
	mongodb.SnapshotStatusRestoring: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "restoring",
	},
}
