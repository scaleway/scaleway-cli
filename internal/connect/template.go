package connect

import (
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/pkg/shlex"
)

// TemplateArgSpec is the shared --template flag spec for connect commands.
func TemplateArgSpec() *core.ArgSpec {
	return &core.ArgSpec{
		Name:  "template",
		Short: "Go text/template rendered into a command line that replaces the default connect command",
	}
}

// RenderCommand renders a Go text/template into a command line using the
// given data. It returns nil when the template is empty so that the caller
// falls back to its default command. The rendered output is split with shlex
// (quotes honored) and executed directly, without any shell.
func RenderCommand(tmpl string, data any) ([]string, error) {
	if strings.TrimSpace(tmpl) == "" {
		return nil, nil
	}

	t, err := template.New("connect").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connect template: %w", err)
	}

	var rendered strings.Builder
	if err := t.Execute(&rendered, data); err != nil {
		return nil, fmt.Errorf("failed to render connect template: %w", err)
	}

	fields, err := shlex.Split(rendered.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse connect command: %w", err)
	}
	if len(fields) == 0 {
		return nil, errors.New("connect template rendered an empty command")
	}

	return fields, nil
}
