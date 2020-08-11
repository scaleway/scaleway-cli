package rdb

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

var nodeTypes = []string{
	"DB-DEV-S",
	"DB-DEV-M",
	"DB-DEV-L",
	"DB-DEV-XL",
	"DB-GP-XS",
	"DB-GP-S",
	"DB-GP-M",
	"DB-GP-L",
	"DB-GP-XL",
}

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(rdb.Instance{}, instanceMarshalerFunc)
	human.RegisterMarshalerFunc(rdb.BackupSchedule{}, backupScheduleMarshalerFunc)

	human.RegisterMarshalerFunc(rdb.InstanceStatus(0), human.EnumMarshalFunc(instanceStatusMarshalSpecs))
	human.RegisterMarshalerFunc(rdb.DatabaseBackupStatus(0), human.EnumMarshalFunc(backupStatusMarshalSpecs))
	human.RegisterMarshalerFunc(rdb.InstanceLogStatus(0), human.EnumMarshalFunc(logStatusMarshalSpecs))
	human.RegisterMarshalerFunc(rdb.NodeTypeStock(0), human.EnumMarshalFunc(nodeTypeStockMarshalSpecs))
	human.RegisterMarshalerFunc(rdb.ACLRuleAction(0), human.EnumMarshalFunc(aclRuleActionMarshalSpecs))

	cmds.Merge(core.NewCommands(
		instanceWaitCommand(),
		instanceConnectCommand(),
		backupWaitCommand(),
	))
	cmds.MustFind("rdb", "backup", "create").Override(backupCreateBuilder)
	cmds.MustFind("rdb", "backup", "export").Override(backupExportBuilder)
	cmds.MustFind("rdb", "backup", "restore").Override(backupRestoreBuilder)

	cmds.MustFind("rdb", "instance", "create").Override(instanceCreateBuilder)
	cmds.MustFind("rdb", "instance", "clone").Override(instanceCloneBuilder)
	cmds.MustFind("rdb", "instance", "create").Override(instanceCreateBuilder)
	cmds.MustFind("rdb", "instance", "upgrade").Override(instanceUpgradeBuilder)

	cmds.MustFind("rdb", "engine", "list").Override(engineListBuilder)

	return cmds
}
