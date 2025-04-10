package rdb_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
)

func Test_ListUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd:       "scw rdb user list instance-id={{ .Instance.ID }}",
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))
}

func Test_CreateUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: fmt.Sprintf(
			"scw rdb user create instance-id={{ $.Instance.Instance.ID }} name=%s password=%s",
			name,
			password,
		),
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))

	t.Run("With password generator", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: fmt.Sprintf(
			"scw rdb user create instance-id={{ $.Instance.Instance.ID }} name=%s generate-password=true",
			name,
		),
		// do not check the golden as the password generated locally and on CI will necessarily be different
		Check:     core.TestCheckExitCode(0),
		AfterFunc: deleteInstance(),
	}))
}

func Test_UpdateUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
			core.ExecBeforeCmd(
				fmt.Sprintf(
					"scw rdb user create instance-id={{ $.Instance.Instance.ID }} name=%s password=%s",
					name,
					password,
				),
			),
		),
		Cmd: fmt.Sprintf(
			"scw rdb user update instance-id={{ $.Instance.Instance.ID }} name=%s password=Newp1ssw0rd! is-admin=true",
			name,
		),
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))

	t.Run("With password generator", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
			core.ExecBeforeCmd(
				fmt.Sprintf(
					"scw rdb user create instance-id={{ $.Instance.Instance.ID }} name=%s password=%s",
					name,
					password,
				),
			),
		),
		Cmd: fmt.Sprintf(
			"scw rdb user update instance-id={{ $.Instance.Instance.ID }} name=%s generate-password=true is-admin=true",
			name,
		),
		// do not check the golden as the password generated locally and on CI will necessarily be different
		Check:     core.TestCheckExitCode(0),
		AfterFunc: deleteInstance(),
	}))
}
