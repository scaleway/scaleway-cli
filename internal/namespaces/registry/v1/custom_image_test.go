package registry

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
)

func Test_ImageList(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
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
						fmt.Sprintf("visibility_%s", registry.ImageVisibilityPublic),
						registry.ImageVisibilityPublic,
					),

					setupImage(
						"busybox:1.30",
						"rg.fr-par.scw.cloud/cli-public-namespace",
						fmt.Sprintf("visibility_%s", registry.ImageVisibilityPrivate),
						registry.ImageVisibilityPrivate,
					),

					setupImage(
						"busybox:1.29",
						"rg.fr-par.scw.cloud/cli-public-namespace",
						fmt.Sprintf("visibility_%s", registry.ImageVisibilityInherit),
						registry.ImageVisibilityInherit,
					),

					setupImage(
						"busybox:1.28",
						"rg.fr-par.scw.cloud/cli-private-namespace",
						fmt.Sprintf("visibility_%s", registry.ImageVisibilityPublic),
						registry.ImageVisibilityPublic,
					),

					setupImage(
						"busybox:1.27",
						"rg.fr-par.scw.cloud/cli-private-namespace",
						fmt.Sprintf("visibility_%s", registry.ImageVisibilityPrivate),
						registry.ImageVisibilityPrivate,
					),

					// namespace_policy: private, image_policy:inherit
					setupImage(
						"busybox:1.26",
						"rg.fr-par.scw.cloud/cli-private-namespace",
						fmt.Sprintf("visibility_%s", registry.ImageVisibilityInherit),
						registry.ImageVisibilityInherit,
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

func setupImage(dockerImage string, namespaceEndpoint string, imageName string, visibility registry.ImageVisibility) core.BeforeFunc {
	remote := fmt.Sprintf("%s/%s:latest", namespaceEndpoint, imageName)
	return core.BeforeFuncCombine(
		core.BeforeFuncOsExec("docker", "pull", dockerImage),
		core.BeforeFuncOsExec("docker", "tag", dockerImage, remote),
		core.BeforeFuncOsExec("docker", "push", remote),
		core.ExecStoreBeforeCmd("ImageListResult", fmt.Sprintf("scw registry image list name=%s", imageName)),
		core.ExecBeforeCmd(fmt.Sprintf(`scw registry image update {{ (index .ImageListResult 0).ID }} visibility=%s`, visibility.String())),
	)
}
