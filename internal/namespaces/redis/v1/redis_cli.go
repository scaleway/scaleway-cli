// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package redis

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		redisRoot(),
		redisCluster(),
		redisNodeType(),
		redisVersion(),
		redisSetting(),
		redisACL(),
		redisEndpoint(),
		redisClusterCreate(),
		redisClusterUpdate(),
		redisClusterGet(),
		redisClusterList(),
		redisClusterMigrate(),
		redisClusterDelete(),
		redisClusterMetrics(),
		redisNodeTypeList(),
		redisVersionList(),
		redisClusterGetCertificate(),
		redisClusterRenewCertificate(),
		redisSettingAdd(),
		redisSettingDelete(),
		redisSettingSet(),
		redisACLSet(),
		redisACLAdd(),
		redisACLDelete(),
		redisACLGet(),
		redisEndpointSet(),
		redisEndpointAdd(),
		redisEndpointDelete(),
		redisEndpointGet(),
		redisEndpointUpdate(),
	)
}
func redisRoot() *core.Command {
	return &core.Command{
		Short:     `Managed Database for Redis™ API`,
		Long:      ``,
		Namespace: "redis",
	}
}

func redisCluster() *core.Command {
	return &core.Command{
		Short:     `Cluster management commands`,
		Long:      `A Redis™ Database Instance, also known as a Redis™ cluster, consists of either one standalone node or a cluster composed of three to six nodes. The cluster uses partitioning to split the keyspace. Each partition is replicated and can be reassigned or elected as the primary when necessary. Standalone mode creates a standalone database provisioned on a single node.`,
		Namespace: "redis",
		Resource:  "cluster",
	}
}

func redisNodeType() *core.Command {
	return &core.Command{
		Short:     `Node Types management commands`,
		Long:      `Nodes are the compute units that make up your Redis™ Database Instance. Different node types are available with varying amounts of RAM and vCPU.`,
		Namespace: "redis",
		Resource:  "node-type",
	}
}

func redisVersion() *core.Command {
	return &core.Command{
		Short:     `Redis™ version management commands`,
		Long:      `The Redis™ database engine versions available at Scaleway for your clusters.`,
		Namespace: "redis",
		Resource:  "version",
	}
}

func redisSetting() *core.Command {
	return &core.Command{
		Short: `Settings management commands`,
		Long: `Advanced settings allow you to tune the behavior of your Redis™ database engine to better fit your needs. Available settings depend on the version of the Redis™ engine. Note that some settings can only be defined upon the Redis™ engine initialization. These are called init settings. You can find a full list of the settings available in the response body of the [list available Redis™ versions](#path-redistm-engine-versions-list-available-redistm-versions) endpoint.

Each advanced setting entry has a default value that users can override. The deletion of a setting entry will restore the setting to default value. Some of the defaults values can be different from the engine's defaults, as we optimize them to the Scaleway platform.`,
		Namespace: "redis",
		Resource:  "setting",
	}
}

func redisACL() *core.Command {
	return &core.Command{
		Short:     `Access Control List (ACL) management commands`,
		Long:      `Network Access Control Lists (ACLs) allow you to manage network inbound traffic by setting up ACL rules.`,
		Namespace: "redis",
		Resource:  "acl",
	}
}

func redisEndpoint() *core.Command {
	return &core.Command{
		Short:     `Endpoints management commands`,
		Long:      `Manage endpoint access to your Redis™ Database Instance through Public or Private Networks.`,
		Namespace: "redis",
		Resource:  "endpoint",
	}
}

func redisClusterCreate() *core.Command {
	return &core.Command{
		Short:     `Create a Redis™ Database Instance`,
		Long:      `Create a new Redis™ Database Instance (Redis™ cluster). You must set the ` + "`" + `zone` + "`" + `, ` + "`" + `project_id` + "`" + `, ` + "`" + `version` + "`" + `, ` + "`" + `node_type` + "`" + `, ` + "`" + `user_name` + "`" + ` and ` + "`" + `password` + "`" + ` parameters. Optionally you can define ` + "`" + `acl_rules` + "`" + `, ` + "`" + `endpoints` + "`" + `, ` + "`" + `tls_enabled` + "`" + ` and ` + "`" + `cluster_settings` + "`" + `.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.CreateClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the Database Instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("ins"),
			},
			{
				Name:       "version",
				Short:      `Redis™ engine version of the Database Instance`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to apply to the Database Instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-type",
				Short:      `Type of node to use for the Database Instance`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-name",
				Short:      `Name of the user created upon Database Instance creation`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password of the user`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cluster-size",
				Short:      `Number of nodes in the Redis™ cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acl-rules.{index}.ip-cidr",
				Short:      `IPv4 network address of the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acl-rules.{index}.description",
				Short:      `Description of the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.id",
				Short:      `UUID of the Private Network to connect to the Database Instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.service-ips.{index}",
				Short:      `Endpoint IPv4 address with a CIDR notation. You must provide at least one IPv4 per node.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tls-enabled",
				Short:      `Defines whether or not TLS is enabled`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cluster-settings.{index}.value",
				Short:      `Value of the setting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cluster-settings.{index}.name",
				Short:      `Name of the setting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.CreateClusterRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.CreateCluster(request)

		},
	}
}

func redisClusterUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a Redis™ Database Instance`,
		Long:      `Update the parameters of a Redis™ Database Instance (Redis™ cluster), including ` + "`" + `name` + "`" + `, ` + "`" + `tags` + "`" + `, ` + "`" + `user_name` + "`" + ` and ` + "`" + `password` + "`" + `.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.UpdateClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Database Instance to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the Database Instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Database Instance tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-name",
				Short:      `Name of the Database Instance user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password of the Database Instance user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.UpdateClusterRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.UpdateCluster(request)

		},
	}
}

