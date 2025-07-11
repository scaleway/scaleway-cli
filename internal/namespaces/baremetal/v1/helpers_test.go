package baremetal_test

import (
	"fmt"
	"os"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}

// createServerAndWait creates a baremetal instance
// register it in the context Meta at metaKey.
func createServerAndWait() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Server",
		"scw baremetal server create type="+offerNameNVME+" zone="+zone+" -w",
	)
}

func createServer(metaKey string, offer string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, "scw baremetal server create zone="+zone+" type="+offer)
}

// deleteServer deletes a server
// previously registered in the context Meta at metaKey.
//
//nolint:unparam
func deleteServer(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd(
		fmt.Sprintf("scw baremetal server delete {{ .%s.ID }} zone=%s", metaKey, zone),
	)
}

// add an ssh key with a given meta key
func addSSH(metaKey string, key string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		ctx.Meta[metaKey] = ctx.ExecuteCmd([]string{
			"scw", "iam", "ssh-key", "create", "public-key=" + key,
		})

		return nil
	}
}

// delete an ssh key with a given meta key
func deleteSSH(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd(fmt.Sprintf("scw iam ssh-key delete {{ .%s.ID }}", metaKey))
}
