package secret_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/scaleway/scaleway-cli/v2/core"
	secret "github.com/scaleway/scaleway-cli/v2/internal/namespaces/secret/v1beta1"
)

// writeInjectTemplate renders tpl with [[ ]] delimiters against Meta into a file stored as TemplatePath.
func writeInjectTemplate(tpl string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		var buf strings.Builder

		err := template.Must(template.New("inject").Delims("[[", "]]").Parse(tpl)).Execute(&buf, ctx.Meta)
		if err != nil {
			return err
		}

		path := filepath.Join(ctx.T.TempDir(), "inject.tpl")
		ctx.Meta["TemplatePath"] = path

		return os.WriteFile(path, []byte(buf.String()), 0o600)
	}
}

func Test_SecretInject(t *testing.T) {
	tests := []struct {
		name         string
		template     string
		withSecret   bool
		wantExitCode int
	}{
		{
			name:       "by-id-fetched-once",
			template:   `{{ secret "[[ .Secret.ID ]]" | jsonField "password" }},{{ with true }}{{ secret "[[ .Secret.ID ]]" | jsonField "password" }}{{ end }}`,
			withSecret: true,
		},
		{
			name:       "by-path-with-revision-and-json-field",
			template:   `{{ secret "cli-test-inject/[[ .Secret.Name ]]@1" | jsonField "password" }}`,
			withSecret: true,
		},
		{
			name:       "secret-behind-value-dependent-branch",
			template:   `{{ if eq (secret "[[ .Secret.ID ]]" | jsonField "password") "hunter2" }}{{ secret "cli-test-inject/[[ .Secret.Name ]]@1" | jsonField "password" }}{{ end }}`,
			withSecret: true,
		},
		{
			name:         "empty-name",
			template:     `{{ secret "cli-test-inject/@1" }}`,
			wantExitCode: 1,
		},
		{
			name:         "unknown-secret-outputs-nothing",
			template:     `before {{ secret "cli-test-inject/missing" }}`,
			wantExitCode: 1,
		},
		{
			name:         "invalid-template",
			template:     `{{ secret `,
			wantExitCode: 1,
		},
	}

	for _, tt := range tests {
		config := &core.TestConfig{
			Commands:   secret.GetCommands(),
			BeforeFunc: writeInjectTemplate(tt.template),
			// Args are rendered as Go templates, so the template goes through a file.
			Cmd: "scw secret inject template=@{{ .TemplatePath }}",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(tt.wantExitCode),
			),
		}
		if tt.withSecret {
			config.BeforeFunc = core.BeforeFuncCombine(
				core.ExecStoreBeforeCmd(
					"Secret",
					"scw secret secret create name=cli-test-inject-"+tt.name+" path=/cli-test-inject",
				),
				createSecretVersion(`{"password":"hunter2"}`),
				config.BeforeFunc,
			)
			config.AfterFunc = deleteSecret()
		}

		t.Run(tt.name, core.Test(config))
	}
}
