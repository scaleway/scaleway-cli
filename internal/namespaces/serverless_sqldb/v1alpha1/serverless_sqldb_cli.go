// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package serverless_sqldb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	serverless_sqldb "github.com/scaleway/scaleway-sdk-go/api/serverless_sqldb/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		sdbSQLRoot(),
		sdbSQLDatabase(),
		sdbSQLBackup(),
		sdbSQLDatabaseCreate(),
		sdbSQLDatabaseGet(),
		sdbSQLDatabaseDelete(),
		sdbSQLDatabaseList(),
		sdbSQLDatabaseUpdate(),
		sdbSQLDatabaseRestore(),
		sdbSQLBackupGet(),
		sdbSQLBackupList(),
		sdbSQLBackupExport(),
	)
}

func sdbSQLRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your Serverless SQL Databases`,
		Long:      `This API allows you to manage your Serverless SQL Databases.`,
		Namespace: "sdb-sql",
	}
}

func sdbSQLDatabase() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "sdb-sql",
		Resource:  "database",
	}
}

func sdbSQLBackup() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "sdb-sql",
		Resource:  "backup",
	}
}

func sdbSQLDatabaseCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new Serverless SQL Database`,
		Long:      `You must provide the following parameters: ` + "`" + `organization_id` + "`" + `, ` + "`" + `project_id` + "`" + `, ` + "`" + `name` + "`" + `, ` + "`" + `cpu_min` + "`" + `, ` + "`" + `cpu_max` + "`" + `. You can also provide ` + "`" + `from_backup_id` + "`" + ` to create a database from a backup.`,
		Namespace: "sdb-sql",
		Resource:  "database",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(serverless_sqldb.CreateDatabaseRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `The name of the Serverless SQL Database to be created.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cpu-min",
				Short:      `The minimum number of CPU units for your Serverless SQL Database.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cpu-max",
				Short:      `The maximum number of CPU units for your Serverless SQL Database.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "from-backup-id",
				Short:      `The ID of the backup to create the database from.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*serverless_sqldb.CreateDatabaseRequest)

			client := core.ExtractClient(ctx)
			api := serverless_sqldb.NewAPI(client)

			return api.CreateDatabase(request)
		},
	}
}

func sdbSQLDatabaseGet() *core.Command {
	return &core.Command{
		Short:     `Get a database information`,
		Long:      `Retrieve information about your Serverless SQL Database. You must provide the ` + "`" + `database_id` + "`" + ` parameter.`,
		Namespace: "sdb-sql",
		Resource:  "database",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(serverless_sqldb.GetDatabaseRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "database-id",
				Short:      `UUID of the Serverless SQL DB database.`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*serverless_sqldb.GetDatabaseRequest)

			client := core.ExtractClient(ctx)
			api := serverless_sqldb.NewAPI(client)

			return api.GetDatabase(request)
		},
	}
}

func sdbSQLDatabaseDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a database`,
		Long:      `Deletes a database. You must provide the ` + "`" + `database_id` + "`" + ` parameter. All data stored in the database will be permanently deleted.`,
		Namespace: "sdb-sql",
		Resource:  "database",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(serverless_sqldb.DeleteDatabaseRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "database-id",
				Short:      `UUID of the Serverless SQL Database.`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*serverless_sqldb.DeleteDatabaseRequest)

			client := core.ExtractClient(ctx)
			api := serverless_sqldb.NewAPI(client)

			return api.DeleteDatabase(request)
		},
	}
}

func sdbSQLDatabaseList() *core.Command {
	return &core.Command{
		Short:     `List your Serverless SQL Databases`,
		Long:      `List all Serverless SQL Databases for a given Scaleway Organization or Scaleway Project. By default, the databases returned in the list are ordered by creation date in ascending order, though this can be modified via the order_by field. For the ` + "`" + `name` + "`" + ` parameter, the value you include will be checked against the whole name string to see if it includes the string you put in the parameter.`,
		Namespace: "sdb-sql",
		Resource:  "database",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(serverless_sqldb.ListDatabasesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Filter by the name of the database`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sorting criteria. One of ` + "`" + `created_at_asc` + "`" + `, ` + "`" + `created_at_desc` + "`" + `, ` + "`" + `name_asc` + "`" + `, ` + "`" + `name_desc` + "`" + ``,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"name_asc",
					"name_desc",
				},
			},
			{
				Name:       "organization-id",
				Short:      `Filter by the UUID of the Scaleway organization`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*serverless_sqldb.ListDatabasesRequest)

			client := core.ExtractClient(ctx)
			api := serverless_sqldb.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListDatabases(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Databases, nil
		},
	}
}

func sdbSQLDatabaseUpdate() *core.Command {
	return &core.Command{
		Short:     `Update database information`,
		Long:      `Update CPU limits of your Serverless SQL Database. You must provide the ` + "`" + `database_id` + "`" + ` parameter.`,
		Namespace: "sdb-sql",
		Resource:  "database",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(serverless_sqldb.UpdateDatabaseRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "database-id",
				Short:      `UUID of the Serverless SQL Database.`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "cpu-min",
				Short:      `The minimum number of CPU units for your Serverless SQL Database.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cpu-max",
				Short:      `The maximum number of CPU units for your Serverless SQL Database.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*serverless_sqldb.UpdateDatabaseRequest)

			client := core.ExtractClient(ctx)
			api := serverless_sqldb.NewAPI(client)

			return api.UpdateDatabase(request)
		},
	}
}

func sdbSQLDatabaseRestore() *core.Command {
	return &core.Command{
		Short:     `Restore a database from a backup`,
		Long:      `Restore a database from a backup. You must provide the ` + "`" + `backup_id` + "`" + ` parameter.`,
		Namespace: "sdb-sql",
		Resource:  "database",
		Verb:      "restore",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(serverless_sqldb.RestoreDatabaseFromBackupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "database-id",
				Short:      `UUID of the Serverless SQL Database.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "backup-id",
				Short:      `UUID of the Serverless SQL Database backup to restore.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*serverless_sqldb.RestoreDatabaseFromBackupRequest)

			client := core.ExtractClient(ctx)
			api := serverless_sqldb.NewAPI(client)

			return api.RestoreDatabaseFromBackup(request)
		},
	}
}

func sdbSQLBackupGet() *core.Command {
	return &core.Command{
		Short:     `Get a database backup information`,
		Long:      `Retrieve information about your Serverless SQL Database backup. You must provide the ` + "`" + `backup_id` + "`" + ` parameter.`,
		Namespace: "sdb-sql",
		Resource:  "backup",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(serverless_sqldb.GetDatabaseBackupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backup-id",
				Short:      `UUID of the Serverless SQL Database backup.`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*serverless_sqldb.GetDatabaseBackupRequest)

			client := core.ExtractClient(ctx)
			api := serverless_sqldb.NewAPI(client)

			return api.GetDatabaseBackup(request)
		},
	}
}

func sdbSQLBackupList() *core.Command {
	return &core.Command{
		Short:     `List your Serverless SQL Database backups`,
		Long:      `List all Serverless SQL Database backups for a given Scaleway Project or Database. By default, the backups returned in the list are ordered by creation date in descending order, though this can be modified via the order_by field.`,
		Namespace: "sdb-sql",
		Resource:  "backup",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(serverless_sqldb.ListDatabaseBackupsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Filter by the UUID of the Scaleway project.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "database-id",
				Short:      `Filter by the UUID of the Serverless SQL Database.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sorting criteria. One of ` + "`" + `created_at_asc` + "`" + `, ` + "`" + `created_at_desc` + "`" + `.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
				},
			},
			{
				Name:       "organization-id",
				Short:      `Filter by the UUID of the Scaleway organization.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*serverless_sqldb.ListDatabaseBackupsRequest)

			client := core.ExtractClient(ctx)
			api := serverless_sqldb.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListDatabaseBackups(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Backups, nil
		},
	}
}

func sdbSQLBackupExport() *core.Command {
	return &core.Command{
		Short:     `Export a database backup`,
		Long:      `Export a database backup providing a download link once the export process is completed. You must provide the ` + "`" + `backup_id` + "`" + ` parameter.`,
		Namespace: "sdb-sql",
		Resource:  "backup",
		Verb:      "export",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(serverless_sqldb.ExportDatabaseBackupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backup-id",
				Short:      `UUID of the Serverless SQL Database backup.`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*serverless_sqldb.ExportDatabaseBackupRequest)

			client := core.ExtractClient(ctx)
			api := serverless_sqldb.NewAPI(client)

			return api.ExportDatabaseBackup(request)
		},
	}
}
