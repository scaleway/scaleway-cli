package serverless_sqldb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/serverless_sqldb/v1alpha1"
)

func Test_ListBackup(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: serverless_sqldb.GetCommands(),
		Cmd:      "scw sdb-sql backup list database-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDatabase(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: serverless_sqldb.GetCommands(),
		Cmd:      "scw sdb-sql database list",
		Check:    core.TestCheckGolden(),
	}))
}
