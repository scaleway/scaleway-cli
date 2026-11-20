package serverless_sqldb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	serverless_sqldb "github.com/scaleway/scaleway-sdk-go/api/serverless_sqldb/v1alpha1"
)

var (
	sdbSQLDatabaseStatusFormat = humanv2.EnumFormats[serverless_sqldb.DatabaseStatus]{
		serverless_sqldb.DatabaseStatusUnknownStatus: humanv2.EnumFormat{
			Color: color.Faint,
		},
		serverless_sqldb.DatabaseStatusError: humanv2.EnumFormat{
			Color: color.FgRed,
		},
		serverless_sqldb.DatabaseStatusCreating: humanv2.EnumFormat{
			Color: color.FgBlue,
		},
		serverless_sqldb.DatabaseStatusReady: humanv2.EnumFormat{
			Color: color.FgGreen,
		},
		serverless_sqldb.DatabaseStatusDeleting: humanv2.EnumFormat{
			Color: color.FgBlue,
		},
		serverless_sqldb.DatabaseStatusRestoring: humanv2.EnumFormat{
			Color: color.FgBlue,
		},
		serverless_sqldb.DatabaseStatusLocked: humanv2.EnumFormat{
			Color: color.FgRed,
		},
	}

	sdbSQLDatabaseBackupStatusFormat = humanv2.EnumFormats[serverless_sqldb.DatabaseBackupStatus]{
		serverless_sqldb.DatabaseBackupStatusUnknownStatus: humanv2.EnumFormat{
			Color: color.Faint,
		},
		serverless_sqldb.DatabaseBackupStatusError: humanv2.EnumFormat{
			Color: color.FgRed,
		},
		serverless_sqldb.DatabaseBackupStatusReady: humanv2.EnumFormat{
			Color: color.FgGreen,
		},
		serverless_sqldb.DatabaseBackupStatusLocked: humanv2.EnumFormat{
			Color: color.FgRed,
		},
	}
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("sdb-sql").Groups = []string{"database"}

	humanv2.RegisterEnumFormat(sdbSQLDatabaseStatusFormat)
	humanv2.RegisterEnumFormat(sdbSQLDatabaseBackupStatusFormat)

	return cmds
}
