package iam_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/iam/v1alpha1"
)

func Test_ListApplication(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: iam.GetCommands(),
		Cmd:      "scw iam application list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListJwt(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: iam.GetCommands(),
		Cmd:      "scw iam jwt list 11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListLog(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: iam.GetCommands(),
		Cmd:      "scw iam log list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPermissionSet(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: iam.GetCommands(),
		Cmd:      "scw iam permission-set list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPolicy(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: iam.GetCommands(),
		Cmd:      "scw iam policy list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListRule(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: iam.GetCommands(),
		Cmd:      "scw iam rule list 11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListScimTokens(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: iam.GetCommands(),
		Cmd:      "scw iam scim-tokens list 11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: iam.GetCommands(),
		Cmd:      "scw iam user list organization-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}
