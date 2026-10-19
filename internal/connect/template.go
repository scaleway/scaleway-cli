package connect

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/pkg/shlex"
	account "github.com/scaleway/scaleway-sdk-go/api/account/v3"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// TemplateArgSpec is the shared --template flag spec for connect commands.
func TemplateArgSpec() *core.ArgSpec {
	return &core.ArgSpec{
		Name:    "template",
		Short:   "Go text/template rendered into a command line that replaces the default connect command",
		Default: core.DefaultValueSetter("k9s"),
	}
}

// RenderCommand renders a Go text/template into a command line using the
// given data. The rendered output is split with shlex (quotes honored) and
// executed directly, without any shell.
func RenderCommand(tmpl string, data any) ([]string, error) {
	if strings.TrimSpace(tmpl) == "" {
		return nil, errors.New("empty connect template")
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

// TemplateData is the base template data for connect commands, embedded by
// each namespace's connect command which adds its own template values.
// OrganizationName and ProjectName are resolved lazily: the API is only called when
// the template references {{.OrganizationName}} or {{.ProjectName}} (Go templates
// invoke no-arg methods at the point of use).
// Name lookups are best-effort: on error the name is left empty.
type TemplateData struct {
	OrganizationID string
	ProjectID      string

	client *scw.Client
}

func NewTemplateData(ctx context.Context) *TemplateData {
	client := core.ExtractClient(ctx)
	data := &TemplateData{client: client}
	data.OrganizationID, _ = client.GetDefaultOrganizationID()
	data.ProjectID, _ = client.GetDefaultProjectID()

	return data
}

// Client returns the client used by the base template data, for embedding
// commands that need their own API lookups.
func (d *TemplateData) Client() *scw.Client {
	return d.client
}

// OrganizationName returns the default organization's name.
func (d *TemplateData) OrganizationName() string {
	org, err := iam.NewAPI(d.client).GetOrganization(
		&iam.GetOrganizationRequest{OrganizationID: d.OrganizationID},
	)
	if err != nil {
		return ""
	}

	return org.Name
}

// ProjectName returns the default project's name.
func (d *TemplateData) ProjectName() string {
	project, err := account.NewProjectAPI(d.client).GetProject(
		&account.ProjectAPIGetProjectRequest{ProjectID: d.ProjectID},
	)
	if err != nil {
		return ""
	}

	return project.Name
}
