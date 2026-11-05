package billing

import (
	"context"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// ChargesView is the table view for lists of *billing.Charge, shared by
// `billing charge list` and the per-resource `charge` verbs.
func ChargesView() *core.View {
	return &core.View{
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
}

func BeginningOfMonth(d *time.Time) *time.Time {
	if d == nil {
		return nil
	}
	start := time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location())

	return &start
}

// ListChargesForResource returns all the charges of a resource, following
// pagination. Since, when non-nil (typically the resource's created_at),
// restricts the result to charges starting after that date.
func ListChargesForResource(
	ctx context.Context,
	resourceID string,
	since *time.Time,
) ([]*billing.Charge, error) {
	request := &billing.FinOpsAPIListChargesRequest{
		ResourceIDs: []string{resourceID},
	}
	if projectID, ok := core.ExtractClient(ctx).GetDefaultProjectID(); ok {
		request.ProjectIDs = append(request.ProjectIDs, projectID)
	}
	if orgID, ok := core.ExtractClient(ctx).GetDefaultOrganizationID(); ok {
		request.OrganizationID = orgID
	}
	if since != nil {
		request.StartDateAfter = BeginningOfMonth(since)
		// Clamp so a charge overlapping the period (e.g. the first one, which
		// may start before the resource's created_at) is still returned.
		clampToTimeRange := true
		request.ClampToTimeRange = &clampToTimeRange
	}

	return listAllCharges(ctx, request)
}

// chargeListBuilder overrides the generated `billing charge list` command:
//   - the generated Run returns the raw ListChargesResponse, which renders as a flat
//     key/value dump; here we return the charges slice so it renders as a table
//     like the other billing list commands.
//   - the API is token-paginated, so we follow NextPageToken to return all charges.
func chargeListBuilder(c *core.Command) *core.Command {
	c.View = ChargesView()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		request := argsI.(*billing.FinOpsAPIListChargesRequest)

		return listAllCharges(ctx, request)
	}

	return c
}

func listAllCharges(
	ctx context.Context,
	request *billing.FinOpsAPIListChargesRequest,
) ([]*billing.Charge, error) {
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
