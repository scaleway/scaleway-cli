package baremetal

import "github.com/scaleway/scaleway-cli/internal/core"

// createServer creates a baremetal instance
// register it in the context Meta at metaKey.
func createServer(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, "scw baremetal server create")
}
