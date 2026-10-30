package redis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/semver"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const redisActionTimeout = 15 * time.Minute

func redisClusterMigrateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("node-type").AutoCompleteFunc = autoCompleteNodeType

	return c
}

type redisEndpointSpecPrivateNetworkSpecCustom struct {
	*redis.EndpointSpecPrivateNetworkSpec
	EnableIpam bool `json:"enable-ipam"`
}

type redisEndpointSpecCustom struct {
	PrivateNetwork *redisEndpointSpecPrivateNetworkSpecCustom `json:"private-network"`
}

type CreateClusterRequestCustom struct {
	*redis.CreateClusterRequest
	Endpoints []*redisEndpointSpecCustom `json:"endpoints"`
}

func clusterCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.AddBefore("endpoints.{index}.private-network.id", &core.ArgSpec{
		Name:     "endpoints.{index}.private-network.enable-ipam",
		Short:    "Will configure your Private Network endpoint with Scaleway IPAM service if true",
		Required: false,
		Default:  core.DefaultValueSetter("false"),
	})

	c.ArgsType = reflect.TypeFor[CreateClusterRequestCustom]()

	c.ArgSpecs.GetByName("version").Default = core.DefaultValueSetter("latest")
	c.ArgSpecs.GetByName("version").AutoCompleteFunc = autoCompleteRedisVersion
	c.ArgSpecs.GetByName("node-type").AutoCompleteFunc = autoCompleteCreateNodeType

	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		api := redis.NewAPI(core.ExtractClient(ctx))
		cluster, err := api.WaitForCluster(&redis.WaitForClusterRequest{
			ClusterID:     respI.(*redis.Cluster).ID,
			Zone:          respI.(*redis.Cluster).Zone,
			Timeout:       new(redisActionTimeout),
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

		customRequest := argsI.(*CreateClusterRequestCustom)
		createClusterRequest := customRequest.CreateClusterRequest

		// Handle default latest version for redis cluster
		if createClusterRequest.Version == "latest" {
			latestVersion, err := getLatestRedisVersion(client, createClusterRequest.Zone)
			if err != nil {
				return nil, err
			}
			createClusterRequest.Version = latestVersion
		}

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
			Timeout:       new(redisActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			// if we get a 404 here, it means the resource was successfully deleted
			notFoundError := &scw.ResourceNotFoundError{}
			if responseError, ok := errors.AsType[*scw.ResponseError](
				err,
			); ok && responseError.StatusCode == http.StatusNotFound ||
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
		ArgsType:  reflect.TypeFor[redis.WaitForClusterRequest](),
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

// redisClusterGetResponse augments the redis cluster with the private
// networks attached to its endpoints. The Redis API does not return the
// private network name, so it is resolved through the VPC API.
type redisClusterGetResponse struct {
	*redis.Cluster
	PrivateNetworks []*vpc.PrivateNetwork `json:"private_networks,omitempty"`
}

func clusterGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		res, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}
		cluster := res.(*redis.Cluster)

		region, err := cluster.Zone.Region()
		if err != nil {
			return nil, err
		}

		vpcAPI := vpc.NewAPI(core.ExtractClient(ctx))
		privateNetworks := make([]*vpc.PrivateNetwork, 0, len(cluster.Endpoints))
		for _, endpoint := range cluster.Endpoints {
			if endpoint.PrivateNetwork == nil {
				continue
			}
			// ponytail: best-effort name resolution. A failing VPC lookup
			// (permissions, 404) drops that PN instead of failing `cluster get`;
			// also a PN shared by two endpoints is listed twice.
			pn, err := vpcAPI.GetPrivateNetwork(
				&vpc.GetPrivateNetworkRequest{
					Region:           region,
					PrivateNetworkID: endpoint.PrivateNetwork.ID,
				},
				scw.WithContext(ctx),
			)
			if err != nil {
				continue
			}
			privateNetworks = append(privateNetworks, pn)
		}

		return redisClusterGetResponse{
			Cluster:         cluster,
			PrivateNetworks: privateNetworks,
		}, nil
	}

	return c
}

// redisPrivateNetworkMarshalerFunc renders an endpoint's private network
// as a compact single line in table cells.
func redisPrivateNetworkMarshalerFunc(i any, _ *human.MarshalOpt) (string, error) {
	pn := i.(redis.PrivateNetwork)

	return fmt.Sprintf("%s (%s)", pn.ID, pn.ProvisioningMode), nil
}

func redisClusterGetResponseMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	// Alias type so human.Marshal does not call this function back.
	type tmp redisClusterGetResponse
	resp := tmp(i.(redisClusterGetResponse))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "Endpoints",
			Title:     "Endpoints",
		},
		{
			FieldName:   "PrivateNetworks",
			Title:       "Private Networks",
			HideIfEmpty: true,
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
				{
					FieldName: "PrivateNetwork",
					Label:     "Private Network",
				},
			},
		},
		"PrivateNetworks": {
			Fields: []*human.MarshalFieldOpt{
				{
					FieldName: "Name",
					Label:     "Name",
				},
				{
					FieldName: "ID",
					Label:     "ID",
				},
				{
					FieldName: "Region",
					Label:     "Region",
				},
			},
		},
	}
	str, err := human.Marshal(resp, opt)
	if err != nil {
		return "", err
	}

	return str, nil
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
	req := request.(*redis.MigrateClusterRequest)

	return completeNodeTypes(ctx, req.Zone, prefix)
}

func CreateRequestZone(req *CreateClusterRequestCustom) scw.Zone {
	// The embedded *redis.CreateClusterRequest is nil during shell completion
	// (reflect.New leaves it unallocated), so guard before promoting .Zone.
	if req.CreateClusterRequest == nil {
		return ""
	}

	return req.CreateClusterRequest.Zone
}

func autoCompleteCreateNodeType(
	ctx context.Context,
	prefix string,
	request any,
) core.AutocompleteSuggestions {
	req := request.(*CreateClusterRequestCustom)

	return completeNodeTypes(ctx, CreateRequestZone(req), prefix)
}

func completeNodeTypes(
	ctx context.Context,
	zone scw.Zone,
	prefix string,
) core.AutocompleteSuggestions {
	suggestions := core.AutocompleteSuggestions(nil)
	client := core.ExtractClient(ctx)
	api := redis.NewAPI(client)

	if completeRedisNoteTypeCache == nil {
		res, err := api.ListNodeTypes(&redis.ListNodeTypesRequest{
			Zone: zone,
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

	return suggestions
}

// Caching ListClusterVersions response for shell completion
var completeRedisVersionCache *redis.ListClusterVersionsResponse

func autoCompleteRedisVersion(
	ctx context.Context,
	prefix string,
	request any,
) core.AutocompleteSuggestions {
	suggestions := core.AutocompleteSuggestions(nil)
	zone := CreateRequestZone(request.(*CreateClusterRequestCustom))

	client := core.ExtractClient(ctx)
	api := redis.NewAPI(client)

	if completeRedisVersionCache == nil {
		res, err := api.ListClusterVersions(&redis.ListClusterVersionsRequest{
			Zone: zone,
		})
		if err != nil {
			return nil
		}
		completeRedisVersionCache = res
	}

	for _, v := range completeRedisVersionCache.Versions {
		if strings.HasPrefix(v.Version, prefix) {
			suggestions = append(suggestions, v.Version)
		}
	}

	return suggestions
}

func getLatestRedisVersion(client *scw.Client, zone scw.Zone) (string, error) {
	api := redis.NewAPI(client)
	resp, err := api.ListClusterVersions(
		&redis.ListClusterVersionsRequest{Zone: zone},
		scw.WithAllPages(),
	)
	if err != nil {
		return "", fmt.Errorf("could not get latest Redis version: %w", err)
	}

	versions := make([]string, 0, len(resp.Versions))
	for _, v := range resp.Versions {
		versions = append(versions, v.Version)
	}

	return semver.LatestVersion(versions)
}

const (
	errorMessagePublicEndpointNotFound  = "public endpoint not found"
	errorMessagePrivateEndpointNotFound = "private endpoint not found"
	errorMessageEndpointNotFound        = "any endpoint is associated on your cluster"
)

func getPublicEndpoint(endpoints []*redis.Endpoint) (*redis.Endpoint, error) {
	for _, e := range endpoints {
		if e.PublicNetwork != nil {
			return e, nil
		}
	}

	return nil, fmt.Errorf("%s", errorMessagePublicEndpointNotFound)
}

func getPrivateEndpoint(endpoints []*redis.Endpoint) (*redis.Endpoint, error) {
	for _, e := range endpoints {
		if e.PrivateNetwork != nil {
			return e, nil
		}
	}

	return nil, fmt.Errorf("%s", errorMessagePrivateEndpointNotFound)
}
