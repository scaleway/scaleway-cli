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
	baseURL := strings.TrimRight(dataSourceURL, "/")
	if strings.HasSuffix(baseURL, cockpitLokiPushPath) {
		return baseURL
	}

	return baseURL + cockpitLokiPushPath
}

// BuildTracesOTLPPushURL returns the OTLP HTTP traces push URL for a Cockpit traces data source.
func BuildTracesOTLPPushURL(dataSourceURL string) string {
	baseURL := strings.TrimRight(dataSourceURL, "/")
	if strings.HasSuffix(baseURL, cockpitTracesOTLPPath) {
		return baseURL
	}

	return baseURL + cockpitTracesOTLPPath
}

// BuildTracesOTLPBaseURL returns the base URL for otelcol.exporter.otlphttp client.endpoint.
func BuildTracesOTLPBaseURL(dataSourceURL string) string {
	baseURL := strings.TrimRight(dataSourceURL, "/")
	if base, ok := strings.CutSuffix(baseURL, cockpitTracesOTLPPath); ok {
		return base
	}

	return baseURL
}

// RenderAlloyConfig renders a Grafana Alloy configuration snippet for stdout.
func RenderAlloyConfig(
	dataSourceType cockpit.DataSourceType,
	dataSourceURL string,
	tokenSecretKey *string,
) (core.RawResult, error) {
	switch dataSourceType {
	case cockpit.DataSourceTypeMetrics:
		return RenderAlloyMetricsConfig(dataSourceURL, tokenSecretKey), nil
	case cockpit.DataSourceTypeLogs:
		return RenderAlloyLogsConfig(dataSourceURL, tokenSecretKey), nil
	case cockpit.DataSourceTypeTraces:
		return RenderAlloyTracesConfig(dataSourceURL, tokenSecretKey), nil
	default:
		return core.RawResult(""), fmt.Errorf(
			"unsupported data source type %q for alloy config",
			dataSourceType,
		)
	}
}

// RenderAlloyMetricsConfig renders an Alloy snippet that scrapes host metrics and remote_writes to Cockpit.
func RenderAlloyMetricsConfig(dataSourceURL string, tokenSecretKey *string) core.RawResult {
	result, err := renderAlloyTemplate("alloy-metrics.tmpl", AlloyTemplateData{
		RemoteWriteURL: BuildPrometheusRemoteWriteURL(dataSourceURL),
		Token:          alloyToken(tokenSecretKey),
	})
	if err != nil {
		return core.RawResult("")
	}

	return result
}

// RenderAlloyLogsConfig renders an Alloy loki.write snippet for Cockpit logs.
func RenderAlloyLogsConfig(dataSourceURL string, tokenSecretKey *string) core.RawResult {
	result, err := renderAlloyTemplate("alloy-logs.tmpl", AlloyTemplateData{
		PushURL: BuildLokiPushURL(dataSourceURL),
		Token:   alloyToken(tokenSecretKey),
	})
	if err != nil {
		return core.RawResult("")
	}

	return result
}

// RenderAlloyTracesConfig renders an Alloy otelcol.exporter.otlphttp snippet for Cockpit traces.
func RenderAlloyTracesConfig(dataSourceURL string, tokenSecretKey *string) core.RawResult {
	result, err := renderAlloyTemplate("alloy-traces.tmpl", AlloyTemplateData{
		OTLPBaseURL:    BuildTracesOTLPBaseURL(dataSourceURL),
		TracesEndpoint: BuildTracesOTLPPushURL(dataSourceURL),
		Token:          alloyToken(tokenSecretKey),
	})
	if err != nil {
		return core.RawResult("")
	}

	return result
}
