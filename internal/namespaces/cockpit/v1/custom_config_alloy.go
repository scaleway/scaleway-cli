package cockpit

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
)

const (
	cockpitConfigTypeAlloy = cockpitConfigType("alloy")

	cockpitLokiPushPath   = "/loki/api/v1/push"
	cockpitTracesOTLPPath = "/otlp/v1/traces"
	alloyTokenPlaceholder = "COCKPIT_TOKEN_SECRET_KEY"
)

// BuildLokiPushURL returns the Loki push URL for a Cockpit logs data source base URL.
func BuildLokiPushURL(dataSourceURL string) string {
	return ensureURLPath(dataSourceURL, cockpitLokiPushPath)
}

// BuildTracesOTLPBaseURL returns the base URL for otelcol.exporter.otlphttp client.endpoint.
func BuildTracesOTLPBaseURL(dataSourceURL string) string {
	baseURL := strings.TrimRight(dataSourceURL, "/")
	if base, ok := strings.CutSuffix(baseURL, cockpitTracesOTLPPath); ok {
		return base
	}

	return baseURL
}

// BuildTracesOTLPPushURL returns the OTLP HTTP traces push URL for a Cockpit traces data source.
func BuildTracesOTLPPushURL(dataSourceURL string) string {
	return BuildTracesOTLPBaseURL(dataSourceURL) + cockpitTracesOTLPPath
}

// RenderAlloyConfig renders a Grafana Alloy configuration snippet for stdout.
func RenderAlloyConfig(
	dataSourceType cockpit.DataSourceType,
	dataSourceURL string,
	tokenSecretKey *string,
) (core.RawResult, error) {
	token := alloyToken(tokenSecretKey)

	switch dataSourceType {
	case cockpit.DataSourceTypeMetrics:
		return renderAlloyTemplate("alloy-metrics.tmpl", alloyTemplateData{
			RemoteWriteURL: BuildPrometheusRemoteWriteURL(dataSourceURL),
			Token:          token,
		})
	case cockpit.DataSourceTypeLogs:
		return renderAlloyTemplate("alloy-logs.tmpl", alloyTemplateData{
			PushURL: BuildLokiPushURL(dataSourceURL),
			Token:   token,
		})
	case cockpit.DataSourceTypeTraces:
		return renderAlloyTemplate("alloy-traces.tmpl", alloyTemplateData{
			OTLPBaseURL:    BuildTracesOTLPBaseURL(dataSourceURL),
			TracesEndpoint: BuildTracesOTLPPushURL(dataSourceURL),
			Token:          token,
		})
	default:
		return core.RawResult(""), fmt.Errorf(
			"unsupported data source type %q for alloy config",
			dataSourceType,
		)
	}
}
