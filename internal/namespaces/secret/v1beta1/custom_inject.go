package secret

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/scaleway/scaleway-cli/v2/core"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"golang.org/x/sync/errgroup"
)

// matches {{ scw://REFERENCE }} with optional surrounding whitespace inside braces
var secretRefRegex = regexp.MustCompile(`\{\{\s*scw://([^\s}]+)\s*\}\}`)

const revisionLatest = "latest"

type secretInjectArgs struct {
	InFile   string
	OutFile  string
	FileMode string
	Region   scw.Region
}

type parsedSecretRef struct {
	secretID   string // UUID-based lookup
	secretName string // name-based lookup
	secretPath string // path prefix for name-based lookup
	revision   string
	field      string
}

func secretInjectCommand() *core.Command {
	return &core.Command{
		Short: `Inject Scaleway secrets into a template file`,
		Long: `Substitute secret references in a template file with their actual values.

References use the syntax ` + "`{{ scw://REFERENCE }}`" + ` where REFERENCE can be:

  UUID                       secret ID, latest revision
  UUID@REVISION              secret ID, specific revision
  UUID:FIELD                 secret ID, latest revision, JSON field
  UUID@REVISION:FIELD        secret ID, specific revision, JSON field
  NAME                       secret name, latest revision
  PATH/NAME                  secret by path and name, latest revision
  PATH/NAME@REVISION:FIELD   path/name with revision and field

REVISION can be an integer, "latest", or "latest_enabled".

FIELD extracts a string field from a JSON-encoded secret.

Examples:

  {{ scw://11111111-1111-1111-1111-111111111111 }}
  {{ scw://11111111-1111-1111-1111-111111111111@2:api-key }}
  {{ scw://my-app/db-password@latest }}`,
		Namespace: "secret", //nolint:goconst
		Resource:  "inject",
		ArgsType:  reflect.TypeOf(secretInjectArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:  "in-file",
				Short: "Template file to read from (default: stdin)",
			},
			{
				Name:  "out-file",
				Short: "Output file (default: stdout)",
			},
			{
				Name:    "file-mode",
				Short:   "Permissions of the output file, in octal (only used with --out-file)",
				Default: core.DefaultValueSetter("0600"),
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Groups: []string{"security"},
		Run:    secretInjectRun,
		Examples: []*core.Example{
			{
				Short: "Render a template from stdin to stdout",
				Raw:   `echo "password: {{ scw://11111111-1111-1111-1111-111111111111 }}" | scw secret inject`,
			},
			{
				Short: "Render a template file to an output file",
				Raw:   `scw secret inject in-file=config.tpl out-file=config.yaml`,
			},
			{
				Short: "Extract a JSON field and write to a file with custom permissions",
				Raw:   `scw secret inject in-file=config.tpl out-file=config.yaml file-mode=0640`,
			},
		},
	}
}

func secretInjectRun(ctx context.Context, argsI any) (any, error) {
	args := argsI.(*secretInjectArgs)

	input, err := readInjectInput(ctx, args.InFile)
	if err != nil {
		return nil, err
	}

	api := secret.NewAPI(core.ExtractClient(ctx))

	rendered, err := renderTemplate(input, api, args.Region)
	if err != nil {
		return nil, err
	}

	if args.OutFile != "" {
		mode, err := strconv.ParseUint(args.FileMode, 8, 32)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid file-mode %q: expected octal value like 0600",
				args.FileMode,
			)
		}

		if err := writeOutputFile(args.OutFile, rendered, os.FileMode(mode)); err != nil {
			return nil, err
		}

		return &core.SuccessResult{Empty: true}, nil
	}

	return core.RawResult(rendered), nil
}

func writeOutputFile(path, content string, mode os.FileMode) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("writing output file: %w", cerr)
		}
	}()

	// Chmod before writing so the secret never lands in a wider-mode file.
	if err := f.Chmod(mode); err != nil {
		return fmt.Errorf("setting output file permissions: %w", err)
	}

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	return nil
}

func readInjectInput(ctx context.Context, inFile string) (string, error) {
	if inFile != "" {
		data, err := os.ReadFile(inFile)
		if err != nil {
			return "", fmt.Errorf("reading template file: %w", err)
		}

		return string(data), nil
	}

	data, err := io.ReadAll(core.ExtractStdin(ctx))
	if err != nil {
		return "", fmt.Errorf("reading stdin: %w", err)
	}

	return string(data), nil
}

