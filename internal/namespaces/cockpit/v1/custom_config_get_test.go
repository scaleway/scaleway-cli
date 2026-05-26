package cockpit_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	cockpit "github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BuildPrometheusRemoteWriteURL(t *testing.T) {
	t.Run("append push path", func(t *testing.T) {
		url := cockpit.BuildPrometheusRemoteWriteURL("https://example.metrics.fr-par.scw.cloud")
		assert.Equal(t, "https://example.metrics.fr-par.scw.cloud/api/v1/push", url)
	})

	t.Run("keep existing push path", func(t *testing.T) {
		url := cockpit.BuildPrometheusRemoteWriteURL(
			"https://example.metrics.fr-par.scw.cloud/api/v1/push",
		)
		assert.Equal(t, "https://example.metrics.fr-par.scw.cloud/api/v1/push", url)
	})

	t.Run("trim trailing slash", func(t *testing.T) {
		url := cockpit.BuildPrometheusRemoteWriteURL("https://example.metrics.fr-par.scw.cloud/")
		assert.Equal(t, "https://example.metrics.fr-par.scw.cloud/api/v1/push", url)
	})
}

func Test_ParseCockpitDataSourceEndpoint(t *testing.T) {
	t.Run("https default port", func(t *testing.T) {
		const host = "000avb0d-34ae-66hh-643b-f9e0n3k17773.logs.cockpit.fr-par.scw.cloud"

		endpoint, err := cockpit.ParseCockpitDataSourceEndpoint("https://" + host)
		require.NoError(t, err)
		assert.Equal(t, host, endpoint.Host)
		assert.Equal(t, 443, endpoint.Port)
	})

	t.Run("explicit port", func(t *testing.T) {
		endpoint, err := cockpit.ParseCockpitDataSourceEndpoint(
			"https://example.logs.cockpit.fr-par.scw.cloud:8443",
		)
		require.NoError(t, err)
		assert.Equal(t, "example.logs.cockpit.fr-par.scw.cloud", endpoint.Host)
		assert.Equal(t, 8443, endpoint.Port)
	})
}

func Test_BuildFluentBitLogsURI(t *testing.T) {
	assert.Equal(t, "/otlp/v1/logs", cockpit.BuildFluentBitLogsURI(
		"https://example.logs.cockpit.fr-par.scw.cloud",
	))
	assert.Equal(t, "/otlp/v1/logs", cockpit.BuildFluentBitLogsURI(
		"https://example.logs.cockpit.fr-par.scw.cloud/otlp/v1/logs",
	))
}

func Test_RenderFluentBitConfig(t *testing.T) {
	endpoint, err := cockpit.ParseCockpitDataSourceEndpoint(
		"https://example.logs.cockpit.fr-par.scw.cloud",
	)
	require.NoError(t, err)

	t.Run("without token", func(t *testing.T) {
		got := cockpit.RenderFluentBitConfig(endpoint, nil)
		want := core.RawResult(
			"# Snippet of Fluent Bit configuration to add to fluent-bit.conf\n" +
				"# Uses a dummy input for testing; replace it with your real log inputs.\n" +
				"[SERVICE]\n" +
				"    Flush        1\n" +
				"    Log_Level    info\n" +
				"\n" +
				"[INPUT]\n" +
				"    Name    dummy\n" +
				"    Tag     dummy.log\n" +
				"    Rate    1\n" +
				"\n" +
				"[OUTPUT]\n" +
				"    Name                 opentelemetry\n" +
				"    Match                dummy.log\n" +
				"    Host                 example.logs.cockpit.fr-par.scw.cloud\n" +
				"    Port                 443\n" +
				"    Logs_uri             /otlp/v1/logs\n" +
				"    Log_response_payload True\n" +
				"    Tls                  On\n" +
				"    Tls.verify           On\n",
		)
		assert.Equal(t, want, got)
	})

	t.Run("with token", func(t *testing.T) {
		token := "my-secret-token"
		got := cockpit.RenderFluentBitConfig(endpoint, &token)
		want := core.RawResult(
			"# Snippet of Fluent Bit configuration to add to fluent-bit.conf\n" +
				"# Uses a dummy input for testing; replace it with your real log inputs.\n" +
				"[SERVICE]\n" +
				"    Flush        1\n" +
				"    Log_Level    info\n" +
				"\n" +
				"[INPUT]\n" +
				"    Name    dummy\n" +
				"    Tag     dummy.log\n" +
				"    Rate    1\n" +
				"\n" +
				"[OUTPUT]\n" +
				"    Name                 opentelemetry\n" +
				"    Match                dummy.log\n" +
				"    Host                 example.logs.cockpit.fr-par.scw.cloud\n" +
				"    Port                 443\n" +
				"    Logs_uri             /otlp/v1/logs\n" +
				"    Log_response_payload True\n" +
				"    Tls                  On\n" +
				"    Tls.verify           On\n" +
				"    header               X-TOKEN my-secret-token\n",
		)
		assert.Equal(t, want, got)
	})
}

func Test_RenderPrometheusRemoteWriteConfig(t *testing.T) {
	const baseURL = "https://example.metrics.fr-par.scw.cloud"

	t.Run("without token", func(t *testing.T) {
		got := cockpit.RenderPrometheusRemoteWriteConfig(baseURL, nil)
		want := core.RawResult(
			"# Snippet of Prometheus configuration to add to prometheus.yml\n" +
				"remote_write:\n" +
				`  - url: "` + baseURL + `/api/v1/push"` + "\n",
		)
		assert.Equal(t, want, got)
	})

	t.Run("with token", func(t *testing.T) {
		token := "my-secret-token"
		got := cockpit.RenderPrometheusRemoteWriteConfig(baseURL, &token)
		want := core.RawResult(
			"# Snippet of Prometheus configuration to add to prometheus.yml\n" +
				"remote_write:\n" +
				`  - url: "` + baseURL + `/api/v1/push"` + "\n" +
				"    headers:\n" +
				"      X-TOKEN: my-secret-token\n",
		)
		assert.Equal(t, want, got)
	})
}
