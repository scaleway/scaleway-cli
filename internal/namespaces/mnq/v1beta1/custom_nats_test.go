package mnq_test

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	mnq "github.com/scaleway/scaleway-cli/v2/internal/namespaces/mnq/v1beta1"
	"github.com/stretchr/testify/assert"
)

func Test_CreateContext(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   mnq.GetCommands(),
		BeforeFunc: createNATSAccount("NATS"),
		Cmd:        "scw mnq nats create-context nats-account-id={{ .NATS.ID }}",
		TmpHomeDir: true,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGoldenAndReplacePatterns(
				core.GoldenReplacement{
					Pattern:     regexp.MustCompile(`cli[\w-]*creds[\w-]*`),
					Replacement: "credential-placeholder",
				},
				core.GoldenReplacement{
					Pattern: regexp.MustCompile(
						"(Select context using `nats context select )cli[\\w-]*`",
					),
					Replacement: "Select context using `nats context select context-placeholder`",
				},
			),
			checkContextFile,
		),
		AfterFunc: core.AfterFuncCombine(
			deleteNATSAccount("NATS"),
			afterFuncDeleteFiles,
		),
	}))
}

func Test_CreateContextWithWrongId(t *testing.T) {
	t.Run("Wrong Account ID", core.Test(&core.TestConfig{
		Commands: mnq.GetCommands(),
		Cmd:      "scw mnq nats create-context nats-account-id=Wrong-id",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}

func Test_CreateContextWithNoAccount(t *testing.T) {
	t.Run("With No Nats Account", core.Test(&core.TestConfig{
		Commands: mnq.GetCommands(),
		Cmd:      "scw mnq nats create-context",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}

func Test_CreateContextNoInteractiveTermAndMultiAccount(t *testing.T) {
	t.Run("Multi Nats Account and no ID", core.Test(&core.TestConfig{
		Commands:   mnq.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(createNATSAccount("NATS"), createNATSAccount("NATS2")),
		Cmd:        "scw mnq nats create-context",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc: core.AfterFuncCombine(deleteNATSAccount("NATS"), deleteNATSAccount("NATS2")),
	}))
}

func beforeFuncCopyConfigToTmpHome() core.BeforeFunc {
	return core.BeforeFuncWhenUpdatingCassette(func(ctx *core.BeforeFuncCtx) error {
		realHomeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		realConfigPath := filepath.Join(realHomeDir, ".config", "scw", "config.yaml")
		if _, err := os.Stat(realConfigPath); os.IsNotExist(err) {
			return nil
		}

		tmpHomeDir := ctx.OverrideEnv["HOME"]
		tmpConfigDir := filepath.Join(tmpHomeDir, ".config", "scw")
		if err := os.MkdirAll(tmpConfigDir, 0o0755); err != nil {
			return err
		}

		data, err := os.ReadFile(realConfigPath)
		if err != nil {
			return err
		}

		return os.WriteFile(filepath.Join(tmpConfigDir, "config.yaml"), data, 0o644)
	})
}

func checkContextFile(t *testing.T, ctx *core.CheckFuncCtx) {
	t.Helper()
	result, isSuccessResult := ctx.Result.(*core.SuccessResult)
	assert.True(t, isSuccessResult, "Expected result to be of type *core.SuccessResult, got %s", reflect.TypeOf(result).String())
	assert.NotNil(t, result)
	expectedContextFile := result.Resource
	if !mnq.FileExists(expectedContextFile) {
		t.Errorf("Expected credentials file not found expected [%s]", expectedContextFile)
	} else {
		ctx.Meta["deleteFiles"] = []string{expectedContextFile}
	}
}

func checkContextFileInXDGConfigHome(t *testing.T, ctx *core.CheckFuncCtx) {
	t.Helper()
	checkContextFile(t, ctx)

	result := ctx.Result.(*core.SuccessResult)
	expectedContextFile := result.Resource
	xdgConfigHome := ctx.OverrideEnv["XDG_CONFIG_HOME"]
	tmpHomeDir := ctx.OverrideEnv["HOME"]

	expectedContextDir := filepath.Join(xdgConfigHome, "nats", "context")
	assert.Contains(t, expectedContextFile, expectedContextDir, "Context file should be in XDG_CONFIG_HOME/nats/context, got: %s", expectedContextFile)

	tmpHomeNatsDir := filepath.Join(tmpHomeDir, ".config", "nats", "context")
	_, err := os.Stat(tmpHomeNatsDir)
	assert.True(t, os.IsNotExist(err), "Files should not be created in HOME/.config/nats/context when XDG_CONFIG_HOME is set: %s", tmpHomeNatsDir)
}

func afterFuncDeleteFiles(ctx *core.AfterFuncCtx) error {
	if ctx.Meta["deleteFiles"] == nil {
		return nil
	}
	filesToDelete := ctx.Meta["deleteFiles"].([]string)
	for _, file := range filesToDelete {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func Test_CreateContextWithXDGConfigHome(t *testing.T) {
	xdgConfigHomeDir := t.TempDir()

	t.Run("XDG_CONFIG_HOME compliance", core.Test(&core.TestConfig{
		Commands: mnq.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			beforeFuncCopyConfigToTmpHome(),
			createNATSAccount("NATS"),
		),
		Cmd:        "scw mnq nats create-context nats-account-id={{ .NATS.ID }}",
		TmpHomeDir: true,
		OverrideEnv: map[string]string{
			"XDG_CONFIG_HOME": xdgConfigHomeDir,
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGoldenAndReplacePatterns(
				core.GoldenReplacement{
					Pattern:     regexp.MustCompile(`cli[\w-]*creds[\w-]*`),
					Replacement: "credential-placeholder",
				},
				core.GoldenReplacement{
					Pattern: regexp.MustCompile(
						"(Select context using `nats context select )cli[\\w-]*`",
					),
					Replacement: "Select context using `nats context select context-placeholder`",
				},
			),
			checkContextFileInXDGConfigHome,
		),
		AfterFunc: core.AfterFuncCombine(
			deleteNATSAccount("NATS"),
			afterFuncDeleteFiles,
		),
	}))
}
