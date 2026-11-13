package domain_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	domain "github.com/scaleway/scaleway-cli/v2/internal/namespaces/domain/v2beta1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
	"github.com/stretchr/testify/require"
)

func Test_RecordImport(t *testing.T) {
	cmds := domain.GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("DryRun", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			path := filepath.Join(ctx.T.TempDir(), "zone.txt")
			require.NoError(
				ctx.T,
				os.WriteFile(path, []byte("cli-import-dry IN A 1.2.3.4\n"), 0o600),
			)
			ctx.Meta["BindFile"] = path

			return nil
		},
		// dry-run never calls the API; any zone name is fine for parse context
		Cmd: "scw dns record import example.privatedns content=@{{ .BindFile }} dry-run=true",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"PN",
				"scw vpc private-network create name=cli-dns-import",
			),
			func(ctx *core.BeforeFuncCtx) error {
				path := filepath.Join(ctx.T.TempDir(), "zone.txt")
				require.NoError(
					ctx.T,
					os.WriteFile(path, []byte("cli-import-simple IN A 1.2.3.4\n"), 0o600),
				)
				ctx.Meta["BindFile"] = path

				return nil
			},
		),
		Cmd: "scw dns record import {{ .PN.ID }}.{{ .PN.VpcID }}.privatedns content=@{{ .BindFile }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(
				"scw dns record delete {{ .PN.ID }}.{{ .PN.VpcID }}.privatedns"+
					" name=cli-import-simple type=A data=1.2.3.4",
			),
			core.ExecAfterCmd("scw vpc private-network delete {{ .PN.ID }}"),
		),
	}))
}
