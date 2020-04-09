package baremetal

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

// createServer creates a baremetal instance
// register it in the context Meta at metaKey.
func createServer(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, "scw baremetal server create")
}

// deleteServer deletes a server
// previously registered in the context Meta at metaKey.
func deleteServer(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw baremetal server delete {{ ." + metaKey + ".ID }}")
}
