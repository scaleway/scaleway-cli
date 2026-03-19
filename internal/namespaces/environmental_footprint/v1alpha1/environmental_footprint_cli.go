// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package environmental_footprint

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	environmental_footprint "github.com/scaleway/scaleway-sdk-go/api/environmental_footprint/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		environmentalFootprintRoot(),
		environmentalFootprintReport(),
		environmentalFootprintData(),
		environmentalFootprintReportList(),
		environmentalFootprintReportGet(),
		environmentalFootprintDataGet(),
	)
}

func environmentalFootprintRoot() *core.Command {
	return &core.Command{
		Short:     `Access and download impact reports and impact data for your Scaleway projects. Our API provides key metrics such as estimated carbon emissions and water usage to help monitor your environmental footprint.`,
		Long:      `Access and download impact reports and impact data for your Scaleway projects. Our API provides key metrics such as estimated carbon emissions and water usage to help monitor your environmental footprint.`,
		Namespace: "environmental-footprint",
	}
}

func environmentalFootprintReport() *core.Command {
	return &core.Command{
		Short:     `Environmental impact report management commands`,
		Long:      `Environmental impact report management commands.`,
		Namespace: "environmental-footprint",
		Resource:  "report",
	}
}

func environmentalFootprintData() *core.Command {
	return &core.Command{
		Short:     `Environmental impact data management commands`,
		Long:      `Environmental impact data management commands.`,
		Namespace: "environmental-footprint",
		Resource:  "data",
	}
}

func environmentalFootprintReportList() *core.Command {
	return &core.Command{
		Short:     `Get available impact reports`,
		Long:      `Returns a list of dates of available impact reports.`,
		Namespace: "environmental-footprint",
		Resource:  "report",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(
			environmental_footprint.UserAPIGetImpactReportAvailabilityRequest{},
		),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "start-date",
				Short:      `Start date of the search period (ISO 8601 format, with time in UTC, ` + "`" + `YYYY-MM-DDTHH:MM:SSZ` + "`" + `). The date is inclusive.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "end-date",
				Short:      `End date of the search period (ISO 8601 format, with time in UTC, ` + "`" + `YYYY-MM-DDTHH:MM:SSZ` + "`" + `). The date is inclusive. Defaults to today's date.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*environmental_footprint.UserAPIGetImpactReportAvailabilityRequest)

			client := core.ExtractClient(ctx)
			api := environmental_footprint.NewUserAPI(client)

			return api.GetImpactReportAvailability(request)
		},
	}
}

func environmentalFootprintReportGet() *core.Command {
	return &core.Command{
		Short:     `Download PDF impact report`,
		Long:      `Download a Scaleway impact PDF report with detailed impact data for your Scaleway projects.`,
		Namespace: "environmental-footprint",
		Resource:  "report",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(environmental_footprint.UserAPIDownloadImpactReportRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "date",
				Short:      `The start date of the period for which you want to download a report (ISO 8601 format, e.g. 2025-05-01T00:00:00Z).`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Type of report to download (e.g. ` + "`" + `monthly` + "`" + `).`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_report_type",
					"monthly",
					"yearly",
				},
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*environmental_footprint.UserAPIDownloadImpactReportRequest)

			client := core.ExtractClient(ctx)
			api := environmental_footprint.NewUserAPI(client)

			return api.DownloadImpactReport(request)
		},
	}
}

func environmentalFootprintDataGet() *core.Command {
	return &core.Command{
		Short:     `Retrieve detailed impact data`,
		Long:      `Retrieve detailed impact data for your Scaleway projects within a specified date range. Filter by project ID, region, zone, service category, and/or product category.`,
		Namespace: "environmental-footprint",
		Resource:  "data",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(environmental_footprint.UserAPIGetImpactDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "start-date",
				Short:      `Start date (inclusive) of the period for which you want to retrieve impact data (ISO 8601 format, e.g. 2025-05-01T00:00:00Z).`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "end-date",
				Short:      `End date (exclusive) of the period for which you want to retrieve impact data (ISO 8601 format, with time in UTC, ` + "`" + `YYYY-MM-DDTHH:MM:SSZ` + "`" + `). Defaults to today's date.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "regions.{index}",
				Short:      `List of regions to filter by (e.g. ` + "`" + `fr-par` + "`" + `). Defaults to all regions.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "zones.{index}",
				Short:      `List of zones to filter by (e.g. ` + "`" + `fr-par-1` + "`" + `). Defaults to all zones.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-ids.{index}",
				Short:      `List of Project IDs to filter by. Defaults to all Projects in the Organization.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "service-categories.{index}",
				Short:      `List of service categories to filter by. Defaults to all service categories.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_service_category",
					"baremetal",
					"compute",
					"storage",
					"network",
					"containers",
					"databases",
				},
			},
			{
				Name:       "product-categories.{index}",
				Short:      `List of product categories to filter by. Defaults to all product categories.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_product_category",
					"apple_silicon",
					"block_storage",
					"dedibox",
					"elastic_metal",
					"instances",
					"object_storage",
					"load_balancer",
					"kubernetes",
					"managed_relational_databases",
					"managed_mongodb",
					"managed_redis",
				},
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*environmental_footprint.UserAPIGetImpactDataRequest)

			client := core.ExtractClient(ctx)
			api := environmental_footprint.NewUserAPI(client)

			return api.GetImpactData(request)
		},
	}
}
