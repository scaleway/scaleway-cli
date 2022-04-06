// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package flexibleip

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		fipRoot(),
		fipIP(),
		fipMac(),
		fipIPCreate(),
		fipIPGet(),
		fipIPList(),
		fipIPUpdate(),
		fipIPDelete(),
		fipIPAttach(),
		fipIPDetach(),
		fipMacCreate(),
		fipMacDuplicate(),
		fipMacDelete(),
	)
}
func fipRoot() *core.Command {
	return &core.Command{
		Short:     `Flexible IP API`,
		Long:      ``,
		Namespace: "fip",
	}
}

func fipIP() *core.Command {
	return &core.Command{
		Short: `Flexible IP management commands`,
		Long: `A Flexible IP can be attached to any server in the same zone.
A server can be linked with multiple Flexible IPs attached to it.
`,
		Namespace: "fip",
		Resource:  "ip",
	}
}

func fipMac() *core.Command {
	return &core.Command{
		Short:     `MAC address management commands`,
		Long:      `MAC address management commands.`,
		Namespace: "fip",
		Resource:  "mac",
	}
}

func fipIPCreate() *core.Command {
	return &core.Command{
		Short:     `Create a Flexible IP`,
		Long:      `Create a Flexible IP.`,
		Namespace: "fip",
		Resource:  "ip",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(flexibleip.CreateFlexibleIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "description",
				Short:      `Description to associate with the Flexible IP, max 255 characters`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to associate to the Flexible IP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-id",
				Short:      `Server ID on which to attach the created Flexible IP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "reverse",
				Short:      `Reverse DNS value`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*flexibleip.CreateFlexibleIPRequest)

			client := core.ExtractClient(ctx)
			api := flexibleip.NewAPI(client)
			return api.CreateFlexibleIP(request)

		},
	}
}

func fipIPGet() *core.Command {
	return &core.Command{
		Short:     `Get a Flexible IP`,
		Long:      `Get a Flexible IP.`,
		Namespace: "fip",
		Resource:  "ip",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(flexibleip.GetFlexibleIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "fip-id",
				Short:      `Flexible IP ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*flexibleip.GetFlexibleIPRequest)

			client := core.ExtractClient(ctx)
			api := flexibleip.NewAPI(client)
			return api.GetFlexibleIP(request)

		},
	}
}

func fipIPList() *core.Command {
	return &core.Command{
		Short:     `List Flexible IPs`,
		Long:      `List Flexible IPs.`,
		Namespace: "fip",
		Resource:  "ip",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(flexibleip.ListFlexibleIPsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `The sort order of the returned Flexible IPs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter Flexible IPs with one or more matching tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status.{index}",
				Short:      `Filter Flexible IPs by status`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "ready", "updating", "attached", "error", "detaching", "locked"},
			},
			{
				Name:       "server-ids.{index}",
				Short:      `Filter Flexible IPs by server IDs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Filter Flexible IPs by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter Flexible IPs by organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*flexibleip.ListFlexibleIPsRequest)

			client := core.ExtractClient(ctx)
			api := flexibleip.NewAPI(client)
			resp, err := api.ListFlexibleIPs(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.FlexibleIPs, nil

		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "IPAddress",
			},
			{
				FieldName: "Status",
			},
			{
				FieldName: "Reverse",
			},
			{
				FieldName: "ServerID",
			},
			{
				FieldName: "Description",
			},
			{
				FieldName: "Tags",
			},
			{
				FieldName: "ProjectID",
			},
			{
				FieldName: "OrganizationID",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "CreatedAt",
			},
		}},
	}
}

func fipIPUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a Flexible IP`,
		Long:      `Update a Flexible IP.`,
		Namespace: "fip",
		Resource:  "ip",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(flexibleip.UpdateFlexibleIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "fip-id",
				Short:      `ID of the Flexible IP to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "description",
				Short:      `Description to associate with the Flexible IP, max 255 characters`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to associate with the Flexible IP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "reverse",
				Short:      `Reverse DNS value`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*flexibleip.UpdateFlexibleIPRequest)

			client := core.ExtractClient(ctx)
			api := flexibleip.NewAPI(client)
			return api.UpdateFlexibleIP(request)

		},
	}
}

func fipIPDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a Flexible IP`,
		Long:      `Delete a Flexible IP.`,
		Namespace: "fip",
		Resource:  "ip",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(flexibleip.DeleteFlexibleIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "fip-id",
				Short:      `ID of the Flexible IP to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*flexibleip.DeleteFlexibleIPRequest)

			client := core.ExtractClient(ctx)
			api := flexibleip.NewAPI(client)
			e = api.DeleteFlexibleIP(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "ip",
				Verb:     "delete",
			}, nil
		},
	}
}

func fipIPAttach() *core.Command {
	return &core.Command{
		Short:     `Attach a Flexible IP to a server`,
		Long:      `Attach a Flexible IP to a server.`,
		Namespace: "fip",
		Resource:  "ip",
		Verb:      "attach",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(flexibleip.AttachFlexibleIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "fips-ids.{index}",
				Short:      `A list of Flexible IP IDs to attach`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-id",
				Short:      `A server ID on which to attach the Flexible IPs`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*flexibleip.AttachFlexibleIPRequest)

			client := core.ExtractClient(ctx)
			api := flexibleip.NewAPI(client)
			return api.AttachFlexibleIP(request)

		},
	}
}

func fipIPDetach() *core.Command {
	return &core.Command{
		Short:     `Detach a Flexible IP from a server`,
		Long:      `Detach a Flexible IP from a server.`,
		Namespace: "fip",
		Resource:  "ip",
		Verb:      "detach",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(flexibleip.DetachFlexibleIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "fips-ids.{index}",
				Short:      `A list of Flexible IP IDs to detach`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*flexibleip.DetachFlexibleIPRequest)

			client := core.ExtractClient(ctx)
			api := flexibleip.NewAPI(client)
			return api.DetachFlexibleIP(request)

		},
	}
}

func fipMacCreate() *core.Command {
	return &core.Command{
		Short:     `Generate a virtual MAC on a given Flexible IP`,
		Long:      `Generate a virtual MAC on a given Flexible IP.`,
		Namespace: "fip",
		Resource:  "mac",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(flexibleip.GenerateMACAddrRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "fip-id",
				Short:      `Flexible IP ID on which to generate a Virtual MAC`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mac-type",
				Short:      `TODO`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_type", "vmware", "xen", "kvm"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*flexibleip.GenerateMACAddrRequest)

			client := core.ExtractClient(ctx)
			api := flexibleip.NewAPI(client)
			return api.GenerateMACAddr(request)

		},
	}
}

func fipMacDuplicate() *core.Command {
	return &core.Command{
		Short:     `Duplicate a Virtual MAC`,
		Long:      `Duplicate a Virtual MAC from a given Flexible IP onto another attached on the same server.`,
		Namespace: "fip",
		Resource:  "mac",
		Verb:      "duplicate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(flexibleip.DuplicateMACAddrRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "fip-id",
				Short:      `Flexible IP ID on which to duplicate the Virtual MAC`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "duplicate-from-fip-id",
				Short:      `Flexible IP ID to duplicate the Virtual MAC from`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*flexibleip.DuplicateMACAddrRequest)

			client := core.ExtractClient(ctx)
			api := flexibleip.NewAPI(client)
			return api.DuplicateMACAddr(request)

		},
	}
}

func fipMacDelete() *core.Command {
	return &core.Command{
		Short:     `Remove a virtual MAC from a Flexible IP`,
		Long:      `Remove a virtual MAC from a Flexible IP.`,
		Namespace: "fip",
		Resource:  "mac",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(flexibleip.DeleteMACAddrRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "fip-id",
				Short:      `Flexible IP ID from which to delete the Virtual MAC`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*flexibleip.DeleteMACAddrRequest)

			client := core.ExtractClient(ctx)
			api := flexibleip.NewAPI(client)
			e = api.DeleteMACAddr(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "mac",
				Verb:     "delete",
			}, nil
		},
	}
}
