// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package k8s

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		k8sRoot(),
		k8sCluster(),
		k8sPool(),
		k8sNode(),
		k8sVersion(),
		k8sClusterList(),
		k8sClusterCreate(),
		k8sClusterGet(),
		k8sClusterUpdate(),
		k8sClusterDelete(),
		k8sClusterUpgrade(),
		k8sClusterListAvailableVersions(),
		k8sClusterResetAdminToken(),
		k8sPoolList(),
		k8sPoolCreate(),
		k8sPoolGet(),
		k8sPoolUpgrade(),
		k8sPoolUpdate(),
		k8sPoolDelete(),
		k8sNodeList(),
		k8sNodeGet(),
		k8sNodeReplace(),
		k8sNodeReboot(),
		k8sVersionList(),
		k8sVersionGet(),
	)
}
func k8sRoot() *core.Command {
	return &core.Command{
		Short:     `Kapsule API`,
		Long:      ``,
		Namespace: "k8s",
	}
}

func k8sCluster() *core.Command {
	return &core.Command{
		Short: `Kapsule cluster management commands`,
		Long: `A cluster is a fully managed Kubernetes cluster.

It is composed of different pools, each pool containing the same kind of nodes.
`,
		Namespace: "k8s",
		Resource:  "cluster",
	}
}

func k8sPool() *core.Command {
	return &core.Command{
		Short: `Kapsule pool management commands`,
		Long: `A pool is a set of identical Nodes. A pool has a name, a size (its current number of nodes), nodes number limits (min, max) and a Scaleway instance type.
Changing those limits increases/decreases the size of a pool. Thus, when autoscaling is enabled, the pool will grow or shrink inside those limits, depending on its load.
A "default pool" is automatically created with every cluster.
`,
		Namespace: "k8s",
		Resource:  "pool",
	}
}

func k8sNode() *core.Command {
	return &core.Command{
		Short: `Kapsule node management commands`,
		Long: `A node (short for worker node) is an abstraction for a Scaleway Instance.
It is part of a pool and is instantiated by Scaleway, making Kubernetes software installed and configured automatically on it.
Please note that Kubernetes nodes cannot be accessed with ssh.
`,
		Namespace: "k8s",
		Resource:  "node",
	}
}

func k8sVersion() *core.Command {
	return &core.Command{
		Short: `Available Kubernetes version commands`,
		Long: `A version is a vanilla Kubernetes version like ` + "`" + `x.y.z` + "`" + `.
It is composed of a major version x, a minor version y and a patch version z.
Scaleway's managed Kubernetes, Kapsule, will at least support the last patch version for the last three minor release.

Also each version have a different set of container runtimes, CNIs, ingresses, feature gates and admission plugins available.
`,
		Namespace: "k8s",
		Resource:  "version",
	}
}

func k8sClusterList() *core.Command {
	return &core.Command{
		Short:     `List all the clusters`,
		Long:      `This method allows to list all the existing Kubernetes clusters in an account.`,
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.ListClustersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `The project ID on which to filter the returned clusters`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `The sort order of the returned clusters`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "name_asc", "name_desc", "status_asc", "status_desc", "version_asc", "version_desc"},
			},
			{
				Name:       "name",
				Short:      `The name on which to filter the returned clusters`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status",
				Short:      `The status on which to filter the returned clusters`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "creating", "ready", "deleting", "deleted", "updating", "locked", "pool_required"},
			},
			{
				Name:       "type",
				Short:      `The type on which to filter the returned clusters`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `The organization ID on which to filter the returned clusters`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.ListClustersRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			resp, err := api.ListClusters(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Clusters, nil

		},
		Examples: []*core.Example{
			{
				Short:    "List all the clusters on your default region",
				ArgsJSON: `null`,
			},
			{
				Short:    "List the ready clusters on your default region",
				ArgsJSON: `{"status":"ready"}`,
			},
			{
				Short:    "List the clusters that match the given name on fr-par ('cluster1' will return 'cluster100' and 'cluster1' but not 'foo')",
				ArgsJSON: `{"name":"cluster1","region":"fr-par"}`,
			},
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Name",
			},
			{
				FieldName: "Status",
			},
			{
				FieldName: "Version",
			},
			{
				FieldName: "Region",
			},
			{
				FieldName: "ProjectID",
			},
			{
				FieldName: "Tags",
			},
			{
				FieldName: "Cni",
			},
			{
				FieldName: "Description",
			},
			{
				FieldName: "ClusterURL",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "Type",
			},
		}},
	}
}