func redisClusterGet() *core.Command {
	return &core.Command{
		Short:     `Get a Redis™ Database Instance`,
		Long:      `Retrieve information about a Redis™ Database Instance (Redis™ cluster). Specify the ` + "`" + `cluster_id` + "`" + ` and ` + "`" + `region` + "`" + ` in your request to get information such as ` + "`" + `id` + "`" + `, ` + "`" + `status` + "`" + `, ` + "`" + `version` + "`" + `, ` + "`" + `tls_enabled` + "`" + `, ` + "`" + `cluster_settings` + "`" + `, ` + "`" + `upgradable_versions` + "`" + ` and ` + "`" + `endpoints` + "`" + ` about your cluster in the response.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.GetClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.GetClusterRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.GetCluster(request)

		},
	}
}

func redisClusterList() *core.Command {
	return &core.Command{
		Short:     `List Redis™ Database Instances`,
		Long:      `List all Redis™ Database Instances (Redis™ cluster) in the specified zone. By default, the Database Instances returned in the list are ordered by creation date in ascending order, though this can be modified via the order_by field. You can define additional parameters for your query, such as ` + "`" + `tags` + "`" + `, ` + "`" + `name` + "`" + `, ` + "`" + `organization_id` + "`" + ` and ` + "`" + `version` + "`" + `.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.ListClustersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tags.{index}",
				Short:      `Filter by Database Instance tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter by Database Instance names`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Criteria to use when ordering the list`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "project-id",
				Short:      `Filter by Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "version",
				Short:      `Filter by Redis™ engine version`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter by Organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2, scw.Zone(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.ListClustersRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListClusters(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Clusters, nil

		},
	}
}

func redisClusterMigrate() *core.Command {
	return &core.Command{
		Short:     `Scale up a Redis™ Database Instance`,
		Long:      `Upgrade your standalone Redis™ Database Instance node, either by upgrading to a bigger node type (vertical scaling) or by adding more nodes to your Database Instance to increase your number of endpoints and distribute cache (horizontal scaling). Note that scaling horizontally your Redis™ Database Instance will not renew its TLS certificate. In order to refresh the TLS certificate, you must use the Renew TLS certificate endpoint.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "migrate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.MigrateClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Database Instance to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "version",
				Short:      `Redis™ engine version of the Database Instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-type",
				Short:      `Type of node to use for the Database Instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cluster-size",
				Short:      `Number of nodes for the Database Instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.MigrateClusterRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.MigrateCluster(request)

		},
	}
}

func redisClusterDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a Redis™ Database Instance`,
		Long:      `Delete a Redis™ Database Instance (Redis™ cluster), specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `cluster_id` + "`" + ` parameters. Deleting a Database Instance is permanent, and cannot be undone. Note that upon deletion all your data will be lost.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.DeleteClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Database Instance to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.DeleteClusterRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.DeleteCluster(request)

		},
	}
}

func redisClusterMetrics() *core.Command {
	return &core.Command{
		Short:     `Get metrics of a Redis™ Database Instance`,
		Long:      `Retrieve the metrics of a Redis™ Database Instance (Redis™ cluster). You can define the period from which to retrieve metrics by specifying the ` + "`" + `start_date` + "`" + ` and ` + "`" + `end_date` + "`" + `.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "metrics",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.GetClusterMetricsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "start-at",
				Short:      `Start date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "end-at",
				Short:      `End date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "metric-name",
				Short:      `Name of the metric to gather`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.GetClusterMetricsRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.GetClusterMetrics(request)

		},
	}
}

