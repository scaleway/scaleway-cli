package k8s

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
	"gopkg.in/yaml.v2"
)

var (
	existingKubeconfigs = `apiVersion: v1
kind: Config
current-context: test@test
clusters:
- name: test
  cluster:
    server: https://test:6443
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJd01ETXdOREV4TVRVd01Gb1hEVE13TURNd05ERXhNVFV3TUZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBS3UvCkxIWXlNUTdFWmRaWmFaRUFsMkJaZDJLVmhPNjR0cGdJN0NUcG56NGRMUzZpdmRQSVJMbGdGdzBBV3ByNkFQeGsKeWp4K09mMkxxUVR4TkdDTFFaOEN2LzF4K3ZSdlM5YXJ3VnVIZ1pML2ZkTTE5Z2oyRWMvZUY3WDhLMEVadVNYMQpHb1ZickFOSDcyb2NkeHdvajZZc29nZHlSRHpXbkRjN2pvVjlsZlZjMlF1UDdHK2FGOEZBSTYrMUZpTDlidFE3Cnl3ZmFDWnE4M1dKcTFpOStqKzRsRUFDaFA1alZ2LzBLenVOMnRud3VaV3ljM3pRbHJHSG1xbGQ4WCtjVEJrMGgKN3hYSTV3U3VPTzMvVzBpcFQ5QVdWZjBoWVpnNVY5VWlKTjVBREZXem15eFQ4bTlVZXR3bTNhVlJtZ3IrNU44YQpGNkdWTWd2TEFteVo1bUpoR0RrQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFFL1JMZU9vbXM3TlF4OUgvNjJIcm9VUTJ5TFcKMXlpaVBVaFBEUUJEZ1p5T251TWN2Sk84SFQrcFZ3OWVPazUxeTNUT01LYmJwQWhxc0I3dEpLQVRWcWJ5Z2tHNgpzWUQrS3kwckNITzdEMTM5RXQ1eEI2dGtDQTJxSlVRMkFFb0kyNy9na3Z3K05WektUU3VSVUF0QXZXQ2diNnVZCmRUa2NpaGNSemhyV1dTV05RN3M2RSt0Y3NJS01ibUJVT0NsMVh6MERPM3RsbzlPLzRMUkE4b3RXUDEweGVLcmsKVWxRVUsvSXFpL2FTN1c4Nm4vSWU1dFVvNTJRZk5NMk0xZkVtUVA2eThYMUZWbDh2cC9WcXgxMWszeS9ZSEcvVgp5N0pIY1hhdVJLSTdieGdPdWdTZXAyZDBsaFl4Smw2NUNQT0NsbGF3Y3U3MERzMzRNS2kzWGtDZTIwST0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
contexts:
- name: test@test
  context:
    cluster: test
    user: test
users:
- name: test
  user:
    token: qotGxuOfD74ajgWir18tMMPicxLIizKg3nt5PKHGbsbNDbGfqNojIdXI
    username: test
    password: test`

	testKubeconfig = &k8s.Kubeconfig{
		APIVersion: "v1",
		Kind:       "config",
		Clusters: []*k8s.KubeconfigClusterWithName{
			{
				Name: "test",
				Cluster: k8s.KubeconfigCluster{
					CertificateAuthorityData: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJd01ETXdOREV4TVRVd01Gb1hEVE13TURNd05ERXhNVFV3TUZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBS3UvCkxIWXlNUTdFWmRaWmFaRUFsMkJaZDJLVmhPNjR0cGdJN0NUcG56NGRMUzZpdmRQSVJMbGdGdzBBV3ByNkFQeGsKeWp4K09mMkxxUVR4TkdDTFFaOEN2LzF4K3ZSdlM5YXJ3VnVIZ1pML2ZkTTE5Z2oyRWMvZUY3WDhLMEVadVNYMQpHb1ZickFOSDcyb2NkeHdvajZZc29nZHlSRHpXbkRjN2pvVjlsZlZjMlF1UDdHK2FGOEZBSTYrMUZpTDlidFE3Cnl3ZmFDWnE4M1dKcTFpOStqKzRsRUFDaFA1alZ2LzBLenVOMnRud3VaV3ljM3pRbHJHSG1xbGQ4WCtjVEJrMGgKN3hYSTV3U3VPTzMvVzBpcFQ5QVdWZjBoWVpnNVY5VWlKTjVBREZXem15eFQ4bTlVZXR3bTNhVlJtZ3IrNU44YQpGNkdWTWd2TEFteVo1bUpoR0RrQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFFL1JMZU9vbXM3TlF4OUgvNjJIcm9VUTJ5TFcKMXlpaVBVaFBEUUJEZ1p5T251TWN2Sk84SFQrcFZ3OWVPazUxeTNUT01LYmJwQWhxc0I3dEpLQVRWcWJ5Z2tHNgpzWUQrS3kwckNITzdEMTM5RXQ1eEI2dGtDQTJxSlVRMkFFb0kyNy9na3Z3K05WektUU3VSVUF0QXZXQ2diNnVZCmRUa2NpaGNSemhyV1dTV05RN3M2RSt0Y3NJS01ibUJVT0NsMVh6MERPM3RsbzlPLzRMUkE4b3RXUDEweGVLcmsKVWxRVUsvSXFpL2FTN1c4Nm4vSWU1dFVvNTJRZk5NMk0xZkVtUVA2eThYMUZWbDh2cC9WcXgxMWszeS9ZSEcvVgp5N0pIY1hhdVJLSTdieGdPdWdTZXAyZDBsaFl4Smw2NUNQT0NsbGF3Y3U3MERzMzRNS2kzWGtDZTIwST0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
					Server:                   "https://test:6443",
				},
			},
		},
		Contexts: []*k8s.KubeconfigContextWithName{
			{
				Name: "test@test",
				Context: k8s.KubeconfigContext{
					Cluster: "test",
					User:    "test",
				},
			},
		},
		Users: []*k8s.KubeconfigUserWithName{
			{
				Name: "test",
				User: k8s.KubeconfigUser{
					Token:    "qotGxuOfD74ajgWir18tMMPicxLIizKg3nt5PKHGbsbNDbGfqNojIdXI",
					Username: "test",
					Password: "test",
				},
			},
		},
	}
)

