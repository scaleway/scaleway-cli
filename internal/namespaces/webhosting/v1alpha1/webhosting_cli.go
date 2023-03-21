// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package webhosting

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/webhosting/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		webhostingRoot(),
		webhostingHosting(),
		webhostingOffer(),
		webhostingHostingCreate(),
		webhostingHostingList(),
		webhostingHostingGet(),
		webhostingHostingUpdate(),
		webhostingHostingDelete(),
		webhostingHostingRestore(),
		webhostingHostingGetDNSRecords(),
		webhostingOfferList(),
	)
}
func webhostingRoot() *core.Command {
	return &core.Command{
		Short:     `Webhosting API`,
		Long:      `Webhosting API.`,
		Namespace: "webhosting",
	}
}

func webhostingHosting() *core.Command {
	return &core.Command{
		Short: `Hosting management commands`,
		Long: `A Scaleway web hosting associated with a domain name.
`,
		Namespace: "webhosting",
		Resource:  "hosting",
	}
}

func webhostingOffer() *core.Command {
	return &core.Command{
		Short: `Offer management commands`,
		Long: `A hosting offer, with a set of features, available for purchase.
`,
		Namespace: "webhosting",
		Resource:  "offer",
	}
}

func webhostingHostingCreate() *core.Command {
	return &core.Command{
		Short:     `Create a hosting`,
		Long:      `Create a hosting.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.CreateHostingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "offer-id",
				Short:      `ID of the selected offer for the hosting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "email",
				Short:      `Contact email of the client for the hosting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The tags of the hosting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain",
				Short:      `The domain name of the hosting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "option-ids.{index}",
				Short:      `IDs of the selected options for the hosting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*webhosting.CreateHostingRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewAPI(client)
			return api.CreateHosting(request)

		},
	}
}

func webhostingHostingList() *core.Command {
	return &core.Command{
		Short:     `List all hostings`,
		Long:      `List all hostings.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.ListHostingsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Define the order of the returned hostings`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "tags.{index}",
				Short:      `Return hostings with these tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "statuses.{index}",
				Short:      `Return hostings with these statuses`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_status", "delivering", "ready", "deleting", "error", "locked"},
			},
			{
				Name:       "domain",
				Short:      `Return hostings with this domain`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Return hostings from this project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Return hostings from this organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*webhosting.ListHostingsRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListHostings(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Hostings, nil

		},
		Examples: []*core.Example{
			{
				Short:    "List all hostings of a given project ID",
				ArgsJSON: `{"organization_id":"a3244331-5d32-4e36-9bf9-b60233e201c7","project_id":"a3244331-5d32-4e36-9bf9-b60233e201c7"}`,
			},
		},
	}
}

func webhostingHostingGet() *core.Command {
	return &core.Command{
		Short:     `Get a hosting`,
		Long:      `Get the details of a Hosting with the given ID.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.GetHostingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `Hosting ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*webhosting.GetHostingRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewAPI(client)
			return api.GetHosting(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get a Hosting with the given ID",
				ArgsJSON: `{"hosting_id":"a3244331-5d32-4e36-9bf9-b60233e201c7"}`,
			},
		},
	}
}

func webhostingHostingUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a hosting`,
		Long:      `Update a hosting.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.UpdateHostingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `Hosting ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "email",
				Short:      `New contact email for the hosting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `New tags for the hosting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "option-ids.{index}",
				Short:      `New options IDs for the hosting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "offer-id",
				Short:      `New offer ID for the hosting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*webhosting.UpdateHostingRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewAPI(client)
			return api.UpdateHosting(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Update the contact email of a given hosting",
				ArgsJSON: `{"email":"foobar@example.com","hosting_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Overwrite tags of a given hosting",
				ArgsJSON: `{"hosting_id":"11111111-1111-1111-1111-111111111111","tags":["foo","bar"]}`,
			},
			{
				Short:    "Overwrite options of a given hosting",
				ArgsJSON: `{"hosting_id":"11111111-1111-1111-1111-111111111111","option_ids":["22222222-2222-2222-2222-222222222222","33333333-3333-3333-3333-333333333333"]}`,
			},
		},
	}
}

func webhostingHostingDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a hosting`,
		Long:      `Delete a hosting with the given ID.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DeleteHostingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `Hosting ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*webhosting.DeleteHostingRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewAPI(client)
			return api.DeleteHosting(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Delete a Hosting with the given ID",
				ArgsJSON: `{"hosting_id":"a3244331-5d32-4e36-9bf9-b60233e201c7"}`,
			},
		},
	}
}

func webhostingHostingRestore() *core.Command {
	return &core.Command{
		Short:     `Restore a hosting`,
		Long:      `Restore a hosting with the given ID.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "restore",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.RestoreHostingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `Hosting ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*webhosting.RestoreHostingRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewAPI(client)
			return api.RestoreHosting(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Restore a Hosting with the given ID",
				ArgsJSON: `{"hosting_id":"a3244331-5d32-4e36-9bf9-b60233e201c7"}`,
			},
		},
	}
}

func webhostingHostingGetDNSRecords() *core.Command {
	return &core.Command{
		Short:     `Get the DNS records`,
		Long:      `The set of DNS record of a specific domain associated to a hosting.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "get-dns-records",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.GetDomainDNSRecordsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Short:      `Domain associated to the DNS records`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*webhosting.GetDomainDNSRecordsRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewAPI(client)
			return api.GetDomainDNSRecords(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get DNS records associated to the given domain",
				ArgsJSON: `{"domain":"foo.com"}`,
			},
		},
	}
}

func webhostingOfferList() *core.Command {
	return &core.Command{
		Short:     `List all offers`,
		Long:      `List all offers.`,
		Namespace: "webhosting",
		Resource:  "offer",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.ListOffersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Define the order of the returned hostings`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"price_asc"},
			},
			{
				Name:       "without-options",
				Short:      `Select only offers, no options`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "only-options",
				Short:      `Select only options`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "hosting-id",
				Short:      `Define a specific hosting id (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*webhosting.ListOffersRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewAPI(client)
			return api.ListOffers(request)

		},
		Examples: []*core.Example{
			{
				Short:    "List all offers available for purchase",
				ArgsJSON: `{"hosting_id":"a3244331-5d32-4e36-9bf9-b60233e201c7","only_options":false,"without_options":false}`,
			},
			{
				Short:    "List only offers, no options",
				ArgsJSON: `{"only_options":false,"without_options":true}`,
			},
			{
				Short:    "List only options",
				ArgsJSON: `{"only_options":true,"without_options":false}`,
			},
		},
	}
}
