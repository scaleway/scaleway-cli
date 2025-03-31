package serverless_sqldb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	serverless_sqldb "github.com/scaleway/scaleway-sdk-go/api/serverless_sqldb/v1alpha1"
)

var (
	sdbSQLDatabaseStatusMarshalSpecs = human.EnumMarshalSpecs{
		serverless_sqldb.DatabaseStatusUnknownStatus: &human.EnumMarshalSpec{
			Attribute: color.Faint,
		},
		serverless_sqldb.DatabaseStatusError: &human.EnumMarshalSpec{
			Attribute: color.FgRed,
		},
		serverless_sqldb.DatabaseStatusCreating: &human.EnumMarshalSpec{
			Attribute: color.FgBlue,
		},
		serverless_sqldb.DatabaseStatusReady: &human.EnumMarshalSpec{
			Attribute: color.FgGreen,
		},
		serverless_sqldb.DatabaseStatusDeleting: &human.EnumMarshalSpec{
			Attribute: color.FgBlue,
		},
		serverless_sqldb.DatabaseStatusRestoring: &human.EnumMarshalSpec{
			Attribute: color.FgBlue,
		},
		serverless_sqldb.DatabaseStatusLocked: &human.EnumMarshalSpec{
			Attribute: color.FgRed,
		},
	}

	sdbSQLDatabaseBackupStatusMarshalSpecs = human.EnumMarshalSpecs{
		serverless_sqldb.DatabaseBackupStatusUnknownStatus: &human.EnumMarshalSpec{
			Attribute: color.Faint,
		},
		serverless_sqldb.DatabaseBackupStatusError: &human.EnumMarshalSpec{
			Attribute: color.FgRed,
		},
		serverless_sqldb.DatabaseBackupStatusReady: &human.EnumMarshalSpec{
			Attribute: color.FgGreen,
		},
		serverless_sqldb.DatabaseBackupStatusLocked: &human.EnumMarshalSpec{
			Attribute: color.FgRed,
		},
	}
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("sdb-sql").Groups = []string{"database"}

	human.RegisterMarshalerFunc(
		serverless_sqldb.DatabaseStatus(""),
		human.EnumMarshalFunc(sdbSQLDatabaseStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		serverless_sqldb.DatabaseBackupStatus(""),
		human.EnumMarshalFunc(sdbSQLDatabaseBackupStatusMarshalSpecs),
	)

	return cmds
}
