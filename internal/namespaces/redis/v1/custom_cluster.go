package redis

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const redisActionTimeout = 15 * time.Minute

func redisClusterMigrateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("node-type").AutoCompleteFunc = autoCompleteNodeType

	return c
}

func clusterCreateBuilder(c *core.Command) *core.Command {
	type redisEndpointSpecPrivateNetworkSpecCustom struct {
		*redis.EndpointSpecPrivateNetworkSpec
		EnableIpam bool `json:"enable-ipam"`
	}

	type redisEndpointSpecCustom struct {
		PrivateNetwork *redisEndpointSpecPrivateNetworkSpecCustom `json:"private-network"`
	}

	type redisCreateClusterRequestCustom struct {
		*redis.CreateClusterRequest
		Endpoints []*redisEndpointSpecCustom `json:"endpoints"`
	}

	c.ArgSpecs.AddBefore("endpoints.{index}.private-network.id", &core.ArgSpec{
		Name:     "endpoints.{index}.private-network.enable-ipam",
		Short:    "Will configure your Private Network endpoint with Scaleway IPAM service if true",
		Required: false,
		Default:  core.DefaultValueSetter("false"),
	})

	c.ArgsType = reflect.TypeOf(redisCreateClusterRequestCustom{})

	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		api := redis.NewAPI(core.ExtractClient(ctx))
		cluster, err := api.WaitForCluster(&redis.WaitForClusterRequest{
			ClusterID:     respI.(*redis.Cluster).ID,
			Zone:          respI.(*redis.Cluster).Zone,
			Timeout:       scw.TimeDurationPtr(redisActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			return nil, err
		}

		return cluster, nil
	}

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		client := core.ExtractClient(ctx)
		api := redis.NewAPI(client)

		customRequest := argsI.(*redisCreateClusterRequestCustom)
		createClusterRequest := customRequest.CreateClusterRequest

		if len(customRequest.Endpoints) == 0 {
			createClusterRequest.Endpoints = append(
				createClusterRequest.Endpoints,
				&redis.EndpointSpec{
					PublicNetwork: &redis.EndpointSpecPublicNetworkSpec{},
				},
			)
		} else {
			for _, customEndpoint := range customRequest.Endpoints {
				if customEndpoint.PrivateNetwork == nil {
					continue
				}
				var ipamConfig *redis.EndpointSpecPrivateNetworkSpecIpamConfig
				if customEndpoint.PrivateNetwork.EnableIpam {
					ipamConfig = &redis.EndpointSpecPrivateNetworkSpecIpamConfig{}
				}
				createClusterRequest.Endpoints = append(
					createClusterRequest.Endpoints,
					&redis.EndpointSpec{
						PrivateNetwork: &redis.EndpointSpecPrivateNetworkSpec{
							ID:         customEndpoint.PrivateNetwork.ID,
							ServiceIPs: customEndpoint.PrivateNetwork.ServiceIPs,
							IpamConfig: ipamConfig,
						},
					},
				)
			}
		}

		cluster, err := api.CreateCluster(createClusterRequest)
		if err != nil {
			return nil, err
		}

		return cluster, nil
	}

	return c
}

func clusterDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		api := redis.NewAPI(core.ExtractClient(ctx))
		cluster, err := api.WaitForCluster(&redis.WaitForClusterRequest{
			ClusterID:     respI.(*redis.Cluster).ID,
			Zone:          respI.(*redis.Cluster).Zone,
			Timeout:       scw.TimeDurationPtr(redisActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			// if we get a 404 here, it means the resource was successfully deleted
			notFoundError := &scw.ResourceNotFoundError{}
			responseError := &scw.ResponseError{}
			if errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound ||
				errors.As(err, &notFoundError) {
				return cluster, nil
			}

			return nil, err
		}

		return cluster, nil
	}

	return c
}

