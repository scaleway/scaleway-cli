package baremetal

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

// createServerAndWait creates a baremetal instance
// register it in the context Meta at metaKey.
func createServerAndWait(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, "scw baremetal server create zone=nl-ams-1 type=GP-BM2-S -w")
}

func createServer(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, "scw baremetal server create zone=nl-ams-1 type=GP-BM2-S")
}

// deleteServer deletes a server
// previously registered in the context Meta at metaKey.
// nolint:unparam
func deleteServer(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw baremetal server delete zone=nl-ams-1 {{ ." + metaKey + ".ID }}")
}

// add an ssh key with a given meta key
func addSSH(metaKey string, key string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		ctx.Meta[metaKey] = ctx.ExecuteCmd([]string{
			"scw", "account", "ssh-key", "add", "public-key=" + key,
		})
		return nil
	}
}

// delete an ssh key with a given meta key
func deleteSSH(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw account ssh-key remove {{ ." + metaKey + ".ID }}")
}
