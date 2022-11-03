package rdb

import (
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_AddACL(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb acl add instance-id={{ .Instance.ID }} rule.0.ip=4.2.3.4",
		Check:      core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			// wait for ACL rule changes
			func(ctx *core.AfterFuncCtx) error {
				time.Sleep(5 * time.Second)
				return nil
			},
			deleteInstance(),
		),
	}))
}

func Test_DeleteACL(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance("PostgreSQL-12"),
			core.ExecBeforeCmd("scw rdb acl add instance-id={{ .Instance.ID }} rule.0.ip=1.2.3.4"),
			func(ctx *core.BeforeFuncCtx) error {
				time.Sleep(5 * time.Second)
				return nil
			},
		),
		Cmd:   "scw rdb acl delete instance-id={{ .Instance.ID }} rule.0.ip=1.2.3.4",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			func(ctx *core.AfterFuncCtx) error {
				time.Sleep(5 * time.Second)
				return nil
			},
			deleteInstance(),
		)}))
}
