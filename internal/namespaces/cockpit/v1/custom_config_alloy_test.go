package cockpit_test

import (
	"testing"

	cockpitSDK "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	cockpit "github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BuildLokiPushURL(t *testing.T) {
	assert.Equal(
		t,
		"https://example.logs.cockpit.fr-par.scw.cloud/loki/api/v1/push",
		cockpit.BuildLokiPushURL("https://example.logs.cockpit.fr-par.scw.cloud"),
	)
	assert.Equal(
		t,
		"https://example.logs.cockpit.fr-par.scw.cloud/loki/api/v1/push",
		cockpit.BuildLokiPushURL("https://example.logs.cockpit.fr-par.scw.cloud/loki/api/v1/push"),
	)
}

func Test_BuildTracesOTLPPushURL(t *testing.T) {
	assert.Equal(
		t,
		"https://example.traces.cockpit.fr-par.scw.cloud/otlp/v1/traces",
		cockpit.BuildTracesOTLPPushURL("https://example.traces.cockpit.fr-par.scw.cloud"),
	)
	assert.Equal(
		t,
		"https://example.traces.cockpit.fr-par.scw.cloud",
		cockpit.BuildTracesOTLPBaseURL("https://example.traces.cockpit.fr-par.scw.cloud/otlp/v1/traces"),
	)
}

func Test_RenderAlloyMetricsConfig(t *testing.T) {
	const baseURL = "https://example.metrics.fr-par.scw.cloud"

	t.Run("without token", func(t *testing.T) {
		got := cockpit.RenderAlloyMetricsConfig(baseURL, nil)
		assert.Contains(t, string(got), `prometheus.remote_write "cockpit"`)
		assert.Contains(t, string(got), baseURL+"/api/v1/push")
		assert.Contains(t, string(got), "COCKPIT_TOKEN_SECRET_KEY")
	})

	t.Run("with token", func(t *testing.T) {
		token := "my-secret-token"
		got := cockpit.RenderAlloyMetricsConfig(baseURL, &token)
		assert.Contains(t, string(got), `X-TOKEN" = "my-secret-token"`)
	})
}

func Test_RenderAlloyLogsConfig(t *testing.T) {
	const baseURL = "https://example.logs.cockpit.fr-par.scw.cloud"

	got := cockpit.RenderAlloyLogsConfig(baseURL, nil)
	assert.Contains(t, string(got), `loki.write "cockpit"`)
	assert.Contains(t, string(got), baseURL+"/loki/api/v1/push")
}

func Test_RenderAlloyTracesConfig(t *testing.T) {
	const baseURL = "https://example.traces.cockpit.fr-par.scw.cloud"

	got := cockpit.RenderAlloyTracesConfig(baseURL, nil)
	assert.Contains(t, string(got), `otelcol.exporter.otlphttp "cockpit"`)
	assert.Contains(t, string(got), `traces_endpoint = "`+baseURL+`/otlp/v1/traces"`)
}

func Test_RenderAlloyConfig(t *testing.T) {
	const metricsURL = "https://example.metrics.fr-par.scw.cloud"

	got, err := cockpit.RenderAlloyConfig(
		cockpitSDK.DataSourceTypeMetrics,
		metricsURL,
		nil,
	)
	require.NoError(t, err)
	assert.Equal(t, cockpit.RenderAlloyMetricsConfig(metricsURL, nil), got)
}
