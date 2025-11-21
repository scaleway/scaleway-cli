package inference

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/inference/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	deploymentActionTimeout = 60 * time.Minute
	deploymentActionCreate  = 1
	deploymentActionDelete  = 2
)

var deploymentStateMarshalSpecs = human.EnumMarshalSpecs{
	inference.DeploymentStatusCreating:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
	inference.DeploymentStatusDeploying: &human.EnumMarshalSpec{Attribute: color.FgBlue},
	inference.DeploymentStatusDeleting:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
	inference.DeploymentStatusError:     &human.EnumMarshalSpec{Attribute: color.FgRed},
	inference.DeploymentStatusReady:     &human.EnumMarshalSpec{Attribute: color.FgGreen},
	inference.DeploymentStatusLocked:    &human.EnumMarshalSpec{Attribute: color.FgRed},
}

func DeploymentMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp inference.Deployment
	deployment := tmp(i.(inference.Deployment))
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

func deploymentDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForDeploymentFunc(deploymentActionDelete)

	return c
}

func deploymentCreateBuilder(c *core.Command) *core.Command {
	type llmInferenceEndpointSpecCustom struct {
		*inference.EndpointSpec
		IsPublic bool `json:"is-public"`
	}

	type llmInferenceCreateDeploymentRequestCustom struct {
		*inference.CreateDeploymentRequest
		Endpoints []*llmInferenceEndpointSpecCustom `json:"endpoints"`
	}

	c.ArgSpecs.AddBefore("endpoints.{index}.private-network.private-network-id", &core.ArgSpec{
		Name:     "endpoints.{index}.is-public",
		Short:    "Will configure your public endpoint if true",
		Required: false,
		Default:  core.DefaultValueSetter("true"),
	})

	c.ArgsType = reflect.TypeOf(llmInferenceCreateDeploymentRequestCustom{})

	c.WaitFunc = waitForDeploymentFunc(deploymentActionCreate)

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		deploymentCreateCustomRequest := argsI.(*llmInferenceCreateDeploymentRequestCustom)
		deploymentRequest := deploymentCreateCustomRequest.CreateDeploymentRequest
		if deploymentCreateCustomRequest.Endpoints == nil {
			publicEndpoint := &inference.EndpointPublicNetworkDetails{}
			endpoint := inference.EndpointSpec{
				PublicNetwork:  publicEndpoint,
				PrivateNetwork: nil,
				DisableAuth:    false,
			}
			deploymentRequest.Endpoints = append(deploymentRequest.Endpoints, &endpoint)

			return runner(ctx, deploymentRequest)
		}
		for _, ep := range deploymentCreateCustomRequest.Endpoints {
			if ep.IsPublic {
				deploymentRequest.Endpoints = append(
					deploymentRequest.Endpoints,
					&inference.EndpointSpec{
						PublicNetwork: &inference.EndpointPublicNetworkDetails{},
						DisableAuth:   ep.DisableAuth,
					},
				)
			}

			if ep.PrivateNetwork != nil {
				deploymentRequest.Endpoints = append(
					deploymentRequest.Endpoints,
					&inference.EndpointSpec{
						PrivateNetwork: &inference.EndpointPrivateNetworkDetails{
							PrivateNetworkID: ep.PrivateNetwork.PrivateNetworkID,
						},
						DisableAuth: ep.DisableAuth,
					},
				)
			}
		}

		return runner(ctx, deploymentRequest)
	}

	return c
}

func waitForDeploymentFunc(action int) core.WaitFunc {
	return func(ctx context.Context, _, respI any) (any, error) {
		deployment, err := inference.NewAPI(core.ExtractClient(ctx)).
			WaitForDeployment(&inference.WaitForDeploymentRequest{
				DeploymentID:  respI.(*inference.Deployment).ID,
				Region:        respI.(*inference.Deployment).Region,
				Timeout:       scw.TimeDurationPtr(deploymentActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})

		switch action {
		case deploymentActionCreate:
			return deployment, err
		case deploymentActionDelete:
			if err != nil {
				// if we get a 404 here, it means the resource was successfully deleted
				notFoundError := &scw.ResourceNotFoundError{}
				responseError := &scw.ResponseError{}
				if errors.As(err, &responseError) &&
					responseError.StatusCode == http.StatusNotFound ||
					errors.As(err, &notFoundError) {
					return fmt.Sprintf(
						"Server %s successfully deleted.",
						respI.(*inference.Deployment).ID,
					), nil
				}
			}
		}

		return nil, err
	}
}
