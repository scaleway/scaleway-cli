package container_test

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	container "github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1beta1"
	registrycmds "github.com/scaleway/scaleway-cli/v2/internal/namespaces/registry/v1"
	containerSDK "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	registrySDK "github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	//go:embed testdata/docker/Dockerfile
	testdataDockerDockerfile string
	//go:embed testdata/docker/index.html
	testdataDockerIndexHTML string
	//go:embed testdata/docker/Dockerfile.build-args
	testdataDockerDockerfileBuildArgs string

	//go:embed testdata/node/index.js
	testDataBuildpackNodeIndexJS string
	//go:embed testdata/node/package.json
	testDataBuildpackNodePackageJSON string
	//go:embed testdata/node/package-lock.json
	testDataBuildpackNodePackageLockJSON string
)

func loadTestdataBeforeFunc(
	path string,
	filename string,
	content string,
) func(ctx *core.BeforeFuncCtx) error {
	return func(_ *core.BeforeFuncCtx) error {
		err := os.WriteFile(filepath.Join(path, filename), []byte(content), 0o600)
		if err != nil {
			return err
		}

		return err
	}
}

func mkdirAllBeforeFunc(path string) func(ctx *core.BeforeFuncCtx) error {
	return func(_ *core.BeforeFuncCtx) error {
		err := os.MkdirAll(path, 0o700)
		if err != nil {
			return err
		}

		return nil
	}
}

func Test_Deploy(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping test on Windows because of flakiness")
	}
	appName := "cli-t-ctnr-deploy"
	path := t.TempDir()

	commands := container.GetCommands()
	commands.Merge(registrycmds.GetCommands())

	simplePath := filepath.Join(path, "simple")
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: commands,
		BeforeFunc: core.BeforeFuncCombine(
			mkdirAllBeforeFunc(simplePath),
			loadTestdataBeforeFunc(simplePath, "index.html", testdataDockerIndexHTML),
			loadTestdataBeforeFunc(simplePath, "Dockerfile", testdataDockerDockerfile),
		),
		Cmd: fmt.Sprintf(
			"scw container deploy name=%s build-source=%s port=80",
			appName+"-s",
			simplePath,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			testDeleteContainersNamespaceAfter(appName+"-s"),
			testDeleteRegistryAfter(appName+"-s"),
		),
		DisableParallel: true,
	}))

	appNameFromPathPath := filepath.Join(path, appName+"-fp")
	t.Run("App name deduced from path", core.Test(&core.TestConfig{
		Commands: commands,
		BeforeFunc: core.BeforeFuncCombine(
			mkdirAllBeforeFunc(appNameFromPathPath),
			loadTestdataBeforeFunc(appNameFromPathPath, "index.html", testdataDockerIndexHTML),
			loadTestdataBeforeFunc(appNameFromPathPath, "Dockerfile", testdataDockerDockerfile),
		),
		Cmd: fmt.Sprintf("scw container deploy build-source=%s port=80", appNameFromPathPath),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			testDeleteContainersNamespaceAfter(appName+"-fp"),
			testDeleteRegistryAfter(appName+"-fp"),
		),
		DisableParallel: true,
	}))

	buildArgsPath := filepath.Join(path, "build-args")
	t.Run("Build args", core.Test(&core.TestConfig{
		Commands: commands,
		BeforeFunc: core.BeforeFuncCombine(
			mkdirAllBeforeFunc(buildArgsPath),
			loadTestdataBeforeFunc(buildArgsPath, "index.html", testdataDockerIndexHTML),
			loadTestdataBeforeFunc(buildArgsPath, "Dockerfile", testdataDockerDockerfileBuildArgs),
		),
		Cmd: fmt.Sprintf(
			"scw container deploy name=%s build-source=%s port=80 build-args.TEST=thisisatest",
			appName+"-ba",
			buildArgsPath,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			testDeleteContainersNamespaceAfter(appName+"-ba"),
			testDeleteRegistryAfter(appName+"-ba"),
		),
		DisableParallel: true,
	}))

	buildpackPath := filepath.Join(path, "bp")
	t.Run("Buildpack", core.Test(&core.TestConfig{
		Commands: commands,
		BeforeFunc: core.BeforeFuncCombine(
			mkdirAllBeforeFunc(buildpackPath),
			loadTestdataBeforeFunc(buildpackPath, "index.js", testDataBuildpackNodeIndexJS),
			loadTestdataBeforeFunc(buildpackPath, "package.json", testDataBuildpackNodePackageJSON),
			loadTestdataBeforeFunc(
				buildpackPath,
				"package-lock.json",
				testDataBuildpackNodePackageLockJSON,
			),
		),
		Cmd: fmt.Sprintf(
			"scw container deploy name=%s build-source=%s port=3000",
			appName+"-bp",
			buildpackPath,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			testDeleteContainersNamespaceAfter(appName+"-bp"),
			testDeleteRegistryAfter(appName+"-bp"),
		),
		DisableParallel: true,
	}))
}

func testDeleteContainersNamespaceAfter(appName string) func(*core.AfterFuncCtx) error {
	return func(ctx *core.AfterFuncCtx) error {
		api := containerSDK.NewAPI(ctx.Client)

		namespaces, err := api.ListNamespaces(&containerSDK.ListNamespacesRequest{
			Name: &appName,
		}, scw.WithAllPages())
		if err != nil {
			return err
		}

		var namespaceID string
		for _, namespace := range namespaces.Namespaces {
			if namespace.Name == appName {
				namespaceID = namespace.ID

				break
			}
		}

		if namespaceID == "" {
			return nil
		}

		return core.ExecAfterCmd(
			fmt.Sprintf("scw container namespace delete %s --wait", namespaceID),
		)(
			ctx,
		)
	}
}

func testDeleteRegistryAfter(appName string) func(*core.AfterFuncCtx) error {
	return func(ctx *core.AfterFuncCtx) error {
		api := registrySDK.NewAPI(ctx.Client)

		registries, err := api.ListNamespaces(&registrySDK.ListNamespacesRequest{
			Name: &appName,
		}, scw.WithAllPages())
		if err != nil {
			return err
		}

		var registryID string
		for _, registry := range registries.Namespaces {
			if registry.Name == appName {
				registryID = registry.ID

				break
			}
		}

		if registryID == "" {
			return nil
		}

		return core.ExecAfterCmd("scw registry namespace delete " + registryID)(ctx)
	}
}
