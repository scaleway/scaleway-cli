package testhelpers

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func PushRegistryImage(dockerImage, namespaceMetaKey, imageName, metaKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		namespaceEndpoint := ctx.Meta.Render(fmt.Sprintf("{{ .%s.Endpoint }}", namespaceMetaKey))
		namespaceID := ctx.Meta.Render(fmt.Sprintf("{{ .%s.ID }}", namespaceMetaKey))
		remote := fmt.Sprintf("%s/%s:latest", namespaceEndpoint, imageName)

		return core.BeforeFuncCombine(
			core.BeforeFuncOsExec("docker", "pull", dockerImage),
			core.BeforeFuncOsExec("docker", "tag", dockerImage, remote),
			core.BeforeFuncOsExec("docker", "push", remote),
			core.ExecStoreBeforeCmd(
				metaKey,
				fmt.Sprintf(
					"scw registry image list namespace-id=%s name=%s",
					namespaceID,
					imageName,
				),
			),
		)(
			ctx,
		)
	}
}
