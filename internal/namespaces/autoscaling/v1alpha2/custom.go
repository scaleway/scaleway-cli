package autoscaling

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	commands := GetGeneratedCommands()

	commands.MustFind("autoscaling", "group", "create").Override(GroupCreateBuilder)
	commands.MustFind("autoscaling", "group", "update").Override(GroupUpdateBuilder)
	commands.MustFind("autoscaling", "group", "list").Override(GroupListBuilder)
	commands.MustFind("autoscaling", "group", "get").Override(GroupGetBuilder)
	commands.MustFind("autoscaling", "group", "delete").Override(GroupDeleteBuilder)

	commands.MustFind("autoscaling", "alerts", "list").Override(AlertsListBuilder)
	commands.MustFind("autoscaling", "logs", "list").Override(LogsListBuilder)
	commands.MustFind("autoscaling", "servers", "list").Override(ServersListBuilder)

	return commands
}
