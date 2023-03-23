package container

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	registrycmds "github.com/scaleway/scaleway-cli/v2/internal/namespaces/registry/v1"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	registry "github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	indexHTML = strings.TrimSpace(`
<!DOCTYPE html>
<html>
<head>
<title>My container</title>
</head>
<body>
<h1>Deployed with scw container deploy</h1>
</body>
</html>
	`)
	nginxDockerfile = strings.TrimSpace(`
FROM nginx:alpine
RUN apk add --no-cache curl git bash
COPY ./index.html /usr/share/nginx/html/index.html
EXPOSE 80
	`)
)

func Test_Deploy(t *testing.T) {
	appName := "cli-test-container-deploy"
	path := t.TempDir()

	commands := GetCommands()
	commands.Merge(registrycmds.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: commands,
		BeforeFunc: core.BeforeFuncCombine(
			func(ctx *core.BeforeFuncCtx) error {
				// Create index.html
				err := os.WriteFile(filepath.Join(path, "index.html"), []byte(indexHTML), 0600)
				if err != nil {
					return err
				}
				return nil
			},
			func(ctx *core.BeforeFuncCtx) error {
				// Create Dockerfile
				err := os.WriteFile(filepath.Join(path, "Dockerfile"), []byte(nginxDockerfile), 0600)
				if err != nil {
					return err
				}
				return nil
			},
		),
		Cmd: fmt.Sprintf("scw container deploy name=%s build-source=%s port=80", appName, path),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			testDeleteContainersNamespaceAfter(appName),
			testDeleteRegistryAfter(appName),
		),
		DisableParallel: true,
	}))

	t.Run("App name deduced from path", core.Test(&core.TestConfig{
		Commands: commands,
		BeforeFunc: core.BeforeFuncCombine(
			func(ctx *core.BeforeFuncCtx) error {
				// Create directory
				err := os.Mkdir(filepath.Join(path, "cli-test-deploy-poney"), 0700)
				if err != nil {
					return err
				}
				return nil
			},
			func(ctx *core.BeforeFuncCtx) error {
				// Create index.html
				err := os.WriteFile(filepath.Join(path, "cli-test-deploy-poney", "index.html"), []byte(indexHTML), 0600)
				if err != nil {
					return err
				}
				return nil
			},
			func(ctx *core.BeforeFuncCtx) error {
				// Create Dockerfile
				err := os.WriteFile(filepath.Join(path, "cli-test-deploy-poney", "Dockerfile"), []byte(nginxDockerfile), 0600)
				if err != nil {
					return err
				}
				return nil
			},
		),
		Cmd: fmt.Sprintf("scw container deploy build-source=%s/cli-test-deploy-poney port=80", path),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			testDeleteContainersNamespaceAfter("app-cli-test-deploy-poney"),
			testDeleteRegistryAfter("app-cli-test-deploy-poney"),
		),
		DisableParallel: true,
	}))
}

func testDeleteContainersNamespaceAfter(appName string) func(*core.AfterFuncCtx) error {
	return func(ctx *core.AfterFuncCtx) error {
		api := container.NewAPI(ctx.Client)

		namespaces, err := api.ListNamespaces(&container.ListNamespacesRequest{
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
			return fmt.Errorf("namespace not found")
		}

		return core.ExecAfterCmd(fmt.Sprintf("scw container namespace delete %s", namespaceID))(ctx)
	}
}

func testDeleteRegistryAfter(appName string) func(*core.AfterFuncCtx) error {
	return func(ctx *core.AfterFuncCtx) error {
		api := registry.NewAPI(ctx.Client)

		registries, err := api.ListNamespaces(&registry.ListNamespacesRequest{
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

		return core.ExecAfterCmd(fmt.Sprintf("scw registry namespace delete %s", registryID))(ctx)
	}
}
