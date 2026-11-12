package cockpit_test

import (
	"testing"

	cockpit "github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	cockpitSDK "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
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
	const (
		tracesHost    = "https://example.traces.cockpit.fr-par.scw.cloud"
		tracesPushURL = tracesHost + "/otlp/v1/traces"
	)

	assert.Equal(t, tracesPushURL, cockpit.BuildTracesOTLPPushURL(tracesHost))
	assert.Equal(t, tracesHost, cockpit.BuildTracesOTLPBaseURL(tracesPushURL))
}

func Test_RenderAlloyConfig(t *testing.T) {
	t.Run("metrics without token", func(t *testing.T) {
		const baseURL = "https://example.metrics.fr-par.scw.cloud"

		got, err := cockpit.RenderAlloyConfig(cockpitSDK.DataSourceTypeMetrics, baseURL, nil)
		require.NoError(t, err)
		assert.Contains(t, string(got), `prometheus.remote_write "cockpit"`)
		assert.Contains(t, string(got), baseURL+"/api/v1/push")
		assert.Contains(t, string(got), "COCKPIT_TOKEN_SECRET_KEY")
	})

	t.Run("metrics with token", func(t *testing.T) {
		const baseURL = "https://example.metrics.fr-par.scw.cloud"
		token := "my-secret-token"

		got, err := cockpit.RenderAlloyConfig(cockpitSDK.DataSourceTypeMetrics, baseURL, &token)
		require.NoError(t, err)
		assert.Contains(t, string(got), `X-TOKEN" = "my-secret-token"`)
	})

	t.Run("logs", func(t *testing.T) {
		const baseURL = "https://example.logs.cockpit.fr-par.scw.cloud"

		got, err := cockpit.RenderAlloyConfig(cockpitSDK.DataSourceTypeLogs, baseURL, nil)
		require.NoError(t, err)
		assert.Contains(t, string(got), `loki.write "cockpit"`)
		assert.Contains(t, string(got), baseURL+"/loki/api/v1/push")
	})

	t.Run("traces", func(t *testing.T) {
		const baseURL = "https://example.traces.cockpit.fr-par.scw.cloud"

		got, err := cockpit.RenderAlloyConfig(cockpitSDK.DataSourceTypeTraces, baseURL, nil)
		require.NoError(t, err)
		assert.Contains(t, string(got), `otelcol.exporter.otlphttp "cockpit"`)
		assert.Contains(t, string(got), `traces_endpoint = "`+baseURL+`/otlp/v1/traces"`)
	})
}
