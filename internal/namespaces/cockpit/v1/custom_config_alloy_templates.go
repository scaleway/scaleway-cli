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

type alloyTemplateData struct {
	RemoteWriteURL string
	PushURL        string
	OTLPBaseURL    string
	TracesEndpoint string
	Token          string
}

func alloyToken(tokenSecretKey *string) string {
	if tokenSecretKey != nil {
		return *tokenSecretKey
	}

	return alloyTokenPlaceholder
}

func renderAlloyTemplate(templateName string, data alloyTemplateData) (core.RawResult, error) {
	var buf bytes.Buffer
	if err := alloyTemplates.ExecuteTemplate(&buf, templateName, data); err != nil {
		return core.RawResult(""), fmt.Errorf("failed to render template %q: %w", templateName, err)
	}

	return core.RawResult(strings.TrimRight(buf.String(), "\n") + "\n"), nil
}
