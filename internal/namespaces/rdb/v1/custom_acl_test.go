package rdb

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_AddACL(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb acl add instance-id={{ .Instance.ID }} rules.0.ip=4.2.3.4",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteInstance(),
	}))
}

func Test_DeleteACL(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance("PostgreSQL-12"),
			core.ExecBeforeCmd("scw rdb acl add instance-id={{ .Instance.ID }} rules.0.ip=1.2.3.4"),
		),
		Cmd:       "scw rdb acl remove instance-id={{ .Instance.ID }} acl-rule-ips.0=1.2.3.4",
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))
}
