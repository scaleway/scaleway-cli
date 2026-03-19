// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package product_catalog

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	product_catalog "github.com/scaleway/scaleway-sdk-go/api/product_catalog/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		productCatalogRoot(),
		productCatalogProduct(),
		productCatalogProductList(),
	)
}

func productCatalogRoot() *core.Command {
	return &core.Command{
		Short:     `Product Catalog API`,
		Long:      ``,
		Namespace: "product-catalog",
	}
}

func productCatalogProduct() *core.Command {
	return &core.Command{
		Short: `Scaleway Product Catalog API`,
		Long: `Scaleway's Product Catalog is an extensive list of the Scaleway products.
The catalog includes details about each product including: description,
locations, prices and properties.`,
		Namespace: "product-catalog",
		Resource:  "product",
	}
}

func productCatalogProductList() *core.Command {
	return &core.Command{
		Short:     `List all available products`,
		Long:      `List all available products in the Scaleway catalog. Returns a complete list of products with their corresponding description, locations, prices and properties. You can define the ` + "`" + `page` + "`" + ` number and ` + "`" + `page_size` + "`" + ` for your query in the request.`,
		Namespace: "product-catalog",
		Resource:  "product",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(
			product_catalog.PublicCatalogAPIListPublicCatalogProductsRequest{},
		),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "product-types.{index}",
				Short:      `The list of filtered product categories.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_product_type",
					"instance",
					"apple_silicon",
					"elastic_metal",
					"dedibox",
					"block_storage",
					"object_storage",
					"managed_inference",
					"generative_apis",
					"load_balancer",
					"secret_manager",
					"key_manager",
					"managed_redis_database",
					"kubernetes",
					"managed_relational_database",
					"managed_mongodb",
				},
			},
			{
				Name:       "global",
				Short:      `Filter global products.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "region",
				Short:      `Filter products by region.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "zone",
				Short:      `Filter products by zone.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "datacenter",
				Short:      `Filter products by datacenter.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status.{index}",
				Short:      `The lists of filtered product status, if empty only products with status public_beta, general_availability, preview, end_of_new_features, end_of_growth, end_of_deployment, end_of_support, end_of_sale, end_of_life or retired will be returned.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_status",
					"public_beta",
					"preview",
					"general_availability",
					"end_of_new_features",
					"end_of_growth",
					"end_of_deployment",
					"end_of_support",
					"end_of_sale",
					"end_of_life",
					"retired",
				},
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*product_catalog.PublicCatalogAPIListPublicCatalogProductsRequest)

			client := core.ExtractClient(ctx)
			api := product_catalog.NewPublicCatalogAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListPublicCatalogProducts(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Products, nil
		},
	}
}