func testIfKubeconfigInFile(t *testing.T, filePath string, suffix string, kubeconfig *k8s.Kubeconfig) {
	kubeconfigBytes, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)
	var existingKubeconfig k8s.Kubeconfig
	err = yaml.Unmarshal(kubeconfigBytes, &existingKubeconfig)
	assert.Nil(t, err)

	found := false
	for _, cluster := range existingKubeconfig.Clusters {
		if cluster.Name == kubeconfig.Clusters[0].Name+suffix {
			assert.Equal(t, kubeconfig.Clusters[0].Cluster, cluster.Cluster)
			found = true
			break
		}
	}
	assert.True(t, found, "cluster not found in kubeconfig for cluster with suffix %s", suffix)

	found = false
	for _, context := range existingKubeconfig.Contexts {
		if context.Name == kubeconfig.Contexts[0].Name+suffix {
			assert.Equal(t, kubeconfig.Contexts[0].Context.Cluster+suffix, context.Context.Cluster)
			assert.Equal(t, kubeconfig.Contexts[0].Context.User+suffix, context.Context.User)
			assert.Equal(t, kubeconfig.Contexts[0].Context.Namespace, context.Context.Namespace)
			found = true
			break
		}
	}

	assert.True(t, found, "context not found in kubeconfig for cluster with suffix %s", suffix)

	found = false
	for _, user := range existingKubeconfig.Users {
		if user.Name == kubeconfig.Users[0].Name+suffix {
			assert.Equal(t, kubeconfig.Users[0].User, user.User)
			found = true
			break
		}
	}

	assert.True(t, found, "user not found in kubeconfig with suffix %s", suffix)
}

func Test_InstallKubeconfig(t *testing.T) {
	////
	// Simple use cases
	////
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		Cmd:        "scw k8s kubeconfig install {{ .Cluster.ID }}",
		BeforeFunc: createClusterAndWaitAndKubeconfig("Cluster", "Kubeconfig", kapsuleVersion),
		Check: core.TestCheckCombine(
			// no golden tests since it's os specific
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				testIfKubeconfigInFile(t, path.Join(os.TempDir(), "cli-test"), "-"+ctx.Meta["Cluster"].(*k8s.Cluster).ID, ctx.Meta["Kubeconfig"].(*k8s.Kubeconfig))
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
		OverrideEnv: map[string]string{
			"KUBECONFIG": path.Join(os.TempDir(), "cli-test"),
		},
	}))

	t.Run("merge", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		Cmd:        "scw k8s kubeconfig install {{ .Cluster.ID }}",
		BeforeFunc: createClusterAndWaitAndKubeconfigAndPopulateFile("Cluster", "Kubeconfig", kapsuleVersion, path.Join(os.TempDir(), "cli-merge-test"), []byte(existingKubeconfigs)),
		Check: core.TestCheckCombine(
			// no golden tests since it's os specific
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				testIfKubeconfigInFile(t, path.Join(os.TempDir(), "cli-merge-test"), "-"+ctx.Meta["Cluster"].(*k8s.Cluster).ID, ctx.Meta["Kubeconfig"].(*k8s.Kubeconfig))
				testIfKubeconfigInFile(t, path.Join(os.TempDir(), "cli-merge-test"), "", testKubeconfig)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
		OverrideEnv: map[string]string{
			"KUBECONFIG": path.Join(os.TempDir(), "cli-merge-test"),
		},
	}))
}
