package secret

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"text/template"
	"text/template/parse"

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

	resolver := &secretResolver{
		ctx:    ctx,
		api:    secret.NewAPI(core.ExtractClient(ctx)),
		region: args.Region,
		cache:  map[string]secretResult{},
	}

	tpl, err := template.New("inject").Funcs(template.FuncMap{
		"secret":    resolver.resolve,
		"jsonField": jsonField,
	}).Parse(text)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	resolver.prefetch(tpl)

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

type secretResolver struct {
	ctx    context.Context //nolint:containedctx // template funcs cannot take a context
	api    *secret.API
	region scw.Region
	mu     sync.Mutex
	cache  map[string]secretResult
}

// prefetch concurrently resolves the literal references of tpl; errors surface only if rendering uses them.
func (r *secretResolver) prefetch(tpl *template.Template) {
	refs := map[string]struct{}{}
	for _, t := range tpl.Templates() {
		if t.Tree != nil {
			collectSecretRefs(t.Root, refs)
		}
	}

	var g errgroup.Group

	g.SetLimit(maxConcurrentSecretFetches)

	for ref := range refs {
		g.Go(func() error {
			_, _ = r.resolve(ref)

			return nil
		})
	}

	_ = g.Wait()
}

// resolve returns the data of the secret version referenced by ref, fetching each reference once.
func (r *secretResolver) resolve(ref string) (string, error) {
	r.mu.Lock()
	res, ok := r.cache[ref]
	r.mu.Unlock()

	if !ok {
		res.value, res.err = r.fetch(ref)

		r.mu.Lock()
		r.cache[ref] = res
		r.mu.Unlock()
	}

	return res.value, res.err
}

func (r *secretResolver) fetch(ref string) (string, error) {
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
		resp, err = r.api.AccessSecretVersion(&secret.AccessSecretVersionRequest{
			Region:   r.region,
			SecretID: id,
			Revision: revision,
		}, scw.WithContext(r.ctx))
	} else {
		path, name := "/", id
		if dir, base, found := strings.CutLast(id, "/"); found {
			path, name = "/"+strings.Trim(dir, "/"), base
		}

		if name == "" {
			return "", fmt.Errorf("invalid secret reference %q", ref)
		}

		resp, err = r.api.AccessSecretVersionByPath(&secret.AccessSecretVersionByPathRequest{
			Region:     r.region,
			SecretPath: path,
			SecretName: name,
			Revision:   revision,
		}, scw.WithContext(r.ctx))
	}

	if err != nil {
		return "", fmt.Errorf("accessing secret %q: %w", ref, err)
	}

	return string(resp.Data), nil
}

// collectSecretRefs adds to refs the string literal argument of every `secret` call under node.
func collectSecretRefs(node parse.Node, refs map[string]struct{}) {
	switch n := node.(type) {
	case *parse.ListNode:
		if n == nil {
			return
		}

		for _, child := range n.Nodes {
			collectSecretRefs(child, refs)
		}
	case *parse.ActionNode:
		collectSecretRefs(n.Pipe, refs)
	case *parse.TemplateNode:
		collectSecretRefs(n.Pipe, refs)
	case *parse.IfNode:
		collectBranchSecretRefs(&n.BranchNode, refs)
	case *parse.RangeNode:
		collectBranchSecretRefs(&n.BranchNode, refs)
	case *parse.WithNode:
		collectBranchSecretRefs(&n.BranchNode, refs)
	case *parse.PipeNode:
		if n == nil {
			return
		}

		for _, cmd := range n.Cmds {
			collectSecretRefs(cmd, refs)
		}
	case *parse.CommandNode:
		if len(n.Args) == 2 {
			fn, isIdent := n.Args[0].(*parse.IdentifierNode)
			ref, isString := n.Args[1].(*parse.StringNode)

			if isIdent && isString && fn.Ident == "secret" {
				refs[ref.Text] = struct{}{}
			}
		}

		for _, arg := range n.Args {
			collectSecretRefs(arg, refs)
		}
	}
}

func collectBranchSecretRefs(n *parse.BranchNode, refs map[string]struct{}) {
	collectSecretRefs(n.Pipe, refs)
	collectSecretRefs(n.List, refs)
	collectSecretRefs(n.ElseList, refs)
}

func jsonField(field, data string) (string, error) {
	value, err := getSecretVersionField([]byte(data), field)
	if err != nil {
		return "", fmt.Errorf("extracting JSON field %q: %w", field, err)
	}

	return string(value), nil
}
