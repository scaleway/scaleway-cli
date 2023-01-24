package rdb

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(rdb.Instance{}, instanceMarshalerFunc)
	human.RegisterMarshalerFunc(rdb.BackupSchedule{}, backupScheduleMarshalerFunc)
	human.RegisterMarshalerFunc(backupDownloadResult{}, backupResultMarshallerFunc)

	human.RegisterMarshalerFunc(rdb.InstanceStatus(""), human.EnumMarshalFunc(instanceStatusMarshalSpecs))
	human.RegisterMarshalerFunc(rdb.DatabaseBackupStatus(""), human.EnumMarshalFunc(backupStatusMarshalSpecs))
	human.RegisterMarshalerFunc(rdb.InstanceLogStatus(""), human.EnumMarshalFunc(logStatusMarshalSpecs))
	human.RegisterMarshalerFunc(rdb.NodeTypeStock(""), human.EnumMarshalFunc(nodeTypeStockMarshalSpecs))
	human.RegisterMarshalerFunc(rdb.ACLRuleAction(""), human.EnumMarshalFunc(aclRuleActionMarshalSpecs))

	cmds.Merge(core.NewCommands(
		instanceWaitCommand(),
		instanceConnectCommand(),
		backupWaitCommand(),
		backupDownloadCommand(),
		engineSettingsCommand(),
	))
	cmds.MustFind("rdb", "acl", "add").Override(aclAddBuilder)
	cmds.MustFind("rdb", "acl", "delete").Override(aclDeleteBuilder)

	cmds.MustFind("rdb", "backup", "create").Override(backupCreateBuilder)
	cmds.MustFind("rdb", "backup", "export").Override(backupExportBuilder)
	cmds.MustFind("rdb", "backup", "restore").Override(backupRestoreBuilder)

	cmds.MustFind("rdb", "instance", "create").Override(instanceCreateBuilder)
	cmds.MustFind("rdb", "instance", "clone").Override(instanceCloneBuilder)
	cmds.MustFind("rdb", "instance", "upgrade").Override(instanceUpgradeBuilder)
	cmds.MustFind("rdb", "instance", "update").Override(instanceUpdateBuilder)

	cmds.MustFind("rdb", "engine", "list").Override(engineListBuilder)

	cmds.MustFind("rdb", "user", "list").Override(userListBuilder)

	cmds.MustFind("rdb", "backup", "list").Override(backupListBuilder)

	return cmds
}
