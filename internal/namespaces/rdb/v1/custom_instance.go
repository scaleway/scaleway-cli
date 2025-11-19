package rdb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	instanceActionTimeout               = 20 * time.Minute
	errorMessagePublicEndpointNotFound  = "public endpoint not found"
	errorMessagePrivateEndpointNotFound = "private endpoint not found"
	errorMessageEndpointNotFound        = "any endpoint is associated on your instance"
)

var instanceStatusMarshalSpecs = human.EnumMarshalSpecs{
	rdbSDK.InstanceStatusUnknown: &human.EnumMarshalSpec{
		Attribute: color.Faint,
		Value:     "unknown",
	},
	rdbSDK.InstanceStatusReady: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
		Value:     "ready",
	},
	rdbSDK.InstanceStatusProvisioning: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "provisioning",
	},
	rdbSDK.InstanceStatusConfiguring: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "configuring",
	},
	rdbSDK.InstanceStatusDeleting: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "deleting",
	},
	rdbSDK.InstanceStatusError: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "error",
	},
	rdbSDK.InstanceStatusAutohealing: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "auto-healing",
	},
	rdbSDK.InstanceStatusLocked: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "locked",
	},
	rdbSDK.InstanceStatusInitializing: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "initialized",
	},
	rdbSDK.InstanceStatusDiskFull: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "disk_full",
	},
	rdbSDK.InstanceStatusBackuping: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "backuping",
	},
}

type serverWaitRequest struct {
	InstanceID string
	Region     scw.Region
	Timeout    time.Duration
}

type CreateInstanceResult struct {
	*rdbSDK.Instance
	Password string `json:"password"`
}

type rdbCreateInstanceRequestCustom struct {
	*rdbSDK.CreateInstanceRequest
	InitEndpoints    []*rdbEndpointSpecCustom `json:"init-endpoints"`
	GeneratePassword bool
}

func createInstanceResultMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	instanceResult := i.(CreateInstanceResult)

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

	instanceStr, err := human.Marshal(instanceResult.Instance, opt)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{
		instanceStr,
		terminal.Style("Password: ", color.Bold) + "\n" + instanceResult.Password,
	}, "\n\n"), nil
}

func instanceMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	// To avoid recursion of human.Marshal we create a dummy type
	type tmp rdbSDK.Instance
	instance := tmp(i.(rdbSDK.Instance))

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

func backupScheduleMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	backupSchedule := i.(rdbSDK.BackupSchedule)

	if opt.TableCell {
		if backupSchedule.Disabled {
			return "Disabled", nil
		}

		return "Enabled", nil
	}

	// To avoid recursion of human.Marshal we create a dummy type
	type LocalBackupSchedule rdbSDK.BackupSchedule
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
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		api := rdbSDK.NewAPI(core.ExtractClient(ctx))

		return api.WaitForInstance(&rdbSDK.WaitForInstanceRequest{
			InstanceID:    respI.(*rdbSDK.Instance).ID,
			Region:        respI.(*rdbSDK.Instance).Region,
			Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

// Caching ListNodeType response for shell completion
var completeListNodeTypeCache *rdbSDK.ListNodeTypesResponse

var completeListEngineCache *rdbSDK.ListDatabaseEnginesResponse

func autoCompleteNodeType(
	ctx context.Context,
	prefix string,
	request any,
) core.AutocompleteSuggestions {
	region := scw.Region("")
	switch req := request.(type) {
	case *rdbSDK.CreateInstanceRequest:
		region = req.Region
	case *rdbSDK.UpgradeInstanceRequest:
		region = req.Region
	}

	suggestions := core.AutocompleteSuggestions(nil)

	client := core.ExtractClient(ctx)
	api := rdbSDK.NewAPI(client)

	if completeListNodeTypeCache == nil {
		res, err := api.ListNodeTypes(&rdbSDK.ListNodeTypesRequest{
			Region: region,
		}, scw.WithAllPages())
		if err != nil {
			return nil
		}
		completeListNodeTypeCache = res
	}

	for _, nodeType := range completeListNodeTypeCache.NodeTypes {
		if strings.HasPrefix(nodeType.Name, prefix) {
			suggestions = append(suggestions, nodeType.Name)
		}
	}

	return suggestions
}

func autoCompleteDatabaseEngines(
	ctx context.Context,
	prefix string,
	request any,
) core.AutocompleteSuggestions {
	var req *rdbCreateInstanceRequestCustom
	switch v := request.(type) {
	case rdbCreateInstanceRequestCustom:
		req = &v
	case *rdbCreateInstanceRequestCustom:
		req = v
	default:
		return nil
	}
	suggestion := core.AutocompleteSuggestions(nil)
	client := core.ExtractClient(ctx)
	api := rdbSDK.NewAPI(client)

	if completeListEngineCache == nil {
		res, err := api.ListDatabaseEngines(&rdbSDK.ListDatabaseEnginesRequest{
			Region: req.Region,
		}, scw.WithAllPages())
		if err != nil {
			return nil
		}
		completeListEngineCache = res
	}

	for _, engine := range completeListEngineCache.Engines {
		if strings.HasPrefix(engine.Name, prefix) {
			suggestion = append(suggestion, engine.Name)
		}
	}

	return suggestion
}

func instanceCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.AddBefore("init-endpoints.{index}.private-network.private-network-id", &core.ArgSpec{
		Name:     "init-endpoints.{index}.load-balancer",
		Short:    "Will configure a load-balancer endpoint along with your private network endpoint if true",
		Required: false,
		Default:  core.DefaultValueSetter("false"),
	})
	c.ArgSpecs.AddBefore("init-endpoints.{index}.private-network.private-network-id", &core.ArgSpec{
		Name:       "init-endpoints.{index}.private-network.enable-ipam",
		Short:      "Will configure your Private Network endpoint with Scaleway IPAM service if true",
		Required:   false,
		OneOfGroup: "config",
	})
	c.ArgSpecs.AddBefore("password", &core.ArgSpec{
		Name:       "generate-password",
		Short:      `Will generate a 21 character-length password that contains a mix of upper/lower case letters, numbers and special symbols`,
		Required:   false,
		Deprecated: false,
		Positional: false,
		Default:    core.DefaultValueSetter("true"),
	})
	c.ArgSpecs.GetByName("password").Required = false
	c.ArgSpecs.GetByName("node-type").Default = core.DefaultValueSetter("DB-DEV-S")
	c.ArgSpecs.GetByName("node-type").AutoCompleteFunc = autoCompleteNodeType
	c.ArgSpecs.GetByName("engine").AutoCompleteFunc = autoCompleteDatabaseEngines

	c.ArgsType = reflect.TypeOf(rdbCreateInstanceRequestCustom{})

	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		api := rdbSDK.NewAPI(core.ExtractClient(ctx))
		instance, err := api.WaitForInstance(&rdbSDK.WaitForInstanceRequest{
			InstanceID:    respI.(CreateInstanceResult).ID,
			Region:        respI.(CreateInstanceResult).Region,
			Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			return nil, err
		}

		result := CreateInstanceResult{
			Instance: instance,
			Password: respI.(CreateInstanceResult).Password,
		}

		return result, nil
	}

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		client := core.ExtractClient(ctx)
		api := rdbSDK.NewAPI(client)

		customRequest := argsI.(*rdbCreateInstanceRequestCustom)
		createInstanceRequest := customRequest.CreateInstanceRequest

		var err error
		createInstanceRequest.NodeType = strings.ToLower(createInstanceRequest.NodeType)
		if customRequest.GeneratePassword && customRequest.Password == "" {
			createInstanceRequest.Password, err = passwordgenerator.GeneratePassword(21, 1, 1, 1, 1)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Your generated password is %s \n", createInstanceRequest.Password)
			fmt.Printf("\n")
		}

		createInstanceRequest.InitEndpoints, err = endpointRequestFromCustom(
			customRequest.InitEndpoints,
		)
		if err != nil {
			return nil, err
		}

		instance, err := api.CreateInstance(createInstanceRequest)
		if err != nil {
			return nil, err
		}

		result := CreateInstanceResult{
			Instance: instance,
			Password: createInstanceRequest.Password,
		}

		return result, nil
	}

	// Waiting for API to accept uppercase node-type
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		args := argsI.(*rdbCreateInstanceRequestCustom)
		args.NodeType = strings.ToLower(args.NodeType)

		return runner(ctx, args)
	}

	return c
}

func instanceGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		res, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}
		instance := res.(*rdbSDK.Instance)

		args := argsI.(*rdbSDK.GetInstanceRequest)

		acls, err := rdbSDK.NewAPI(core.ExtractClient(ctx)).
			ListInstanceACLRules(&rdbSDK.ListInstanceACLRulesRequest{
				Region:     args.Region,
				InstanceID: args.InstanceID,
			}, scw.WithAllPages())
		if err != nil {
			return nil, err
		}

		return struct {
			*rdbSDK.Instance
			ACLs []*rdbSDK.ACLRule `json:"acls"`
		}{
			instance,
			acls.Rules,
		}, nil
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName: "Endpoint",
				Title:     "Endpoint",
			},
			{
				FieldName: "Volume",
				Title:     "Volume",
			},
			{
				FieldName: "BackupSchedule",
				Title:     "Backup schedule",
			},
			{
				FieldName:   "Settings",
				Title:       "Settings",
				HideIfEmpty: true,
			},
			{
				FieldName:   "ACLs",
				Title:       "ACLs",
				HideIfEmpty: true,
			},
		},
	}

	return c
}

