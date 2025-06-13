package rdb

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type rdbEndpointCustomResult struct {
	Endpoints []*rdb.Endpoint
	Success   core.SuccessResult
}

func rdbEndpointCustomResultMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	result := i.(rdbEndpointCustomResult)
	messageStr, err := result.Success.MarshalHuman()
	if err != nil {
		return "", err
	}
	if len(result.Endpoints) == 0 {
		return messageStr, nil
	}

	endpointsStr, err := human.Marshal(result.Endpoints, opt)
	if err != nil {
		return "", err
	}

	return messageStr + "\n" + endpointsStr, nil
}

type rdbEndpointSpecCustom struct {
	PrivateNetwork *rdbEndpointSpecPrivateNetworkCustom `json:"private-network"`
	LoadBalancer   bool                                 `json:"load-balancer"`
}

type rdbEndpointSpecPrivateNetworkCustom struct {
	*rdb.EndpointSpecPrivateNetwork
	EnableIpam bool `json:"enable-ipam"`
}

func endpointCreateBuilder(c *core.Command) *core.Command {
	type rdbCreateEndpointRequestCustom struct {
		*rdb.CreateEndpointRequest
		LoadBalancer   bool                                 `json:"load-balancer"`
		PrivateNetwork *rdbEndpointSpecPrivateNetworkCustom `json:"private-network"`
	}

	c.ArgsType = reflect.TypeOf(rdbCreateEndpointRequestCustom{})

	c.ArgSpecs = core.ArgSpecs{
		{
			Name:       "instance-id",
			Short:      `UUID of the Database Instance to which you want to add an endpoint`,
			Required:   true,
			Positional: true,
		},
		{
			Name:     "private-network.private-network-id",
			Short:    `UUID of the Private Network to be connected to the Database Instance`,
			Required: false,
		},
		{
			Name:     "private-network.service-ip",
			Short:    `Endpoint IPv4 address with a CIDR notation. Refer to the official Scaleway documentation to learn more about IP and subnet limitations.`,
			Required: false,
		},
		{
			Name:     "private-network.enable-ipam",
			Short:    "Will configure your Private Network endpoint with Scaleway IPAM service if true",
			Required: false,
			Default:  core.DefaultValueSetter("true"),
		},
		{
			Name:     "load-balancer",
			Short:    "Will configure a public Load-Balancer endpoint",
			Required: false,
			Default:  core.DefaultValueSetter("true"),
		},
		core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
	}

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)
		args := argsI.(*rdbCreateEndpointRequestCustom)

		customEndpointSpec := &rdbEndpointSpecCustom{
			PrivateNetwork: args.PrivateNetwork,
			LoadBalancer:   args.LoadBalancer,
		}
		endpointRequest, err := endpointRequestFromCustom(
			[]*rdbEndpointSpecCustom{customEndpointSpec},
		)
		if err != nil {
			return nil, err
		}

		endpoint, err := api.CreateEndpoint(&rdb.CreateEndpointRequest{
			Region:       args.Region,
			InstanceID:   args.InstanceID,
			EndpointSpec: endpointRequest[0],
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("could not create endpoint: %w", err)
		}

		return &rdbEndpointCustomResult{
			Endpoints: []*rdb.Endpoint{endpoint},
			Success: core.SuccessResult{
				Message: fmt.Sprintf("Endpoint %s successfully added", endpoint.ID),
			},
		}, nil
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI any) (any, error) {
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)
		args := argsI.(*rdbCreateEndpointRequestCustom)

		instance, err := api.WaitForInstance(&rdb.WaitForInstanceRequest{
			InstanceID: args.InstanceID,
			Region:     args.Region,
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		newEndpointID := respI.(*rdbEndpointCustomResult).Endpoints[0].ID
		// The ID in the CreateEndpoint response is empty for public endpoints, so we have to fetch it among the WaitForInstance response's endpoints
		if newEndpointID == "" {
			for _, endpoint := range instance.Endpoints {
				if endpoint.LoadBalancer != nil {
					newEndpointID = endpoint.ID
				}
			}
		}

		return &rdbEndpointCustomResult{
			Endpoints: instance.Endpoints,
			Success: core.SuccessResult{
				Message: fmt.Sprintf("Endpoint %s successfully added", newEndpointID),
			},
		}, nil
	}

	return c
}

