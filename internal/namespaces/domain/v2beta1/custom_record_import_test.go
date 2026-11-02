package domain_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	domain "github.com/scaleway/scaleway-cli/v2/internal/namespaces/domain/v2beta1"
	"github.com/stretchr/testify/require"
)

// Private DNS zone used when recording cassettes (must have default name servers).
const testRecordImportDNSZone = "3d910824-710e-4ece-92ad-3a96f59fe943.b1972fa8-5349-4d37-8777-f4f205bbd322.privatedns"

func Test_RecordImport(t *testing.T) {
	cmds := domain.GetCommands()

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
		Cmd: "scw dns record import " + testRecordImportDNSZone +
			" content=@{{ .BindFile }} dry-run=true",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			path := filepath.Join(ctx.T.TempDir(), "zone.txt")
			require.NoError(
				ctx.T,
				os.WriteFile(path, []byte("cli-import-simple IN A 1.2.3.4\n"), 0o600),
			)
			ctx.Meta["BindFile"] = path

			return nil
		},
		Cmd: "scw dns record import " + testRecordImportDNSZone + " content=@{{ .BindFile }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.ExecAfterCmd(
			"scw dns record delete " + testRecordImportDNSZone +
				" name=cli-import-simple type=A data=1.2.3.4",
		),
	}))
}
