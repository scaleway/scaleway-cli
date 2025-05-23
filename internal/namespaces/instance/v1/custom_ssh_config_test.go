package instance_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/sshconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SSHConfigInstall(t *testing.T) {
	t.Run("Install config and create default", core.Test(&core.TestConfig{
		TmpHomeDir: true,
		Commands:   instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", testServerCommand("stopped=true ip=new")),
		Args:       []string{"scw", "instance", "ssh", "install-config"},
		Check: core.TestCheckCombine(
			core.TestCheckGoldenAndReplacePatterns(
				core.GoldenReplacement{
					Pattern:     regexp.MustCompile("generated to .*scaleway.config"),
					Replacement: "generated to /tmp/scw/.ssh/scaleway.config",
				},
			),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				server := ctx.Meta["Server"].(*instance.ServerWithWarningsResponse).Server

				configPath := sshconfig.ConfigFilePath(ctx.Meta["HOME"].(string))
				content, err := os.ReadFile(configPath)
				require.NoError(t, err)
				assert.Contains(t, string(content), "Host "+server.Name)

				included, err := sshconfig.ConfigIsIncluded(ctx.Meta["HOME"].(string))
				require.NoError(t, err)
				assert.True(t, included)
			},
		),
		AfterFunc: deleteServer("Server"),
	}))

	t.Run("Install config and include", core.Test(&core.TestConfig{
		TmpHomeDir: true,
		Commands:   instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			func(ctx *core.BeforeFuncCtx) error {
				homeDir := ctx.Meta["HOME"].(string)
				configPath := sshconfig.DefaultConfigFilePath(homeDir)
				err := os.Mkdir(filepath.Join(homeDir, ".ssh"), 0o700)
				require.NoError(t, err)
				err = os.WriteFile(configPath, []byte(`Host myhost`), 0o600)
				assert.NoError(t, err)

				return nil
			},
			core.ExecStoreBeforeCmd("Server", testServerCommand("stopped=true ip=new")),
		),
		Args: []string{"scw", "instance", "ssh", "install-config"},
		Check: core.TestCheckCombine(
			core.TestCheckGoldenAndReplacePatterns(
				core.GoldenReplacement{
					Pattern:     regexp.MustCompile("generated to .*scaleway.config"),
					Replacement: "generated to /tmp/scw/.ssh/scaleway.config",
				},
			),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				server := ctx.Meta["Server"].(*instance.ServerWithWarningsResponse).Server

				defaultConfigPath := sshconfig.DefaultConfigFilePath(ctx.Meta["HOME"].(string))
				content, err := os.ReadFile(defaultConfigPath)
				require.NoError(t, err)
				assert.Contains(t, string(content), "Include scaleway.config")
				assert.Contains(t, string(content), "Host myhost")

				configPath := sshconfig.ConfigFilePath(ctx.Meta["HOME"].(string))
				content, err = os.ReadFile(configPath)
				require.NoError(t, err)
				assert.Contains(t, string(content), "Host "+server.Name)

				included, err := sshconfig.ConfigIsIncluded(ctx.Meta["HOME"].(string))
				require.NoError(t, err)
				assert.True(t, included)
			},
		),
		AfterFunc: deleteServer("Server"),
	}))
}
