package redis

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("redis").Groups = []string{"database"}

	human.RegisterMarshalerFunc(redis.Cluster{}, redisClusterGetMarshalerFunc)

	cmds.Merge(core.NewCommands(clusterWaitCommand(), clusterConnectCommand()))
	cmds.MustFind("redis", "cluster", "create").Override(clusterCreateBuilder)
	cmds.MustFind("redis", "cluster", "delete").Override(clusterDeleteBuilder)
	cmds.MustFind("redis", "acl", "add").Override(ACLAddListBuilder)
	cmds.MustFind("redis", "setting", "add").Override(redisSettingAddBuilder)
	cmds.MustFind("redis", "cluster", "migrate").Override(redisClusterMigrateBuilder)

	cmds.Merge(core.NewCommands(redisACLUpdateCommand()))
	cmds.Merge(core.NewCommands(redisVersionSettingsCommand()))

	cmds.MustFind("redis", "cluster", "list").
		Override(core.ListViewBuilder("ID", "Name", "ProjectID", "Status", "Version", "Endpoints", "Tags", "NodeType", "CreatedAt", "UpdatedAt", "TLSEnabled", "ClusterSettings", "ACLRules", "ClusterSize", "Zone", "UserName", "UpgradableVersions"))
	cmds.MustFind("redis", "node-type", "list").
		Override(core.ListViewBuilder("Name", "StockStatus", "Description", "Vcpus", "Memory", "Disabled", "Beta", "Zone"))
	cmds.MustFind("redis", "version", "list").
		Override(core.ListViewBuilder("Version", "EndOfLifeAt", "AvailableSettings", "LogoURL", "ReleasedAt", "Zone"))

	return cmds
}
