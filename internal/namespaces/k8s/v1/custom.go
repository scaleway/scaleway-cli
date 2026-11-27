package k8s

import (
	"context"
	"errors"
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
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
		clusterConnectCommand(),
		k8sClusterWaitCommand(),
		k8sNodeWaitCommand(),
		k8sPoolWaitCommand(),
		k8sPoolAddExternalNodeCommand(),
		k8sPoolSetTaintCommand(),
		k8sPoolRemoveTaintCommand(),
		k8sPoolSetStartupTaintCommand(),
		k8sPoolRemoveStartupTaintCommand(),
		k8sPoolSetLabelCommand(),
		k8sPoolRemoveLabelCommand(),
	))

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[k8s.ClusterStatus]{
		k8s.ClusterStatusCreating:     humanv2.EnumFormat{Color: color.FgBlue},
		k8s.ClusterStatusReady:        humanv2.EnumFormat{Color: color.FgGreen},
		k8s.ClusterStatusPoolRequired: humanv2.EnumFormat{Color: color.FgRed},
		k8s.ClusterStatusLocked:       humanv2.EnumFormat{Color: color.FgRed},
		k8s.ClusterStatusUpdating:     humanv2.EnumFormat{Color: color.FgBlue},
	})

	humanv2.RegisterLayout(reflect.TypeFor[k8s.Cluster](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "AutoscalerConfig", Title: "Autoscaler configuration"},
			{Path: "AutoUpgrade", Title: "Auto-upgrade settings"},
			{Path: "OpenIDConnectConfig", Title: "Open ID Connect configuration"},
		},
	})

	humanv2.RegisterLayout(reflect.TypeFor[k8s.ListClusterAvailableTypesResponse](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "ClusterTypes", Title: "Cluster Types"},
		},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[k8s.NodeStatus]{
		k8s.NodeStatusCreating:      humanv2.EnumFormat{Color: color.FgBlue},
		k8s.NodeStatusRebooting:     humanv2.EnumFormat{Color: color.FgBlue},
		k8s.NodeStatusReady:         humanv2.EnumFormat{Color: color.FgGreen},
		k8s.NodeStatusNotReady:      humanv2.EnumFormat{Color: color.FgYellow},
		k8s.NodeStatusCreationError: humanv2.EnumFormat{Color: color.FgRed},
		k8s.NodeStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[k8s.PoolStatus]{
		k8s.PoolStatusScaling:   humanv2.EnumFormat{Color: color.FgBlue},
		k8s.PoolStatusReady:     humanv2.EnumFormat{Color: color.FgGreen},
		k8s.PoolStatusLocked:    humanv2.EnumFormat{Color: color.FgRed},
		k8s.PoolStatusUpgrading: humanv2.EnumFormat{Color: color.FgBlue},
		k8s.PoolStatusWarning:   humanv2.EnumFormat{Color: color.FgHiYellow},
	})

	humanv2.RegisterLayout(reflect.TypeFor[k8s.Pool](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "UpgradePolicy", Title: "Upgrade Policy"},
			{Path: "Labels", Title: "Node labels"},
			{Path: "Taints", Title: "Node taints"},
			{Path: "StartupTaints", Title: "Node startup taints"},
		},
	})

	humanv2.RegisterLayout(reflect.TypeFor[k8s.Version](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "AvailableKubeletArgs", Title: "Available Kubelet Arguments"},
			{Path: "AvailableCnis", Title: "Available CNIs"},
			{Path: "AvailableFeatureGates", Title: "Available Feature Gates"},
			{Path: "AvailableAdmissionPlugins", Title: "Available Admission Plugins"},
		},
	})

	cmds.MustFind("k8s", "cluster", "list-available-versions").
		Override(clusterAvailableVersionsListBuilder)
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

	cmds.Merge(core.NewCommands(clusterCockpitCommand()))

	return cmds
}

func extractSecretKey(ctx context.Context) (string, error) {
	config, _ := scw.LoadConfigFromPath(core.ExtractConfigPath(ctx))
	profileName := core.ExtractProfileName(ctx)

	switch {
	// Environment variable check
	case core.ExtractEnv(ctx, scw.ScwSecretKeyEnv) != "":
		return core.ExtractEnv(ctx, scw.ScwSecretKeyEnv), nil
	// There is no config file
	case config == nil:
		return "", errors.New("config not provided")
	// Config file with profile name
	case config.Profiles[profileName] != nil && config.Profiles[profileName].SecretKey != nil:
		return *config.Profiles[profileName].SecretKey, nil
	// Default config
	case config.SecretKey != nil:
		return *config.SecretKey, nil
	}

	return "", errors.New("unable to find secret key")
}
