package vpcgw

import (
	"context"
	"errors"
	"net/http"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var gatewayNetworkStatusFormat = humanv2.EnumFormats[vpcgw.GatewayNetworkStatus]{
	vpcgw.GatewayNetworkStatusAttaching:     humanv2.EnumFormat{Color: color.FgBlue},
	vpcgw.GatewayNetworkStatusConfiguring:   humanv2.EnumFormat{Color: color.FgBlue},
	vpcgw.GatewayNetworkStatusCreated:       humanv2.EnumFormat{Color: color.FgGreen},
	vpcgw.GatewayNetworkStatusDetaching:     humanv2.EnumFormat{Color: color.FgBlue},
	vpcgw.GatewayNetworkStatusReady:         humanv2.EnumFormat{Color: color.FgGreen},
	vpcgw.GatewayNetworkStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},
}

func gatewayNetworkCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		getResp := respI.(*vpcgw.GatewayNetwork)
		api := vpcgw.NewAPI(core.ExtractClient(ctx))

		return api.WaitForGatewayNetwork(&vpcgw.WaitForGatewayNetworkRequest{
			GatewayNetworkID: getResp.ID,
			Zone:             getResp.Zone,
			Timeout:          new(gatewayActionTimeout),
			RetryInterval:    core.DefaultRetryInterval,
		})
	}

	return c
}

func gatewayNetworkDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, _ any) (any, error) {
		getResp := argsI.(*vpcgw.DeleteGatewayNetworkRequest)
		api := vpcgw.NewAPI(core.ExtractClient(ctx))
		gwNetwork, err := api.WaitForGatewayNetwork(&vpcgw.WaitForGatewayNetworkRequest{
			GatewayNetworkID: getResp.GatewayNetworkID,
			Zone:             getResp.Zone,
			Timeout:          new(gatewayActionTimeout),
			RetryInterval:    core.DefaultRetryInterval,
		})
		if err != nil {
			notFoundError := &scw.ResourceNotFoundError{}
			if responseError, ok := errors.AsType[*scw.ResponseError](
				err,
			); ok && responseError.StatusCode == http.StatusNotFound ||
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
