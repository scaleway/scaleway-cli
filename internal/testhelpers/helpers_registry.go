package testhelpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/registry/v1"
)

var (
	ErrResultIsNot                    = errors.New("result is not")
	ErrExpectedExactly1ImageWithName  = errors.New("expected exactly 1 image with name")
	ErrUnexpectedNumberOfTagsForImage = errors.New("unexpected number of tags for image")
	ErrExpectedExactly1TagWithName    = errors.New("expected exactly 1 tag with name")
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
		imageName, _, _ := strings.Cut(dockerImageName, ":")
		namespaceID := ctx.Meta.Render(fmt.Sprintf("{{ .%s.ID }}", namespaceMetaKey))
		imageListResult := core.ExecBeforeCmdWithResult(ctx, fmt.Sprintf(
			"scw registry image list namespace-id=%s name=%s",
			namespaceID,
			imageName,
		))

		// Select the image
		imageList, ok := imageListResult.([]registry.CustomImage)
		if !ok {
			return fmt.Errorf("%w []registry.CustomImage but %T", ErrResultIsNot, imageListResult)
		}
		if len(imageList) != 1 {
			return fmt.Errorf(
				"%w %q, got %d",
				ErrExpectedExactly1ImageWithName,
				imageName,
				len(imageList),
			)
		}
		image := imageList[0]

		// Build image identifier and store it in Meta
		if len(image.Tags) != 1 {
			return fmt.Errorf(
				"%w %s: expected 1, got %d",
				ErrUnexpectedNumberOfTagsForImage,
				image.Name,
				len(image.Tags),
			)
		}
		ctx.Meta[imageMetaKey] = fmt.Sprintf("%s:%s", image.FullName, image.Tags[0])

		return nil
	}
}

func StoreTagIDInMeta(namespaceMetaKey, dockerImageName, tagMetaKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		// List images
		dockerImageNameSplit := strings.Split(dockerImageName, ":")
		imageName := dockerImageNameSplit[0]
		namespaceID := ctx.Meta.Render(fmt.Sprintf("{{ .%s.ID }}", namespaceMetaKey))
		imageListResult := core.ExecBeforeCmdWithResult(
			ctx,
			fmt.Sprintf("scw registry image list namespace-id=%s name=%s",
				namespaceID,
				imageName,
			),
		)

		// Select the image
		imageList, ok := imageListResult.([]registry.CustomImage)
		if !ok {
			return fmt.Errorf("%w []registry.CustomImage but %T", ErrResultIsNot, imageListResult)
		}
		if len(imageList) != 1 {
			return fmt.Errorf(
				"%w %q, got %d",
				ErrExpectedExactly1ImageWithName,
				imageName,
				len(imageList),
			)
		}
		image := imageList[0]

		// List tags for the image
		tagName := dockerImageNameSplit[1]
		tagListResult := core.ExecBeforeCmdWithResult(ctx, fmt.Sprintf(
			"scw registry tag list image-id=%s name=%s",
			image.ID,
			tagName,
		))

		// Select the tag
		tagList, ok := tagListResult.([]registry.CustomTag)
		if !ok {
			return fmt.Errorf("%w []registry.CustomTag but %T", ErrResultIsNot, tagListResult)
		}
		if len(tagList) != 1 {
			return fmt.Errorf(
				"%w %q, got %d",
				ErrExpectedExactly1TagWithName,
				tagName,
				len(tagList),
			)
		}
		tag := tagList[0]

		ctx.Meta[tagMetaKey] = tag.ID

		return nil
	}
}
