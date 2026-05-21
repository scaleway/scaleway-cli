package cockpit_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	cockpit "github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	"github.com/stretchr/testify/assert"
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
