package llm_inference

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	llm_inference "github.com/scaleway/scaleway-sdk-go/api/llm_inference/v1beta1"
)

func endpointCreateBuilder(c *core.Command) *core.Command {
	type llmInferenceEndpointSpecCustom struct {
		*llm_inference.EndpointSpec
		IsPublic bool `json:"is-public"`
	}

	type createEndpointRequestCustom struct {
		*llm_inference.CreateEndpointRequest
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
			publicEndpoint := &llm_inference.EndpointSpecPublic{}
			endpointToCreate := llm_inference.EndpointSpec{
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