func redisNodeTypeList() *core.Command {
	return &core.Command{
		Short:     `List available node types`,
		Long:      `List all available node types. By default, the node types returned in the list are ordered by creation date in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "redis",
		Resource:  "node-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.ListNodeTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "include-disabled-types",
				Short:      `Defines whether or not to include disabled types`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2, scw.Zone(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.ListNodeTypesRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListNodeTypes(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.NodeTypes, nil

		},
	}
}

func redisVersionList() *core.Command {
	return &core.Command{
		Short:     `List available Redis™ versions`,
		Long:      `List the Redis™ database engine versions available. You can define additional parameters for your query, such as ` + "`" + `include_disabled` + "`" + `, ` + "`" + `include_beta` + "`" + `, ` + "`" + `include_deprecated` + "`" + ` and ` + "`" + `version` + "`" + `.`,
		Namespace: "redis",
		Resource:  "version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.ListClusterVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "include-disabled",
				Short:      `Defines whether or not to include disabled Redis™ engine versions`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "include-beta",
				Short:      `Defines whether or not to include beta Redis™ engine versions`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "include-deprecated",
				Short:      `Defines whether or not to include deprecated Redis™ engine versions`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "version",
				Short:      `List Redis™ engine versions that match a given name pattern`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2, scw.Zone(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.ListClusterVersionsRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListClusterVersions(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Versions, nil

		},
	}
}

func redisClusterGetCertificate() *core.Command {
	return &core.Command{
		Short:     `Get the TLS certificate of a cluster`,
		Long:      `Retrieve information about the TLS certificate of a Redis™ Database Instance (Redis™ cluster). Details like name and content are returned in the response.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "get-certificate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.GetClusterCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.GetClusterCertificateRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.GetClusterCertificate(request)

		},
	}
}

func redisClusterRenewCertificate() *core.Command {
	return &core.Command{
		Short:     `Renew the TLS certificate of a cluster`,
		Long:      `Renew a TLS certificate for a Redis™ Database Instance (Redis™ cluster). Renewing a certificate means that you will not be able to connect to your Database Instance using the previous certificate. You will also need to download and update the new certificate for all database clients.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "renew-certificate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.RenewClusterCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.RenewClusterCertificateRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.RenewClusterCertificate(request)

		},
	}
}

func redisSettingAdd() *core.Command {
	return &core.Command{
		Short:     `Add advanced settings`,
		Long:      `Add an advanced setting to a Redis™ Database Instance (Redis™ cluster). You must set the ` + "`" + `name` + "`" + ` and the ` + "`" + `value` + "`" + ` of each setting.`,
		Namespace: "redis",
		Resource:  "setting",
		Verb:      "add",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.AddClusterSettingsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Database Instance you want to add settings to`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "settings.{index}.value",
				Short:      `Value of the setting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "settings.{index}.name",
				Short:      `Name of the setting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.AddClusterSettingsRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.AddClusterSettings(request)

		},
	}
}

func redisSettingDelete() *core.Command {
	return &core.Command{
		Short:     `Delete advanced setting`,
		Long:      `Delete an advanced setting in a Redis™ Database Instance (Redis™ cluster). You must specify the names of the settings you want to delete in the request body.`,
		Namespace: "redis",
		Resource:  "setting",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.DeleteClusterSettingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Database Instance where the settings must be set`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "setting-name",
				Short:      `Setting name to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.DeleteClusterSettingRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.DeleteClusterSetting(request)

		},
	}
}

func redisSettingSet() *core.Command {
	return &core.Command{
		Short:     `Set advanced settings`,
		Long:      `Update an advanced setting for a Redis™ Database Instance (Redis™ cluster). Settings added upon database engine initalization can only be defined once, and cannot, therefore, be updated.`,
		Namespace: "redis",
		Resource:  "setting",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.SetClusterSettingsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Database Instance where the settings must be set`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "settings.{index}.value",
				Short:      `Value of the setting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "settings.{index}.name",
				Short:      `Name of the setting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.SetClusterSettingsRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.SetClusterSettings(request)

		},
	}
}

func redisACLSet() *core.Command {
	return &core.Command{
		Short:     `Set ACL rules for a cluster`,
		Long:      `Replace all the ACL rules of a Redis™ Database Instance (Redis™ cluster).`,
		Namespace: "redis",
		Resource:  "acl",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.SetACLRulesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Database Instance where the ACL rules have to be set`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acl-rules.{index}.ip-cidr",
				Short:      `IPv4 network address of the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acl-rules.{index}.description",
				Short:      `Description of the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.SetACLRulesRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.SetACLRules(request)

		},
	}
}

