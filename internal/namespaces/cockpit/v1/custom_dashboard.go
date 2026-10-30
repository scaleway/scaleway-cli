package cockpit

import (
	"context"
	"fmt"
	"net/url"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/skratchdot/open-golang/open"
)

// OpenDashboard opens the Grafana product dashboard for a resource and
// returns the "Opening <url>" message to print.
//
// dashboardName is the product dashboard (e.g. "scw-mongodb-overview").
// vars maps each dashboard variable to its value (e.g. {"region": "fr-par",
// "cl_name": "mydb"}); empty values are skipped. Each product sets exactly the
// variables its dashboard expects.
func OpenDashboard(
	ctx context.Context,
	dashboardName string,
	vars map[string]string,
) (any, error) {
	// Fetch the product dashboard to get its (org-specific) base URL.
	api := cockpit.NewGlobalAPI(core.ExtractClient(ctx))
	dash, err := api.GetGrafanaProductDashboard(
		&cockpit.GlobalAPIGetGrafanaProductDashboardRequest{
			DashboardName: dashboardName,
			ProjectID:     core.GetProjectIDFromContext(ctx),
		},
		scw.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("orgId", "1")
	q.Set("from", "now-1h")
	q.Set("to", "now")
	q.Set("timezone", "browser")
	for key, value := range vars {
		if value != "" {
			q.Set("var-"+key, value)
		}
	}
	fullURL := dash.URL + "?" + q.Encode()

	if err := open.Start(fullURL); err != nil {
		return nil, &core.CliError{
			Err:     err,
			Message: "Failed to open cockpit url",
			Details: "You can open it: " + fullURL,
			Hint:    "You may not have a default browser configured",
			Code:    1,
		}
	}

	return fmt.Sprintf("Opening %s\n", fullURL), nil
}
