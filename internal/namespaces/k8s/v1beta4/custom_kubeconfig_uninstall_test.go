package k8s

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
	"gopkg.in/yaml.v2"
)

func testIfKubeconfigNotInFile(t *testing.T, filePath string, suffix string, kubeconfig *k8s.Kubeconfig) {
	kubeconfigBytes, err := ioutil.ReadFile(filePath)
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
		if user.Name == kubeconfig.Users[0].Name+suffix {
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
		Cmd:        "scw k8s kubeconfig uninstall {{ .Cluster.ID }}",
		BeforeFunc: createClusterAndWaitAndInstallKubeconfig("Cluster", "Kubeconfig", kapsuleVersion),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				testIfKubeconfigNotInFile(t, "/tmp/cli-uninstall-test", "-"+ctx.Meta["Cluster"].(*k8s.Cluster).ID, ctx.Meta["Kubeconfig"].(*k8s.Kubeconfig))
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
		OverrideEnv: map[string]string{
			"KUBECONFIG": "/tmp/cli-uninstall-test",
		},
	}))
	t.Run("empty file", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		Cmd:        "scw k8s kubeconfig uninstall {{ .EmptyCluster.ID }}",
		BeforeFunc: createClusterAndWaitAndKubeconfig("EmptyCluster", "Kubeconfig", kapsuleVersion),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				_, err := os.Stat("/tmp/emptyfile")
				assert.True(t, os.IsNotExist(err))
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("EmptyCluster"),
		OverrideEnv: map[string]string{
			"KUBECONFIG": "/tmp/emptyfile",
		},
	}))
	t.Run("uninstall-merge", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		Cmd:        "scw k8s kubeconfig uninstall {{ .Cluster.ID }}",
		BeforeFunc: createClusterAndWaitAndKubeconfigAndPopulateFileAndInstall("Cluster", "Kubeconfig", kapsuleVersion, "/tmp/cli-uninstall-merge-test", []byte(existingKubeconfigs)),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				testIfKubeconfigNotInFile(t, "/tmp/cli-uninstall-merge-test", "-"+ctx.Meta["Cluster"].(*k8s.Cluster).ID, ctx.Meta["Kubeconfig"].(*k8s.Kubeconfig))
				testIfKubeconfigInFile(t, "/tmp/cli-uninstall-merge-test", "", testKubeconfig)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
		OverrideEnv: map[string]string{
			"KUBECONFIG": "/tmp/cli-uninstall-merge-test",
		},
	}))
}
