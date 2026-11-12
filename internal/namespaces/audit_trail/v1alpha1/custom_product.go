package audit_trail

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	audit_trail "github.com/scaleway/scaleway-sdk-go/api/audit_trail/v1alpha1"
)

// productRow flattens a Product into one row per logged method. Keeping
// every cell short prevents the table's width-based column shrinking from
// dropping the methods (a single "all methods" cell is hundreds of chars
// wide and would disappear on a normal terminal).
type productRow struct {
	Name    string `json:"name"`
	Service string `json:"service"`
	Method  string `json:"method"`
}

func productListBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "Name",
				FieldName: "Name",
			},
			{
				Label:     "Service",
				FieldName: "Service",
			},
			{
				Label:     "Method",
				FieldName: "Method",
			},
		},
	}

	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
			originalRes, err := runner(ctx, argsI)
			if err != nil {
				return nil, err
			}

			productsResponse := originalRes.(*audit_trail.ListProductsResponse)

			rows := make([]productRow, 0, len(productsResponse.Products))
			for _, product := range productsResponse.Products {
				if len(product.Services) == 0 {
					rows = append(rows, productRow{Name: product.Name})

					continue
				}
				for _, service := range product.Services {
					if len(service.Methods) == 0 {
						rows = append(rows, productRow{Name: product.Name, Service: service.Name})

						continue
					}
					for _, method := range service.Methods {
						rows = append(rows, productRow{
							Name:    product.Name,
							Service: service.Name,
							Method:  method,
						})
					}
				}
			}

			return rows, nil
		},
	)

	return c
}
