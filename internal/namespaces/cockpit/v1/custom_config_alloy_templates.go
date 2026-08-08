package cockpit

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/scaleway/scaleway-cli/v2/core"
)

//go:embed templates/alloy-*.tmpl
var alloyTemplatesFS embed.FS

var alloyTemplates = template.Must(template.ParseFS(alloyTemplatesFS, "templates/alloy-*.tmpl"))

func init() {
	for _, name := range []string{
		"alloy-metrics.tmpl",
		"alloy-logs.tmpl",
		"alloy-traces.tmpl",
	} {
		if alloyTemplates.Lookup(name) == nil {
			panic("missing cockpit alloy template " + name)
		}
	}
}

// AlloyTemplateData holds values injected into Grafana Alloy configuration templates.
type AlloyTemplateData struct {
	// RemoteWriteURL is the Prometheus remote write endpoint for metrics data sources.
	RemoteWriteURL string
	// PushURL is the Loki HTTP push endpoint for logs data sources.
	PushURL string
	// OTLPBaseURL is the OTLP HTTP client endpoint base URL for traces data sources.
	OTLPBaseURL string
	// TracesEndpoint is the OTLP HTTP traces push endpoint for traces data sources.
	TracesEndpoint string
	// Token is the Cockpit secret key value for the X-TOKEN request header.
	Token string
}

func alloyToken(tokenSecretKey *string) string {
	if tokenSecretKey != nil {
		return *tokenSecretKey
	}

	return alloyTokenPlaceholder
}

func renderAlloyTemplate(templateName string, data AlloyTemplateData) (core.RawResult, error) {
	var buf bytes.Buffer
	if err := alloyTemplates.ExecuteTemplate(&buf, templateName, data); err != nil {
		return core.RawResult(""), fmt.Errorf("failed to render template %q: %w", templateName, err)
	}

	return core.RawResult(strings.TrimRight(buf.String(), "\n") + "\n"), nil
}