func k8sClusterCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new cluster`,
		Long:      `This method allows to create a new Kubernetes cluster on an account.`,
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.CreateClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "type",
				Short:      `The type of the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `The name of the cluster`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("k8s"),
			},
			{
				Name:       "description",
				Short:      `The description of the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The tags associated with the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "version",
				Short:      `The Kubernetes version of the cluster`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cni",
				Short:      `The Container Network Interface (CNI) plugin that will run in the cluster`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_cni", "cilium", "calico", "weave", "flannel", "kilo"},
			},
			{
				Name:       "enable-dashboard",
				Short:      `The enablement of the Kubernetes Dashboard in the cluster`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "ingress",
				Short:      `The Ingress Controller that will run in the cluster`,
				Required:   false,
				Deprecated: true,
				Positional: false,
				EnumValues: []string{"unknown_ingress", "none", "nginx", "traefik", "traefik2"},
			},
			{
				Name:       "pools.{index}.name",
				Short:      `The name of the pool`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pools.{index}.node-type",
				Short:      `The node type is the type of Scaleway Instance wanted for the pool`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pools.{index}.placement-group-id",
				Short:      `The placement group ID in which all the nodes of the pool will be created`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pools.{index}.autoscaling",
				Short:      `The enablement of the autoscaling feature for the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pools.{index}.size",
				Short:      `The size (number of nodes) of the pool`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pools.{index}.min-size",
				Short:      `The minimun size of the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pools.{index}.max-size",
				Short:      `The maximum size of the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pools.{index}.container-runtime",
				Short:      `The container runtime for the nodes of the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_runtime", "docker", "containerd", "crio"},
			},
			{
				Name:       "pools.{index}.autohealing",
				Short:      `The enablement of the autohealing feature for the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pools.{index}.tags.{index}",
				Short:      `The tags associated with the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pools.{index}.kubelet-args.{key}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pools.{index}.upgrade-policy.max-unavailable",
				Short:      `The maximum number of nodes that can be not ready at the same time`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pools.{index}.upgrade-policy.max-surge",
				Short:      `The maximum number of nodes to be created during the upgrade`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pools.{index}.zone",
				Short:      `The Zone in which the Pool's node will be spawn in`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.scale-down-disabled",
				Short:      `Disable the cluster autoscaler`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.scale-down-delay-after-add",
				Short:      `How long after scale up that scale down evaluation resumes`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.estimator",
				Short:      `Type of resource estimator to be used in scale up`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_estimator", "binpacking"},
			},
			{
				Name:       "autoscaler-config.expander",
				Short:      `Type of node group expander to be used in scale up`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_expander", "random", "most_pods", "least_waste", "priority", "price"},
			},
			{
				Name:       "autoscaler-config.ignore-daemonsets-utilization",
				Short:      `Ignore DaemonSet pods when calculating resource utilization for scaling down`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.balance-similar-node-groups",
				Short:      `Detect similar node groups and balance the number of nodes between them`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.expendable-pods-priority-cutoff",
				Short:      `Pods with priority below cutoff will be expendable`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.scale-down-unneeded-time",
				Short:      `How long a node should be unneeded before it is eligible for scale down`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.scale-down-utilization-threshold",
				Short:      `Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale down`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.max-graceful-termination-sec",
				Short:      `Maximum number of seconds the cluster autoscaler waits for pod termination when trying to scale down a node`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-upgrade.enable",
				Short:      `Whether or not auto upgrade is enabled for the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-upgrade.maintenance-window.start-hour",
				Short:      `The start hour of the 2-hour maintenance window`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-upgrade.maintenance-window.day",
				Short:      `The day of the week for the maintenance window`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"any", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"},
			},
			{
				Name:       "feature-gates.{index}",
				Short:      `List of feature gates to enable`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "admission-plugins.{index}",
				Short:      `List of admission plugins to enable`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.issuer-url",
				Short:      `URL of the provider which allows the API server to discover public signing keys`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.client-id",
				Short:      `A client id that all tokens must be issued for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.username-claim",
				Short:      `JWT claim to use as the user name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.username-prefix",
				Short:      `Prefix prepended to username`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.groups-claim.{index}",
				Short:      `JWT claim to use as the user's group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.groups-prefix",
				Short:      `Prefix prepended to group claims`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.required-claim.{index}",
				Short:      `Multiple key=value pairs that describes a required claim in the ID Token`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "apiserver-cert-sans.{index}",
				Short:      `Additional Subject Alternative Names for the Kubernetes API server certificate`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.CreateClusterRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.CreateCluster(request)

		},
		Examples: []*core.Example{
			{
				Short: "Create a Kubernetes cluster named foo with cilium as CNI, in version 1.17.4 and with a pool named default composed of 3 DEV1-M",
				Raw:   `scw k8s cluster create name=foo version=1.17.4 pools.0.size=3 pools.0.node-type=DEV1-M pools.0.name=default`,
			},
			{
				Short: "Create a Kubernetes cluster named bar, tagged, calico as CNI, in version 1.17.4 and with a tagged pool named default composed of 2 RENDER-S and autohealing and autoscaling enabled (between 1 and 10 nodes)",
				Raw:   `scw k8s cluster create name=bar version=1.17.4 tags.0=tag1 tags.1=tag2 cni=cilium pools.0.size=2 pools.0.node-type=RENDER-S pools.0.min-size=1 pools.0.max-size=10 pools.0.autohealing=true pools.0.autoscaling=true pools.0.tags.0=pooltag1 pools.0.tags.1=pooltag2 pools.0.name=default`,
			},
		},
	}
}