func renderTemplate(input string, api *secret.API, region scw.Region) (string, error) {
	matches := secretRefRegex.FindAllStringSubmatch(input, -1)

	// Deduplicate and parse all references before fetching.
	type pendingRef struct {
		rawRef string
		ref    *parsedSecretRef
	}

	seen := make(map[string]struct{}, len(matches))
	pending := make([]pendingRef, 0, len(matches))

	for _, m := range matches {
		rawRef := m[1]
		if _, ok := seen[rawRef]; ok {
			continue
		}

		seen[rawRef] = struct{}{}

		ref, err := parseSecretRef(rawRef)
		if err != nil {
			return "", fmt.Errorf("invalid reference %q: %w", rawRef, err)
		}

		pending = append(pending, pendingRef{rawRef: rawRef, ref: ref})
	}

	// Fetch all unique secrets concurrently.
	var mu sync.Mutex

	cache := make(map[string]string, len(pending))
	g, _ := errgroup.WithContext(context.Background())

	for _, p := range pending {
		g.Go(func() error {
			value, err := resolveSecretRef(api, p.ref, region)
			if err != nil {
				return fmt.Errorf("resolving {{ scw://%s }}: %w", p.rawRef, err)
			}

			mu.Lock()
			cache[p.rawRef] = value
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return "", err
	}

	rendered := secretRefRegex.ReplaceAllStringFunc(input, func(match string) string {
		sub := secretRefRegex.FindStringSubmatch(match)

		return cache[sub[1]]
	})

	return rendered, nil
}

// parseSecretRef parses a reference of the form:
//
//	(UUID|PATH/NAME)[@REVISION][:FIELD]
func parseSecretRef(raw string) (*parsedSecretRef, error) {
	ref := &parsedSecretRef{revision: revisionLatest}

	// Strip optional :FIELD (last colon, since neither UUID nor revision contain one)
	if idx := strings.LastIndex(raw, ":"); idx != -1 {
		ref.field = raw[idx+1:]
		raw = raw[:idx]

		if ref.field == "" {
			return nil, errors.New("empty field name after ':'")
		}
	}

	// Strip optional @REVISION
	if idx := strings.Index(raw, "@"); idx != -1 {
		ref.revision = raw[idx+1:]
		raw = raw[:idx]

		if ref.revision == "" {
			return nil, errors.New("empty revision after '@'")
		}
	}

	if raw == "" {
		return nil, errors.New("empty secret identifier")
	}

	if isSecretUUID(raw) {
		ref.secretID = raw
	} else {
		// PATH/NAME: last path segment is the name, the rest is the path.
		if idx := strings.LastIndex(raw, "/"); idx != -1 {
			ref.secretPath = "/" + strings.TrimPrefix(raw[:idx], "/")
			ref.secretName = raw[idx+1:]
		} else {
			ref.secretPath = "/"
			ref.secretName = raw
		}

		if ref.secretName == "" {
			return nil, errors.New("empty secret name in path reference")
		}
	}

	return ref, nil
}

func isSecretUUID(s string) bool {
	if len(s) != 36 {
		return false
	}

	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
				return false
			}
		}
	}

	return true
}

func resolveSecretRef(api *secret.API, ref *parsedSecretRef, region scw.Region) (string, error) {
	var (
		resp *secret.AccessSecretVersionResponse
		err  error
	)

	if ref.secretID != "" {
		resp, err = api.AccessSecretVersion(&secret.AccessSecretVersionRequest{
			Region:   region,
			SecretID: ref.secretID,
			Revision: ref.revision,
		})
	} else {
		resp, err = api.AccessSecretVersionByPath(&secret.AccessSecretVersionByPathRequest{
			Region:     region,
			SecretName: ref.secretName,
			SecretPath: ref.secretPath,
			Revision:   ref.revision,
		})
	}

	if err != nil {
		return "", err
	}

	data := resp.Data
	if ref.field != "" {
		data, err = getSecretVersionField(data, ref.field)
		if err != nil {
			return "", err
		}
	}

	return string(data), nil
}
