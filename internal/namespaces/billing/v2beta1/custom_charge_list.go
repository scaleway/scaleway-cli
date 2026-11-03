package billing

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// chargeListBuilder overrides the generated `billing charge list` command:
//   - the generated Run returns the raw ListChargesResponse, which renders as a flat
//     key/value dump; here we return the charges slice so it renders as a table
//     like the other billing list commands.
//   - the API is token-paginated, so we follow NextPageToken to return all charges.
func chargeListBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Fields: []*core.ViewField{
			{Label: "Resource Name", FieldName: "ResourceName"},
			{Label: "Resource ID", FieldName: "ResourceID"},
			{Label: "Price", FieldName: "Price"},
			{Label: "SKU", FieldName: "Sku"},
			{Label: "Start Date", FieldName: "StartDate"},
			{Label: "End Date", FieldName: "EndDate"},
			{Label: "Invoice ID", FieldName: "InvoiceID"},
			{Label: "Project Name", FieldName: "ProjectName"},
			{Label: "Project ID", FieldName: "ProjectID"},
		},
	}

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		request := argsI.(*billing.FinOpsAPIListChargesRequest)
		api := billing.NewFinOpsAPI(core.ExtractClient(ctx))

		charges := []*billing.Charge(nil)
		for {
			resp, err := api.ListCharges(request, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}
			charges = append(charges, resp.Charges...)
			if resp.NextPageToken == nil || *resp.NextPageToken == "" {
				return charges, nil
			}
			request.PageToken = resp.NextPageToken
		}
	}

	return c
}
