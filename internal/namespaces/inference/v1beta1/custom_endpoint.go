package inference

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	inference "github.com/scaleway/scaleway-sdk-go/api/inference/v1beta1"
)

func endpointCreateBuilder(c *core.Command) *core.Command {
	type llmInferenceEndpointSpecCustom struct {
		*inference.EndpointSpec
		IsPublic bool `json:"is-public"`
	}

	type createEndpointRequestCustom struct {
		*inference.CreateEndpointRequest
		Endpoint *llmInferenceEndpointSpecCustom `json:"endpoint"`
	}

	c.ArgSpecs.AddBefore("endpoint.private-network.private-network-id", &core.ArgSpec{
		Name:     "endpoint.is-public",
		Short:    "Will configure your public endpoint if true",
		Required: false,
		Default:  core.DefaultValueSetter("false"),
	})

	c.ArgsType = reflect.TypeOf(createEndpointRequestCustom{})

	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		createEndpointCustomRequest := argsI.(*createEndpointRequestCustom)
		createEndpointreq := createEndpointCustomRequest.CreateEndpointRequest
		endpoint := createEndpointCustomRequest.Endpoint
		if endpoint.IsPublic {
			publicEndpoint := &inference.EndpointSpecPublic{}
			endpointToCreate := inference.EndpointSpec{
				Public:         publicEndpoint,
				PrivateNetwork: nil,
				DisableAuth:    endpoint.DisableAuth,
			}
			createEndpointreq.Endpoint = &endpointToCreate
		}
		return runner(ctx, createEndpointreq)
	}
	return c
}
