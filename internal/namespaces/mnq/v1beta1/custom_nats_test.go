package mnq

import (
	"os"
	"regexp"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_CreateContext(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createNATSAccount("NATS"),
		Cmd:        "scw mnq nats create-context nats-account-id={{ .NATS.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGoldenAndReplacePatterns(
				core.GoldenReplacement{
					Pattern:     regexp.MustCompile(`cli[\w-]*creds[\w-]*`),
					Replacement: "credential-placeholder",
				},
				core.GoldenReplacement{
					Pattern:     regexp.MustCompile("(Select context using `nats context select )cli[\\w-]*`"),
					Replacement: "Select context using `nats context select context-placeholder`",
				},
			),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				result := ctx.Result.(*core.SuccessResult)
				expectedContextFile := result.Resource
				if !fileExists(expectedContextFile) {
					t.Errorf("Expected credentials file not found expected [%s] ", expectedContextFile)
				} else {
					ctx.Meta["deleteFiles"] = []string{expectedContextFile}
				}
			},
		),
		AfterFunc: core.AfterFuncCombine(deleteNATSAccount("NATS"), func(ctx *core.AfterFuncCtx) error {
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
		Commands: GetCommands(),
		Cmd:      "scw mnq nats create-context nats-account-id=Wrong-id",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}

func Test_CreateContextWithNoAccount(t *testing.T) {
	t.Run("With No Nats Account", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw mnq nats create-context",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}

func Test_CreateContextNoInteractiveTermAndMultiAccount(t *testing.T) {
	t.Run("Multi Nats Account and no ID", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(createNATSAccount("NATS"), createNATSAccount("NATS2")),
		Cmd:        "scw mnq nats create-context",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc: core.AfterFuncCombine(deleteNATSAccount("NATS"), deleteNATSAccount("NATS2")),
	}))
}
