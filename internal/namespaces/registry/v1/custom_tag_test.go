package registry_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/registry/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
)

func Test_RegistryTagDelete(t *testing.T) {
	registryNamespaceMetaKey := "RegistryNamespace"
	helloWorldImage := "hello-world:latest"
	helloWorldImageMetaKey := "HelloWorldImage"
	tagIDMetaKey := "TagID"

	t.Run("simple", core.Test(&core.TestConfig{
		Commands: registry.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				registryNamespaceMetaKey,
				fmt.Sprintf("scw registry namespace create name=%s is-public=false",
					core.GetRandomName("test-rg-tag-delete"),
				),
			),
			core.BeforeFuncWhenUpdatingCassette(
				core.BeforeFuncCombine(
					core.ExecBeforeCmd("scw registry login"),
					testhelpers.PushRegistryImage(helloWorldImage, registryNamespaceMetaKey),
				),
			),
			testhelpers.StoreImageIdentifierInMeta(
				registryNamespaceMetaKey,
				helloWorldImage,
				helloWorldImageMetaKey,
			),
			testhelpers.StoreTagIDInMeta(registryNamespaceMetaKey, helloWorldImage, tagIDMetaKey),
		),
		Cmd: fmt.Sprintf("scw registry tag delete {{ .%s }}", tagIDMetaKey),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			return core.ExecAfterCmd(
				fmt.Sprintf(
					"scw registry namespace delete {{ .%s.ID }}",
					registryNamespaceMetaKey,
				),
			)(
				ctx,
			)
		},
	}))

	t.Run("timeout-ok", core.Test(&core.TestConfig{
		Commands: registry.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				registryNamespaceMetaKey,
				fmt.Sprintf("scw registry namespace create name=%s is-public=false",
					core.GetRandomName("test-rg-tag-delete"),
				),
			),
			core.BeforeFuncWhenUpdatingCassette(
				core.BeforeFuncCombine(
					core.ExecBeforeCmd("scw registry login"),
					testhelpers.PushRegistryImage(helloWorldImage, registryNamespaceMetaKey),
				),
			),
			testhelpers.StoreImageIdentifierInMeta(
				registryNamespaceMetaKey,
				helloWorldImage,
				helloWorldImageMetaKey,
			),
			testhelpers.StoreTagIDInMeta(registryNamespaceMetaKey, helloWorldImage, tagIDMetaKey),
		),
		Cmd: fmt.Sprintf("scw registry tag delete {{ .%s }} timeout=1s", tagIDMetaKey),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			return core.ExecAfterCmd(
				fmt.Sprintf(
					"scw registry namespace delete {{ .%s.ID }}",
					registryNamespaceMetaKey,
				),
			)(
				ctx,
			)
		},
	}))
}
