package k8s

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
)

// GetCommands returns cluster commands.
//
// This function:
// - Gets the generated commands
// - Register handwritten marshalers
// - Apply handwritten overrides (of Command.Run)
func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(k8s.ClusterStatus(0), human.BindAttributesMarshalFunc(clusterStatusAttributes))
	human.RegisterMarshalerFunc(k8s.PoolStatus(0), human.BindAttributesMarshalFunc(poolStatusAttributes))
	human.RegisterMarshalerFunc(k8s.NodeStatus(0), human.BindAttributesMarshalFunc(nodeStatusAttributes))

	cmds.MustFind("k8s", "cluster", "list-available-versions").Override(clusterAvailableVersionsListBuilder)
	cmds.MustFind("k8s", "cluster", "create").Override(clusterCreateBuilder)
	cmds.MustFind("k8s", "cluster", "update").Override(clusterUpdateBuilder)
	cmds.MustFind("k8s", "cluster", "upgrade").Override(clusterUpgradeBuilder)
	cmds.MustFind("k8s", "cluster", "delete").Override(clusterDeleteBuilder)

	cmds.MustFind("k8s", "pool", "create").Override(poolCreateBuilder)
	cmds.MustFind("k8s", "pool", "update").Override(poolUpdateBuilder)
	cmds.MustFind("k8s", "pool", "upgrade").Override(poolUpgradeBuilder)
	cmds.MustFind("k8s", "pool", "delete").Override(poolDeleteBuilder)

	cmds.MustFind("k8s", "node", "reboot").Override(nodeRebootBuilder)

	return cmds
}
