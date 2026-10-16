package redis

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	"github.com/scaleway/scaleway-cli/v2/internal/semver"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
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
	Endpoints  []*redisEndpointSpecCustom `json:"endpoints"`
	CliConnect bool
}

var connectSeeAlsos = []*core.SeeAlso{
	{
		Command: "scw redis cluster connect",
		Short:   "Connect to a Redis cluster using locally installed redis-cli",
	},
}

func clusterUpdateBuilder(c *core.Command) *core.Command {
	type updateClusterRequestCustom struct {
		*redis.UpdateClusterRequest
		CliConnect bool
	}

	c.ArgSpecs.AddBefore("password", connect.CliConnectArgSpec())
	c.ArgsType = reflect.TypeFor[updateClusterRequestCustom]()

	c.SeeAlsos = connectSeeAlsos

	c.Run = func(ctx context.Context, args any) (any, error) {
		req := args.(*updateClusterRequestCustom)
		region, err := req.Zone.Region()
		if err != nil {
			return nil, err
		}

		// The secret is keyed by user name: when --user-name is not given the
		// cluster keeps its current user, so resolve it to persist at the right path.
		name := ""
		if req.CliConnect {
			if req.UserName != nil {
				name = *req.UserName
			} else {
				cluster, err := redis.NewAPI(core.ExtractClient(ctx)).
					GetCluster(&redis.GetClusterRequest{
						ClusterID: req.ClusterID,
						Zone:      req.Zone,
					})
				if err != nil {
					return nil, err
				}
				name = cluster.UserName
			}
		}

		return connect.RunUserUpdate(
			ctx,
			req.CliConnect,
			connect.Redis,
			region,
			string(req.Zone),
			req.ClusterID,
			name,
			&req.Password,
			func(ctx context.Context) (any, error) {
				return redis.NewAPI(core.ExtractClient(ctx)).
					UpdateCluster(req.UpdateClusterRequest, scw.WithContext(ctx))
			},
		)
	}

	return c
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

	c.ArgSpecs.AddBefore("password", connect.CliConnectArgSpec())
	c.ArgSpecs.GetByName("password").Required = false
	c.ArgSpecs.GetByName("user-name").Default = core.DefaultValueSetter("admin")

	c.SeeAlsos = []*core.SeeAlso{
		{
			Command: "scw redis cluster connect",
			Short:   "Connect to a Redis cluster using locally installed redis-cli",
		},
		{
			Command: "scw redis cluster update",
			Short:   "Change the password of a user of the cluster",
		},
	}

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

		generatedPassword, err := passwordgenerator.GenerateIfEmpty(&createClusterRequest.Password)
		if err != nil {
			return nil, err
		}

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

		if generatedPassword != "" {
			interactive.Printf("Your generated password is %s \n\n", generatedPassword)
		}

		if customRequest.CliConnect {
			region, err := cluster.Zone.Region()
			if err != nil {
				return nil, err
			}

			res := connect.Resource{
				Product:  connect.Redis,
				Location: string(cluster.Zone),
				ID:       cluster.ID,
				Username: createClusterRequest.UserName,
			}
			if err := res.SaveSecret(
				ctx,
				region,
				createClusterRequest.UserName,
				createClusterRequest.Password,
			); err != nil {
				return nil, err
			}
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

	type deleteClusterRequestCustom struct {
		*redis.DeleteClusterRequest
		CliConnect bool
	}

	c.ArgsType = reflect.TypeFor[deleteClusterRequestCustom]()
	c.ArgSpecs.AddBefore("zone", connect.CliConnectDeleteArgSpec())

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		args := argsI.(*deleteClusterRequestCustom)

		client := core.ExtractClient(ctx)
		api := redis.NewAPI(client)
		cluster, err := api.DeleteCluster(args.DeleteClusterRequest, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		if args.CliConnect {
			region, err := args.Zone.Region()
			if err != nil {
				return cluster, err
			}

			resource := connect.Resource{
				Product:  connect.Redis,
				Location: string(args.Zone),
				ID:       args.ClusterID,
			}
			if err := resource.DeleteSecret(ctx, region); err != nil {
				return cluster, err
			}
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
		return "", fmt.Errorf("could not get latest Redis version: %s", err)
	}

	versions := make([]string, 0, len(resp.Versions))
	for _, v := range resp.Versions {
		versions = append(versions, v.Version)
	}

	return semver.LatestVersion(versions)
}

type clusterConnectArgs struct {
	Zone           scw.Zone
	PrivateNetwork bool
	ClusterID      string
	Template       string
}

const (
	errorMessagePublicEndpointNotFound  = "public endpoint not found"
	errorMessagePrivateEndpointNotFound = "private endpoint not found"
	errorMessageEndpointNotFound        = "any endpoint is associated on your cluster"
	errorMessageRedisCliNotFound        = "redis-cli is not installed. Please install redis-cli to use this command"
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

func checkRedisCliInstalled() error {
	cmd := exec.Command("redis-cli", "--version") //nolint:gosec
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s", errorMessageRedisCliNotFound)
	}

	return nil
}

func getRedisZones() []scw.Zone {
	// Get zones dynamically from the Redis API SDK
	// We create a minimal client just to access the Zones() method
	// which doesn't require authentication
	client := &scw.Client{}
	api := redis.NewAPI(client)

	return api.Zones()
}

func clusterConnectCommand() *core.Command {
	return &core.Command{
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "connect",
		Short:     "Connect to a Redis cluster using locally installed redis-cli",
		Long: `Connect to a Redis cluster using locally installed redis-cli. The command will check if redis-cli is installed, download the certificate if TLS is enabled, and prompt for the password.
Use --template to run an arbitrary tool: the Go text/template is rendered with
{{.Username}}, {{.Password}} (from the persisted secret), {{.Host}}, {{.Port}}, {{.CertPath}} and {{.TLS}},
and the result is executed as the connect command.`,
		ArgsType: reflect.TypeFor[clusterConnectArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "private-network",
				Short:    `Connect by the private network endpoint attached.`,
				Required: false,
				Default:  core.DefaultValueSetter("false"),
			},
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster`,
				Required:   true,
				Positional: true,
			},
			connect.TemplateArgSpec(),
			core.ZoneArgSpec(getRedisZones()...),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*clusterConnectArgs)

			if args.Template == "" {
				if err := checkRedisCliInstalled(); err != nil {
					return nil, err
				}
			}

			client := core.ExtractClient(ctx)
			api := redis.NewAPI(client)
			cluster, err := api.GetCluster(&redis.GetClusterRequest{
				Zone:      args.Zone,
				ClusterID: args.ClusterID,
			})
			if err != nil {
				return nil, err
			}

			if len(cluster.Endpoints) == 0 {
				return nil, fmt.Errorf("%s", errorMessageEndpointNotFound)
			}

			var endpoint *redis.Endpoint
			switch {
			case args.PrivateNetwork:
				endpoint, err = getPrivateEndpoint(cluster.Endpoints)
				if err != nil {
					return nil, err
				}
			default:
				endpoint, err = getPublicEndpoint(cluster.Endpoints)
				if err != nil {
					return nil, err
				}
			}

			if len(endpoint.IPs) == 0 {
				return nil, errors.New("endpoint has no IP addresses")
			}

			port := endpoint.Port

			var certPath string
			if cluster.TLSEnabled {
				certResp, err := api.GetClusterCertificate(&redis.GetClusterCertificateRequest{
					Zone:      args.Zone,
					ClusterID: args.ClusterID,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to get certificate: %w", err)
				}

				certContent, err := io.ReadAll(certResp.Content)
				if err != nil {
					return nil, fmt.Errorf("failed to read certificate content: %w", err)
				}

				tmpDir := os.TempDir()
				certPath = filepath.Join(tmpDir, fmt.Sprintf("redis-cert-%s.crt", args.ClusterID))
				if err := os.WriteFile(certPath, certContent, 0o600); err != nil {
					return nil, fmt.Errorf("failed to write certificate: %w", err)
				}
				defer func() {
					if err := os.Remove(certPath); err != nil {
						core.ExtractLogger(ctx).Debugf("failed to remove certificate file: %v", err)
					}
				}()
			}

			region, err := cluster.Zone.Region()
			if err != nil {
				return nil, err
			}

			password, err := connect.Resource{
				Product:  connect.Redis,
				Location: string(cluster.Zone),
				ID:       cluster.ID,
				Username: cluster.UserName,
			}.FetchOrPrompt(ctx, region, cluster.UserName)
			if err != nil {
				return nil, err
			}

			hostStr := endpoint.IPs[0].String()

			// Fail fast on an unreachable endpoint instead of letting redis-cli hang on the dial
			if args.Template == "" {
				if err := core.CheckReachable(hostStr, port); err != nil {
					return nil, err
				}
			}

			cmdArgs := []string{
				"redis-cli",
				"-h", hostStr,
				"-p", strconv.FormatUint(uint64(port), 10),
				"--pass", password,
			}

			if cluster.TLSEnabled {
				cmdArgs = append(cmdArgs, "--tls", "--cert", certPath)
			}

			if cluster.UserName != "" {
				cmdArgs = append(cmdArgs, "--user", cluster.UserName)
			}

			templateArgs, err := connect.RenderCommand(args.Template, map[string]any{
				"Username": cluster.UserName,
				"Password": password,
				"Host":     hostStr,
				"Port":     strconv.FormatUint(uint64(port), 10),
				"CertPath": certPath,
				"TLS":      cluster.TLSEnabled,
			})
			if err != nil {
				return nil, err
			}
			if templateArgs != nil {
				cmdArgs = templateArgs
			}

			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			core.ExtractLogger(ctx).Debugf("executing: %s\n", cmd.Args)

			if err := cmd.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					return nil, &core.CliError{Empty: true, Code: exitError.ExitCode()}
				}

				return nil, err
			}

			return &core.SuccessResult{
				Empty: true,
			}, nil
		},
		Examples: []*core.Example{
			{
				Short: "Connect to a Redis cluster",
				Raw:   `scw redis cluster connect 11111111-1111-1111-1111-111111111111`,
			},
			{
				Short: "Connect to a Redis cluster via private network",
				Raw:   `scw redis cluster connect 11111111-1111-1111-1111-111111111111 private-network=true`,
			},
			{
				Short: "Connect with a custom tool using a template",
				Raw:   `scw redis cluster connect 11111111-1111-1111-1111-111111111111 --template "my-redis-tool -h {{.Host}} -p {{.Port}}"`,
			},
		},
	}
}
