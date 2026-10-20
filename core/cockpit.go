package core

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"

	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/skratchdot/open-golang/open"
)

// OpenCockpitDashboard opens the Grafana product dashboard for a resource and
// returns the "Opening <url>" message to print.
//
// dashboardName is the product dashboard (e.g. "scw-mongodb-overview"),
// resourceVar the dashboard variable that receives the resource's name
// (e.g. "cl_name"), and resource the fetched resource (any struct with Name
// and Region/Zone fields).
func OpenCockpitDashboard(
	ctx context.Context,
	dashboardName, resourceVar string,
	resource any,
) (any, error) {
	name, region, err := extractNameAndRegion(resource)
	if err != nil {
		return nil, err
	}

	// Fetch the product dashboard to get its (org-specific) base URL.
	api := cockpit.NewGlobalAPI(ExtractClient(ctx))
	dash, err := api.GetGrafanaProductDashboard(
		&cockpit.GlobalAPIGetGrafanaProductDashboardRequest{
			DashboardName: dashboardName,
			ProjectID:     GetProjectIDFromContext(ctx),
		},
		scw.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}

	// var-Datasource is intentionally omitted: the dashboard falls back to its
	// default datasource (ponytail: add it if a dashboard needs it).
	q := url.Values{}
	q.Set("orgId", "1")
	q.Set("from", "now-1h")
	q.Set("to", "now")
	q.Set("timezone", "browser")
	if region != "" {
		q.Set("var-region", region)
	}
	q.Set("var-"+resourceVar, name)
	fullURL := dash.URL + "?" + q.Encode()

	if err := open.Start(fullURL); err != nil {
		return nil, &CliError{
			Err:     err,
			Message: "Failed to open cockpit url",
			Details: "You can open it: " + fullURL,
			Hint:    "You may not have a default browser configured",
			Code:    1,
		}
	}

	return fmt.Sprintf("Opening %s\n", fullURL), nil
}

// extractNameAndRegion pulls the Name and Region (or Zone) fields out of a resource
// struct via reflection. Region is returned as a string, converting Zone when needed.
func extractNameAndRegion(respI any) (name string, region string, err error) {
	v := reflect.ValueOf(respI)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return "", "", errors.New("resource is nil")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return "", "", fmt.Errorf("expected a resource struct, got %s", v.Kind())
	}

	if f := v.FieldByName("Name"); f.IsValid() && f.Kind() == reflect.String {
		name = f.String()
	}
	if name == "" {
		return "", "", errors.New("resource has no usable Name field")
	}

	if f := v.FieldByName("Region"); f.IsValid() {
		region = fmt.Sprint(f.Interface())
	}
	if region == "" {
		if f := v.FieldByName("Zone"); f.IsValid() {
			if z, ok := reflect.TypeAssert[scw.Zone](f); ok {
				r, err := z.Region()
				if err != nil {
					return "", "", err
				}
				region = string(r)
			}
		}
	}

	return name, region, nil
}
