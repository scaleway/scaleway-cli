package rdb

import (
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_WaitLog(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["startDate"] = time.Now().Format(time.RFC3339)
				return nil
			},
			createInstance("PostgreSQL-12"),
			func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["endDate"] = time.Now().Format(time.RFC3339)
				return nil
			},
			core.ExecStoreBeforeCmd("Logs", "scw rdb log prepare {{ .Instance.ID }} start-date={{ .startDate }} end-date={{ .endDate }}"),
		),
		Cmd:       "scw rdb log wait {{ .Logs[0].ID }}",
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))
}
