package instance

import (
	"context"
	"fmt"
	"net"
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var privateNetworkInterfaceStateMarshalSpecs = human.EnumMarshalSpecs{
	instance.PrivateNetworkInterfaceStatusAttaching: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.PrivateNetworkInterfaceStatusAvailable: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
	},
	instance.PrivateNetworkInterfaceStatusDetaching: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.PrivateNetworkInterfaceStatusSyncing: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.PrivateNetworkInterfaceStatusUnknownStatus: &human.EnumMarshalSpec{
		Attribute: color.Faint,
	},
}

func privateNetworkInterfaceCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("private-network-id").Required = true

	return c
}

func privateNetworkInterfaceUpdateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("private-network-interface-id").Positional = true

	return c
}

func privateNetworkInterfaceDeleteBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("private-network-interface-id").Positional = true

	return c
}

func privateNetworkInterfaceGetBuilder(c *core.Command) *core.Command {
	c.Long += " You can get private network interfaces either by ID or by MAC address. When using the latter, it is recommended to also provide the server-id and/or the private-network-id to narrow search results, depending on the volume of PNIs in the scope."

	type customPrivateNetworkInterfaceGetRequest struct {
		*instance.GetPrivateNetworkInterfaceRequest
		ServerID         *string
		PrivateNetworkID *string
	}

	c.ArgSpecs.GetByName("private-network-interface-id").Short = "The Private Network Interface ID or MAC address"
	c.ArgSpecs.GetByName("private-network-interface-id").Positional = true
	c.ArgSpecs.AddBefore("zone", &core.ArgSpec{
		Name:  "server-id",
		Short: "The ID of the server to filter list results when getting a PNI by MAC address.",
	})
	c.ArgSpecs.AddBefore("zone", &core.ArgSpec{
		Name:  "private-network-id",
		Short: "The ID of the private network to filter list results when getting a PNI by MAC address.",
	})

	c.ArgsType = reflect.TypeFor[customPrivateNetworkInterfaceGetRequest]()

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		args := argsI.(*customPrivateNetworkInterfaceGetRequest)
		request := args.GetPrivateNetworkInterfaceRequest

		if isMacAddress(args.PrivateNetworkInterfaceID) {
			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			listReq := &instance.ListPrivateNetworkInterfacesRequest{
				Zone: args.Zone,
			}

			if args.ServerID != nil {
				listReq.ServerIDs = []string{*args.ServerID}
			}

			if args.PrivateNetworkID != nil {
				listReq.PrivateNetworkIDs = []string{*args.PrivateNetworkID}
			}

			pnis, err := api.ListPrivateNetworkInterfaces(
				listReq,
				scw.WithContext(ctx),
				scw.WithAllPages(),
			)
			if err != nil {
				return nil, err
			}

			for _, pni := range pnis.PrivateNetworkInterfaces {
				if pni.MacAddress == args.PrivateNetworkInterfaceID {
					request.PrivateNetworkInterfaceID = pni.ID
				}
			}
		}

		return runner(ctx, request)
	}

	return c
}

func isMacAddress(address string) bool {
	_, err := net.ParseMAC(address)

	return err == nil
}

func privateNetworkInterfaceListBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}

		pniList, ok := rawResp.(*instance.ListPrivateNetworkInterfacesResponse)
		if !ok {
			return "", fmt.Errorf(
				"expected response of type *instance.ListPrivateNetworkInterfacesResponse, got %T",
				rawResp,
			)
		}

		return pniList.PrivateNetworkInterfaces, nil
	}

	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "ID",
				FieldName: "ID",
			},
			{
				Label:     "SERVER ID",
				FieldName: "ServerID",
			},
			{
				Label:     "PRIVATE NETWORK ID",
				FieldName: "PrivateNetworkID",
			},
			{
				Label:     "MAC ADDRESS",
				FieldName: "MacAddress",
			},
			{
				Label:     "TAGS",
				FieldName: "Tags",
			},
			{
				Label:     "STATUS",
				FieldName: "Status",
			},
			{
				Label:     "IP IDS",
				FieldName: "IPIDs",
			},
			{
				Label:     "PROJECT ID",
				FieldName: "ProjectID",
			},
			{
				Label:     "CREATED AT",
				FieldName: "CreatedAt",
			},
			{
				Label:     "UPDATED AT",
				FieldName: "UpdatedAt",
			},
		},
	}

	return c
}
