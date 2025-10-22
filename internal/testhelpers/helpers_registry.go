package testhelpers

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func PushRegistryImage(dockerImage, namespaceMetaKey, imageName string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		namespaceEndpoint := ctx.Meta.Render(fmt.Sprintf("{{ .%s.Endpoint }}", namespaceMetaKey))
		remote := fmt.Sprintf("%s/%s:latest", namespaceEndpoint, imageName)

		return core.BeforeFuncCombine(
			core.BeforeFuncOsExec("docker", "pull", dockerImage),
			core.BeforeFuncOsExec("docker", "tag", dockerImage, remote),
			core.BeforeFuncOsExec("docker", "push", remote),
		)(
			ctx,
		)
	}
}

func StoreImageInMeta(metaKey, namespaceMetaKey, imageName string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		namespaceID := ctx.Meta.Render(fmt.Sprintf("{{ .%s.ID }}", namespaceMetaKey))

		return core.ExecStoreBeforeCmd(
			metaKey,
			fmt.Sprintf("scw registry image list namespace-id=%s name=%s", namespaceID, imageName),
		)(
			ctx,
		)
	}
}
