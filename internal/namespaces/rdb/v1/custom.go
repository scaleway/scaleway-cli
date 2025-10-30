package rdb

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("rdb").Groups = []string{"database"}

	human.RegisterMarshalerFunc(rdb.Instance{}, instanceMarshalerFunc)
	human.RegisterMarshalerFunc(rdb.BackupSchedule{}, backupScheduleMarshalerFunc)
	human.RegisterMarshalerFunc(backupDownloadResult{}, backupResultMarshallerFunc)
	human.RegisterMarshalerFunc(CreateInstanceResult{}, createInstanceResultMarshalerFunc)
	human.RegisterMarshalerFunc(CustomACLResult{}, rdbACLCustomResultMarshalerFunc)
	human.RegisterMarshalerFunc(rdbEndpointCustomResult{}, rdbEndpointCustomResultMarshalerFunc)

	human.RegisterMarshalerFunc(
		rdb.InstanceStatus(""),
		human.EnumMarshalFunc(instanceStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		rdb.DatabaseBackupStatus(""),
		human.EnumMarshalFunc(backupStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		rdb.InstanceLogStatus(""),
		human.EnumMarshalFunc(logStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		rdb.NodeTypeStock(""),
		human.EnumMarshalFunc(nodeTypeStockMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		rdb.ACLRuleAction(""),
		human.EnumMarshalFunc(aclRuleActionMarshalSpecs),
	)

	cmds.Merge(core.NewCommands(
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

	return cmds
}
