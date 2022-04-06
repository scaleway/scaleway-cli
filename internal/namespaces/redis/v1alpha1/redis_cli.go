// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package redis

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1alpha1"
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
	)
}
func redisRoot() *core.Command {
	return &core.Command{
		Short:     `Database Redis API`,
		Long:      ``,
		Namespace: "redis",
	}
}

func redisCluster() *core.Command {
	return &core.Command{
		Short: `Cluster management commands`,
		Long: `A Redis cluster is composed of one or more Nodes depending of the cluster_size setting.
`,
		Namespace: "redis",
		Resource:  "cluster",
	}
}

func redisNodeType() *core.Command {
	return &core.Command{
		Short: `Node Types management commands`,
		Long: `Node types powering your cluster.
`,
		Namespace: "redis",
		Resource:  "node-type",
	}
}

func redisVersion() *core.Command {
	return &core.Command{
		Short:     `Redis Version management commands`,
		Long:      `Redis versions powering your cluster.`,
		Namespace: "redis",
		Resource:  "version",
	}
}

func redisSetting() *core.Command {
	return &core.Command{
		Short: `Settings management commands`,
		Long: `Cluster Settings are tunables of Redis Cluster. Available settings depend on the Redis version.
`,
		Namespace: "redis",
		Resource:  "setting",
	}
}

func redisACL() *core.Command {
	return &core.Command{
		Short: `Access Control List (ACL) management commands`,
		Long: `Network Access Control List allows to control network inbound traffic allowed by setting up ACL rules. ACL rules could be created, edited, deleted.
`,
		Namespace: "redis",
		Resource:  "acl",
	}
}

func redisEndpoint() *core.Command {
	return &core.Command{
		Short: `Endpoints management commands`,
		Long: `Manage endpoints to access to your cluster trough Public or Private networks
`,
		Namespace: "redis",
		Resource:  "endpoint",
	}
}

