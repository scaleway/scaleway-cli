package feedback

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"log"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
)

const githubURL = "https://github.com/scaleway/scaleway-cli/issues/new"

type issueTemplate string

const (
	bug     = issueTemplate("bug")
	feature = issueTemplate("feature")
	linux   = "linux"
	darwin  = "darwin"
	Windows = "windows"
)

type issue struct {
	BuildInfo     *core.BuildInfo
	IssueTemplate issueTemplate
}

const bugBodyTemplate = `
## Description:

## How to reproduce:

### Command attempted

### Expected Behavior

### Actual Behavior

## More info

## Version

{{ .BuildInfoStr }}
`

const featureBodyTemplate = `
## Description

## How this functionality would be exposed

## References

## Version

{{ .BuildInfoStr }}
`

func (i issue) getURL() string {
	u, err := url.Parse(githubURL)
	if err != nil {
		log.Fatal(err)
	}

	params := url.Values{}
	switch i.IssueTemplate {
	case feature:
		params.Add("labels", "enhancement")
		params.Add("issueTemplate", "feature_request.md")
		renderedBody, _ := i.renderTemplate(featureBodyTemplate)
		params.Add("body", renderedBody)
	case bug:
		params.Add("labels", "bug")
		params.Add("issueTemplate", "bug_report.md")
		renderedBody, _ := i.renderTemplate(bugBodyTemplate)
		params.Add("body", renderedBody)
	}

	u.RawQuery = params.Encode()

	return u.String()
}

func (i issue) openInBrowser(ctx context.Context) error {
	var err error
	var openCmd *exec.Cmd

	switch runtime.GOOS {
	case linux:
		openCmd = exec.Command("xdg-open", i.getURL()) //nolint:gosec
	case Windows:
		openCmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", i.getURL()) //nolint:gosec
	case darwin:
		openCmd = exec.Command("open", i.getURL()) //nolint:gosec
	default:
		return errors.New("unsupported platform")
	}

	exitCode, err := core.ExecCmd(ctx, openCmd)
	if exitCode != 0 {
		return &core.CliError{Empty: true, Code: exitCode}
	}

	return err
}

func (i issue) renderTemplate(bodyTemplate string) (string, error) {
	tmpl, err := template.New("configuration").Parse(bodyTemplate)
	if err != nil {
		return "", err
	}
	buildInfoStr, err := human.Marshal(i.BuildInfo, nil)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, struct {
		BuildInfoStr string
	}{
		BuildInfoStr: buildInfoStr,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