func k8sClusterGet() *core.Command {
	return &core.Command{
		Short:     `Get a cluster`,
		Long:      `This method allows to get details about a specific Kubernetes cluster.`,
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.GetClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `The ID of the requested cluster`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.GetClusterRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.GetCluster(request)

		},
		Examples: []*core.Example{
			{
				Short: "Get a given cluster",
				Raw:   `scw k8s cluster get 11111111-1111-1111-111111111111`,
			},
		},
	}
}

func k8sClusterUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a cluster`,
		Long:      `This method allows to update a specific Kubernetes cluster. Note that this method is not made to upgrade a Kubernetes cluster.`,
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.UpdateClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `The ID of the cluster to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `The new name of the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `The new description of the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The new tags associated with the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.scale-down-disabled",
				Short:      `Disable the cluster autoscaler`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.scale-down-delay-after-add",
				Short:      `How long after scale up that scale down evaluation resumes`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.estimator",
				Short:      `Type of resource estimator to be used in scale up`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_estimator", "binpacking"},
			},
			{
				Name:       "autoscaler-config.expander",
				Short:      `Type of node group expander to be used in scale up`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_expander", "random", "most_pods", "least_waste", "priority", "price"},
			},
			{
				Name:       "autoscaler-config.ignore-daemonsets-utilization",
				Short:      `Ignore DaemonSet pods when calculating resource utilization for scaling down`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.balance-similar-node-groups",
				Short:      `Detect similar node groups and balance the number of nodes between them`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.expendable-pods-priority-cutoff",
				Short:      `Pods with priority below cutoff will be expendable`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.scale-down-unneeded-time",
				Short:      `How long a node should be unneeded before it is eligible for scale down`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.scale-down-utilization-threshold",
				Short:      `Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale down`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaler-config.max-graceful-termination-sec",
				Short:      `Maximum number of seconds the cluster autoscaler waits for pod termination when trying to scale down a node`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-dashboard",
				Short:      `The new value of the Kubernetes Dashboard enablement`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "ingress",
				Short:      `The new Ingress Controller for the cluster`,
				Required:   false,
				Deprecated: true,
				Positional: false,
				EnumValues: []string{"unknown_ingress", "none", "nginx", "traefik", "traefik2"},
			},
			{
				Name:       "auto-upgrade.enable",
				Short:      `Whether or not auto upgrade is enabled for the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-upgrade.maintenance-window.start-hour",
				Short:      `The start hour of the 2-hour maintenance window`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-upgrade.maintenance-window.day",
				Short:      `The day of the week for the maintenance window`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"any", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"},
			},
			{
				Name:       "feature-gates.{index}",
				Short:      `List of feature gates to enable`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "admission-plugins.{index}",
				Short:      `List of admission plugins to enable`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.issuer-url",
				Short:      `URL of the provider which allows the API server to discover public signing keys`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.client-id",
				Short:      `A client id that all tokens must be issued for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.username-claim",
				Short:      `JWT claim to use as the user name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.username-prefix",
				Short:      `Prefix prepended to username`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.groups-claim.{index}",
				Short:      `JWT claim to use as the user's group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.groups-prefix",
				Short:      `Prefix prepended to group claims`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "open-id-connect-config.required-claim.{index}",
				Short:      `Multiple key=value pairs that describes a required claim in the ID Token`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "apiserver-cert-sans.{index}",
				Short:      `Additional Subject Alternative Names for the Kubernetes API server certificate`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.UpdateClusterRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.UpdateCluster(request)

		},
		Examples: []*core.Example{
			{
				Short: "Enable dashboard on a given cluster",
				Raw:   `scw k8s cluster update 11111111-1111-1111-111111111111 enable-dashboard=true`,
			},
			{
				Short: "Add TTLAfterFinished and ServiceNodeExclusion as feature gates on a given cluster",
				Raw:   `scw k8s cluster update 11111111-1111-1111-111111111111 feature-gates.0=TTLAfterFinished feature-gates.1=ServiceNodeExclusion`,
			},
		},
	}
}

