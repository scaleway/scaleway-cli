// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package billing

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/billing/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		billingRoot(),
		billingInvoice(),
		billingDiscount(),
		billingInvoiceList(),
		billingInvoiceDownload(),
	)
}
func billingRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to query your consumption`,
		Long:      `This API allows you to query your consumption.`,
		Namespace: "billing",
	}
}

func billingInvoice() *core.Command {
	return &core.Command{
		Short:     `Invoices management commands`,
		Long:      `Invoices management commands.`,
		Namespace: "billing",
		Resource:  "invoice",
	}
}

func billingDiscount() *core.Command {
	return &core.Command{
		Short:     `Discounts management commands`,
		Long:      `Discounts management commands.`,
		Namespace: "billing",
		Resource:  "discount",
	}
}

func billingInvoiceList() *core.Command {
	return &core.Command{
		Short:     `List invoices`,
		Long:      `List all your invoices, filtering by ` + "`" + `start_date` + "`" + ` and ` + "`" + `invoice_type` + "`" + `. Each invoice has its own ID.`,
		Namespace: "billing",
		Resource:  "invoice",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.ListInvoicesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "started-after",
				Short:      `Invoice's ` + "`" + `start_date` + "`" + ` is greater or equal to ` + "`" + `started_after` + "`" + ``,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "started-before",
				Short:      `Invoice's ` + "`" + `start_date` + "`" + ` precedes ` + "`" + `started_before` + "`" + ``,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "invoice-type",
				Short:      `Invoice type. It can either be ` + "`" + `periodic` + "`" + ` or ` + "`" + `purchase` + "`" + ``,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_type", "periodic", "purchase"},
			},
			{
				Name:       "order-by",
				Short:      `How invoices are ordered in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"invoice_number_desc", "invoice_number_asc", "start_date_desc", "start_date_asc", "issued_date_desc", "issued_date_asc", "due_date_desc", "due_date_asc", "total_untaxed_desc", "total_untaxed_asc", "total_taxed_desc", "total_taxed_asc", "invoice_type_desc", "invoice_type_asc"},
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID to filter for, only invoices from this Organization will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*billing.ListInvoicesRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListInvoices(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Invoices, nil

		},
	}
}

func billingInvoiceDownload() *core.Command {
	return &core.Command{
		Short:     `Download an invoice`,
		Long:      `Download a specific invoice, specified by its ID.`,
		Namespace: "billing",
		Resource:  "invoice",
		Verb:      "download",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.DownloadInvoiceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "invoice-id",
				Short:      `Invoice ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "file-type",
				Short:      `Wanted file type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"pdf"},
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*billing.DownloadInvoiceRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)
			return api.DownloadInvoice(request)

		},
	}
}
