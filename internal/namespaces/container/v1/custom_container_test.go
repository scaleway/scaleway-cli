package container_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	container "github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/registry/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	containerSDK "github.com/scaleway/scaleway-sdk-go/api/container/v1"
	"github.com/stretchr/testify/assert"
)

func createNamespace(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		metaKey,
		fmt.Sprintf("scw container namespace create name=%s -w", core.GetRandomName("ns")),
	)
}

func deleteNamespace(metaKey string) core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		return core.ExecAfterCmd("scw container namespace delete {{ ." + metaKey + ".ID }}")(ctx)
	}
}

func Test_CreateContainer(t *testing.T) {
	commands := container.GetCommands()

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: commands,
		BeforeFunc: core.BeforeFuncCombine(
			createNamespace("Namespace"),
		),
		Cmd: fmt.Sprintf(
			"scw container container create namespace-id={{ .Namespace.ID }} name=%s image=hello-world:latest",
			core.GetRandomName("test"),
		),
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				c := ctx.Result.(*containerSDK.Container)
				assert.Equal(t, containerSDK.ContainerStatusCreating, c.Status)
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteNamespace("Namespace"),
		),
		DisableParallel: true,
	}))
}

func Test_UpdateContainer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: container.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createNamespace("Namespace"),
			core.ExecStoreBeforeCmd("Container", fmt.Sprintf(
				"scw container container create namespace-id={{ .Namespace.ID }} name=%s image=nginx:latest port=80 -w",
				core.GetRandomName("test"),
			)),
		),
		Cmd: "scw container container update {{ .Container.ID }} tags.0=new_tag port=8080 mvcpu-limit=1500",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				if c, ok := ctx.Result.(*containerSDK.Container); ok {
					assert.Equal(t, []string{"new_tag"}, c.Tags)
					assert.Equal(t, uint32(8080), c.Port, "unexpected port number")
					assert.Equal(t, uint32(1500), c.MvcpuLimit, "unexpected CPU limit")
				} else {
					t.Fail()
				}
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteNamespace("Namespace"),
		),
	}))

	lighttpdImage := "sebp/lighttpd:latest"
	nginxImage := "nginx:1.29.2-alpine"
	lighttpdImageMetaKey := "LighttpdImage"
	nginxImageMetaKey := "NginxImage"
	registryNamespaceMetaKey := "RegistryNamespace"
	containerNamespaceMetaKey := "ContainerNamespace"
	containerMetaKey := "Container"

	t.Run("RegistryImage", core.Test(&core.TestConfig{
		Commands: core.NewCommandsMerge(
			container.GetCommands(),
			registry.GetCommands(),
		),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				registryNamespaceMetaKey,
				fmt.Sprintf("scw registry namespace create name=%s is-public=false",
					core.GetRandomName("test-ctn-update-rg-img"),
				),
			),
			core.BeforeFuncWhenUpdatingCassette(
				core.BeforeFuncCombine(
					core.ExecBeforeCmd("scw registry login"),
					testhelpers.PushRegistryImage(lighttpdImage, registryNamespaceMetaKey),
					testhelpers.PushRegistryImage(nginxImage, registryNamespaceMetaKey),
				),
			),
			testhelpers.StoreImageIdentifierInMeta(
				registryNamespaceMetaKey,
				lighttpdImage,
				lighttpdImageMetaKey,
			),
			testhelpers.StoreImageIdentifierInMeta(
				registryNamespaceMetaKey,
				nginxImage,
				nginxImageMetaKey,
			),
			createNamespace(containerNamespaceMetaKey),
			core.ExecStoreBeforeCmd(containerMetaKey, fmt.Sprintf(
				"scw container container create namespace-id={{ .%s.ID }} name=%s image={{ .%s }} port=80 -w",
				containerNamespaceMetaKey,
				core.GetRandomName("test"),
				lighttpdImageMetaKey,
			)),
			// NB: after this step, the container with the sebp/lighttpd image will deploy but stay in error state because it has no content to serve
		),
		Cmd: fmt.Sprintf(
			"scw container container update {{ .%s.ID }} image={{ .%s }} port=80 -w",
			containerMetaKey,
			nginxImageMetaKey,
		),
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				c := ctx.Result.(*containerSDK.Container)
				// Check image
				expectedImageName := ctx.Meta.Render(fmt.Sprintf("{{ .%s }}", nginxImageMetaKey))
				assert.Equal(t, expectedImageName, c.Image)
				// Check status
				assert.Equal(t, containerSDK.ContainerStatusReady, c.Status)
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteNamespace(containerNamespaceMetaKey),
			func(ctx *core.AfterFuncCtx) error {
				return core.ExecAfterCmd(
					fmt.Sprintf(
						"scw registry namespace delete {{ .%s.ID }}",
						registryNamespaceMetaKey,
					),
				)(
					ctx,
				)
			},
		),
	}))
}
