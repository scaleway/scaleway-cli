package registry_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/registry/v1"
	registrySDK "github.com/scaleway/scaleway-sdk-go/api/registry/v1"
)

func Test_ImageList(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: registry.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("PublicNamespace", "scw registry namespace create name=cli-public-namespace is-public=true"),
			core.ExecStoreBeforeCmd("PrivateNamespace", "scw registry namespace create name=cli-private-namespace is-public=false"),
			core.BeforeFuncWhenUpdatingCassette(
				core.ExecBeforeCmd("scw registry login"),
			),
			// We need to push 6 images with different hashes, we choose small images
			core.BeforeFuncWhenUpdatingCassette(
				core.BeforeFuncCombine(
					setupImage(
						"busybox:1.31",
						"rg.fr-par.scw.cloud/cli-public-namespace",
						fmt.Sprintf("visibility_%s", registrySDK.ImageVisibilityPublic),
						registrySDK.ImageVisibilityPublic,
					),

					setupImage(
						"busybox:1.30",
						"rg.fr-par.scw.cloud/cli-public-namespace",
						fmt.Sprintf("visibility_%s", registrySDK.ImageVisibilityPrivate),
						registrySDK.ImageVisibilityPrivate,
					),

					setupImage(
						"busybox:1.29",
						"rg.fr-par.scw.cloud/cli-public-namespace",
						fmt.Sprintf("visibility_%s", registrySDK.ImageVisibilityInherit),
						registrySDK.ImageVisibilityInherit,
					),

					setupImage(
						"busybox:1.28",
						"rg.fr-par.scw.cloud/cli-private-namespace",
						fmt.Sprintf("visibility_%s", registrySDK.ImageVisibilityPublic),
						registrySDK.ImageVisibilityPublic,
					),

					setupImage(
						"busybox:1.27",
						"rg.fr-par.scw.cloud/cli-private-namespace",
						fmt.Sprintf("visibility_%s", registrySDK.ImageVisibilityPrivate),
						registrySDK.ImageVisibilityPrivate,
					),

					// namespace_policy: private, image_policy:inherit
					setupImage(
						"busybox:1.26",
						"rg.fr-par.scw.cloud/cli-private-namespace",
						fmt.Sprintf("visibility_%s", registrySDK.ImageVisibilityInherit),
						registrySDK.ImageVisibilityInherit,
					),
				),
			),
		),
		Cmd: "scw registry image list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw registry namespace delete {{ .PublicNamespace.ID }}"),
			core.ExecAfterCmd("scw registry namespace delete {{ .PrivateNamespace.ID }}"),
		),
	}))
}

func setupImage(dockerImage string, namespaceEndpoint string, imageName string, visibility registrySDK.ImageVisibility) core.BeforeFunc {
	remote := fmt.Sprintf("%s/%s:latest", namespaceEndpoint, imageName)
	return core.BeforeFuncCombine(
		core.BeforeFuncOsExec("docker", "pull", dockerImage),
		core.BeforeFuncOsExec("docker", "tag", dockerImage, remote),
		core.BeforeFuncOsExec("docker", "push", remote),
		core.ExecStoreBeforeCmd("ImageListResult", "scw registry image list name="+imageName),
		core.ExecBeforeCmd("scw registry image update {{ (index .ImageListResult 0).ID }} visibility="+visibility.String()),
	)
}
