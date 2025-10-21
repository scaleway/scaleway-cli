package k8s_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
	api "github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1/types"
	k8sSDK "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	existingKubeconfig = `apiVersion: v1
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
    user: test-token
users:
- name: test-token
  user:
    token: qotGxuOfD74ajgWir18tMMPicxLIizKg3nt5PKHGbsbNDbGfqNojIdXI
- name: test-cert
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvekNDQWVlZ0F3SUJBZ0lJZERQak80Umphdzh3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TURBek1qTXdPRFEyTkRoYUZ3MHlNVEF6TWpNd09EUTJORGhhTUVFeApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1TWXdKQVlEVlFRREV4MXJkV0psTFdGd2FYTmxjblpsCmNpMXJkV0psYkdWMExXTnNhV1Z1ZERDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUIKQU54VjByQ0lBemNsdXIyV1VNb3NqOW1LQmlkclYzcnB5RmNwdnltMmtFVjZaOVo2TTBSRXpyTHo1c3BaWndCTwo1bHZrbEdzL2RJVndFK2pBd2tNWWNRRWlOaTQ2bHU4UFNSei9HVTFkek5mOEF2TXpnRWZER0xUY2x3eUs4di9kCklLenhTUnVOUFFseDZoTUw1bFpDeVBBZ3hqejNEdDZGWmUxUnVUdURWTUhnOWZIaHNwOFZTYnVCbWFYTTU2T0IKLzNZQXJLMXZOTlY0enRlQ3libFZnVUd3QUdKQ09zTlE0d0l4R0xSdjN5TVhtK3V3YVpGeTFxSEh6ZlpXclRpQQpKQ2lQNFVCbDV3bnUzeEhNaFZaemI0RnNCLzBmVEl1WHQ0ZjQ5L201KzdpM01vMEdrMjJNMjAvQldzNURZVmo1CnptSVVxcU9kK09UekdkcjgvcTRsdnQ4Q0F3RUFBYU1uTUNVd0RnWURWUjBQQVFIL0JBUURBZ1dnTUJNR0ExVWQKSlFRTU1Bb0dDQ3NHQVFVRkJ3TUNNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUNuVFFCWlhsbm1aVWpDNEdscwpKdTZWWEYxN040ckpzWkNVOVR3SEhETlpZNGo5YlZHbG8wZzAva3JTajBId3hNNVU1NXl2WUJDaWNpZzVkSS96Cnd2ZENUQm5FQWIxRWtuZVR1ZkVPYzFtNzBtSzg0dnd3WWZtRVNkY1NXMmJieHBuUFNpak5BdnlTekZTTmZZZDEKMy9FZlRlQjQ0VFNGRGZQVk83YnpKYXBpYVJCTlZocVJQSncwc0lJWGM1Q0hiQzFEMHU5Mk4zRnhCa3JKcFN2UAp1QXBQT2dyNUgwUk5rOEk2TTBjd0FBc1RqdUkxd2Z4MjhJU0FWcmZLUjU4d1Eza1NsZzlUTTQrN01VMFA4eUZHClJXRkIrVFZiMTExYTRDc2RSbWMzQnZtcnFEbjZ2Ny9LOTJ4c0hNeDdBd3FObk1XUDQ4QStoVFNFVFh3U1Btb3cKL040RAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBM0ZYU3NJZ0ROeVc2dlpaUXlpeVAyWW9HSjJ0WGV1bklWeW0vS2JhUVJYcG4xbm96ClJFVE9zdlBteWxsbkFFN21XK1NVYXo5MGhYQVQ2TURDUXhoeEFTSTJManFXN3c5SkhQOFpUVjNNMS93Qzh6T0EKUjhNWXROeVhESXJ5LzkwZ3JQRkpHNDA5Q1hIcUV3dm1Wa0xJOENER1BQY08zb1ZsN1ZHNU80TlV3ZUQxOGVHeQpueFZKdTRHWnBjem5vNEgvZGdDc3JXODAxWGpPMTRMSnVWV0JRYkFBWWtJNncxRGpBakVZdEcvZkl4ZWI2N0JwCmtYTFdvY2ZOOWxhdE9JQWtLSS9oUUdYbkNlN2ZFY3lGVm5OdmdXd0gvUjlNaTVlM2gvajMrYm43dUxjeWpRYVQKYll6YlQ4RmF6a05oV1BuT1loU3FvNTM0NVBNWjJ2eityaVcrM3dJREFRQUJBb0lCQVFESDRsdldwaTAwbEZmSwpzbGpzY0d5M2p3MXlLV0VkTW9UNi9mWmNJekRTdHU4SWxhZDRvV3RhMFFWb1FKNittdFZFUENPZy85bjNTK3ZqCjFTcm1yMytrNWFKOVljMlhaaWlQMDZUaW1OdkNmTzg0TGxxTHY2UGtQOUlRSU9XOTFKOVdCVGFyZGdBUFYzWmcKZlFVaThFZFdBSVdXdlJLU01EWjlpd3dkdjFEZTZFUmt4Z0Y2R0NTSXQ2Ri80RS81Uk1VbkJObU1ycjZHWHR5NgorK0cxWExCcWxRdExYVm1yRDAyVW05Y1Yyb2QwOEczTzdUM3VqUWl6ZjR6emx4LzVWWGk0ZTFkVEViY05PRU53Cnlwd3lSajBCdFh2TXVwUXZvdUZRM3I5UmVQL0g5dmp1Q3NiZGF3T1pGQkFDb3J0UVJxcnFodENZMERRK2tiM1AKQWV1SjNnb0JBb0dCQU56QU0rbzBqckhKWGZJVFppSWEwZkV4QkIyNHBpTE5NTmoxaHIxSEJFWjI2eXFucFg5UApTSkRIbXhWREo2UXROZFQxZ3Y5L1MwelF3S2ZzYVZ5M2VYNW9OcU5hVEZERGhPSEoxWDZZUElwREZGTUgwNGV2CnRXV1ZNd21MVU9mdmhQR3NYOW5rdFRlVmxueTlnSUZOK0dkWFRTSlgyOEVIaHE2NGg4ay9IVXFCQW9HQkFQK0UKb2pVcUN3RWJ0UTIzcGRwYjNGRnNwdXZWM0F1aENiNnNmcUlxQ1ovVDRlUXJSSWtPU2luYmlva042ZFR5MVhuNQp6cGlJTEhOQ20wYkl2cVpJZmhkdERsUlcwcGQxbmlGZ2R3c3FacjdFUlFlN29XSHZkbVRwa1NaQ2p1M04zb1NjCjRPSmUwVmxBdWdwMjRsbms2bisxb01ySjJRUjNqQkxQVGUvZ3dKbGZBb0dCQUtBcUhBQ3J6WFNWQThLbDdJNkcKSXhqNlZXQXpIdWRWTlVIVk1zT1dDVFlQQmlWV3FhOHJHUjFpbGRUaGVwdVY2ZDd2bXZKQnE2SzZPMjRiQzM4bgo1OUNkVURkSlJ1RzZXbWx3QmFUcVU5S0ZSUFBSVTlxNDA4WTJjR2RXVzRkTXM0cWRaSlkxYUg1QjNJUDVBb25PCnhwSkVOMFRadGluaGlnaXUvbVkza3NzQkFvR0JBUFZDb0ZnYmhQaUpXZDVTMnRXZnV2aEZMR3ZPbVNwb1p1d28Kc2x5QnNUOUNwOTdWVVRHbEQ3Ymh6allEcnVFQ1BicVk5NThkaGwwVUgrdHZvT0FIVVZDM0V6d05JcExUQ1BmTQptamVUZVkrKzRPdXRSQmkzTzVOZFJqL05QMWd2ZFZraEpCTGxKRmxoY2JHOXIwTE9JZkIzckdFNkloN1JpUmc4CjkvZzZhV1JOQW9HQWJSSG94d1B4MUVlRnRxVDlUdXowZWZUR3RwQTB0bkhDZTN2b2x1L201eEZ6N3BwS29HeUQKRkNPVm5jMmZ3LzQwYUFGTGdHYlFLMFBqTzFCbWZ3cjFvb09aT1hZYnExUXo3Q1cvN3A1OUFkR0VrWXFzdWZZcAp6OXlMd1dBUEdybm9jVjBVQXZ2SHcvbC9OK29NZEdpdmVTdDhRb3RHclgzdm9PTmVsWThCZDRNPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=`
)

func Test_InstallKubeconfig(t *testing.T) {
	////
	// Simple use case
	// no flags, no env, using default path
	////
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   k8s.GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: core.BeforeFuncCombine(
			createCluster("install-kubeconfig-simple", true),
			fetchClusterKubeconfigMetadata(true),
		),
		Cmd: "scw k8s kubeconfig install {{ ." + clusterMetaKey + ".ID }}",
		Check: core.TestCheckCombine(
			// no golden tests since it's os specific
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				kubeconfigPath := path.Join(ctx.Meta["HOME"].(string), ".kube", "config")
				t.Logf("using kubeconfigPath: %s", kubeconfigPath)

				fileKubeconfig, err := os.ReadFile(kubeconfigPath)
				require.NoError(t, err)

				var finalKubeconfig api.Config
				err = yaml.Unmarshal(fileKubeconfig, &finalKubeconfig)
				require.NoError(t, err)

				suffix := "-" + ctx.Meta[clusterMetaKey].(*k8sSDK.Cluster).ID
				assertKubeconfigClusterHasSuffix(t, finalKubeconfig, suffix)
				assertKubeconfigContextHasSuffix(t, finalKubeconfig, suffix)

				authInfo := api.NamedAuthInfo{
					Name: "cli-config-00000000",
					AuthInfo: api.AuthInfo{
						Exec: &api.ExecConfig{
							APIVersion:      "client.authentication.k8s.io/v1",
							Command:         "scw",
							Args:            []string{"k8s", "exec-credential"},
							InstallHint:     k8s.InstallHint,
							InteractiveMode: "Never",
						},
					},
				}
				assertKubeconfigAuthUserContains(t, finalKubeconfig, authInfo)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
		OverrideEnv: map[string]string{
			scw.ScwAccessKeyEnv: "", // Hide keys in test env
			scw.ScwSecretKeyEnv: "", // Hide keys in test env
		},
	}))

	////
	// Merge will use install on top of an existing kubeconfig
	// no envs, no flags, using default path
	////
	t.Run("merge", core.Test(&core.TestConfig{
		Commands:   k8s.GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: core.BeforeFuncCombine(
			writeKubeconfigFile([]byte(existingKubeconfig)),
			createCluster("install-kubeconfig-merge", true),
			fetchClusterKubeconfigMetadata(true),
		),
		Cmd: "scw k8s kubeconfig install {{ ." + clusterMetaKey + ".ID }}",
		Check: core.TestCheckCombine(
			// no golden tests since it's os specific
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				kubeconfigPath := path.Join(ctx.Meta["HOME"].(string), ".kube", "config")
				t.Logf("using kubeconfigPath: %s", kubeconfigPath)

				fileKubeconfig, err := os.ReadFile(kubeconfigPath)
				require.NoError(t, err)

				var finalKubeconfig api.Config
				err = yaml.Unmarshal(fileKubeconfig, &finalKubeconfig)
				require.NoError(t, err)

				assert.NotEqualf(t,
					[]byte(existingKubeconfig),
					fileKubeconfig,
					"expected kubeconfig file to be merged",
				)

				suffix := "-" + ctx.Meta[clusterMetaKey].(*k8sSDK.Cluster).ID
				assertKubeconfigClusterHasSuffix(t, finalKubeconfig, suffix)
				assertKubeconfigContextHasSuffix(t, finalKubeconfig, suffix)

				authInfo := api.NamedAuthInfo{
					Name: "cli-config-00000000",
					AuthInfo: api.AuthInfo{
						Exec: &api.ExecConfig{
							APIVersion:      "client.authentication.k8s.io/v1",
							Command:         "scw",
							Args:            []string{"k8s", "exec-credential"},
							InstallHint:     k8s.InstallHint,
							InteractiveMode: "Never",
						},
					},
				}
				assertKubeconfigAuthUserContains(t, finalKubeconfig, authInfo)

				// check currentContext changed
				found := strings.HasSuffix(finalKubeconfig.CurrentContext, suffix)
				assert.True(t, found, "current context does not contains suffix %s", suffix)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
		OverrideEnv: map[string]string{
			scw.ScwAccessKeyEnv: "", // Hide keys in test env
			scw.ScwSecretKeyEnv: "", // Hide keys in test env
		},
	}))

	////
	// Merge and keep currentContext of existing kubeconfig
	// no envs, no flags, using default path
	////
	t.Run("merge-keep", core.Test(&core.TestConfig{
		Commands:   k8s.GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: core.BeforeFuncCombine(
			writeKubeconfigFile([]byte(existingKubeconfig)),
			createCluster("install-kubeconfig-merge-keep", true),
			fetchClusterKubeconfigMetadata(true),
		),
		Cmd: "scw k8s kubeconfig install {{ ." + clusterMetaKey + ".ID }} keep-current-context=true",
		Check: core.TestCheckCombine(
			// no golden tests since it's os specific
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				kubeconfigPath := path.Join(ctx.Meta["HOME"].(string), ".kube", "config")
				t.Logf("using kubeconfigPath: %s", kubeconfigPath)

				fileKubeconfig, err := os.ReadFile(kubeconfigPath)
				require.NoError(t, err)

				var finalKubeconfig api.Config
				err = yaml.Unmarshal(fileKubeconfig, &finalKubeconfig)
				require.NoError(t, err)

				assert.NotEqualf(
					t,
					[]byte(existingKubeconfig),
					fileKubeconfig,
					"expected kubeconfig file to be merged",
				)

				suffix := "-" + ctx.Meta[clusterMetaKey].(*k8sSDK.Cluster).ID
				assertKubeconfigClusterHasSuffix(t, finalKubeconfig, suffix)
				assertKubeconfigContextHasSuffix(t, finalKubeconfig, suffix)

				authInfo := api.NamedAuthInfo{
					Name: "cli-config-00000000",
					AuthInfo: api.AuthInfo{
						Exec: &api.ExecConfig{
							APIVersion:      "client.authentication.k8s.io/v1",
							Command:         "scw",
							Args:            []string{"k8s", "exec-credential"},
							InstallHint:     k8s.InstallHint,
							InteractiveMode: "Never",
						},
					},
				}
				assertKubeconfigAuthUserContains(t, finalKubeconfig, authInfo)

				// check keep-current-context=true
				assert.Equal(t, "test@test", finalKubeconfig.CurrentContext)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
		OverrideEnv: map[string]string{
			scw.ScwAccessKeyEnv: "", // Hide keys in test env
			scw.ScwSecretKeyEnv: "", // Hide keys in test env
		},
	}))

	////
	// Merge will use install on top of an existing kubeconfig
	// no flags, using KUBECONFIG env
	////
	t.Run("merge-kubeconfig", func(t *testing.T) {
		f, err := os.CreateTemp(t.TempDir(), "kubeconfig")
		require.NoError(t, err)
		require.NotNil(t, f)
		defer os.Remove(f.Name()) // clean up

		_, err = f.WriteString(existingKubeconfig)
		require.NoError(t, err)

		err = f.Close()
		require.NoError(t, err)

		testConfig := &core.TestConfig{
			Commands:   k8s.GetCommands(),
			TmpHomeDir: true,
			BeforeFunc: core.BeforeFuncCombine(
				writeKubeconfigFile([]byte(existingKubeconfig)),
				createCluster("install-kubeconfig-merge-kubeconfig", true),
				fetchClusterKubeconfigMetadata(true),
			),
			Cmd: "scw k8s kubeconfig install {{ ." + clusterMetaKey + ".ID }}",
			Check: core.TestCheckCombine(
				// no golden tests since it's os specific
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()

					kubeconfigPath := ctx.OverrideEnv["KUBECONFIG"]
					t.Logf("using kubeconfigPath: %s", kubeconfigPath)

					fileKubeconfig, err := os.ReadFile(kubeconfigPath)
					require.NoError(t, err)

					var finalKubeconfig api.Config
					err = yaml.Unmarshal(fileKubeconfig, &finalKubeconfig)
					require.NoError(t, err)

					assert.NotEqualf(
						t,
						[]byte(existingKubeconfig),
						fileKubeconfig,
						"expected kubeconfig file to be merged",
					)

					suffix := "-" + ctx.Meta[clusterMetaKey].(*k8sSDK.Cluster).ID
					assertKubeconfigClusterHasSuffix(t, finalKubeconfig, suffix)
					assertKubeconfigContextHasSuffix(t, finalKubeconfig, suffix)

					authInfo := api.NamedAuthInfo{
						Name: "cli-config-00000000",
						AuthInfo: api.AuthInfo{
							Exec: &api.ExecConfig{
								APIVersion:      "client.authentication.k8s.io/v1",
								Command:         "scw",
								Args:            []string{"k8s", "exec-credential"},
								InstallHint:     k8s.InstallHint,
								InteractiveMode: "Never",
							},
						},
					}
					assertKubeconfigAuthUserContains(t, finalKubeconfig, authInfo)

					// check currentContext changed
					found := strings.HasSuffix(finalKubeconfig.CurrentContext, suffix)
					assert.True(t, found, "current context does not contains suffix %s", suffix)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteCluster(),
			OverrideEnv: map[string]string{
				"KUBECONFIG":        f.Name(),
				scw.ScwAccessKeyEnv: "", // Hide keys in test env
				scw.ScwSecretKeyEnv: "", // Hide keys in test env
			},
		}
		core.Test(testConfig)(t)
	})
}

func assertKubeconfigClusterHasSuffix(
	t *testing.T,
	kubeconfig api.Config,
	suffix string,
) {
	t.Helper()

	found := 0
	for _, cluster := range kubeconfig.Clusters {
		if strings.HasSuffix(cluster.Name, suffix) {
			found++
		}
	}
	assert.Equal(
		t,
		1,
		found,
		"%d clusters with suffix %s found in kubeconfig",
		found,
		suffix,
	)
	if found == 0 {
		config, err := yaml.Marshal(kubeconfig)
		require.NoError(t, err)
		t.Logf("kubeconfig: %s", config)
	}
}

func assertKubeconfigContextHasSuffix(
	t *testing.T,
	kubeconfig api.Config,
	suffix string,
) {
	t.Helper()

	found := false
	for _, context := range kubeconfig.Contexts {
		if strings.HasSuffix(context.Name, suffix) {
			assert.True(
				t,
				strings.HasSuffix(context.Context.Cluster, suffix),
				"context cluster does not match cluster pattern",
			)
			found = true

			break
		}
	}
	assert.True(t, found, "no context not found in kubeconfig for cluster with suffix %s", suffix)
	if !found {
		config, err := yaml.Marshal(kubeconfig)
		require.NoError(t, err)
		t.Logf("kubeconfig: %s", config)
	}
}

func assertKubeconfigAuthUserContains(
	t *testing.T,
	kubeconfig api.Config,
	authUser api.NamedAuthInfo,
) {
	t.Helper()

	found := false
	for _, user := range kubeconfig.AuthInfos {
		if user.Name == authUser.Name {
			assert.Equal(t, user.AuthInfo, authUser.AuthInfo)
			found = true

			break
		}
	}
	assert.True(t, found, "user not found in kubeconfig with name %s", authUser.Name)
}
