package rdb

import (
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("rdb").Groups = []string{"database"}

	humanv2.RegisterLayout(reflect.TypeFor[rdb.Instance](), instanceLayout)
	humanv2.RegisterFormatFunc(backupScheduleView)
	humanv2.RegisterFormatFunc(createInstanceResultView)
	human.RegisterMarshalerFunc(rdb.BackupSchedule{}, backupScheduleMarshalerFunc)
	human.RegisterMarshalerFunc(backupDownloadResult{}, backupResultMarshallerFunc)
	human.RegisterMarshalerFunc(CustomACLResult{}, rdbACLCustomResultMarshalerFunc)
	human.RegisterMarshalerFunc(rdbEndpointCustomResult{}, rdbEndpointCustomResultMarshalerFunc)
	human.RegisterMarshalerFunc(logDownloadResult{}, logDownloadResultMarshalerFunc)
	humanv2.RegisterFormatFunc(rdbACLResultView)
	humanv2.RegisterFormatFunc(rdbEndpointCustomResultView)
	humanv2.RegisterFormatFunc(backupDownloadResultView)

	humanv2.RegisterEnumFormat(instanceStatusFormat)
	humanv2.RegisterEnumFormat(backupStatusFormat)
	humanv2.RegisterEnumFormat(logStatusFormat)
	humanv2.RegisterEnumFormat(nodeTypeStockFormat)
	humanv2.RegisterEnumFormat(aclRuleActionFormat)
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
	cmds.MustFind("rdb", "instance", "list").Override(instanceListBuilder)

	cmds.MustFind("rdb", "acl", "add").Override(aclAddBuilder)
	cmds.MustFind("rdb", "acl", "delete").Override(aclDeleteBuilder)
	cmds.MustFind("rdb", "acl", "set").Override(aclSetBuilder)

	cmds.MustFind("rdb", "backup", "create").Override(backupCreateBuilder)
	cmds.MustFind("rdb", "backup", "export").Override(backupExportBuilder)
	cmds.MustFind("rdb", "backup", "list").Override(backupListBuilder)
	cmds.MustFind("rdb", "backup", "restore").Override(backupRestoreBuilder)

	cmds.MustFind("rdb", "endpoint", "create").Override(endpointCreateBuilder)
	cmds.MustFind("rdb", "endpoint", "delete").Override(endpointDeleteBuilder)
	cmds.MustFind("rdb", "endpoint", "list").Override(endpointListBuilder)
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
