package k8s

import (
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/ghodss/yaml"
	api "github.com/kubernetes-client/go-base/config/api"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
)

// testIfKubeconfigNotInFile checks if the given kubeconfig is not in the given file
// it tests if the user, cluster and context of the kubeconfig file are not in the given file
func testIfKubeconfigNotInFile(t *testing.T, filePath string, suffix string, kubeconfig api.Config) {
	kubeconfigBytes, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	var existingKubeconfig k8s.Kubeconfig
	err = yaml.Unmarshal(kubeconfigBytes, &existingKubeconfig)
	assert.Nil(t, err)

	found := false
	for _, cluster := range existingKubeconfig.Clusters {
		if cluster.Name == kubeconfig.Clusters[0].Name+suffix {
			found = true
			break
		}
	}
	assert.False(t, found, "cluster found in kubeconfig for cluster with suffix %s", suffix)

	found = false
	for _, context := range existingKubeconfig.Contexts {
		if context.Name == kubeconfig.Contexts[0].Name+suffix {
			found = true
			break
		}
	}
	assert.False(t, found, "context found in kubeconfig for cluster with suffix %s", suffix)

	found = false
	for _, user := range existingKubeconfig.Users {
		if user.Name == kubeconfig.AuthInfos[0].Name+suffix {
			found = true
			break
		}
	}
	assert.False(t, found, "user found in kubeconfig with suffix %s", suffix)
}

func Test_UninstallKubeconfig(t *testing.T) {
	////
	// Simple use cases
	////
	t.Run("uninstall", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createClusterAndWaitAndInstallKubeconfig("uninstall-kubeconfig", "Cluster", "Kubeconfig", kapsuleVersion),
		Cmd:        "scw k8s kubeconfig uninstall {{ .Cluster.ID }}",
		Check: core.TestCheckCombine(
			// no golden tests since it's os specific
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				testIfKubeconfigNotInFile(t, path.Join(os.TempDir(), "cli-uninstall-test"), "-"+ctx.Meta["Cluster"].(*k8s.Cluster).ID, ctx.Meta["Kubeconfig"].(api.Config))
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
		OverrideEnv: map[string]string{
			"KUBECONFIG": path.Join(os.TempDir(), "cli-uninstall-test"),
		},
	}))
	t.Run("empty file", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createClusterAndWaitAndKubeconfig("uninstall-kubeconfig-empty", "EmptyCluster", "Kubeconfig", kapsuleVersion),
		Cmd:        "scw k8s kubeconfig uninstall {{ .EmptyCluster.ID }}",
		Check: core.TestCheckCombine(
			// no golden tests since it's os specific
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				_, err := os.Stat(path.Join(os.TempDir(), "emptyfile"))
				assert.True(t, os.IsNotExist(err))
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("EmptyCluster"),
		OverrideEnv: map[string]string{
			"KUBECONFIG": path.Join(os.TempDir(), "emptyfile"),
		},
	}))
	t.Run("uninstall-merge", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createClusterAndWaitAndKubeconfigAndPopulateFileAndInstall("uninstall-kubeconfig-merge", "Cluster", "Kubeconfig", kapsuleVersion, path.Join(os.TempDir(), "cli-uninstall-merge-test"), []byte(existingKubeconfig)),
		Cmd:        "scw k8s kubeconfig uninstall {{ .Cluster.ID }}",
		Check: core.TestCheckCombine(
			// no golden tests since it's os specific
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				testIfKubeconfigNotInFile(t, path.Join(os.TempDir(), "cli-uninstall-merge-test"), "-"+ctx.Meta["Cluster"].(*k8s.Cluster).ID, ctx.Meta["Kubeconfig"].(api.Config))
				testIfKubeconfigInFile(t, path.Join(os.TempDir(), "cli-uninstall-merge-test"), "", testKubeconfig)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
		OverrideEnv: map[string]string{
			"KUBECONFIG": path.Join(os.TempDir(), "cli-uninstall-merge-test"),
		},
	}))
}
