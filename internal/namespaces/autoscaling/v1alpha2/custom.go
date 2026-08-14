package autoscaling

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	autoscaling "github.com/scaleway/scaleway-sdk-go/api/autoscaling/v1alpha2"
)

func GetCommands() *core.Commands {
	commands := GetGeneratedCommands()

	human.RegisterMarshalerFunc(autoscaling.Group{}, groupMarshalerFunc)

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
