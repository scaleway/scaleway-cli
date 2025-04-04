package vpcgw

import (
	"context"
	"errors"
	"net/http"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var gatewayNetworkStatusMarshalSpecs = human.EnumMarshalSpecs{
	vpcgw.GatewayNetworkStatusAttaching:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
	vpcgw.GatewayNetworkStatusConfiguring:   &human.EnumMarshalSpec{Attribute: color.FgBlue},
	vpcgw.GatewayNetworkStatusCreated:       &human.EnumMarshalSpec{Attribute: color.FgGreen},
	vpcgw.GatewayNetworkStatusDetaching:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
	vpcgw.GatewayNetworkStatusReady:         &human.EnumMarshalSpec{Attribute: color.FgGreen},
	vpcgw.GatewayNetworkStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
}

func gatewayNetworkCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI interface{}) (interface{}, error) {
		getResp := respI.(*vpcgw.GatewayNetwork)
		api := vpcgw.NewAPI(core.ExtractClient(ctx))

		return api.WaitForGatewayNetwork(&vpcgw.WaitForGatewayNetworkRequest{
			GatewayNetworkID: getResp.ID,
			Zone:             getResp.Zone,
			Timeout:          scw.TimeDurationPtr(gatewayActionTimeout),
			RetryInterval:    core.DefaultRetryInterval,
		})
	}

	return c
}

func gatewayNetworkDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, _ interface{}) (interface{}, error) {
		getResp := argsI.(*vpcgw.DeleteGatewayNetworkRequest)
		api := vpcgw.NewAPI(core.ExtractClient(ctx))
		gwNetwork, err := api.WaitForGatewayNetwork(&vpcgw.WaitForGatewayNetworkRequest{
			GatewayNetworkID: getResp.GatewayNetworkID,
			Zone:             getResp.Zone,
			Timeout:          scw.TimeDurationPtr(gatewayActionTimeout),
			RetryInterval:    core.DefaultRetryInterval,
		})
		if err != nil {
			notFoundError := &scw.ResourceNotFoundError{}
			responseError := &scw.ResponseError{}
			if errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound ||
				errors.As(err, &notFoundError) {
				return &core.SuccessResult{
					Resource: "gateway-network",
					Verb:     "delete",
				}, nil
			}

			return nil, err
		}

		return gwNetwork, nil
	}

	return c
}
