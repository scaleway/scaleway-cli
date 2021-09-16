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
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	instanceActionTimeout = 20 * time.Minute
)

var (
	instanceStatusMarshalSpecs = human.EnumMarshalSpecs{
		rdb.InstanceStatusUnknown:      &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
		rdb.InstanceStatusReady:        &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
		rdb.InstanceStatusProvisioning: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "provisioning"},
		rdb.InstanceStatusConfiguring:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "configuring"},
		rdb.InstanceStatusDeleting:     &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "deleting"},
		rdb.InstanceStatusError:        &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
		rdb.InstanceStatusAutohealing:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "auto-healing"},
		rdb.InstanceStatusLocked:       &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "locked"},
		rdb.InstanceStatusInitializing: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "initialized"},
		rdb.InstanceStatusDiskFull:     &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "disk_full"},
		rdb.InstanceStatusBackuping:    &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "backuping"},
	}
)

type serverWaitRequest struct {
	InstanceID string
	Region     scw.Region
}

func instanceMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	// To avoid recursion of human.Marshal we create a dummy type
	type tmp rdb.Instance
	instance := tmp(i.(rdb.Instance))

	// Sections
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "Endpoint",
		},
		{
			FieldName: "Volume",
		},
		{
			FieldName: "BackupSchedule",
		},
		{
			FieldName: "Settings",
		},
	}

	str, err := human.Marshal(instance, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func backupScheduleMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	backupSchedule := i.(rdb.BackupSchedule)

	if opt.TableCell {
		if backupSchedule.Disabled {
			return "Disabled", nil
		}
		return "Enabled", nil
	}

	// To avoid recursion of human.Marshal we create a dummy type
	type LocalBackupSchedule rdb.BackupSchedule
	type tmp struct {
		LocalBackupSchedule
		Frequency *scw.Duration `json:"frequency"`
		Retention *scw.Duration `json:"retention"`
	}

	localBackupSchedule := tmp{
		LocalBackupSchedule: LocalBackupSchedule(backupSchedule),
		Frequency: &scw.Duration{
			Seconds: int64(backupSchedule.Frequency) * 3600,
		},
		Retention: &scw.Duration{
			Seconds: int64(backupSchedule.Retention) * 24 * 3600,
		},
	}

	str, err := human.Marshal(localBackupSchedule, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func instanceCloneBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForInstance(&rdb.WaitForInstanceRequest{
			InstanceID:    respI.(*rdb.Instance).ID,
			Region:        respI.(*rdb.Instance).Region,
			Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func instanceCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("node-type").Default = core.DefaultValueSetter("DB-DEV-S")
	c.ArgSpecs.GetByName("node-type").EnumValues = nodeTypes

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForInstance(&rdb.WaitForInstanceRequest{
			InstanceID:    respI.(*rdb.Instance).ID,
			Region:        respI.(*rdb.Instance).Region,
			Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	// Waiting for API to accept uppercase node-type
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		args := argsI.(*rdb.CreateInstanceRequest)
		args.NodeType = strings.ToLower(args.NodeType)
		return runner(ctx, args)
	}

	return c
}

func instanceUpgradeBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("node-type").EnumValues = nodeTypes

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForInstance(&rdb.WaitForInstanceRequest{
			InstanceID:    respI.(*rdb.Instance).ID,
			Region:        respI.(*rdb.Instance).Region,
			Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func instanceWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for an instance to reach a stable state`,
		Long:      `Wait for an instance to reach a stable state. This is similar to using --wait flag.`,
		Namespace: "rdb",
		Resource:  "instance",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(serverWaitRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := rdb.NewAPI(core.ExtractClient(ctx))
			return api.WaitForInstance(&rdb.WaitForInstanceRequest{
				Region:        argsI.(*serverWaitRequest).Region,
				InstanceID:    argsI.(*serverWaitRequest).InstanceID,
				Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-id",
				Short:      `ID of the instance you want to wait for.`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for an instance to reach a stable state",
				ArgsJSON: `{"instance_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

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
			"--host", instance.Endpoints[0].IP.String(),
			"--port", fmt.Sprintf("%d", instance.Endpoints[0].Port),
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
			"--host", instance.Endpoints[0].IP.String(),
			"--port", fmt.Sprintf("%d", instance.Endpoints[0].Port),
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
		Short:     "Connect to an instance using locally installed CLI",
		Long:      "Connect to an instance using locally installed CLI such as psql or mysql.",
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
				Name:    "database",
				Short:   "Name of the database",
				Default: core.DefaultValueSetter("rdb"),
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
			//cmd.Stdin = os.Stdin
			core.ExtractLogger(ctx).Debugf("executing: %s\n", cmd.Args)
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
