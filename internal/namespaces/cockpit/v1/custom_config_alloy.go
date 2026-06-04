package cockpit

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
)

const (
	cockpitConfigTypeAlloy = cockpitConfigType("alloy")

	cockpitLokiPushPath    = "/loki/api/v1/push"
	cockpitTracesOTLPPath  = "/otlp/v1/traces"
	alloyTokenPlaceholder  = "COCKPIT_TOKEN_SECRET_KEY"
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
	if strings.HasSuffix(baseURL, cockpitTracesOTLPPath) {
		return strings.TrimSuffix(baseURL, cockpitTracesOTLPPath)
	}

	return baseURL
}

func renderAlloyTokenHeaders(tokenSecretKey *string) []string {
	token := alloyTokenPlaceholder
	if tokenSecretKey != nil {
		token = *tokenSecretKey
	}

	return []string{
		"        headers = {",
		`            "X-TOKEN" = "` + token + `",`,
		"        }",
	}
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
		return core.RawResult(""), fmt.Errorf("unsupported data source type %q for alloy config", dataSourceType)
	}
}

// RenderAlloyMetricsConfig renders an Alloy snippet that scrapes host metrics and remote_writes to Cockpit.
func RenderAlloyMetricsConfig(dataSourceURL string, tokenSecretKey *string) core.RawResult {
	remoteWriteURL := BuildPrometheusRemoteWriteURL(dataSourceURL)

	lines := []string{
		"// Snippet of Grafana Alloy configuration to add to cockpit.alloy",
		"// Collects host metrics and pushes them to your Cockpit metrics data source.",
		"// Run with: alloy run cockpit.alloy",
		"",
		`prometheus.exporter.unix "node" {`,
		"    set_collectors = [",
		`        "uname",`,
		`        "cpu",`,
		`        "cpufreq",`,
		`        "loadavg",`,
		`        "meminfo",`,
		`        "filesystem",`,
		`        "netdev",`,
		"    ]",
		"}",
		"",
		`prometheus.scrape "node" {`,
		`    scrape_interval = "60s"`,
		`    scrape_timeout  = "4s"`,
		"",
		"    targets    = prometheus.exporter.unix.node.targets",
		"    forward_to = [prometheus.remote_write.cockpit.receiver]",
		"}",
		"",
		`prometheus.remote_write "cockpit" {`,
		"    endpoint {",
		`        url = "` + remoteWriteURL + `"`,
	}
	lines = append(lines, renderAlloyTokenHeaders(tokenSecretKey)...)
	lines = append(lines,
		"    }",
		"}",
		"",
	)

	return core.RawResult(strings.Join(lines, "\n"))
}

// RenderAlloyLogsConfig renders an Alloy loki.write snippet for Cockpit logs.
func RenderAlloyLogsConfig(dataSourceURL string, tokenSecretKey *string) core.RawResult {
	pushURL := BuildLokiPushURL(dataSourceURL)

	lines := []string{
		"// Snippet of Grafana Alloy configuration to add to cockpit.alloy",
		"// Forwards logs to your Cockpit logs data source. Add loki.source.* components",
		"// and forward them to loki.write.cockpit.receiver.",
		"// Run with: alloy run cockpit.alloy",
		"",
		`loki.write "cockpit" {`,
		"    endpoint {",
		`        url = "` + pushURL + `"`,
	}
	lines = append(lines, renderAlloyTokenHeaders(tokenSecretKey)...)
	lines = append(lines,
		"    }",
		"}",
		"",
	)

	return core.RawResult(strings.Join(lines, "\n"))
}

// RenderAlloyTracesConfig renders an Alloy otelcol.exporter.otlphttp snippet for Cockpit traces.
func RenderAlloyTracesConfig(dataSourceURL string, tokenSecretKey *string) core.RawResult {
	baseURL := BuildTracesOTLPBaseURL(dataSourceURL)
	tracesURL := BuildTracesOTLPPushURL(dataSourceURL)

	lines := []string{
		"// Snippet of Grafana Alloy configuration to add to cockpit.alloy",
		"// Exports traces to your Cockpit traces data source over OTLP HTTP.",
		"// Forward otelcol.* pipeline outputs to otelcol.exporter.otlphttp.cockpit.input.",
		"// Run with: alloy run cockpit.alloy",
		"",
		`otelcol.exporter.otlphttp "cockpit" {`,
		"    client {",
		`        endpoint        = "` + baseURL + `"`,
		`        traces_endpoint = "` + tracesURL + `"`,
	}
	lines = append(lines, renderAlloyTokenHeaders(tokenSecretKey)...)
	lines = append(lines,
		"    }",
		"}",
		"",
	)

	return core.RawResult(strings.Join(lines, "\n"))
}
