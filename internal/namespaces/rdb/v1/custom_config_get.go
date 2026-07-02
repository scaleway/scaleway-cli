package rdb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type rdbConfigGetRequest struct {
	InstanceID     string
	Type           rdbConfigType
	User           string
	Db             string //nolint: stylecheck
	PrivateNetwork bool
	Region         scw.Region
}

func rdbConfigRoot() *core.Command {
	return &core.Command{
		Short:     "Database client configuration snippets",
		Long:      "Generate ready-to-use database client configuration snippets.",
		Namespace: "rdb",
		Resource:  "config",
	}
}

func rdbConfigGetCommand() *core.Command {
	return &core.Command{
		Namespace: "rdb",
		Resource:  "config",
		Verb:      "get",
		Short:     "Generate a database client configuration snippet",
		Long: `Generate a ready-to-use database client configuration snippet for a Database Instance.

Supported languages:
  - php: PHP connection snippet (pg_connect for PostgreSQL, PDO for MySQL).
  - node: Node.js connection snippet (pg or mysql2).
  - typescript: TypeScript connection snippet (pg or mysql2).
  - python: Python connection snippet (psycopg2 or mysql.connector).
  - go: Go connection snippet (database/sql with pgx or go-sql-driver/mysql).
  - rust: Rust connection snippet (sqlx).

Use private-network=true when the instance has no public endpoint and is reachable only
from resources attached to its Private Network.
Replace YOUR_PASSWORD with your database user password, or use environment variables
where suggested in the snippet.`,
		ArgsType: reflect.TypeOf(rdbConfigGetRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-id",
				Short:      "ID of the Database Instance",
				Required:   true,
				Positional: true,
			},
			{
				Name:     "type",
				Short:    "Configuration template type",
				Required: true,
				EnumValues: []string{
					string(rdbConfigTypePHP),
					string(rdbConfigTypeNode),
					string(rdbConfigTypeTypeScript),
					string(rdbConfigTypePython),
					string(rdbConfigTypeGo),
					string(rdbConfigTypeRust),
				},
			},
			{
				Name:  "user",
				Short: "Database user to connect as",
			},
			{
				Name:  "db",
				Short: "Database name to connect to (defaults to rdb)",
			},
			{
				Name:  "private-network",
				Short: "Use the Private Network endpoint instead of the public one",
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Examples: []*core.Example{
			{
				Short:    "Generate a PHP connection snippet",
				ArgsJSON: `{"instance_id":"11111111-1111-1111-1111-111111111111","type":"php"}`,
			},
			{
				Short:    "Generate a Node.js connection snippet for a custom user and database",
				ArgsJSON: `{"instance_id":"11111111-1111-1111-1111-111111111111","type":"node","user":"myuser","db":"mydb"}`,
			},
			{
				Short: "Generate a Python connection snippet using the Private Network endpoint",
				ArgsJSON: `{"instance_id":"11111111-1111-1111-1111-111111111111",` +
					`"type":"python","private_network":true}`,
			},
			{
				Short:    "Generate a Go connection snippet",
				ArgsJSON: `{"instance_id":"11111111-1111-1111-1111-111111111111","type":"go"}`,
			},
			{
				Short:    "Generate a Rust connection snippet",
				ArgsJSON: `{"instance_id":"11111111-1111-1111-1111-111111111111","type":"rust"}`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw rdb user get-url",
				Short:   "Get a connection URL for a database user",
			},
			{
				Command: "scw rdb instance connect",
				Short:   "Open a database shell session",
			},
			{
				Command: "scw rdb instance get-certificate",
				Short:   "Get the TLS certificate of a Database Instance",
			},
		},
		Run: rdbConfigGetRun,
	}
}

func rdbConfigGetRun(ctx context.Context, argsI any) (any, error) {
	args := argsI.(*rdbConfigGetRequest)

	client := core.ExtractClient(ctx)
	api := rdb.NewAPI(client)

	info, err := resolveRDBConnectionInfo(ctx, api, &rdbConnectionParams{
		InstanceID:     args.InstanceID,
		Region:         args.Region,
		User:           args.User,
		Db:             args.Db,
		PrivateNetwork: args.PrivateNetwork,
	})
	if err != nil {
		return nil, err
	}

	return renderRDBConfig(args.Type, info)
}