func endpointDeleteBuilder(c *core.Command) *core.Command {
	type rdbDeleteEndpointRequestCustom struct {
		*rdb.DeleteEndpointRequest
		InstanceID string `json:"instance-id"`
	}

	c.ArgsType = reflect.TypeOf(rdbDeleteEndpointRequestCustom{})

	c.ArgSpecs.GetByName("endpoint-id").Positional = true
	c.ArgSpecs.AddBefore("endpoint-id", &core.ArgSpec{
		Name:       "instance-id",
		Short:      `UUID of the Database Instance from which you want to delete an endpoint`,
		Required:   true,
		Positional: false,
	})

	c.Run = func(ctx context.Context, argsI any) (i any, e error) {
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)
		args := argsI.(*rdbDeleteEndpointRequestCustom)

		err := api.DeleteEndpoint(&rdb.DeleteEndpointRequest{
			Region:     args.Region,
			EndpointID: args.EndpointID,
		})
		if err != nil {
			return nil, err
		}

		return &core.SuccessResult{
			Resource: "endpoint",
			Verb:     "delete",
		}, nil
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI any) (any, error) {
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)
		args := argsI.(*rdbDeleteEndpointRequestCustom)

		_, err := api.WaitForInstance(&rdb.WaitForInstanceRequest{
			InstanceID: args.InstanceID,
			Region:     args.Region,
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		return respI, nil
	}

	return c
}

func endpointGetBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("endpoint-id").Positional = true

	return c
}

func endpointListCommand() *core.Command {
	type rdbEndpointListCustomArgs struct {
		InstanceID string
		Region     scw.Region
	}

	return &core.Command{
		Namespace: "rdb",
		Resource:  "endpoint",
		Verb:      "list",
		Short:     "Lists a Database Instance's endpoints",
		Long:      "Lists all public and private endpoints of a Database Instance",
		ArgsType:  reflect.TypeOf(rdbEndpointListCustomArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-id",
				Short:      `ID of the Database Instance`,
				Required:   true,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			client := core.ExtractClient(ctx)
			api := rdb.NewAPI(client)
			args := argsI.(*rdbEndpointListCustomArgs)

			instance, err := api.GetInstance(&rdb.GetInstanceRequest{
				Region:     args.Region,
				InstanceID: args.InstanceID,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to get instance: %w", err)
			}

			return instance.Endpoints, nil
		},
	}
}

func endpointRequestFromCustom(
	customEndpoints []*rdbEndpointSpecCustom,
) ([]*rdb.EndpointSpec, error) {
	endpoints := []*rdb.EndpointSpec(nil)
	for _, customEndpoint := range customEndpoints {
		switch {
		case customEndpoint.PrivateNetwork != nil && customEndpoint.PrivateNetwork.EndpointSpecPrivateNetwork != nil:
			ipamConfig := &rdb.EndpointSpecPrivateNetworkIpamConfig{}
			if !customEndpoint.PrivateNetwork.EnableIpam ||
				customEndpoint.PrivateNetwork.ServiceIP != nil {
				ipamConfig = nil
			}
			endpoints = append(endpoints, &rdb.EndpointSpec{
				PrivateNetwork: &rdb.EndpointSpecPrivateNetwork{
					PrivateNetworkID: customEndpoint.PrivateNetwork.PrivateNetworkID,
					ServiceIP:        customEndpoint.PrivateNetwork.ServiceIP,
					IpamConfig:       ipamConfig,
				},
			})
		case customEndpoint.LoadBalancer:
			endpoints = append(endpoints, &rdb.EndpointSpec{
				LoadBalancer: &rdb.EndpointSpecLoadBalancer{},
			})
		default:
			return nil, errors.New("endpoint must be either a load-balancer or a private network")
		}
	}

	return endpoints, nil
}
