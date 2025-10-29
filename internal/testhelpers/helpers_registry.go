package testhelpers

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/registry/v1"
)

func PushRegistryImage(dockerImage, namespaceMetaKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		namespaceEndpoint := ctx.Meta.Render(fmt.Sprintf("{{ .%s.Endpoint }}", namespaceMetaKey))
		remote := fmt.Sprintf("%s/%s", namespaceEndpoint, dockerImage)

		return core.BeforeFuncCombine(
			core.BeforeFuncOsExec("docker", "pull", dockerImage),
			core.BeforeFuncOsExec("docker", "tag", dockerImage, remote),
			core.BeforeFuncOsExec("docker", "push", remote),
		)(
			ctx,
		)
	}
}

func StoreImageIdentifierInMeta(
	namespaceMetaKey, dockerImageName, imageMetaKey string,
) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		// List images
		imageName := strings.Split(dockerImageName, ":")[0]
		namespaceID := ctx.Meta.Render(fmt.Sprintf("{{ .%s.ID }}", namespaceMetaKey))
		args := strings.Split(
			fmt.Sprintf(
				"scw registry image list namespace-id=%s name=%s",
				namespaceID,
				imageName,
			),
			" ",
		)
		imageListResult := ctx.ExecuteCmd(args)

		// Select the image
		imageList, ok := imageListResult.([]registry.CustomImage)
		if !ok {
			return fmt.Errorf("result is not []registry.CustomImage but %T", imageListResult)
		}
		if len(imageList) != 1 {
			return fmt.Errorf(
				"expected exactly 1 image with name %q, got %d",
				imageName,
				len(imageList),
			)
		}
		image := imageList[0]

		// Build image identifier and store it in Meta
		if len(image.Tags) != 1 {
			return fmt.Errorf(
				"unexpected number of tags for image %s: expected 1, got %d",
				image.Name,
				len(image.Tags),
			)
		}
		ctx.Meta[imageMetaKey] = fmt.Sprintf("%s:%s", image.FullName, image.Tags[0])

		return nil
	}
}