func k8sClusterDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a cluster`,
		Long:      `This method allows to delete a specific cluster and all its associated pools and nodes. Note that this method will not delete any Load Balancers or Block Volumes that are associated with the cluster.`,
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.DeleteClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `The ID of the cluster to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "with-additional-resources",
				Short:      `Set true if you want to delete all volumes (including retain volume type) and loadbalancers whose name start with cluster ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.DeleteClusterRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.DeleteCluster(request)

		},
		Examples: []*core.Example{
			{
				Short: "Delete a given cluster",
				Raw:   `scw k8s cluster delete 11111111-1111-1111-111111111111`,
			},
		},
	}
}

func k8sClusterUpgrade() *core.Command {
	return &core.Command{
		Short:     `Upgrade a cluster`,
		Long:      `This method allows to upgrade a specific Kubernetes cluster and/or its associated pools to a specific and supported Kubernetes version.`,
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "upgrade",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.UpgradeClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `The ID of the cluster to upgrade`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "version",
				Short:      `The new Kubernetes version of the cluster`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "upgrade-pools",
				Short:      `The enablement of the pools upgrade`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.UpgradeClusterRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.UpgradeCluster(request)

		},
		Examples: []*core.Example{
			{
				Short: "Upgrade a given cluster to Kubernetes version 1.17.4 (without upgrading the pools)",
				Raw:   `scw k8s cluster upgrade 11111111-1111-1111-111111111111 version=1.17.4`,
			},
			{
				Short: "Upgrade a given cluster to Kubernetes version 1.17.4 (and upgrade the pools)",
				Raw:   `scw k8s cluster upgrade 11111111-1111-1111-111111111111 version=1.17.4 upgrade-pools=true`,
			},
		},
	}
}

func k8sClusterListAvailableVersions() *core.Command {
	return &core.Command{
		Short:     `List available versions for a cluster`,
		Long:      `This method allows to list the versions that a specific Kubernetes cluster is allowed to upgrade to. Note that it will be every patch version greater than the actual one as well a one minor version ahead of the actual one. Upgrades skipping a minor version will not work.`,
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "list-available-versions",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.ListClusterAvailableVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `The ID of the cluster which the available Kuberentes versions will be listed from`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.ListClusterAvailableVersionsRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.ListClusterAvailableVersions(request)

		},
		Examples: []*core.Example{
			{
				Short: "List all available versions for a given cluster to upgrade to",
				Raw:   `scw k8s cluster list-available-versions 11111111-1111-1111-111111111111`,
			},
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "Name",
			},
			{
				FieldName: "Label",
			},
			{
				FieldName: "AvailableIngresses",
			},
			{
				FieldName: "AvailableContainerRuntimes",
			},
		}},
	}
}

func k8sClusterResetAdminToken() *core.Command {
	return &core.Command{
		Short:     `Reset the admin token of a cluster`,
		Long:      `This method allows to reset the admin token for a specific Kubernetes cluster. This will invalidate the old admin token (which will not be usable after) and create a new one. Note that the redownload of the kubeconfig will be necessary to keep interacting with the cluster (if the old admin token was used).`,
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "reset-admin-token",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.ResetClusterAdminTokenRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `The ID of the cluster of which the admin token will be renewed`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.ResetClusterAdminTokenRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			e = api.ResetClusterAdminToken(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "cluster",
				Verb:     "reset-admin-token",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short: "Reset the admin token for a given cluster",
				Raw:   `scw k8s cluster reset-admin-token 11111111-1111-1111-111111111111`,
			},
		},
	}
}

func k8sPoolList() *core.Command {
	return &core.Command{
		Short:     `List all the pools in a cluster`,
		Long:      `This method allows to list all the existing pools for a specific Kubernetes cluster.`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.ListPoolsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `The ID of the cluster from which the pools will be listed from`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `The sort order of the returned pools`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "name_asc", "name_desc", "status_asc", "status_desc", "version_asc", "version_desc"},
			},
			{
				Name:       "name",
				Short:      `The name on which to filter the returned pools`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status",
				Short:      `The status on which to filter the returned pools`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "ready", "deleting", "deleted", "scaling", "warning", "locked", "upgrading"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.ListPoolsRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			resp, err := api.ListPools(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Pools, nil

		},
		Examples: []*core.Example{
			{
				Short: "List all pools for a given cluster",
				Raw:   `scw k8s pool list cluster-id=11111111-1111-1111-111111111111`,
			},
			{
				Short: "List all scaling pools for a given cluster",
				Raw:   `scw k8s pool list cluster-id=11111111-1111-1111-111111111111 status=scaling`,
			},
			{
				Short: "List all pools for a given cluster that contain the word foo in the pool name",
				Raw:   `scw k8s pool list cluster-id=11111111-1111-1111-111111111111 name=foo`,
			},
			{
				Short: "List all pools for a given cluster and order them by ascending creation date",
				Raw:   `scw k8s pool list cluster-id=11111111-1111-1111-111111111111 order-by=created_at_asc`,
			},
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Name",
			},
			{
				FieldName: "Status",
			},
			{
				FieldName: "NodeType",
			},
			{
				FieldName: "Size",
			},
			{
				FieldName: "MinSize",
			},
			{
				FieldName: "MaxSize",
			},
			{
				FieldName: "Autoscaling",
			},
			{
				FieldName: "Autohealing",
			},
			{
				FieldName: "Version",
			},
			{
				FieldName: "Tags",
			},
			{
				FieldName: "ContainerRuntime",
			},
			{
				FieldName: "ClusterID",
			},
			{
				FieldName: "Zone",
			},
			{
				FieldName: "Region",
			},
			{
				FieldName: "KubeletArgs",
			},
			{
				FieldName: "PlacementGroupID",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "UpdatedAt",
			},
		}},
	}
}

func k8sPoolCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new pool in a cluster`,
		Long:      `This method allows to create a new pool in a specific Kubernetes cluster.`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.CreatePoolRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `The ID of the cluster in which the pool will be created`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `The name of the pool`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("pool"),
			},
			{
				Name:       "node-type",
				Short:      `The node type is the type of Scaleway Instance wanted for the pool`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "placement-group-id",
				Short:      `The placement group ID in which all the nodes of the pool will be created`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autoscaling",
				Short:      `The enablement of the autoscaling feature for the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "size",
				Short:      `The size (number of nodes) of the pool`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "min-size",
				Short:      `The minimun size of the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "max-size",
				Short:      `The maximum size of the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "container-runtime",
				Short:      `The container runtime for the nodes of the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_runtime", "docker", "containerd", "crio"},
			},
			{
				Name:       "autohealing",
				Short:      `The enablement of the autohealing feature for the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The tags associated with the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "kubelet-args.{key}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "upgrade-policy.max-unavailable",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "upgrade-policy.max-surge",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "zone",
				Short:      `The Zone in which the Pool's node will be spawn in`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.CreatePoolRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.CreatePool(request)

		},
		Examples: []*core.Example{
			{
				Short: "Create a pool named bar with 2 DEV1-XL on a given cluster",
				Raw:   `scw k8s pool create cluster-id=11111111-1111-1111-111111111111 name=bar node-type=DEV1-XL size=2`,
			},
			{
				Short: "Create a pool named fish with 5 GP1-L with autoscaling enabled within 0 and 10 nodes, autohealing enabled, and containerd as the container runtime on a given cluster",
				Raw:   `scw k8s pool create cluster-id=11111111-1111-1111-111111111111 name=fish node-type=GP1-L size=5 min-size=0 max-size=10 autoscaling=true autohealing=true container-runtime=containerd`,
			},
			{
				Short: "Create a tagged pool named turtle with 1 GP1-S which is using the already created placement group 2222222222222-2222-222222222222 for all the nodes in the pool on a given cluster",
				Raw:   `scw k8s pool create cluster-id=11111111-1111-1111-111111111111 name=turtle node-type=GP1-S size=1 placement-group-id=2222222222222-2222-222222222222 tags.0=turtle tags.1=placement-group`,
			},
		},
	}
}

