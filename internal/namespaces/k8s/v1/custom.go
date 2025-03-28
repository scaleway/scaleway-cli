package k8s

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
)

// GetCommands returns cluster commands.
//
// This function:
// - Gets the generated commands
// - Register handwritten marshalers
// - Apply handwritten overrides (of Command.Run)
func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("k8s").Groups = []string{"container"}

	cmds.Merge(core.NewCommands(
		k8sExecCredentialCommand(),
		k8sKubeconfigCommand(),
		k8sKubeconfigGetCommand(),
		k8sKubeconfigInstallCommand(),
		k8sKubeconfigUninstallCommand(),
		k8sClusterWaitCommand(),
		k8sNodeWaitCommand(),
		k8sPoolWaitCommand(),
	))

	human.RegisterMarshalerFunc(k8s.Version{}, versionMarshalerFunc)
	human.RegisterMarshalerFunc(k8s.Cluster{}, clusterMarshalerFunc)
	human.RegisterMarshalerFunc(k8s.ClusterStatus(""), human.EnumMarshalFunc(clusterStatusMarshalSpecs))
	human.RegisterMarshalerFunc(k8s.PoolStatus(""), human.EnumMarshalFunc(poolStatusMarshalSpecs))
	human.RegisterMarshalerFunc(k8s.NodeStatus(""), human.EnumMarshalFunc(nodeStatusMarshalSpecs))
	human.RegisterMarshalerFunc(k8s.ListClusterAvailableTypesResponse{}, clusterAvailableTypesListMarshalerFunc)

	cmds.MustFind("k8s", "cluster", "list-available-versions").Override(clusterAvailableVersionsListBuilder)
	cmds.MustFind("k8s", "cluster", "create").Override(clusterCreateBuilder)
	cmds.MustFind("k8s", "cluster", "get").Override(clusterGetBuilder)
	cmds.MustFind("k8s", "cluster", "update").Override(clusterUpdateBuilder)
	cmds.MustFind("k8s", "cluster", "upgrade").Override(clusterUpgradeBuilder)
	cmds.MustFind("k8s", "cluster", "delete").Override(clusterDeleteBuilder)
	cmds.MustFind("k8s", "pool", "create").Override(poolCreateBuilder)
	cmds.MustFind("k8s", "pool", "update").Override(poolUpdateBuilder)
	cmds.MustFind("k8s", "pool", "upgrade").Override(poolUpgradeBuilder)
	cmds.MustFind("k8s", "pool", "delete").Override(poolDeleteBuilder)

	cmds.MustFind("k8s", "node", "reboot").Override(nodeRebootBuilder)

	cmds.MustFind("k8s", "version", "list").Override(versionListBuilder)

	return cmds
}
