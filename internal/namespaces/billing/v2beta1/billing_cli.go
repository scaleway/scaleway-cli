// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package billing

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2beta1"
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
		billingConsumption(),
		billingDiscount(),
		billingConsumptionList(),
		billingConsumptionListTaxes(),
		billingInvoiceList(),
		billingInvoiceExport(),
		billingInvoiceGet(),
		billingInvoiceDownload(),
		billingDiscountList(),
	)
}

func billingRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage and query your Scaleway billing and consumption`,
		Long:      `This API allows you to manage and query your Scaleway billing and consumption.`,
		Namespace: "billing",
	}
}

func billingInvoice() *core.Command {
	return &core.Command{
		Short:     `Invoice management commands`,
		Long:      `Invoice management commands.`,
		Namespace: "billing",
		Resource:  "invoice",
	}
}

func billingConsumption() *core.Command {
	return &core.Command{
		Short:     `Consumption management commands`,
		Long:      `Consumption management commands.`,
		Namespace: "billing",
		Resource:  "consumption",
	}
}

func billingDiscount() *core.Command {
	return &core.Command{
		Short:     `Discount management commands`,
		Long:      `Discount management commands.`,
		Namespace: "billing",
		Resource:  "discount",
	}
}

func billingConsumptionList() *core.Command {
	return &core.Command{
		Short:     `Get monthly consumption`,
		Long:      `Consumption allows you to retrieve your past or current consumption cost, by project or category.`,
		Namespace: "billing",
		Resource:  "consumption",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.ListConsumptionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order consumptions list in the response by their update date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"updated_at_desc",
					"updated_at_asc",
					"category_name_desc",
					"category_name_asc",
				},
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "category-name",
				Short:      `Filter by name of a Category as they are shown in the invoice (Compute, Network, Observability)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "billing-period",
				Short:      `Filter by the billing period in the YYYY-MM format. If it is empty the current billing period will be used as default`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*billing.ListConsumptionsRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListConsumptions(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Consumptions, nil
		},
	}
}

func billingConsumptionListTaxes() *core.Command {
	return &core.Command{
		Short:     `Get monthly consumption taxes`,
		Long:      `Consumption Tax allows you to retrieve your past or current tax charges, by project or category.`,
		Namespace: "billing",
		Resource:  "consumption",
		Verb:      "list-taxes",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.ListTaxesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order consumed taxes list in the response by their update date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"updated_at_desc",
					"updated_at_asc",
					"category_name_desc",
					"category_name_asc",
				},
			},
			{
				Name:       "billing-period",
				Short:      `Filter by the billing period in the YYYY-MM format. If it is empty the current billing period will be used as default`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*billing.ListTaxesRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListTaxes(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Taxes, nil
		},
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
				Name:       "billing-period-start-after",
				Short:      `Return only invoice with start date greater than billing_period_start`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "billing-period-start-before",
				Short:      `Return only invoice with start date less than billing_period_start`,
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
				EnumValues: []string{
					"unknown_type",
					"periodic",
					"purchase",
				},
			},
			{
				Name:       "order-by",
				Short:      `How invoices are ordered in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"invoice_number_desc",
					"invoice_number_asc",
					"start_date_desc",
					"start_date_asc",
					"issued_date_desc",
					"issued_date_asc",
					"due_date_desc",
					"due_date_asc",
					"total_untaxed_desc",
					"total_untaxed_asc",
					"total_taxed_desc",
					"total_taxed_asc",
					"invoice_type_desc",
					"invoice_type_asc",
				},
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID. If specified, only invoices from this Organization will be returned`,
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

func billingInvoiceExport() *core.Command {
	return &core.Command{
		Short:     `Export invoices`,
		Long:      `Export invoices in a CSV file.`,
		Namespace: "billing",
		Resource:  "invoice",
		Verb:      "export",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.ExportInvoicesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "billing-period-start-after",
				Short:      `Return only invoice with start date greater than billing_period_start`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "billing-period-start-before",
				Short:      `Return only invoice with start date less than billing_period_start`,
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
				EnumValues: []string{
					"unknown_type",
					"periodic",
					"purchase",
				},
			},
			{
				Name:       "page",
				Short:      `Page number`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "page-size",
				Short:      `Positive integer lower or equal to 100 to select the number of items to return`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("20"),
			},
			{
				Name:       "order-by",
				Short:      `How invoices are ordered in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"invoice_number_desc",
					"invoice_number_asc",
					"start_date_desc",
					"start_date_asc",
					"issued_date_desc",
					"issued_date_asc",
					"due_date_desc",
					"due_date_asc",
					"total_untaxed_desc",
					"total_untaxed_asc",
					"total_taxed_desc",
					"total_taxed_asc",
					"invoice_type_desc",
					"invoice_type_asc",
				},
			},
			{
				Name:       "file-type",
				Short:      `File format for exporting the invoice list`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("CSV"),
				EnumValues: []string{
					"csv",
				},
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID. If specified, only invoices from this Organization will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*billing.ExportInvoicesRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)

			return api.ExportInvoices(request)
		},
	}
}

func billingInvoiceGet() *core.Command {
	return &core.Command{
		Short:     `Get an invoice`,
		Long:      `Get a specific invoice, specified by its ID.`,
		Namespace: "billing",
		Resource:  "invoice",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.GetInvoiceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "invoice-id",
				Short:      `Invoice ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*billing.GetInvoiceRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)

			return api.GetInvoice(request)
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
				Positional: true,
			},
			{
				Name:       "file-type",
				Short:      `File type. PDF by default`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"pdf",
				},
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

func billingDiscountList() *core.Command {
	return &core.Command{
		Short: `List discounts`,
		Long: `List all discounts for your Organization and usable categories, products, offers, references, regions and zones where the discount can be applied. As a reseller:
- If you do not specify an ` + "`" + `organization_id` + "`" + ` you will list the discounts applied to your own Organization and your customers
- If you indicate your ` + "`" + `organization_id` + "`" + ` you will list only the discounts applied to your Organization
- If you indicate ` + "`" + `the organization_id` + "`" + ` of one of your customers, you will list the discounts applied to their Organization.`,
		Namespace: "billing",
		Resource:  "discount",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.ListDiscountsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order discounts in the response by their description`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"creation_date_desc",
					"creation_date_asc",
					"start_date_desc",
					"start_date_asc",
					"stop_date_desc",
					"stop_date_asc",
				},
			},
			{
				Name:       "organization-id",
				Short:      `ID of the organization`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*billing.ListDiscountsRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListDiscounts(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Discounts, nil
		},
	}
}
