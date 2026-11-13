package connect

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/pkg/shlex"
	"github.com/scaleway/scaleway-sdk-go/api/account/v3"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

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
// Organization and Project are resolved lazily: the API is only called when
// the template references {{.Organization}} or {{.Project}} (Go templates
// invoke no-arg methods at the point of use).
type TemplateData struct {
	client *scw.Client
}

func NewTemplateData(ctx context.Context) *TemplateData {
	client := core.ExtractClient(ctx)
	data := &TemplateData{client: client}

	return data
}

// Organization returns the default organization.
func (d *TemplateData) Organization() (*iam.Organization, error) {
	org, exists := d.client.GetDefaultOrganizationID()
	if !exists {
		return nil, errors.New("organization not found")
	}

	return iam.NewAPI(d.client).GetOrganization(
		&iam.GetOrganizationRequest{OrganizationID: org},
	)
}

// Project returns the default project.
func (d *TemplateData) Project() (*account.Project, error) {
	project, exists := d.client.GetDefaultProjectID()
	if !exists {
		return nil, errors.New("project not found")
	}

	return account.NewProjectAPI(d.client).GetProject(
		&account.ProjectAPIGetProjectRequest{ProjectID: project},
	)
}
