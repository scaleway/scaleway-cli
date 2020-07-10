package rdb

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"reflect"
	"runtime"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type instanceConnectArgs struct {
	Region     scw.Region
	InstanceID string
	Username   string
	Database   *string
	CliDB      *string
}

type engineFamily string

const (
	Unknown        = engineFamily("Unknown")
	PostgreSQL     = engineFamily("PostgreSQL")
	MySQL          = engineFamily("MySQL")
	postgreSQLHint = `
psql supports password file to avoid typing your password manually.
Learn more at: https://www.postgresql.org/docs/current/libpq-pgpass.html`
	mySQLHint = `
mysql supports loading your password from a file to avoid typing them manually.
Learn more at: https://dev.mysql.com/doc/refman/8.0/en/option-files.html`
)

func passwordFileExist(ctx context.Context, family engineFamily) bool {
	passwordFilePath := ""
	switch family {
	case PostgreSQL:
		switch runtime.GOOS {
		case "windows":
			passwordFilePath = path.Join(core.ExtractUserHomeDir(ctx), core.ExtractEnv(ctx, "APPDATA"), "postgresql", "pgpass.conf")
		default:
			passwordFilePath = path.Join(core.ExtractUserHomeDir(ctx), ".pgpass")
		}
	case MySQL:
		passwordFilePath = path.Join(core.ExtractUserHomeDir(ctx), ".my.cnf")
	default:
		return false
	}
	if passwordFilePath == "" {
		return false
	}
	_, err := os.Stat(passwordFilePath)
	return err == nil
}

func passwordFileHint(family engineFamily) string {
	switch family {
	case PostgreSQL:
		return postgreSQLHint
	case MySQL:
		return mySQLHint
	default:
		return ""
	}
}

func detectEngineFamily(instance *rdb.Instance) (engineFamily, error) {
	if instance == nil {
		return Unknown, fmt.Errorf("instance engine is nil")
	}
	if strings.HasPrefix(instance.Engine, string(PostgreSQL)) {
		return PostgreSQL, nil
	}
	if strings.HasPrefix(instance.Engine, string(MySQL)) {
		return MySQL, nil
	}
	return Unknown, fmt.Errorf("unknown engine: %s", instance.Engine)
}

func createConnectCommandLineArgs(instance *rdb.Instance, family engineFamily, args *instanceConnectArgs) ([]string, error) {
	database := "rdb"
	if args.Database != nil {
		database = *args.Database
	}

	switch family {
	case PostgreSQL:
		clidb := "psql"
		if args.CliDB != nil {
			clidb = *args.CliDB
		}

		// psql -h 51.159.25.206 --port 13917 -d rdb -U username
		return []string{
			clidb,
			"--host", instance.Endpoint.IP.String(),
			"--port", fmt.Sprintf("%d", instance.Endpoint.Port),
			"--username", args.Username,
			"--dbname", database,
		}, nil
	case MySQL:
		clidb := "mysql"
		if args.CliDB != nil {
			clidb = *args.CliDB
		}

		// mysql -h 195.154.69.163 --port 12210 -p -u username
		return []string{
			clidb,
			"--host", instance.Endpoint.IP.String(),
			"--port", fmt.Sprintf("%d", instance.Endpoint.Port),
			"--database", database,
			"--user", args.Username,
		}, nil
	}

	return nil, fmt.Errorf("unrecognize database engine: %s", instance.Engine)
}

func instanceConnectCommand() *core.Command {
	return &core.Command{
		Namespace: "rdb",
		Resource:  "instance",
		Verb:      "connect",
		Short:     "Connect to an instance using local database cli",
		Long:      "Connect to an instance using local database cli such as psql or mysql.",
		ArgsType:  reflect.TypeOf(instanceConnectArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-id",
				Short:      `UUID of the instance`,
				Required:   true,
				Positional: true,
			},
			{
				Name:     "username",
				Short:    "Name of the user to connect with to the database",
				Required: true,
			},
			{
				Name:  "database",
				Short: "Name of the database",
			},
			{
				Name:  "cli-db",
				Short: "Command line tool to use, default to psql/mysql",
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*instanceConnectArgs)

			client := core.ExtractClient(ctx)
			api := rdb.NewAPI(client)
			instance, err := api.GetInstance(&rdb.GetInstanceRequest{
				Region:     args.Region,
				InstanceID: args.InstanceID,
			})
			if err != nil {
				return nil, err
			}

			engineFamily, err := detectEngineFamily(instance)
			if err != nil {
				return nil, err
			}

			cmdArgs, err := createConnectCommandLineArgs(instance, engineFamily, args)
			if err != nil {
				return nil, err
			}

			if !passwordFileExist(ctx, engineFamily) {
				interactive.Println(passwordFileHint(engineFamily))
			}

			// Run command
			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec
			cmd.Stdin = os.Stdin
			exitCode, err := core.ExecCmd(ctx, cmd)

			if err != nil {
				return nil, err
			}
			if exitCode != 0 {
				return nil, &core.CliError{Empty: true, Code: exitCode}
			}

			return &core.SuccessResult{
				Empty: true, // the program will output the success message
			}, nil
		},
	}
}
