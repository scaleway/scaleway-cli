package kafka

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("kafka", "cluster", "create").Override(clusterCreateBuilder)

	return cmds
}
