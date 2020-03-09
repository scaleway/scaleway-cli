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

	cmds.MustFind("k8s", "cluster", "list-available-versions").Override(clusterAvailableVersionsListBuilder)
	cmds.MustFind("k8s", "cluster", "create").Override(clusterCreateBuilder)
	cmds.MustFind("k8s", "cluster", "update").Override(clusterUpdateBuilder)
	cmds.MustFind("k8s", "cluster", "upgrade").Override(clusterUpgradeBuilder)
	cmds.MustFind("k8s", "cluster", "delete").Override(clusterDeleteBuilder)

	return cmds
}
