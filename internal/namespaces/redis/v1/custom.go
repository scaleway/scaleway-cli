package redis

import (
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("redis").Groups = []string{"database"}

	humanv2.RegisterLayout(reflect.TypeFor[redis.Cluster](), clusterLayout)

	cmds.Merge(core.NewCommands(clusterWaitCommand(), clusterConnectCommand()))
	cmds.MustFind("redis", "cluster", "create").Override(clusterCreateBuilder)
	cmds.MustFind("redis", "cluster", "delete").Override(clusterDeleteBuilder)
	cmds.MustFind("redis", "acl", "add").Override(ACLAddListBuilder)
	cmds.MustFind("redis", "setting", "add").Override(redisSettingAddBuilder)
	cmds.MustFind("redis", "cluster", "migrate").Override(redisClusterMigrateBuilder)

	cmds.Merge(core.NewCommands(redisACLUpdateCommand()))
	cmds.Merge(core.NewCommands(redisVersionSettingsCommand()))

	cmds.Merge(core.NewCommands(clusterCockpitCommand()))

	return cmds
}
