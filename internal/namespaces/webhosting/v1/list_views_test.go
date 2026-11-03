package webhosting_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/webhosting/v1"
)

func Test_ListBackup(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: webhosting.GetCommands(),
		Cmd:      "scw webhosting backup list hosting-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListControlPanel(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: webhosting.GetCommands(),
		Cmd:      "scw webhosting control-panel list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDatabase(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: webhosting.GetCommands(),
		Cmd:      "scw webhosting database list hosting-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDatabaseUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: webhosting.GetCommands(),
		Cmd:      "scw webhosting database-user list hosting-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListFtpAccount(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: webhosting.GetCommands(),
		Cmd:      "scw webhosting ftp-account list hosting-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListHosting(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: webhosting.GetCommands(),
		Cmd:      "scw webhosting hosting list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListMailAccount(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: webhosting.GetCommands(),
		Cmd:      "scw webhosting mail-account list hosting-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListWebsite(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: webhosting.GetCommands(),
		Cmd:      "scw webhosting website list hosting-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}
