package mnq_test

import (
	"os"
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
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				result, isSuccessResult := ctx.Result.(*core.SuccessResult)
				assert.True(
					t,
					isSuccessResult,
					"Expected result to be of type *core.SuccessResult, got %s",
					reflect.TypeOf(result).String(),
				)
				assert.NotNil(t, result)
				expectedContextFile := result.Resource
				if !mnq.FileExists(expectedContextFile) {
					t.Errorf(
						"Expected credentials file not found expected [%s] ",
						expectedContextFile,
					)
				} else {
					ctx.Meta["deleteFiles"] = []string{expectedContextFile}
				}
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteNATSAccount("NATS"),
			func(ctx *core.AfterFuncCtx) error {
				if ctx.Meta["deleteFiles"] == nil {
					return nil
				}
				filesToDelete := ctx.Meta["deleteFiles"].([]string)
				for _, file := range filesToDelete {
					err := os.Remove(file)
					if err != nil {
						t.Errorf("Failed to delete the file : %s", err)
					}
				}

				return nil
			},
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
