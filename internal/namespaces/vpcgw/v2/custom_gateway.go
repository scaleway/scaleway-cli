package vpcgw

import (
	"context"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
)

const (
	gatewayActionTimeout = 60 * time.Minute
)

var gatewayStatusMarshalSpecs = human.EnumMarshalSpecs{
	vpcgw.GatewayStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
	vpcgw.GatewayStatusAllocating:    &human.EnumMarshalSpec{Attribute: color.FgBlue},
	vpcgw.GatewayStatusConfiguring:   &human.EnumMarshalSpec{Attribute: color.FgBlue},
	vpcgw.GatewayStatusDeleting:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	vpcgw.GatewayStatusFailed:        &human.EnumMarshalSpec{Attribute: color.FgRed},
	vpcgw.GatewayStatusRunning:       &human.EnumMarshalSpec{Attribute: color.FgGreen},
	vpcgw.GatewayStatusStopped:       &human.EnumMarshalSpec{Attribute: color.FgRed},
	vpcgw.GatewayStatusLocked:        &human.EnumMarshalSpec{Attribute: color.FgRed},
	vpcgw.GatewayStatusStopping:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
}

func gatewayCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		getResp := respI.(*vpcgw.Gateway)
		api := vpcgw.NewAPI(core.ExtractClient(ctx))

		return api.WaitForGateway(&vpcgw.WaitForGatewayRequest{
			GatewayID:     getResp.ID,
			Zone:          getResp.Zone,
			Timeout:       new(gatewayActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func gatewayMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp vpcgw.Gateway
	vpcgtw := tmp(i.(vpcgw.Gateway))
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "IPv4",
			Title:     "IPv4",
		},
		{
			FieldName: "GatewayNetworks",
			Title:     "GatewayNetworks",
		},
	}
	str, err := human.Marshal(vpcgtw, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

// Custom delete gateway command with interactive IP deletion confirmation
type customDeleteGatewayRequest struct {
	Zone      scw.Zone
	GatewayID string
	DeleteIP  bool
	WithIP    string // "prompt", "true", "false"
}

const (
	withIPPrompt = "prompt"
	withIPTrue   = "true"
	withIPFalse  = "false"
)

func gatewayDeleteBuilder(c *core.Command) *core.Command {
	c.ArgsType = reflect.TypeOf(customDeleteGatewayRequest{})
	c.ArgSpecs = core.ArgSpecs{
		{
			Name:       "gateway-id",
			Short:      "ID of the gateway to delete",
			Required:   true,
			Positional: true,
		},
		{
			Name:    "with-ip",
			Short:   "Delete the IP attached to the gateway",
			Default: core.DefaultValueSetter(withIPPrompt),
			EnumValues: []string{
				withIPPrompt,
				withIPTrue,
				withIPFalse,
			},
		},
		core.ZoneArgSpec(
			scw.ZoneFrPar1,
			scw.ZoneFrPar2,
			scw.ZoneNlAms1,
			scw.ZoneNlAms2,
			scw.ZoneNlAms3,
			scw.ZonePlWaw1,
			scw.ZonePlWaw2,
			scw.ZonePlWaw3,
		),
	}
	c.Run = func(ctx context.Context, argsI any) (any, error) {
		args := argsI.(*customDeleteGatewayRequest)

		client := core.ExtractClient(ctx)
		api := vpcgw.NewAPI(client)

		// Get gateway info to check if it has an IP
		gateway, err := api.GetGateway(&vpcgw.GetGatewayRequest{
			Zone:      args.Zone,
			GatewayID: args.GatewayID,
		})
		if err != nil {
			return nil, err
		}

		// Determine if we should delete the IP
		deleteIP, err := shouldDeleteIP(ctx, gateway, args.WithIP)
		if err != nil {
			return nil, err
		}

		request := &vpcgw.DeleteGatewayRequest{
			Zone:      args.Zone,
			GatewayID: args.GatewayID,
			DeleteIP:  deleteIP,
		}

		return api.DeleteGateway(request)
	}

	return c
}

func shouldDeleteIP(
	ctx context.Context,
	gateway *vpcgw.Gateway,
	withIP string,
) (bool, error) {
	switch withIP {
	case withIPTrue:
		return true, nil
	case withIPFalse:
		return false, nil
	case withIPPrompt:
		// Only prompt user if gateway has an IP
		if gateway.IPv4 == nil {
			return false, nil
		}

		return interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
			Prompt:       "Do you also want to delete the IP attached to this gateway?",
			DefaultValue: false,
			Ctx:          ctx,
		})
	default:
		return false, nil
	}
}
