package llm_inference

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	llm_inference "github.com/scaleway/scaleway-sdk-go/api/llm_inference/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	deploymentActionTimeout = 40 * time.Minute
)

var (
	deployementStateMarshalSpecs = human.EnumMarshalSpecs{
		llm_inference.DeploymentStatusCreating:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		llm_inference.DeploymentStatusDeploying: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		llm_inference.DeploymentStatusDeleting:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		llm_inference.DeploymentStatusError:     &human.EnumMarshalSpec{Attribute: color.FgRed},
		llm_inference.DeploymentStatusReady:     &human.EnumMarshalSpec{Attribute: color.FgGreen},
		llm_inference.DeploymentStatusLocked:    &human.EnumMarshalSpec{Attribute: color.FgRed},
	}
)

func DeploymentMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp llm_inference.Deployment
	deployment := tmp(i.(llm_inference.Deployment))
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "Endpoints",
			Title:     "Endpoints",
		},
	}
	str, err := human.Marshal(deployment, opt)
	if err != nil {
		return "", err
	}
	return str, nil
}

func deploymentCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("node-type").AutoCompleteFunc = autocompleteDeploymentNodeType
	type llmInferenceEndpointSpecCustom struct {
		*llm_inference.EndpointSpec
		IsPublic bool `json:"is-public"`
	}

	type llmInferenceCreateDeploymentRequestCustom struct {
		*llm_inference.CreateDeploymentRequest
		Endpoints []*llmInferenceEndpointSpecCustom `json:"endpoints"`
	}

	c.ArgSpecs.AddBefore("endpoints.{index}.private-network.private-network-id", &core.ArgSpec{
		Name:     "endpoints.{index}.is-public",
		Short:    "Will configure your public endpoint if true",
		Required: false,
		Default:  core.DefaultValueSetter("false"),
	})

	c.ArgsType = reflect.TypeOf(llmInferenceCreateDeploymentRequestCustom{})

	c.WaitFunc = func(ctx context.Context, _, respI interface{}) (interface{}, error) {
		api := llm_inference.NewAPI(core.ExtractClient(ctx))
		return api.WaitForDeployment(&llm_inference.WaitForDeploymentRequest{
			DeploymentId:  respI.(*llm_inference.Deployment).ID,
			Region:        respI.(*llm_inference.Deployment).Region,
			Status:        respI.(*llm_inference.Deployment).Status,
			Timeout:       scw.TimeDurationPtr(deploymentActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		deploymentCreateCustomRequest := argsI.(*llmInferenceCreateDeploymentRequestCustom)
		deploymentRequest := deploymentCreateCustomRequest.CreateDeploymentRequest
		if deploymentCreateCustomRequest.Endpoints == nil {
			publicEndpoint := &llm_inference.EndpointSpecPublic{}
			endpoint := llm_inference.EndpointSpec{
				Public:         publicEndpoint,
				PrivateNetwork: nil,
				DisableAuth:    false,
			}
			deploymentRequest.Endpoints = append(deploymentRequest.Endpoints, &endpoint)
			return runner(ctx, deploymentRequest)
		}
		for _, endpoint := range deploymentCreateCustomRequest.Endpoints {
			publicEndpoint := &llm_inference.EndpointSpecPublic{}
			if !endpoint.IsPublic {
				publicEndpoint = nil
			}
			privateNetwork := &llm_inference.EndpointSpecPrivateNetwork{}
			if endpoint.EndpointSpec == nil {
				privateNetwork = nil
			} else {
				privateNetwork.PrivateNetworkID = endpoint.PrivateNetwork.PrivateNetworkID
			}
			endpoint := llm_inference.EndpointSpec{
				Public:         publicEndpoint,
				PrivateNetwork: privateNetwork,
				DisableAuth:    endpoint.DisableAuth,
			}
			deploymentRequest.Endpoints = append(deploymentRequest.Endpoints, &endpoint)
		}

		return runner(ctx, deploymentRequest)
	}

	return c
}

func deploymentDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI interface{}) (interface{}, error) {
		api := llm_inference.NewAPI(core.ExtractClient(ctx))
		deployment, err := api.WaitForDeployment(&llm_inference.WaitForDeploymentRequest{
			DeploymentId:  respI.(*llm_inference.Deployment).ID,
			Region:        respI.(*llm_inference.Deployment).Region,
			Status:        respI.(*llm_inference.Deployment).Status,
			Timeout:       scw.TimeDurationPtr(deploymentActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			notFoundError := &scw.ResourceNotFoundError{}
			responseError := &scw.ResponseError{}
			if errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound || errors.As(err, &notFoundError) {
				return &core.SuccessResult{
					Resource: "deployment",
					Verb:     "delete",
				}, nil
			}
			return nil, err
		}
		return deployment, nil
	}
	return c
}

var completeListNodeTypesCache *llm_inference.ListNodeTypesResponse

func autocompleteDeploymentNodeType(ctx context.Context, prefix string, request any) core.AutocompleteSuggestions {
	req := request.(*llm_inference.CreateDeploymentRequest)
	suggestions := core.AutocompleteSuggestions(nil)

	client := core.ExtractClient(ctx)
	api := llm_inference.NewAPI(client)

	if completeListNodeTypesCache == nil {
		res, err := api.ListNodeTypes(&llm_inference.ListNodeTypesRequest{
			Region: req.Region,
		})
		if err != nil {
			return nil
		}
		completeListNodeTypesCache = res
	}

	for _, nodeType := range completeListNodeTypesCache.NodeTypes {
		if strings.HasPrefix(nodeType.Name, prefix) {
			suggestions = append(suggestions, nodeType.Name)
		}
	}

	return suggestions
}
