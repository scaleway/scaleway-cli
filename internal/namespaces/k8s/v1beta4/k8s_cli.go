// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package k8s

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
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
		k8sClusterList(),
		k8sClusterCreate(),
		k8sClusterGet(),
		k8sClusterUpdate(),
		k8sClusterDelete(),
		k8sClusterUpgrade(),
		k8sClusterListAvailableVersions(),
		k8sClusterResetAdminToken(),
	)
}
func k8sRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your Kapsule clusters`,
		Long:      ``,
		Namespace: "k8s",
	}
}

func k8sCluster() *core.Command {
	return &core.Command{
		Short: `A cluster is a Kubernetes Kapsule cluster`,
		Long: `A cluster is a fully managed Kubernetes cluster.

It is composed of different pools, each pool containing the same kind of nodes.
`,
		Namespace: "k8s",
		Resource:  "cluster",
	}
}

func k8sClusterList() *core.Command {
	return &core.Command{
		Short:     `List all the clusters`,
		Long:      `This method allows to list all the existing Kubernetes clusters in an account.`,
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(k8s.ListClustersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `The sort order of the returned clusters`,
				Required:   false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "name_asc", "name_desc", "status_asc", "status_desc", "version_asc", "version_desc"},
			},
			{
				Name:     "name",
				Short:    `The name on which to filter the returned clusters`,
				Required: false,
			},
			{
				Name:       "status",
				Short:      `The status on which to filter the returned clusters`,
				Required:   false,
				EnumValues: []string{"unknown", "creating", "ready", "deleting", "deleted", "updating", "warning", "error", "locked"},
			},
			{
				Name:     "organization-id",
				Short:    `The organization ID on which to filter the returned clusters`,
				Required: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
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
				Short:   "List all the clusters on your default region",
				Request: `null`,
			},
			{
				Short:   "List the ready clusters on your default region",
				Request: `{"status":"ready"}`,
			},
			{
				Short:   "List the clusters that match the given name on fr-par ('cluster1' will return 'cluster100' and 'cluster1' but not 'foo')",
				Request: `{"name":"cluster1","region":"fr-par"}`,
			},
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "id",
			},
			{
				FieldName: "name",
			},
			{
				FieldName: "status",
			},
			{
				FieldName: "version",
			},
			{
				FieldName: "region",
			},
			{
				FieldName: "organization_id",
			},
			{
				FieldName: "tags",
			},
			{
				FieldName: "cni",
			},
			{
				FieldName: "description",
			},
			{
				FieldName: "cluster_url",
			},
			{
				FieldName: "created_at",
			},
			{
				FieldName: "updated_at",
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
		ArgsType:  reflect.TypeOf(k8s.CreateClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "name",
				Short:    `The name of the cluster`,
				Required: true,
			},
			{
				Name:     "description",
				Short:    `The description of the cluster`,
				Required: false,
			},
			{
				Name:     "tags.{index}",
				Short:    `The tags associated with the cluster`,
				Required: false,
			},
			{
				Name:     "version",
				Short:    `The Kubernetes version of the cluster`,
				Required: true,
			},
			{
				Name:       "cni",
				Short:      `The Container Network Interface (CNI) plugin that will run in the cluster`,
				Required:   true,
				EnumValues: []string{"unknown_cni", "cilium", "calico", "weave", "flannel"},
			},
			{
				Name:     "enable-dashboard",
				Short:    `The enablement of the Kubernetes Dashboard in the cluster`,
				Required: false,
			},
			{
				Name:       "ingress",
				Short:      `The Ingress Controller that will run in the cluster`,
				Required:   false,
				EnumValues: []string{"unknown_ingress", "none", "nginx", "traefik"},
			},
			{
				Name:     "default-pool-config.node-type",
				Short:    `The node type is the type of Scaleway Instance wanted for the pool`,
				Required: true,
			},
			{
				Name:     "default-pool-config.placement-group-id",
				Short:    `The placement group ID in which all the nodes of the pool will be created`,
				Required: false,
			},
			{
				Name:     "default-pool-config.autoscaling",
				Short:    `The enablement of the autoscaling feature for the pool`,
				Required: false,
			},
			{
				Name:     "default-pool-config.size",
				Short:    `The size (number of nodes) of the pool`,
				Required: true,
			},
			{
				Name:     "default-pool-config.min-size",
				Short:    `The minimun size of the pool`,
				Required: false,
			},
			{
				Name:     "default-pool-config.max-size",
				Short:    `The maximum size of the pool`,
				Required: false,
			},
			{
				Name:       "default-pool-config.container-runtime",
				Short:      `The container runtime for the nodes of the pool`,
				Required:   false,
				EnumValues: []string{"unknown_runtime", "docker", "containerd", "crio"},
			},
			{
				Name:     "default-pool-config.autohealing",
				Short:    `The enablement of the autohealing feature for the pool`,
				Required: false,
			},
			{
				Name:     "default-pool-config.tags.{index}",
				Short:    `The tags associated with the pool`,
				Required: false,
			},
			{
				Name:     "autoscaler-config.scale-down-disabled",
				Short:    `Disable the cluster autoscaler`,
				Required: false,
			},
			{
				Name:     "autoscaler-config.scale-down-delay-after-add",
				Short:    `How long after scale up that scale down evaluation resumes`,
				Required: false,
			},
			{
				Name:       "autoscaler-config.estimator",
				Short:      `Type of resource estimator to be used in scale up`,
				Required:   false,
				EnumValues: []string{"unknown_estimator", "binpacking", "oldbinpacking"},
			},
			{
				Name:       "autoscaler-config.expander",
				Short:      `Type of node group expander to be used in scale up`,
				Required:   false,
				EnumValues: []string{"unknown_expander", "random", "most_pods", "least_waste", "priority"},
			},
			{
				Name:     "autoscaler-config.ignore-daemonsets-utilization",
				Short:    `Ignore DaemonSet pods when calculating resource utilization for scaling down`,
				Required: false,
			},
			{
				Name:     "autoscaler-config.balance-similar-node-groups",
				Short:    `Detect similar node groups and balance the number of nodes between them`,
				Required: false,
			},
			{
				Name:     "autoscaler-config.expendable-pods-priority-cutoff",
				Short:    `Pods with priority below cutoff will be expendable`,
				Required: false,
			},
			{
				Name:     "auto-upgrade.enable",
				Short:    `Whether or not auto upgrade is enabled for the cluster`,
				Required: false,
			},
			{
				Name:     "auto-upgrade.maintenance-window.start-hour",
				Short:    `The start hour of the 2-hour maintenance window`,
				Required: false,
			},
			{
				Name:       "auto-upgrade.maintenance-window.day",
				Short:      `The day of the week for the maintenance window`,
				Required:   false,
				EnumValues: []string{"any", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"},
			},
			{
				Name:     "feature-gates.{index}",
				Short:    `List of feature gates to enable`,
				Required: false,
			},
			{
				Name:     "admission-plugins.{index}",
				Short:    `List of admission plugins to enable`,
				Required: false,
			},
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.CreateClusterRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.CreateCluster(request)

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
		ArgsType:  reflect.TypeOf(k8s.GetClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "cluster-id",
				Short:    `The ID of the requested cluster`,
				Required: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.GetClusterRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.GetCluster(request)

		},
		Examples: []*core.Example{
			{
				Short:   "Get the cluster with id 11111111-1111-1111-111111111111",
				Request: `{"cluster-id":"11111111-1111-1111-111111111111"}`,
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
		ArgsType:  reflect.TypeOf(k8s.UpdateClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "cluster-id",
				Short:    `The ID of the cluster to update`,
				Required: true,
			},
			{
				Name:     "name",
				Short:    `The new name of the cluster`,
				Required: false,
			},
			{
				Name:     "description",
				Short:    `The new description of the cluster`,
				Required: false,
			},
			{
				Name:     "tags",
				Short:    `The new tags associated with the cluster`,
				Required: false,
			},
			{
				Name:     "autoscaler-config.scale-down-disabled",
				Short:    `Disable the cluster autoscaler`,
				Required: false,
			},
			{
				Name:     "autoscaler-config.scale-down-delay-after-add",
				Short:    `How long after scale up that scale down evaluation resumes`,
				Required: false,
			},
			{
				Name:       "autoscaler-config.estimator",
				Short:      `Type of resource estimator to be used in scale up`,
				Required:   false,
				EnumValues: []string{"unknown_estimator", "binpacking", "oldbinpacking"},
			},
			{
				Name:       "autoscaler-config.expander",
				Short:      `Type of node group expander to be used in scale up`,
				Required:   false,
				EnumValues: []string{"unknown_expander", "random", "most_pods", "least_waste", "priority"},
			},
			{
				Name:     "autoscaler-config.ignore-daemonsets-utilization",
				Short:    `Ignore DaemonSet pods when calculating resource utilization for scaling down`,
				Required: false,
			},
			{
				Name:     "autoscaler-config.balance-similar-node-groups",
				Short:    `Detect similar node groups and balance the number of nodes between them`,
				Required: false,
			},
			{
				Name:     "autoscaler-config.expendable-pods-priority-cutoff",
				Short:    `Pods with priority below cutoff will be expendable`,
				Required: false,
			},
			{
				Name:     "enable-dashboard",
				Short:    `The new value of the Kubernetes Dashboard enablement`,
				Required: false,
			},
			{
				Name:       "ingress",
				Short:      `The new Ingress Controller for the cluster`,
				Required:   false,
				EnumValues: []string{"unknown_ingress", "none", "nginx", "traefik"},
			},
			{
				Name:     "auto-upgrade.enable",
				Short:    `Whether or not auto upgrade is enabled for the cluster`,
				Required: false,
			},
			{
				Name:     "auto-upgrade.maintenance-window.start-hour",
				Short:    `The start hour of the 2-hour maintenance window`,
				Required: false,
			},
			{
				Name:       "auto-upgrade.maintenance-window.day",
				Short:      `The day of the week for the maintenance window`,
				Required:   false,
				EnumValues: []string{"any", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"},
			},
			{
				Name:     "feature-gates",
				Short:    `List of feature gates to enable`,
				Required: false,
			},
			{
				Name:     "admission-plugins",
				Short:    `List of admission plugins to enable`,
				Required: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.UpdateClusterRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.UpdateCluster(request)

		},
		Examples: []*core.Example{
			{
				Short:   "Enable dashboard on cluster 11111111-1111-1111-111111111111",
				Request: `{"cluster-id":"11111111-1111-1111-111111111111","enable-dashboard":true}`,
			},
			{
				Short:   "Add TTLAfterFinished Feature gates on cluster 11111111-1111-1111-111111111111",
				Request: `{"cluster-id":"11111111-1111-1111-111111111111","feature-gates":["TTLAfterFinished"]}`,
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
		ArgsType:  reflect.TypeOf(k8s.DeleteClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "cluster-id",
				Short:    `The ID of the cluster to delete`,
				Required: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.DeleteClusterRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.DeleteCluster(request)

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
		ArgsType:  reflect.TypeOf(k8s.UpgradeClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "cluster-id",
				Short:    `The ID of the cluster to upgrade`,
				Required: true,
			},
			{
				Name:     "version",
				Short:    `The new Kubernetes version of the cluster`,
				Required: false,
			},
			{
				Name:     "upgrade-pools",
				Short:    `The enablement of the pools upgrade`,
				Required: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.UpgradeClusterRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.UpgradeCluster(request)

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
		ArgsType:  reflect.TypeOf(k8s.ListClusterAvailableVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "cluster-id",
				Short:    `The ID of the cluster which the available Kuberentes versions will be listed from`,
				Required: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.ListClusterAvailableVersionsRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			return api.ListClusterAvailableVersions(request)

		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "name",
			},
			{
				FieldName: "label",
			},
			{
				FieldName: "available_ingresses",
			},
			{
				FieldName: "available_container_runtimes",
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
		ArgsType:  reflect.TypeOf(k8s.ResetClusterAdminTokenRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "cluster-id",
				Short:    `The ID of the cluster of which the admin token will be renewed`,
				Required: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*k8s.ResetClusterAdminTokenRequest)

			client := core.ExtractClient(ctx)
			api := k8s.NewAPI(client)
			e = api.ResetClusterAdminToken(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{}, nil
		},
	}
}
