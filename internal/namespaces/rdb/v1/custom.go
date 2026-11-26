package rdb

import (
	"fmt"
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("rdb").Groups = []string{"database"}

	// TODO
	// human.RegisterMarshalerFunc(rdb.Instance{}, instanceMarshalerFunc)
	// human.RegisterMarshalerFunc(rdb.BackupSchedule{}, backupScheduleMarshalerFunc)
	// human.RegisterMarshalerFunc(backupDownloadResult{}, backupResultMarshallerFunc)
	// human.RegisterMarshalerFunc(CreateInstanceResult{}, createInstanceResultMarshalerFunc)
	// human.RegisterMarshalerFunc(CustomACLResult{}, rdbACLCustomResultMarshalerFunc)
	// human.RegisterMarshalerFunc(rdbEndpointCustomResult{}, rdbEndpointCustomResultMarshalerFunc)
	// human.RegisterMarshalerFunc(logDownloadResult{}, logDownloadResultMarshalerFunc)

	// human.RegisterMarshalerFunc(
	// 	rdb.InstanceStatus(""),
	// 	human.EnumMarshalFunc(instanceStatusMarshalSpecs),
	// )
	// human.RegisterMarshalerFunc(
	// 	rdb.DatabaseBackupStatus(""),
	// 	human.EnumMarshalFunc(backupStatusMarshalSpecs),
	// )
	// human.RegisterMarshalerFunc(
	// 	rdb.InstanceLogStatus(""),
	// 	human.EnumMarshalFunc(logStatusMarshalSpecs),
	// )
	// human.RegisterMarshalerFunc(
	// 	rdb.NodeTypeStock(""),
	// 	human.EnumMarshalFunc(nodeTypeStockMarshalSpecs),
	// )
	// human.RegisterMarshalerFunc(
	// 	rdb.ACLRuleAction(""),
	// 	human.EnumMarshalFunc(aclRuleActionMarshalSpecs),
	// )

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[rdb.ACLRuleAction]{
		rdb.ACLRuleActionAllow: humanv2.EnumFormat{Color: color.FgGreen, Value: "allow"},
		rdb.ACLRuleActionDeny:  humanv2.EnumFormat{Color: color.FgRed, Value: "deny"},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[rdb.DatabaseBackupStatus]{
		rdb.DatabaseBackupStatusUnknown:   humanv2.EnumFormat{Color: color.Faint, Value: "unknown"},
		rdb.DatabaseBackupStatusCreating:  humanv2.EnumFormat{Color: color.FgBlue, Value: "creating"},
		rdb.DatabaseBackupStatusReady:     humanv2.EnumFormat{Color: color.FgGreen, Value: "ready"},
		rdb.DatabaseBackupStatusRestoring: humanv2.EnumFormat{Color: color.FgBlue, Value: "restoring"},
		rdb.DatabaseBackupStatusDeleting:  humanv2.EnumFormat{Color: color.FgBlue, Value: "deleting"},
		rdb.DatabaseBackupStatusError:     humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
		rdb.DatabaseBackupStatusExporting: humanv2.EnumFormat{Color: color.FgBlue, Value: "exporting"},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[rdb.InstanceStatus]{
		rdb.InstanceStatusUnknown:      humanv2.EnumFormat{Color: color.Faint, Value: "unknown"},
		rdb.InstanceStatusReady:        humanv2.EnumFormat{Color: color.FgGreen, Value: "ready"},
		rdb.InstanceStatusProvisioning: humanv2.EnumFormat{Color: color.FgBlue, Value: "provisioning"},
		rdb.InstanceStatusConfiguring:  humanv2.EnumFormat{Color: color.FgBlue, Value: "configuring"},
		rdb.InstanceStatusDeleting:     humanv2.EnumFormat{Color: color.FgBlue, Value: "deleting"},
		rdb.InstanceStatusError:        humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
		rdb.InstanceStatusAutohealing:  humanv2.EnumFormat{Color: color.FgBlue, Value: "auto-healing"},
		rdb.InstanceStatusLocked:       humanv2.EnumFormat{Color: color.FgRed, Value: "locked"},
		rdb.InstanceStatusInitializing: humanv2.EnumFormat{Color: color.FgBlue, Value: "initialized"},
		rdb.InstanceStatusDiskFull:     humanv2.EnumFormat{Color: color.FgRed, Value: "disk_full"},
		rdb.InstanceStatusBackuping:    humanv2.EnumFormat{Color: color.FgBlue, Value: "backuping"},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[rdb.InstanceLogStatus]{
		rdb.InstanceLogStatusUnknown:  humanv2.EnumFormat{Color: color.Faint, Value: "unknown"},
		rdb.InstanceLogStatusReady:    humanv2.EnumFormat{Color: color.FgGreen, Value: "ready"},
		rdb.InstanceLogStatusCreating: humanv2.EnumFormat{Color: color.FgBlue, Value: "creating"},
		rdb.InstanceLogStatusError:    humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[rdb.NodeTypeStock]{
		rdb.NodeTypeStockAvailable:  humanv2.EnumFormat{Color: color.FgGreen, Value: "available"},
		rdb.NodeTypeStockUnknown:    humanv2.EnumFormat{Color: color.Faint, Value: "unknown"},
		rdb.NodeTypeStockLowStock:   humanv2.EnumFormat{Color: color.FgYellow, Value: "low stock"},
		rdb.NodeTypeStockOutOfStock: humanv2.EnumFormat{Color: color.FgRed, Value: "out of stock"},
	})

	humanv2.RegisterLayout(reflect.TypeFor[CustomACLResult](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Main: true, Path: "Success"},
			{Path: "Rules"},
		},
	})

	humanv2.RegisterLayout(reflect.TypeFor[rdb.Instance](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "Endpoint"},
			{Path: "Volume"},
			{Path: "BackupSchedule"},
			{Path: "Settings"},
		},
	})

	humanv2.RegisterLayout(reflect.TypeFor[CreateInstanceResult](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "Endpoint"},
			{Path: "Volume"},
			{Path: "BackupSchedule"},
			{Path: "Settings"},
			{Path: "Password"},
		},
	})

	humanv2.RegisterLayout(reflect.TypeFor[rdbEndpointCustomResult](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Main: true, Path: "Success"},
			{Path: "Endpoints", Fields: []humanv2.Field{
				{Path: "ID"},
				{Path: "IP"},
				{Path: "Port"},
				{Path: "Name"},
				{Path: "Hostname"},
			}},
		},
	})

	humanv2.RegisterLayout(reflect.TypeFor[[]*rdb.Endpoint](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Main: true, Fields: []humanv2.Field{
				{Path: "ID"},
				{Path: "IP"},
				{Path: "Port"},
				{Path: "Name"},
				{Path: "Hostname"},
			}},
		},
	})

	humanv2.RegisterFormatFunc(func(r backupDownloadResult) (string, error) {
		size, err := humanv2.Format(r.Size)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf(
			"Backup downloaded to %s successfully (%s written)",
			r.FileName,
			size,
		), nil
	})

	humanv2.RegisterFormatFunc(func(s rdb.BackupSchedule) (string, error) {
		if s.Disabled {
			return "Disabled", nil
		}

		return "Enabled", nil
	})

	// TODO
	// func rdbEndpointCustomResultMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	// 	result := i.(rdbEndpointCustomResult)
	// 	messageStr, err := result.Success.MarshalHuman()
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	if len(result.Endpoints) == 0 {
	// 		return messageStr, nil
	// 	}

	// 	endpointsStr, err := human.Marshal(result.Endpoints, opt)
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	return messageStr + "\n" + endpointsStr, nil
	// }

	// func backupScheduleMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	// 	backupSchedule := i.(rdbSDK.BackupSchedule)

	// 	if opt.TableCell {
	// 		if backupSchedule.Disabled {
	// 			return "Disabled", nil
	// 		}

	// 		return "Enabled", nil
	// 	}

	// 	// To avoid recursion of human.Marshal we create a dummy type
	// 	type LocalBackupSchedule rdbSDK.BackupSchedule
	// 	type tmp struct {
	// 		LocalBackupSchedule
	// 		Frequency *scw.Duration `json:"frequency"`
	// 		Retention *scw.Duration `json:"retention"`
	// 	}

	// 	localBackupSchedule := tmp{
	// 		LocalBackupSchedule: LocalBackupSchedule(backupSchedule),
	// 		Frequency: &scw.Duration{
	// 			Seconds: int64(backupSchedule.Frequency) * 3600,
	// 		},
	// 		Retention: &scw.Duration{
	// 			Seconds: int64(backupSchedule.Retention) * 24 * 3600,
	// 		},
	// 	}

	// 	str, err := human.Marshal(localBackupSchedule, opt)
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	return str, nil
	// }

	// func logDownloadResultMarshalerFunc(i any, _ *human.MarshalOpt) (string, error) {
	// 	result := i.(logDownloadResult)
	// 	sizeStr, err := human.Marshal(result.Size, nil)
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	return fmt.Sprintf(
	// 		"Log downloaded to %s successfully (%s written)",
	// 		result.FileName,
	// 		sizeStr,
	// 	), nil
	// }
	//
	cmds.Merge(core.NewCommands(
		rdbConfigRoot(),
		rdbConfigGetCommand(),
		aclEditCommand(),
		backupDownloadCommand(),
		backupWaitCommand(),
		databaseGetURLCommand(),
		endpointListCommand(),
		engineSettingsCommand(),
		instanceConnectCommand(),
		instanceWaitCommand(),
		userGetURLCommand(),
		instanceEditSettingsCommand(),
	))
	cmds.MustFind("rdb", "acl", "add").Override(aclAddBuilder)
	cmds.MustFind("rdb", "acl", "delete").Override(aclDeleteBuilder)
	cmds.MustFind("rdb", "acl", "set").Override(aclSetBuilder)

	cmds.MustFind("rdb", "backup", "create").Override(backupCreateBuilder)
	cmds.MustFind("rdb", "backup", "export").Override(backupExportBuilder)
	cmds.MustFind("rdb", "backup", "list").Override(backupListBuilder)
	cmds.MustFind("rdb", "backup", "restore").Override(backupRestoreBuilder)

	cmds.MustFind("rdb", "endpoint", "create").Override(endpointCreateBuilder)
	cmds.MustFind("rdb", "endpoint", "delete").Override(endpointDeleteBuilder)
	cmds.MustFind("rdb", "endpoint", "get").Override(endpointGetBuilder)

	cmds.MustFind("rdb", "engine", "list").Override(engineListBuilder)

	cmds.MustFind("rdb", "instance", "create").Override(instanceCreateBuilder)
	cmds.MustFind("rdb", "instance", "clone").Override(instanceCloneBuilder)
	cmds.MustFind("rdb", "instance", "delete").Override(instanceDeleteBuilder)
	cmds.MustFind("rdb", "instance", "get").Override(instanceGetBuilder)
	cmds.MustFind("rdb", "instance", "update").Override(instanceUpdateBuilder)
	cmds.MustFind("rdb", "instance", "upgrade").Override(instanceUpgradeBuilder)

	// Make database create idempotent
	cmds.MustFind("rdb", "database", "create").Override(databaseCreateBuilder)

	cmds.MustFind("rdb", "user", "create").Override(userCreateBuilder)
	cmds.MustFind("rdb", "user", "list").Override(userListBuilder)
	cmds.MustFind("rdb", "user", "update").Override(userUpdateBuilder)

	cmds.MustFind("rdb", "log", "prepare").Override(logPrepareBuilder)
	cmds.Merge(core.NewCommands(
		logDownloadCommand(),
		instanceCockpitCommand(),
	))

	return cmds
}
