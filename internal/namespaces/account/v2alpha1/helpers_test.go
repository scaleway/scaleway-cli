package account

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

func addSSHKey(metaKey string, key string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		ctx.Meta[metaKey] = ctx.ExecuteCmd([]string{
			"scw", "account", "ssh-key", "add", "public-key=" + key,
		})
		return nil
	}
}