func k8sPoolGet() *core.Command {
	return &core.Command{
		Short:     `Get a pool in a cluster`,
		Long:      `This method allows to get details about a specific pool.`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.GetPoolRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pool-id",
				Short:      `The ID of the requested pool`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.GetPoolRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.GetPool(request)

		},
		Examples: []*core.Example{
			{
				Short: "Get a given pool",
				Raw:   `scw k8s pool get 11111111-1111-1111-111111111111`,
			},
		},
	}
}

func k8sPoolUpgrade() *core.Command {
	return &core.Command{
		Short:     `Upgrade a pool in a cluster`,
		Long:      `This method allows to upgrade the Kubernetes version of a specific pool. Note that this will work when the targeted version is the same than the version of the cluster.`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "upgrade",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.UpgradePoolRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pool-id",
				Short:      `The ID of the pool to upgrade`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "version",
				Short:      `The new Kubernetes version for the pool`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.UpgradePoolRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.UpgradePool(request)

		},
		Examples: []*core.Example{
			{
				Short: "Upgrade a given pool to the Kubernetes version 1.17.4",
				Raw:   `scw k8s pool upgrade 11111111-1111-1111-111111111111 version=1.17.4`,
			},
		},
	}
}

func k8sPoolUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a pool in a cluster`,
		Long:      `This method allows to update some attributes of a specific pool such as the size, the autoscaling enablement, the tags, ...`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.UpdatePoolRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pool-id",
				Short:      `The ID of the pool to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "autoscaling",
				Short:      `The new value for the enablement of autoscaling for the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "size",
				Short:      `The new size for the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "min-size",
				Short:      `The new minimun size for the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "max-size",
				Short:      `The new maximum size for the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "autohealing",
				Short:      `The new value for the enablement of autohealing for the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The new tags associated with the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "kubelet-args.value.{key}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "upgrade-policy.max-unavailable",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "upgrade-policy.max-surge",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.UpdatePoolRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.UpdatePool(request)

		},
		Examples: []*core.Example{
			{
				Short: "Enable autoscaling on a given pool",
				Raw:   `scw k8s pool update 11111111-1111-1111-111111111111 autoscaling=true`,
			},
			{
				Short: "Reduce the size and max size of a given pool to 4",
				Raw:   `scw k8s pool update 11111111-1111-1111-111111111111 size=4 max-size=4`,
			},
			{
				Short: "Change the tags of the given pool",
				Raw:   `scw k8s pool update 11111111-1111-1111-111111111111 tags.0=my tags.1=new tags.2=pool`,
			},
		},
	}
}

func k8sPoolDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a pool in a cluster`,
		Long:      `This method allows to delete a specific pool from a cluster, deleting all the nodes associated with it.`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.DeletePoolRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pool-id",
				Short:      `The ID of the pool to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.DeletePoolRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.DeletePool(request)

		},
		Examples: []*core.Example{
			{
				Short: "Delete a given pool",
				Raw:   `scw k8s pool delete 11111111-1111-1111-111111111111`,
			},
		},
	}
}

func k8sNodeList() *core.Command {
	return &core.Command{
		Short:     `List all the nodes in a cluster`,
		Long:      `This method allows to list all the existing nodes for a specific Kubernetes cluster.`,
		Namespace: "k8s",
		Resource:  "node",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.ListNodesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `The cluster ID from which the nodes will be listed from`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pool-id",
				Short:      `The pool ID on which to filter the returned nodes`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `The sort order of the returned nodes`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "name",
				Short:      `The name on which to filter the returned nodes`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status",
				Short:      `The status on which to filter the returned nodes`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "creating", "not_ready", "ready", "deleting", "deleted", "locked", "rebooting", "creation_error", "upgrading"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.ListNodesRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			resp, err := api.ListNodes(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Nodes, nil

		},
		Examples: []*core.Example{
			{
				Short: "List all the nodes in the given cluster",
				Raw:   `scw k8s node list 11111111-1111-1111-111111111111`,
			},
			{
				Short: "List all the nodes in the pool 2222222222222-2222-222222222222 in the given cluster",
				Raw:   `scw k8s node list 11111111-1111-1111-111111111111 pool-id=2222222222222-2222-222222222222`,
			},
			{
				Short: "List all ready nodes in the given cluster",
				Raw:   `scw k8s node list 11111111-1111-1111-111111111111 status=ready`,
			},
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Name",
			},
			{
				FieldName: "Status",
			},
			{
				FieldName: "PublicIPV4",
			},
			{
				FieldName: "PublicIPV6",
			},
			{
				FieldName: "PoolID",
			},
			{
				FieldName: "ClusterID",
			},
			{
				FieldName: "Region",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "UpdatedAt",
			},
		}},
	}
}

func k8sNodeGet() *core.Command {
	return &core.Command{
		Short:     `Get a node in a cluster`,
		Long:      `This method allows to get details about a specific Kubernetes node.`,
		Namespace: "k8s",
		Resource:  "node",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.GetNodeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "node-id",
				Short:      `The ID of the requested node`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.GetNodeRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.GetNode(request)

		},
		Examples: []*core.Example{
			{
				Short: "Get a given node",
				Raw:   `scw k8s node get 11111111-1111-1111-111111111111`,
			},
		},
	}
}

func k8sNodeReplace() *core.Command {
	return &core.Command{
		Short:     `Replace a node in a cluster`,
		Long:      `This method allows to replace a specific node. The node will be set cordoned, meaning that scheduling will be disabled. Then the existing pods on the node will be drained and reschedule onto another schedulable node. Then the node will be deleted, and a new one will be created after the deletion. Note that when there is not enough space to reschedule all the pods (in a one node cluster for instance), you may experience some disruption of your applications.`,
		Namespace: "k8s",
		Resource:  "node",
		Verb:      "replace",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.ReplaceNodeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "node-id",
				Short:      `The ID of the node to replace`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.ReplaceNodeRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.ReplaceNode(request)

		},
		Examples: []*core.Example{
			{
				Short: "Replace a given node",
				Raw:   `scw k8s node replace 11111111-1111-1111-111111111111`,
			},
		},
	}
}

func k8sNodeReboot() *core.Command {
	return &core.Command{
		Short:     `Reboot a node in a cluster`,
		Long:      `This method allows to reboot a specific node. This node will frist be cordoned, meaning that scheduling will be disabled. Then the existing pods on the node will be drained and reschedule onto another schedulable node. Note that when there is not enough space to reschedule all the pods (in a one node cluster for instance), you may experience some disruption of your applications.`,
		Namespace: "k8s",
		Resource:  "node",
		Verb:      "reboot",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.RebootNodeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "node-id",
				Short:      `The ID of the node to reboot`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.RebootNodeRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.RebootNode(request)

		},
		Examples: []*core.Example{
			{
				Short: "Reboot a given node",
				Raw:   `scw k8s node reboot 11111111-1111-1111-111111111111`,
			},
		},
	}
}

func k8sVersionList() *core.Command {
	return &core.Command{
		Short:     `List all available versions`,
		Long:      `This method allows to list all available versions for the creation of a new Kubernetes cluster.`,
		Namespace: "k8s",
		Resource:  "version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.ListVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.ListVersionsRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.ListVersions(request)

		},
		Examples: []*core.Example{
			{
				Short: "List all available Kubernetes version in Kapsule",
				Raw:   `scw k8s version list`,
			},
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "Name",
			},
			{
				FieldName: "AvailableCnis",
			},
			{
				FieldName: "AvailableIngresses",
			},
			{
				FieldName: "AvailableContainerRuntimes",
			},
			{
				FieldName: "AvailableFeatureGates",
			},
			{
				FieldName: "AvailableAdmissionPlugins",
			},
			{
				FieldName: "AvailableKubeletArgs",
			},
		}},
	}
}

func k8sVersionGet() *core.Command {
	return &core.Command{
		Short:     `Get details about a specific version`,
		Long:      `This method allows to get a specific Kubernetes version and the details about the version.`,
		Namespace: "k8s",
		Resource:  "version",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(k8s.GetVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "version-name",
				Short:      `The requested version name`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.GetVersionRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.GetVersion(request)

		},
		Examples: []*core.Example{
			{
				Short: "Get the Kubernetes version 1.18.0",
				Raw:   `scw k8s version get 1.18.0`,
			},
		},
	}
}