func instanceUpgradeInterceptor(
	ctx context.Context,
	argsI any,
	runner core.CommandRunner,
) (any, error) {
	req := argsI.(*rdbSDK.UpgradeInstanceRequest)
	api := rdbSDK.NewAPI(core.ExtractClient(ctx))

	instance, err := api.GetInstance(&rdbSDK.GetInstanceRequest{
		Region:     req.Region,
		InstanceID: req.InstanceID,
	})
	if err != nil {
		return nil, err
	}

	if !needsUpgrade(req, instance) {
		return &core.SuccessResult{Message: "Nothing to do!"}, nil
	}

	return runner(ctx, argsI)
}

func needsUpgrade(req *rdbSDK.UpgradeInstanceRequest, instance *rdbSDK.Instance) bool {
	if req.NodeType != nil && *req.NodeType != "" {
		if !strings.EqualFold(instance.NodeType, *req.NodeType) {
			return true
		}
	}

	if req.EnableHa != nil && *req.EnableHa && !instance.IsHaCluster {
		return true
	}

	if instance.Volume != nil {
		if req.VolumeType != nil && *req.VolumeType != instance.Volume.Type {
			return true
		}
		if req.VolumeSize != nil && *req.VolumeSize > uint64(instance.Volume.Size) {
			return true
		}
	}

	if req.UpgradableVersionID != nil && *req.UpgradableVersionID != "" {
		return true
	}

	if req.MajorUpgradeWorkflow != nil && req.MajorUpgradeWorkflow.UpgradableVersionID != "" {
		return true
	}

	return false
}

func instanceUpgradeBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("node-type").AutoCompleteFunc = autoCompleteNodeType

	c.Interceptor = instanceUpgradeInterceptor

	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		if _, ok := respI.(*core.SuccessResult); ok {
			return respI, nil
		}

		instance := respI.(*rdbSDK.Instance)
		api := rdbSDK.NewAPI(core.ExtractClient(ctx))

		return api.WaitForInstance(&rdbSDK.WaitForInstanceRequest{
			InstanceID:    instance.ID,
			Region:        instance.Region,
			Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func instanceUpdateBuilder(_ *core.Command) *core.Command {
	type rdbUpdateInstanceRequestCustom struct {
		*rdbSDK.UpdateInstanceRequest
		Settings []*rdbSDK.InstanceSetting
	}

	return &core.Command{
		Short:     `Update an instance`,
		Long:      `Update an instance.`,
		Namespace: "rdb",
		Resource:  "instance",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(rdbUpdateInstanceRequestCustom{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backup-schedule-frequency",
				Short:      `In hours`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "backup-schedule-retention",
				Short:      `In days`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-backup-schedule-disabled",
				Short:      `Whether or not the backup schedule is disabled`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "instance-id",
				Short:      `UUID of the instance to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of a given instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "logs-policy.max-age-retention",
				Short:      `Max age (in day) of remote logs to keep on the database instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "logs-policy.total-disk-retention",
				Short:      `Max disk size of remote logs to keep on the database instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "backup-same-region",
				Short:      `Store logical backups in the same region as the database instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "settings.{index}.name",
				Short:      `Setting name of a given instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "settings.{index}.value",
				Short:      `Setting value of a given instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			customRequest := args.(*rdbUpdateInstanceRequestCustom)

			updateInstanceRequest := customRequest.UpdateInstanceRequest

			client := core.ExtractClient(ctx)
			api := rdbSDK.NewAPI(client)

			getResp, err := api.GetInstance(&rdbSDK.GetInstanceRequest{
				Region:     customRequest.Region,
				InstanceID: customRequest.InstanceID,
			})
			if err != nil {
				return nil, err
			}

			if customRequest.Settings != nil {
				settings := getResp.Settings
				changes := customRequest.Settings

				for _, change := range changes {
					matched := false
					for _, setting := range settings {
						if change.Name == setting.Name {
							setting.Value = change.Value
							matched = true

							break
						}
					}
					if !matched {
						settings = append(settings, change)
					}
				}

				_, err = api.SetInstanceSettings(&rdbSDK.SetInstanceSettingsRequest{
					Region:     updateInstanceRequest.Region,
					InstanceID: updateInstanceRequest.InstanceID,
					Settings:   settings,
				})
				if err != nil {
					return nil, err
				}
			}

			updateInstanceResponse, err := api.UpdateInstance(updateInstanceRequest)
			if err != nil {
				return nil, err
			}

			return updateInstanceResponse, nil
		},
		WaitFunc: func(ctx context.Context, _, respI any) (any, error) {
			api := rdbSDK.NewAPI(core.ExtractClient(ctx))

			return api.WaitForInstance(&rdbSDK.WaitForInstanceRequest{
				InstanceID:    respI.(*rdbSDK.Instance).ID,
				Region:        respI.(*rdbSDK.Instance).Region,
				Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
		},
		Examples: []*core.Example{
			{
				Short: "Update instance name",
				Raw:   "scw rdb instance update 11111111-1111-1111-1111-111111111111 name=foo --wait",
			},
			{
				Short: "Update instance tags",
				Raw:   "scw rdb instance update 11111111-1111-1111-1111-111111111111 tags.0=a --wait",
			},
			{
				Short: "Set a timezone",
				Raw:   "scw rdb instance update 11111111-1111-1111-1111-111111111111 settings.0.name=timezone settings.0.value=UTC --wait",
			},
		},
	}
}

func instanceDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		api := rdbSDK.NewAPI(core.ExtractClient(ctx))
		instance, err := api.WaitForInstance(&rdbSDK.WaitForInstanceRequest{
			InstanceID:    respI.(*rdbSDK.Instance).ID,
			Region:        respI.(*rdbSDK.Instance).Region,
			Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			// if we get a 404 here, it means the resource was successfully deleted
			notFoundError := &scw.ResourceNotFoundError{}
			responseError := &scw.ResponseError{}
			if errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound ||
				errors.As(err, &notFoundError) {
				return instance, nil
			}

			return nil, err
		}

		return instance, nil
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
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(serverWaitRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			api := rdbSDK.NewAPI(core.ExtractClient(ctx))

			return api.WaitForInstance(&rdbSDK.WaitForInstanceRequest{
				Region:        argsI.(*serverWaitRequest).Region,
				InstanceID:    argsI.(*serverWaitRequest).InstanceID,
				Timeout:       scw.TimeDurationPtr(argsI.(*serverWaitRequest).Timeout),
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
			core.WaitTimeoutArgSpec(instanceActionTimeout),
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
	Region         scw.Region
	PrivateNetwork bool
	InstanceID     string
	Username       string
	Database       *string
	CliDB          *string
}

type engineFamily string

const (
	Unknown        = engineFamily("Unknown")
	PostgreSQL     = engineFamily("PostgreSQL")
	MySQL          = engineFamily("MySQL")
	postgreSQLHint = `
psql supports password file (.pgpass) to avoid typing your password manually.
Create ~/.pgpass (Linux/macOS) or %APPDATA%\postgresql\pgpass.conf (Windows) with:
  hostname:port:database:username:password
Learn more at: https://www.postgresql.org/docs/current/libpq-pgpass.html`
	mySQLHint = `
mysql supports mysql_config_editor for secure password storage.
Use: mysql_config_editor set --login-path=scw --host=HOST --user=USER --password
Or create ~/.mylogin.cnf with connection credentials.
Learn more at: https://dev.mysql.com/doc/refman/8.0/en/mysql-config-editor.html`
)

func passwordFileExist(ctx context.Context, family engineFamily) bool {
	var passwordFilePath string
	switch family {
	case PostgreSQL:
		switch runtime.GOOS {
		case "windows":
			passwordFilePath = path.Join(
				core.ExtractUserHomeDir(ctx),
				core.ExtractEnv(ctx, "APPDATA"),
				"postgresql",
				"pgpass.conf",
			)
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

func detectEngineFamily(instance *rdbSDK.Instance) (engineFamily, error) {
	if instance == nil {
		return Unknown, errors.New("instance engine is nil")
	}
	if strings.HasPrefix(instance.Engine, string(PostgreSQL)) {
		return PostgreSQL, nil
	}
	if strings.HasPrefix(instance.Engine, string(MySQL)) {
		return MySQL, nil
	}

	return Unknown, fmt.Errorf("unknown engine: %s", instance.Engine)
}

func getPublicEndpoint(endpoints []*rdbSDK.Endpoint) (*rdbSDK.Endpoint, error) {
	for _, e := range endpoints {
		if e.LoadBalancer != nil {
			return e, nil
		}
	}

	return nil, fmt.Errorf("%s", errorMessagePublicEndpointNotFound)
}

func getPrivateEndpoint(endpoints []*rdbSDK.Endpoint) (*rdbSDK.Endpoint, error) {
	for _, e := range endpoints {
		if e.PrivateNetwork != nil {
			return e, nil
		}
	}

	return nil, fmt.Errorf("%s", errorMessagePrivateEndpointNotFound)
}

func createConnectCommandLineArgs(
	endpoint *rdbSDK.Endpoint,
	family engineFamily,
	args *instanceConnectArgs,
) ([]string, error) {
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
			"--host", endpoint.IP.String(),
			"--port", strconv.FormatUint(uint64(endpoint.Port), 10),
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
			"--host", endpoint.IP.String(),
			"--port", strconv.FormatUint(uint64(endpoint.Port), 10),
			"--database", database,
			"--user", args.Username,
		}, nil
	}

	return nil, fmt.Errorf("unrecognize database engine: %s", family)
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
				Name:     "private-network",
				Short:    `Connect by the private network endpoint attached.`,
				Required: false,
				Default:  core.DefaultValueSetter("false"),
			},
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
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*instanceConnectArgs)

			client := core.ExtractClient(ctx)
			api := rdbSDK.NewAPI(client)
			instance, err := api.GetInstance(&rdbSDK.GetInstanceRequest{
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

			if len(instance.Endpoints) == 0 {
				return nil, fmt.Errorf("%s", errorMessageEndpointNotFound)
			}

			var endpoint *rdbSDK.Endpoint
			switch {
			case args.PrivateNetwork:
				endpoint, err = getPrivateEndpoint(instance.Endpoints)
				if err != nil {
					return nil, err
				}
			default:
				endpoint, err = getPublicEndpoint(instance.Endpoints)
				if err != nil {
					return nil, err
				}
			}

			cmdArgs, err := createConnectCommandLineArgs(endpoint, engineFamily, args)
			if err != nil {
				return nil, err
			}

			if !passwordFileExist(ctx, engineFamily) {
				interactive.Println(passwordFileHint(engineFamily))
			}

			// Run command
			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec
			// cmd.Stdin = os.Stdin
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

func instanceEditSettingsCommand() *core.Command {
	type editSettingsArgs struct {
		InstanceID string     `arg:"positional,required"`
		Region     scw.Region `arg:"required"`
		Mode       editor.MarshalMode
	}

	return &core.Command{
		Namespace: "rdb",
		Resource:  "setting",
		Verb:      "edit",
		Short:     "Edit Database Instance settings in your default editor",
		Long: `This command opens the current settings of your RDB instance in your $EDITOR.
You can modify the values and save the file to apply the new configuration.`,
		ArgsType: reflect.TypeOf(editSettingsArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-id",
				Short:      "ID of the instance",
				Required:   true,
				Positional: true,
			},
			editor.MarshalModeArgSpec(), // --mode=yaml|json
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Examples: []*core.Example{
			{
				Short: "Edit instance settings in YAML",
				Raw:   "scw rdb setting edit 12345678-1234-1234-1234-123456789abc --region=fr-par --mode=yaml",
			},
			{
				Short: "Edit instance settings in JSON",
				Raw:   "scw rdb setting edit 12345678-1234-1234-1234-123456789abc --region=fr-par --mode=json",
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*editSettingsArgs)

			client := core.ExtractClient(ctx)
			api := rdbSDK.NewAPI(client)

			instance, err := api.GetInstance(&rdbSDK.GetInstanceRequest{
				InstanceID: args.InstanceID,
				Region:     args.Region,
			})
			if err != nil {
				return nil, err
			}

			initialRequest := &rdbSDK.SetInstanceSettingsRequest{
				Region:     args.Region,
				InstanceID: args.InstanceID,
				Settings:   instance.Settings,
			}

			editedRequestRaw, err := editor.UpdateResourceEditor(
				initialRequest,
				&rdbSDK.SetInstanceSettingsRequest{
					Region:     args.Region,
					InstanceID: args.InstanceID,
				},
				&editor.Config{
					PutRequest:  true,
					MarshalMode: args.Mode,
				},
			)
			if err != nil {
				return nil, err
			}

			editedRequest := editedRequestRaw.(*rdbSDK.SetInstanceSettingsRequest)

			if reflect.DeepEqual(initialRequest.Settings, editedRequest.Settings) {
				return &core.SuccessResult{Message: "No changes detected."}, nil
			}

			_, err = api.SetInstanceSettings(editedRequest)
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{Message: "Settings successfully updated."}, nil
		},
	}
}
