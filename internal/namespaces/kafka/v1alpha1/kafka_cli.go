// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package kafka

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	kafka "github.com/scaleway/scaleway-sdk-go/api/kafka/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		kafkaRoot(),
		kafkaNodeType(),
		kafkaVersion(),
		kafkaCluster(),
		kafkaEndpoint(),
		kafkaUsers(),
		kafkaNodeTypeList(),
		kafkaVersionList(),
		kafkaClusterList(),
		kafkaClusterGet(),
		kafkaClusterCreate(),
		kafkaClusterUpdate(),
		kafkaClusterDelete(),
		kafkaClusterGetCa(),
		kafkaClusterRenewCa(),
		kafkaUsersList(),
		kafkaUsersUpdate(),
	)
}

func kafkaRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your Clusters for Apache Kafka®.`,
		Long:      `This API allows you to manage your Clusters for Apache Kafka®.`,
		Namespace: "kafka",
	}
}

func kafkaNodeType() *core.Command {
	return &core.Command{
		Short: `Node types management commands`,
		Long: `Two node type ranges are available:

* **Shared:** a complete and highly reliable node range with shared resources, made for scaling from development to production needs, at affordable prices.
* **Dedicated:** Kafka nodes with dedicated vCPU for the most demanding workloads and mission-critical applications.`,
		Namespace: "kafka",
		Resource:  "node-type",
	}
}

func kafkaVersion() *core.Command {
	return &core.Command{
		Short:     `Kafka version management commands`,
		Long:      `A version of Apache Kafka®.`,
		Namespace: "kafka",
		Resource:  "version",
	}
}

func kafkaCluster() *core.Command {
	return &core.Command{
		Short:     `Cluster management commands`,
		Long:      `A Kafka cluster is composed of one or multiple dedicated compute nodes running a single Kafka broker.`,
		Namespace: "kafka",
		Resource:  "cluster",
	}
}

func kafkaEndpoint() *core.Command {
	return &core.Command{
		Short:     `Endpoints management commands`,
		Long:      `Cluster endpoints enable connection to your cluster.`,
		Namespace: "kafka",
		Resource:  "endpoint",
	}
}

func kafkaUsers() *core.Command {
	return &core.Command{
		Short:     `Kafka users management commands`,
		Long:      `Kafka users enable authentication to your cluster.`,
		Namespace: "kafka",
		Resource:  "users",
	}
}

func kafkaNodeTypeList() *core.Command {
	return &core.Command{
		Short:     `List available node types`,
		Long:      `List available node types.`,
		Namespace: "kafka",
		Resource:  "node-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(kafka.ListNodeTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "include-disabled-types",
				Short:      `Defines whether or not to include disabled types`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*kafka.ListNodeTypesRequest)

			client := core.ExtractClient(ctx)
			api := kafka.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListNodeTypes(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.NodeTypes, nil
		},
	}
}

func kafkaVersionList() *core.Command {
	return &core.Command{
		Short:     `List Kafka versions`,
		Long:      `List all available versions of Kafka at the current time.`,
		Namespace: "kafka",
		Resource:  "version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(kafka.ListVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "version",
				Short:      `Kafka version to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*kafka.ListVersionsRequest)

			client := core.ExtractClient(ctx)
			api := kafka.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListVersions(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Versions, nil
		},
	}
}

func kafkaClusterList() *core.Command {
	return &core.Command{
		Short:     `List Kafka clusters`,
		Long:      `List all Kafka clusters in the specified region. By default, the clusters returned in the list are ordered by creation date in ascending order, though this can be modified via the order_by field. You can define additional parameters for your query, such as ` + "`" + `tags` + "`" + ` and ` + "`" + `name` + "`" + `. For the ` + "`" + `name` + "`" + ` parameter, the value you include will be checked against the whole name string to see if it includes the string you put in the parameter.`,
		Namespace: "kafka",
		Resource:  "cluster",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(kafka.ListClustersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tags.{index}",
				Short:      `List Kafka cluster with a given tag`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Lists Kafka clusters that match a name pattern`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Criteria to use when ordering Kafka cluster listings`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"name_asc",
					"name_desc",
					"status_asc",
					"status_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID of the Kafka cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*kafka.ListClustersRequest)

			client := core.ExtractClient(ctx)
			api := kafka.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListClusters(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Clusters, nil
		},
	}
}

func kafkaClusterGet() *core.Command {
	return &core.Command{
		Short:     `Get a Kafka cluster`,
		Long:      `Retrieve information about a given Kafka cluster, specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `cluster_id` + "`" + ` parameters. Its full details, including name, status, IP address and port, are returned in the response object.`,
		Namespace: "kafka",
		Resource:  "cluster",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(kafka.GetClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Kafka Cluster`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*kafka.GetClusterRequest)

			client := core.ExtractClient(ctx)
			api := kafka.NewAPI(client)

			return api.GetCluster(request)
		},
	}
}

func kafkaClusterCreate() *core.Command {
	return &core.Command{
		Short:     `Create a Kafka cluster`,
		Long:      `Create a new Kafka cluster.`,
		Namespace: "kafka",
		Resource:  "cluster",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(kafka.CreateClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the Kafka cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("kafk"),
			},
			{
				Name:       "version",
				Short:      `Version of Kafka`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to apply to the Kafka cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-amount",
				Short:      `Number of nodes to use for the Kafka cluster`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-type",
				Short:      `Type of node to use for the Kafka cluster`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volume.size-bytes",
				Short:      `Volume size`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volume.type",
				Short:      `Type of volume where data is stored`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"sbs_5k",
					"sbs_15k",
				},
			},
			{
				Name:       "endpoints.{index}.private-network.private-network-id",
				Short:      `UUID of the Private Network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-name",
				Short:      `Username for the kafka user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password for the kafka user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mono-az.zone",
				Short:      `Zone is the zone on which the cluster nodes are deployed.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*kafka.CreateClusterRequest)

			client := core.ExtractClient(ctx)
			api := kafka.NewAPI(client)

			return api.CreateCluster(request)
		},
	}
}

func kafkaClusterUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a Kafka cluster`,
		Long:      `Update the parameters of a Kafka cluster.`,
		Namespace: "kafka",
		Resource:  "cluster",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(kafka.UpdateClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Kafka Clusters to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the Kafka Cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of a Kafka Cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "version",
				Short:      `Version of Kafka`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*kafka.UpdateClusterRequest)

			client := core.ExtractClient(ctx)
			api := kafka.NewAPI(client)

			return api.UpdateCluster(request)
		},
	}
}

func kafkaClusterDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a Kafka cluster`,
		Long:      `Delete a given Kafka cluster, specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `cluster_id` + "`" + ` parameters. Deleting a Kafka cluster is permanent, and cannot be undone. Note that upon deletion all your data will be lost.`,
		Namespace: "kafka",
		Resource:  "cluster",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(kafka.DeleteClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Kafka Cluster to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*kafka.DeleteClusterRequest)

			client := core.ExtractClient(ctx)
			api := kafka.NewAPI(client)

			return api.DeleteCluster(request)
		},
	}
}

func kafkaClusterGetCa() *core.Command {
	return &core.Command{
		Short:     `Get a Kafka cluster's certificate authority`,
		Long:      `Retrieve certificate authority for a given Kafka cluster, specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `cluster_id` + "`" + ` parameters. The response object contains the certificate in PEM format. The certificate is required to validate the sever from the client side during TLS connection.`,
		Namespace: "kafka",
		Resource:  "cluster",
		Verb:      "get-ca",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(kafka.GetClusterCertificateAuthorityRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Kafka Cluster`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*kafka.GetClusterCertificateAuthorityRequest)

			client := core.ExtractClient(ctx)
			api := kafka.NewAPI(client)

			return api.GetClusterCertificateAuthority(request)
		},
	}
}

func kafkaClusterRenewCa() *core.Command {
	return &core.Command{
		Short:     `Renew the Kafka cluster's certificate authority`,
		Long:      `Request to renew the certificate authority for a given Kafka cluster, specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `cluster_id` + "`" + ` parameters. The certificate authority will be renewed within a few minutes.`,
		Namespace: "kafka",
		Resource:  "cluster",
		Verb:      "renew-ca",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(kafka.RenewClusterCertificateAuthorityRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Kafka Cluster`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*kafka.RenewClusterCertificateAuthorityRequest)

			client := core.ExtractClient(ctx)
			api := kafka.NewAPI(client)
			e = api.RenewClusterCertificateAuthority(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "cluster",
				Verb:     "renew-ca",
			}, nil
		},
	}
}

func kafkaUsersList() *core.Command {
	return &core.Command{
		Short:     `List kafka resources`,
		Long:      `List kafka resources.`,
		Namespace: "kafka",
		Resource:  "users",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(kafka.ListUsersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"name_asc",
					"name_desc",
				},
			},
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cluster-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*kafka.ListUsersRequest)

			client := core.ExtractClient(ctx)
			api := kafka.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListUsers(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Users, nil
		},
	}
}

func kafkaUsersUpdate() *core.Command {
	return &core.Command{
		Short:     `Update kafka resources`,
		Long:      `Update kafka resources.`,
		Namespace: "kafka",
		Resource:  "users",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(kafka.UpdateUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `ID of the cluster in which to update the user's password`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Username of the Kafka cluster user`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `New password for the Kafka cluster user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*kafka.UpdateUserRequest)

			client := core.ExtractClient(ctx)
			api := kafka.NewAPI(client)

			return api.UpdateUser(request)
		},
	}
}