func redisClusterCreate() *core.Command {
	return &core.Command{
		Short:     `Create a cluster`,
		Long:      `Create a cluster.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.CreateClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("ins"),
			},
			{
				Name:       "version",
				Short:      `Redis version of the cluster`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to apply to the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-type",
				Short:      `Type of node to use for the cluster`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-name",
				Short:      `Name of the user created when the cluster is created`,
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
				Short:      `Number of nodes for the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acl-rules.{index}.ip",
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
				Short:      `UUID of the private network to be connected to the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.service-ips.{index}",
				Short:      `Endpoint IPv4 adress with a CIDR notation. You must provide at least one IPv4 per node. Check documentation about IP and subnet limitation.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tls-enabled",
				Short:      `Whether or not TLS is enabled`,
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
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Short:     `Update a cluster`,
		Long:      `Update a cluster.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.UpdateClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of a given cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-name",
				Short:      `Name of the cluster user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password of the cluster user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Short:     `Get a cluster`,
		Long:      `Get a cluster.`,
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
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Short:     `List clusters`,
		Long:      `List clusters.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.ListClustersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tags.{index}",
				Short:      `Tags of the clusters to filter upon`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the clusters to filter upon`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Criteria to use when ordering cluster listing`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "project-id",
				Short:      `Project ID to list the cluster of`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.ListClustersRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			resp, err := api.ListClusters(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Clusters, nil

		},
	}
}

func redisClusterMigrate() *core.Command {
	return &core.Command{
		Short:     `Migrate a cluster`,
		Long:      `Migrate a cluster.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "migrate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.MigrateClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "version",
				Short:      `Redis version of the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-type",
				Short:      `Type of node to use for the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cluster-size",
				Short:      `Number of nodes for the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Short:     `Delete a cluster`,
		Long:      `Delete a cluster.`,
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.DeleteClusterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Short:     `Get metrics of a cluster`,
		Long:      `Get metrics of a cluster.`,
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
				Name:       "start-date",
				Short:      `Start date to gather metrics from`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "end-date",
				Short:      `End date to gather metrics from`,
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
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Long:      `List available node types.`,
		Namespace: "redis",
		Resource:  "node-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.ListNodeTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "include-disabled-types",
				Short:      `Whether or not to include disabled types`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.ListNodeTypesRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			resp, err := api.ListNodeTypes(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.NodeTypes, nil

		},
	}
}

func redisVersionList() *core.Command {
	return &core.Command{
		Short:     `List available Redis versions`,
		Long:      `List available Redis versions.`,
		Namespace: "redis",
		Resource:  "version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.ListVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "include-disabled",
				Short:      `Whether or not to include disabled Redis versions`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "include-beta",
				Short:      `Whether or not to include beta Redis versions`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "include-deprecated",
				Short:      `Whether or not to include deprecated Redis versions`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "version-name",
				Short:      `List Redis versions that match a given name pattern`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.ListVersionsRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			resp, err := api.ListVersions(request, scw.WithAllPages())
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
		Long:      `Get the TLS certificate of a cluster.`,
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
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Long:      `Renew the TLS certificate of a cluster.`,
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
			core.ZoneArgSpec(scw.ZoneFrPar2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.RenewClusterCertificateRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			e = api.RenewClusterCertificate(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "cluster",
				Verb:     "renew-certificate",
			}, nil
		},
	}
}

func redisSettingAdd() *core.Command {
	return &core.Command{
		Short:     `Add cluster settings`,
		Long:      `Add cluster settings.`,
		Namespace: "redis",
		Resource:  "setting",
		Verb:      "add",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.AddClusterSettingsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster you want to add settings to`,
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
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Short:     `Delete a cluster setting`,
		Long:      `Delete a cluster setting.`,
		Namespace: "redis",
		Resource:  "setting",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.DeleteClusterSettingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster where the settings has to be set`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "settings-name",
				Short:      `Setting name to delete`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.DeleteClusterSettingRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			e = api.DeleteClusterSetting(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "setting",
				Verb:     "delete",
			}, nil
		},
	}
}

func redisSettingSet() *core.Command {
	return &core.Command{
		Short:     `Set cluster settings`,
		Long:      `Set cluster settings.`,
		Namespace: "redis",
		Resource:  "setting",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.SetClusterSettingsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster where the settings has to be set`,
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
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Short:     `Set ACL rules for a given cluster`,
		Long:      `Set ACL rules for a given cluster.`,
		Namespace: "redis",
		Resource:  "acl",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.SetACLRulesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster where the ACL rules has to be set`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acl-rules.{index}.ip",
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
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Short:     `Add ACL rules for a given cluster`,
		Long:      `Add ACL rules for a given cluster.`,
		Namespace: "redis",
		Resource:  "acl",
		Verb:      "add",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.AddACLRulesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster you want to add acl rules to`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acl-rules.{index}.ip",
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
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Short:     `Delete an ACL rule for a given cluster`,
		Long:      `Delete an ACL rule for a given cluster.`,
		Namespace: "redis",
		Resource:  "acl",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.DeleteACLRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "acl-id",
				Short:      `UUID of the acl rule you want to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.DeleteACLRuleRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			e = api.DeleteACLRule(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "acl",
				Verb:     "delete",
			}, nil
		},
	}
}

func redisACLGet() *core.Command {
	return &core.Command{
		Short:     `Get an ACL rule`,
		Long:      `Get an ACL rule.`,
		Namespace: "redis",
		Resource:  "acl",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.GetACLRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "acl-id",
				Short:      `UUID of the acl rule you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Short:     `Set endpoints for a given cluster`,
		Long:      `Set endpoints for a given cluster.`,
		Namespace: "redis",
		Resource:  "endpoint",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.SetEndpointsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster where the endpoints has to be set`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.id",
				Short:      `UUID of the private network to be connected to the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.service-ips.{index}",
				Short:      `Endpoint IPv4 adress with a CIDR notation. You must provide at least one IPv4 per node. Check documentation about IP and subnet limitation.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Short:     `Add endpoints for a given cluster`,
		Long:      `Add endpoints for a given cluster.`,
		Namespace: "redis",
		Resource:  "endpoint",
		Verb:      "add",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.AddEndpointsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster you want to add endpoints to`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.id",
				Short:      `UUID of the private network to be connected to the cluster`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.service-ips.{index}",
				Short:      `Endpoint IPv4 adress with a CIDR notation. You must provide at least one IPv4 per node. Check documentation about IP and subnet limitation.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
		Short:     `Delete an endpoint for a given cluster`,
		Long:      `Delete an endpoint for a given cluster.`,
		Namespace: "redis",
		Resource:  "endpoint",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.DeleteEndpointRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "endpoint-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.DeleteEndpointRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			e = api.DeleteEndpoint(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "endpoint",
				Verb:     "delete",
			}, nil
		},
	}
}

func redisEndpointGet() *core.Command {
	return &core.Command{
		Short:     `Get an endpoint`,
		Long:      `Get an endpoint.`,
		Namespace: "redis",
		Resource:  "endpoint",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(redis.GetEndpointRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "endpoint-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*redis.GetEndpointRequest)

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			return api.GetEndpoint(request)

		},
	}
}