func clusterWaitCommand() *core.Command {
	return &core.Command{
		Short:     "Wait for a Redis cluster to reach a stable state",
		Long:      "Wait for a Redis cluster to reach a stable state. This is similar to using --wait flag.",
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(redis.WaitForClusterRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			api := redis.NewAPI(core.ExtractClient(ctx))

			return api.WaitForCluster(&redis.WaitForClusterRequest{
				Zone:          argsI.(*redis.WaitForClusterRequest).Zone,
				ClusterID:     argsI.(*redis.WaitForClusterRequest).ClusterID,
				Timeout:       argsI.(*redis.WaitForClusterRequest).Timeout,
				RetryInterval: core.DefaultRetryInterval,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      "ID of the cluster you want to wait for",
				Required:   true,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
			),
			core.WaitTimeoutArgSpec(redisActionTimeout),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a Redis cluster to reach a stable state",
				ArgsJSON: `{"cluster-id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func ACLAddListBuilder(c *core.Command) *core.Command {
	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
			originalResp, err := runner(ctx, argsI)
			if err != nil {
				return nil, err
			}
			ACLAddResponse := originalResp.(*redis.AddACLRulesResponse)

			return ACLAddResponse.ACLRules, nil
		},
	)

	return c
}

func redisSettingAddBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("settings.{index}.name").AutoCompleteFunc = autoCompleteSettingsName

	return c
}

func redisClusterGetMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp redis.Cluster
	redisClusterResponse := tmp(i.(redis.Cluster))
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "Endpoints",
			Title:     "Endpoints",
		},
		{
			FieldName: "ACLRules",
			Title:     "ACLRules",
		},
	}
	opt.SubOptions = map[string]*human.MarshalOpt{
		"Endpoints": {
			Fields: []*human.MarshalFieldOpt{
				{
					FieldName: "ID",
					Label:     "ID",
				},
				{
					FieldName: "Port",
					Label:     "Port",
				},
				{
					FieldName: "IPs",
					Label:     "IPs",
				},
			},
		},
	}
	str, err := human.Marshal(redisClusterResponse, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

var completeClusterCache *redis.Cluster

var completeClusterVersionCache *redis.ListClusterVersionsResponse

func autoCompleteSettingsName(
	ctx context.Context,
	prefix string,
	request any,
) core.AutocompleteSuggestions {
	suggestions := core.AutocompleteSuggestions(nil)
	req := request.(*redis.AddClusterSettingsRequest)
	client := core.ExtractClient(ctx)
	api := redis.NewAPI(client)
	if req.ClusterID != "" {
		if completeClusterCache == nil {
			res, err := api.GetCluster(&redis.GetClusterRequest{
				ClusterID: req.ClusterID,
			})
			if err != nil {
				return nil
			}
			completeClusterCache = res
		}
		if completeClusterVersionCache == nil {
			res, err := api.ListClusterVersions(&redis.ListClusterVersionsRequest{
				Zone:    completeClusterCache.Zone,
				Version: &completeClusterCache.Version,
			})
			if err != nil {
				return nil
			}
			completeClusterVersionCache = res
		}

		for _, version := range completeClusterVersionCache.Versions {
			for _, settingName := range version.AvailableSettings {
				if strings.HasPrefix(settingName.Name, prefix) {
					suggestions = append(suggestions, settingName.Name)
				}
			}
		}
	}

	return suggestions
}

var completeRedisNoteTypeCache *redis.ListNodeTypesResponse

func autoCompleteNodeType(
	ctx context.Context,
	prefix string,
	request any,
) core.AutocompleteSuggestions {
	suggestions := core.AutocompleteSuggestions(nil)
	req := request.(*redis.MigrateClusterRequest)
	client := core.ExtractClient(ctx)
	api := redis.NewAPI(client)
	if req.Zone != "" {
		if completeRedisNoteTypeCache == nil {
			res, err := api.ListNodeTypes(&redis.ListNodeTypesRequest{
				Zone: req.Zone,
			})
			if err != nil {
				return nil
			}
			completeRedisNoteTypeCache = res
		}
		for _, nodeType := range completeRedisNoteTypeCache.NodeTypes {
			if strings.HasPrefix(nodeType.Name, prefix) {
				suggestions = append(suggestions, nodeType.Name)
			}
		}
	}

	return suggestions
}
