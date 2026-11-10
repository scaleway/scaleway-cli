package secret

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"text/template"

	"github.com/scaleway/scaleway-cli/v2/core"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
	"golang.org/x/sync/errgroup"
)

const (
	revisionLatest = "latest"
	// maxConcurrentSecretFetches bounds parallel API calls during prefetch.
	maxConcurrentSecretFetches = 10
)

type secretInjectArgs struct {
	Template string
	Region   scw.Region
}

func secretInjectCommand() *core.Command {
	return &core.Command{
		Short: `Inject Scaleway secrets into a template`,
		Long: `Render a Go template (text/template) in which the ` + "`secret`" + ` function resolves secret references.

- ` + "`{{ secret \"REFERENCE\" }}`" + `: secret value
- ` + "`{{ secret \"REFERENCE\" | jsonField \"KEY\" }}`" + `: string field of a JSON-encoded secret

REFERENCE is ` + "`(UUID|[PATH/]NAME)[@REVISION]`" + `, REVISION being an integer, "latest" (default) or "latest_enabled".

The rendered template is written to stdout.`,
		Namespace: "secret", //nolint:goconst
		Resource:  "inject",
		ArgsType:  reflect.TypeFor[secretInjectArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:        "template",
				Short:       "Template to render (default: stdin)",
				CanLoadFile: true,
			},
			// scw.AllRegions includes regions without Secret Manager.
			secretVersionAccess().ArgSpecs.GetByName("region"),
		},
		Groups: []string{"security"},
		Run:    secretInjectRun,
		Examples: []*core.Example{
			{
				Short: "Render a template from stdin",
				Raw:   `echo 'password: {{ secret "11111111-1111-1111-1111-111111111111" }}' | scw secret inject`,
			},
			{
				Short: "Render a template file into another file",
				Raw:   `scw secret inject template=@config.tpl > config.yaml`,
			},
			{
				Short: "Extract a JSON field from a secret referenced by path, name and revision",
				Raw:   `scw secret inject template='{{ secret "my-app/db@2" | jsonField "password" }}'`,
			},
		},
	}
}

func secretInjectRun(ctx context.Context, argsI any) (any, error) {
	args := argsI.(*secretInjectArgs)

	text := args.Template
	if text == "" {
		data, err := io.ReadAll(core.ExtractStdin(ctx))
		if err != nil {
			return nil, fmt.Errorf("reading stdin: %w", err)
		}

		text = string(data)
	}

	api := secret.NewAPI(core.ExtractClient(ctx))
	fetch := func(ref string) (string, error) {
		return fetchSecret(ctx, api, args.Region, ref)
	}

	// First pass only collects references, with placeholder values.
	refs := map[string]struct{}{}

	tpl, err := template.New("inject").Funcs(template.FuncMap{
		"secret": func(ref string) string {
			refs[ref] = struct{}{}

			return ""
		},
		"jsonField": func(string, string) string { return "" },
	}).Parse(text)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	// Placeholders may break value-dependent logic; the second pass reports real errors.
	_ = tpl.Execute(io.Discard, nil)

	results := fetchSecrets(refs, fetch)

	tpl.Funcs(template.FuncMap{
		"secret": func(ref string) (string, error) {
			// Missed by the first pass when behind a branch depending on a secret value.
			res, ok := results[ref]
			if !ok {
				res.value, res.err = fetch(ref)
				results[ref] = res
			}

			return res.value, res.err
		},
		"jsonField": jsonField,
	})

	// Render fully before writing so a failure never outputs a partial template.
	var out strings.Builder
	if err := tpl.Execute(&out, nil); err != nil {
		return nil, fmt.Errorf("rendering template: %w", err)
	}

	return core.RawResult(out.String()), nil
}

type secretResult struct {
	value string
	err   error
}

// fetchSecrets concurrently fetches refs; errors surface only if rendering uses them.
func fetchSecrets(refs map[string]struct{}, fetch func(string) (string, error)) map[string]secretResult {
	var (
		mu sync.Mutex
		g  errgroup.Group
	)

	results := make(map[string]secretResult, len(refs))

	g.SetLimit(maxConcurrentSecretFetches)

	for ref := range refs {
		g.Go(func() error {
			value, err := fetch(ref)

			mu.Lock()
			results[ref] = secretResult{value: value, err: err}
			mu.Unlock()

			return nil
		})
	}

	_ = g.Wait()

	return results
}

func fetchSecret(ctx context.Context, api *secret.API, region scw.Region, ref string) (string, error) {
	id, revision, hasRevision := strings.Cut(ref, "@")
	if !hasRevision {
		revision = revisionLatest
	}

	if id == "" || revision == "" {
		return "", fmt.Errorf("invalid secret reference %q", ref)
	}

	var (
		resp *secret.AccessSecretVersionResponse
		err  error
	)

	if validation.IsUUID(id) {
		resp, err = api.AccessSecretVersion(&secret.AccessSecretVersionRequest{
			Region:   region,
			SecretID: id,
			Revision: revision,
		}, scw.WithContext(ctx))
	} else {
		path, name := "/", id
		if dir, base, found := strings.CutLast(id, "/"); found {
			path, name = "/"+strings.Trim(dir, "/"), base
		}

		if name == "" {
			return "", fmt.Errorf("invalid secret reference %q", ref)
		}

		resp, err = api.AccessSecretVersionByPath(&secret.AccessSecretVersionByPathRequest{
			Region:     region,
			SecretPath: path,
			SecretName: name,
			Revision:   revision,
		}, scw.WithContext(ctx))
	}

	if err != nil {
		return "", fmt.Errorf("accessing secret %q: %w", ref, err)
	}

	return string(resp.Data), nil
}

func jsonField(field, data string) (string, error) {
	value, err := getSecretVersionField([]byte(data), field)
	if err != nil {
		return "", fmt.Errorf("extracting JSON field %q: %w", field, err)
	}

	return string(value), nil
}
