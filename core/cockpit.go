package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"text/template"

	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/skratchdot/open-golang/open"
)

func commandHasCockpit(cmd *Command) bool {
	return cmd.CockpitDashboard != ""
}

// ensureClient makes sure an authenticated client is on the meta, mirroring the
// client creation done in cobraRun. It is needed for the non-runnable (namespace-level)
// --cockpit path, where cobraRun never runs. The runnable path already has a client,
// so it should not call this.
func ensureClient(ctx context.Context, allowAnonymous bool) error {
	meta := ExtractMeta(ctx)
	if !allowAnonymous && !meta.isClientFromBootstrapConfig {
		client, err := meta.Platform.CreateClient(
			meta.httpClient,
			ExtractConfigPath(ctx),
			ExtractProfileName(ctx),
		)
		if err != nil {
			return createClientError(err)
		}
		meta.Client = client
	}

	return nil
}

// runCockpit opens the Grafana dashboard for the current resource.
// respI is the command response (used to resolve the resource name/region and any
// dashboard-name template), or nil for a product-level dashboard.
func runCockpit(ctx context.Context, cmd *Command, respI any) (any, error) {
	// Resolve the dashboard name, evaluating any template against the response.
	dashboardName := cmd.CockpitDashboard
	if strings.Contains(dashboardName, "{{") {
		if respI == nil {
			return nil, &CliError{
				Err:     errors.New("dashboard name requires a resource"),
				Message: "This command cannot resolve a cockpit dashboard",
			}
		}
		tmpl, err := template.New("dashboard").Parse(dashboardName)
		if err != nil {
			return nil, err
		}
		buf := bytes.NewBuffer(nil)
		if err := tmpl.Execute(buf, respI); err != nil {
			return nil, err
		}
		dashboardName = buf.String()
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

	// Resolve the identifying resource variable, if any.
	resourceName, region := "", ""
	if cmd.CockpitResourceVar != "" {
		if respI == nil {
			return nil, &CliError{
				Err:     errors.New("resource is required for this dashboard"),
				Message: "Pass a resource to open its cockpit dashboard",
			}
		}
		resourceName, region, err = extractNameAndRegion(respI)
		if err != nil {
			return nil, err
		}
	}

	// Build the Grafana URL. var-Datasource is intentionally omitted: the dashboard
	// falls back to its default datasource (ponytail: add it if a dashboard needs it).
	q := url.Values{}
	q.Set("orgId", "1")
	q.Set("from", "now-1h")
	q.Set("to", "now")
	q.Set("timezone", "browser")
	if region != "" {
		q.Set("var-region", region)
	}
	if cmd.CockpitResourceVar != "" {
		q.Set("var-"+cmd.CockpitResourceVar, resourceName)
	}
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
