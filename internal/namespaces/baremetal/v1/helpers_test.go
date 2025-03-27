package baremetal_test

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
)

// createServerAndWait creates a baremetal instance
// register it in the context Meta at metaKey.
func createServerAndWait(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		metaKey,
		"scw baremetal server create zone="+region+" type="+offerName+" -w",
	)
}

func createServerAndWaitDefault(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, "scw baremetal server create type="+offerName+" -w")
}

func createServer(metaKey string, offerType string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		metaKey,
		"scw baremetal server create zone="+region+" type="+offerType,
	)
}

// deleteServer deletes a server
// previously registered in the context Meta at metaKey.
//
//nolint:unparam
func deleteServer(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd(
		fmt.Sprintf("scw baremetal server delete zone="+region+" {{ .%s.ID }}", metaKey),
	)
}

func deleteServerDefault(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd(fmt.Sprintf("scw baremetal server delete {{ .%s.ID }}", metaKey))
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