func redisACLAdd() *core.Command {
	return &core.Command{
		Short:     `Add ACL rules for a cluster`,
		Long:      `Add an additional ACL rule to a Redis™ Database Instance (Redis™ cluster).`,
		Namespace: "redis",
		Resource:  "acl",
		Verb:      "add",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.AddACLRulesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Database Instance you want to add ACL rules to`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acl-rules.{index}.ip-cidr",
				Short:      `IPv4 network address of the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acl-rules.{index}.description",
				Short:      `Description of the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.AddACLRulesRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.AddACLRules(request)

		},
	}
}

func redisACLDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an ACL rule for a cluster`,
		Long:      `Delete an ACL rule of a Redis™ Database Instance (Redis™ cluster). You must specify the ` + "`" + `acl_id` + "`" + ` of the rule you want to delete in your request.`,
		Namespace: "redis",
		Resource:  "acl",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.DeleteACLRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "acl-id",
				Short:      `UUID of the ACL rule you want to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.DeleteACLRuleRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.DeleteACLRule(request)

		},
	}
}

func redisACLGet() *core.Command {
	return &core.Command{
		Short:     `Get an ACL rule`,
		Long:      `Retrieve information about an ACL rule of a Redis™ Database Instance (Redis™ cluster). You must specify the ` + "`" + `acl_id` + "`" + ` of the rule in your request.`,
		Namespace: "redis",
		Resource:  "acl",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.GetACLRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "acl-id",
				Short:      `UUID of the ACL rule you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.GetACLRuleRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.GetACLRule(request)

		},
	}
}

func redisEndpointSet() *core.Command {
	return &core.Command{
		Short:     `Set endpoints for a cluster`,
		Long:      `Update an endpoint for a Redis™ Database Instance (Redis™ cluster). You must specify the ` + "`" + `cluster_id` + "`" + ` and the ` + "`" + `endpoints` + "`" + ` parameters in your request.`,
		Namespace: "redis",
		Resource:  "endpoint",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.SetEndpointsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Database Instance where the endpoints have to be set`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.id",
				Short:      `UUID of the Private Network to connect to the Database Instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.service-ips.{index}",
				Short:      `Endpoint IPv4 address with a CIDR notation. You must provide at least one IPv4 per node.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.SetEndpointsRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.SetEndpoints(request)

		},
	}
}

func redisEndpointAdd() *core.Command {
	return &core.Command{
		Short:     `Add endpoints for a cluster`,
		Long:      `Add a new endpoint for a Redis™ Database Instance (Redis™ cluster). You can add ` + "`" + `private_network` + "`" + ` or ` + "`" + `public_network` + "`" + ` specifications to the body of the request.`,
		Namespace: "redis",
		Resource:  "endpoint",
		Verb:      "add",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.AddEndpointsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the Database Instance you want to add endpoints to`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.id",
				Short:      `UUID of the Private Network to connect to the Database Instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.service-ips.{index}",
				Short:      `Endpoint IPv4 address with a CIDR notation. You must provide at least one IPv4 per node.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.AddEndpointsRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.AddEndpoints(request)

		},
	}
}

func redisEndpointDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an endpoint for a cluster`,
		Long:      `Delete the endpoint of a Redis™ Database Instance (Redis™ cluster). You must specify the ` + "`" + `region` + "`" + ` and ` + "`" + `endpoint_id` + "`" + ` parameters of the endpoint you want to delete. Note that might need to update any environment configurations that point to the deleted endpoint.`,
		Namespace: "redis",
		Resource:  "endpoint",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.DeleteEndpointRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "endpoint-id",
				Short:      `UUID of the endpoint you want to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.DeleteEndpointRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.DeleteEndpoint(request)

		},
	}
}

func redisEndpointGet() *core.Command {
	return &core.Command{
		Short:     `Get an endpoint`,
		Long:      `Retrieve information about a Redis™ Database Instance (Redis™ cluster) endpoint. Full details about the endpoint, like ` + "`" + `ips` + "`" + `, ` + "`" + `port` + "`" + `, ` + "`" + `private_network` + "`" + ` and ` + "`" + `public_network` + "`" + ` specifications are returned in the response.`,
		Namespace: "redis",
		Resource:  "endpoint",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.GetEndpointRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "endpoint-id",
				Short:      `UUID of the endpoint you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.GetEndpointRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.GetEndpoint(request)

		},
	}
}

func redisEndpointUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an endpoint`,
		Long:      `Update information about a Redis™ Database Instance (Redis™ cluster) endpoint. Full details about the endpoint, like ` + "`" + `ips` + "`" + `, ` + "`" + `port` + "`" + `, ` + "`" + `private_network` + "`" + ` and ` + "`" + `public_network` + "`" + ` specifications are returned in the response.`,
		Namespace: "redis",
		Resource:  "endpoint",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.UpdateEndpointRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "endpoint-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network.id",
				Short:      `UUID of the Private Network to connect to the Database Instance`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network.service-ips.{index}",
				Short:      `Endpoint IPv4 address with a CIDR notation. You must provide at least one IPv4 per node.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZonePlWaw1, scw.ZonePlWaw2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.UpdateEndpointRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.UpdateEndpoint(request)

		},
	}
}
