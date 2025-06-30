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

	cmds.Merge(core.NewCommands(clusterWaitCommand()))
	cmds.MustFind("redis", "cluster", "create").Override(clusterCreateBuilder)
	cmds.MustFind("redis", "cluster", "delete").Override(clusterDeleteBuilder)
	cmds.MustFind("redis", "acl", "add").Override(ACLAddListBuilder)
	cmds.MustFind("redis", "setting", "add").Override(redisSettingAddBuilder)
	cmds.MustFind("redis", "cluster", "migrate").Override(redisClusterMigrateBuilder)

	cmds.Merge(core.NewCommands(redisACLUpdateCommand()))

	return cmds
}
