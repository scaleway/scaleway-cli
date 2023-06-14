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
		Short:     `Web Hosting API`,
		Long:      `Web Hosting API.`,
		Namespace: "webhosting",
	}
}

func webhostingHosting() *core.Command {
	return &core.Command{
		Short:     `Hosting management commands`,
		Long:      `With a Scaleway Web Hosting plan, you can manage your domain, configure your web hosting services, manage your emails and more. Create, list, update and delete your Web Hosting plans with these calls.`,
		Namespace: "webhosting",
		Resource:  "hosting",
	}
}

func webhostingOffer() *core.Command {
	return &core.Command{
		Short:     `Offer management commands`,
		Long:      `Web Hosting offers represent the different types of Web Hosting plan available to order at Scaleway.`,
		Namespace: "webhosting",
		Resource:  "offer",
	}
}

func webhostingHostingCreate() *core.Command {
	return &core.Command{
		Short:     `Order a Web Hosting plan`,
		Long:      `Order a Web Hosting plan, specifying the offer type required via the ` + "`" + `offer_id` + "`" + ` parameter.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.CreateHostingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "offer-id",
				Short:      `ID of the selected offer for the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "email",
				Short:      `Contact email for the Web Hosting client`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags for the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain",
				Short:      `Domain name to link to the Web Hosting plan. You must already own this domain name, and have completed the DNS validation process beforehand`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "option-ids.{index}",
				Short:      `IDs of any selected additional options for the Web Hosting plan`,
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
		Short:     `List all Web Hosting plans`,
		Long:      `List all of your existing Web Hosting plans. Various filters are available to limit the results, including filtering by domain, status, tag and Project ID.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.ListHostingsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order for Web Hosting plans in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to filter for, only Web Hosting plans with matching tags will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "statuses.{index}",
				Short:      `Statuses to filter for, only Web Hosting plans with matching statuses will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_status", "delivering", "ready", "deleting", "error", "locked", "migrating"},
			},
			{
				Name:       "domain",
				Short:      `Domain to filter for, only Web Hosting plans associated with this domain will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Project ID to filter for, only Web Hosting plans from this Project will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID to filter for, only Web Hosting plans from this Organization will be returned`,
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
		Short:     `Get a Web Hosting plan`,
		Long:      `Get the details of one of your existing Web Hosting plans, specified by its ` + "`" + `hosting_id` + "`" + `.`,
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
		Short:     `Update a Web Hosting plan`,
		Long:      `Update the details of one of your existing Web Hosting plans, specified by its ` + "`" + `hosting_id` + "`" + `. You can update parameters including the contact email address, tags, options and offer.`,
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
				Short:      `New contact email for the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `New tags for the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "option-ids.{index}",
				Short:      `IDs of the new options for the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "offer-id",
				Short:      `ID of the new offer for the Web Hosting plan`,
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
		Short:     `Delete a Web Hosting plan`,
		Long:      `Delete a Web Hosting plan, specified by its ` + "`" + `hosting_id` + "`" + `. Note that deletion is not immediate: it will take place at the end of the calendar month, after which time your Web Hosting plan and all its data (files and emails) will be irreversibly lost.`,
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
		Short:     `Restore a Web Hosting plan`,
		Long:      `When you [delete a Web Hosting plan](#path-hostings-delete-a-hosting), definitive deletion does not take place until the end of the calendar month. In the time between initiating the deletion, and definitive deletion at the end of the month, you can choose to **restore** the Web Hosting plan, using this endpoint and specifying its ` + "`" + `hosting_id` + "`" + `.`,
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
		Short:     `Get DNS records`,
		Long:      `Get the set of DNS records of a specified domain associated with a Web Hosting plan.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "get-dns-records",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.GetDomainDNSRecordsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Short:      `Domain associated with the DNS records`,
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
		Long:      `List the different Web Hosting offers, and their options, available to order from Scaleway.`,
		Namespace: "webhosting",
		Resource:  "offer",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.ListOffersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of offers in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"price_asc"},
			},
			{
				Name:       "without-options",
				Short:      `Defines whether the response should consist of offers only, without options`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "only-options",
				Short:      `Defines whether the response should consist of options only, without offers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "hosting-id",
				Short:      `ID of a Web Hosting plan, to check compatibility with returned offers (in case of wanting to update the plan)`,
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
