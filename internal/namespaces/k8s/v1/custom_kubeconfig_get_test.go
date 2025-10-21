package k8s_test

import (
	"regexp"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
	api "github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1/types"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetKubeconfig(t *testing.T) {
	////
	// Simple use case
	////
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createCluster("get-kubeconfig-simple", true),
			fetchClusterKubeconfigMetadata(true),
		),
		Cmd: "scw k8s kubeconfig get {{ ." + clusterMetaKey + ".ID }}",
		OverrideEnv: map[string]string{
			scw.ScwAccessKeyEnv: "", // Hide keys in test env
			scw.ScwSecretKeyEnv: "", // Hide keys in test env
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				var cliKubeconfig api.Config
				err := yaml.Unmarshal([]byte(ctx.Result.(string)), &cliKubeconfig)
				require.NoError(t, err)

				require.Lenf(t, cliKubeconfig.Clusters, 1,
					"expected 1 cluster, got %d", len(cliKubeconfig.Clusters))

				require.Lenf(t, cliKubeconfig.AuthInfos, 1,
					"expected 1 user, got %d", len(cliKubeconfig.AuthInfos))

				authInfos := []api.NamedAuthInfo{
					{
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
					},
				}
				assert.Equal(t, authInfos, cliKubeconfig.AuthInfos)

				assert.NotEmptyf(t, cliKubeconfig.CurrentContext,
					"expected current context not empty")

				validNamedContext := api.NamedContext{
					Name: cliKubeconfig.CurrentContext,
					Context: api.Context{
						Cluster:  cliKubeconfig.Clusters[0].Name,
						AuthInfo: cliKubeconfig.AuthInfos[0].Name,
					},
				}
				assert.Equal(t, validNamedContext, cliKubeconfig.Contexts[0])
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
	}))

	////
	// Default use case with flags set
	////
	t.Run("with_flags", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createCluster("get-kubeconfig-with-flags", true),
			fetchClusterKubeconfigMetadata(true),
		),
		Cmd: "scw k8s --profile=default kubeconfig get {{ ." + clusterMetaKey + ".ID }}",
		OverrideEnv: map[string]string{
			scw.ScwAccessKeyEnv: "", // Hide keys in test env
			scw.ScwSecretKeyEnv: "", // Hide keys in test env
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				var cliKubeconfig api.Config
				err := yaml.Unmarshal([]byte(ctx.Result.(string)), &cliKubeconfig)
				require.NoError(t, err)

				require.Lenf(t, cliKubeconfig.Clusters, 1,
					"expected 1 cluster, got %d", len(cliKubeconfig.Clusters))

				require.Lenf(t, cliKubeconfig.AuthInfos, 1,
					"expected 1 user, got %d", len(cliKubeconfig.AuthInfos))

				authInfos := []api.NamedAuthInfo{
					{
						Name: "cli-config-2e6b12a6",
						AuthInfo: api.AuthInfo{
							Exec: &api.ExecConfig{
								APIVersion: "client.authentication.k8s.io/v1",
								Command:    "scw",
								Args: []string{
									"--profile",
									"default",
									"k8s",
									"exec-credential",
								},
								InstallHint:     k8s.InstallHint,
								InteractiveMode: api.NeverExecInteractiveMode,
							},
						},
					},
				}
				assert.Equal(t, authInfos, cliKubeconfig.AuthInfos)

				assert.NotEmptyf(t, cliKubeconfig.CurrentContext,
					"expected current context not empty")

				validNamedContext := api.NamedContext{
					Name: cliKubeconfig.CurrentContext,
					Context: api.Context{
						Cluster:  cliKubeconfig.Clusters[0].Name,
						AuthInfo: cliKubeconfig.AuthInfos[0].Name,
					},
				}
				assert.Equal(t, validNamedContext, cliKubeconfig.Contexts[0])
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
	}))
	////
	// Default use case with envs presents
	////

	t.Run("with_envs", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createCluster("get-kubeconfig-with-envs", true),
			fetchClusterKubeconfigMetadata(true),
		),
		Cmd: "scw k8s kubeconfig get {{ ." + clusterMetaKey + ".ID }}",
		OverrideEnv: map[string]string{
			scw.ScwActiveProfileEnv: "default",
			scw.ScwAccessKeyEnv:     "", // Hide keys in test env
			scw.ScwSecretKeyEnv:     "", // Hide keys in test env
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				var cliKubeconfig api.Config
				err := yaml.Unmarshal([]byte(ctx.Result.(string)), &cliKubeconfig)
				require.NoError(t, err)

				require.Lenf(t, cliKubeconfig.Clusters, 1,
					"expected 1 cluster, got %d", len(cliKubeconfig.Clusters))

				require.Lenf(t, cliKubeconfig.AuthInfos, 1,
					"expected 1 user, got %d", len(cliKubeconfig.AuthInfos))

				authInfos := []api.NamedAuthInfo{
					{
						Name: "cli-config-3b1b8942",
						AuthInfo: api.AuthInfo{
							Exec: &api.ExecConfig{
								APIVersion: "client.authentication.k8s.io/v1",
								Command:    "scw",
								Args:       []string{"k8s", "exec-credential"},
								Env: []api.ExecEnvVar{
									{Name: "SCW_PROFILE", Value: "default"},
								},
								InstallHint:     k8s.InstallHint,
								InteractiveMode: api.NeverExecInteractiveMode,
							},
						},
					},
				}
				assert.Equal(t, authInfos, cliKubeconfig.AuthInfos)

				assert.NotEmptyf(t, cliKubeconfig.CurrentContext,
					"expected current context not empty")

				validNamedContext := api.NamedContext{
					Name: cliKubeconfig.CurrentContext,
					Context: api.Context{
						Cluster:  cliKubeconfig.Clusters[0].Name,
						AuthInfo: cliKubeconfig.AuthInfos[0].Name,
					},
				}
				assert.Equal(t, validNamedContext, cliKubeconfig.Contexts[0])
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
	}))

	////
	// Default use case with flags and envs set
	////
	t.Run("with_flags_and_envs", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createCluster("get-kubeconfig-with-flags-and-envs", true),
			fetchClusterKubeconfigMetadata(true),
		),
		Cmd: "scw --profile=default k8s kubeconfig get {{ ." + clusterMetaKey + ".ID }}",
		OverrideEnv: map[string]string{
			scw.ScwActiveProfileEnv: "default2",
			scw.ScwAccessKeyEnv:     "", // Hide keys in test env
			scw.ScwSecretKeyEnv:     "", // Hide keys in test env
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				var cliKubeconfig api.Config
				err := yaml.Unmarshal([]byte(ctx.Result.(string)), &cliKubeconfig)
				require.NoError(t, err)

				require.Lenf(t, cliKubeconfig.Clusters, 1,
					"expected 1 cluster, got %d", len(cliKubeconfig.Clusters))

				require.Lenf(t, cliKubeconfig.AuthInfos, 1,
					"expected 1 user, got %d", len(cliKubeconfig.AuthInfos))

				authInfos := []api.NamedAuthInfo{
					{
						Name: "cli-config-90703896",
						AuthInfo: api.AuthInfo{
							Exec: &api.ExecConfig{
								APIVersion: "client.authentication.k8s.io/v1",
								Command:    "scw",
								Args: []string{
									"--profile",
									"default",
									"k8s",
									"exec-credential",
								},
								Env: []api.ExecEnvVar{
									{Name: "SCW_PROFILE", Value: "default2"},
								},
								InstallHint:     k8s.InstallHint,
								InteractiveMode: api.NeverExecInteractiveMode,
							},
						},
					},
				}
				assert.Equal(t, authInfos, cliKubeconfig.AuthInfos)

				assert.NotEmptyf(t, cliKubeconfig.CurrentContext,
					"expected current context not empty")

				validNamedContext := api.NamedContext{
					Name: cliKubeconfig.CurrentContext,
					Context: api.Context{
						Cluster:  cliKubeconfig.Clusters[0].Name,
						AuthInfo: cliKubeconfig.AuthInfos[0].Name,
					},
				}
				assert.Equal(t, validNamedContext, cliKubeconfig.Contexts[0])
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster(),
	}))

	////
	// Legacy use case
	////
	t.Run("legacy", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		Cmd:      "scw k8s kubeconfig get {{ ." + clusterMetaKey + ".ID }} auth-method=legacy",
		BeforeFunc: core.BeforeFuncCombine(
			createCluster("get-kubeconfig-legacy", true),
			fetchClusterKubeconfigMetadata(false),
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				config, err := yaml.Marshal(ctx.Meta["Kubeconfig"].(api.Config))
				require.NoError(t, err)
				assert.Equal(t, ctx.Result.(string), string(config))
			},
		),
		AfterFunc: deleteCluster(),
	}))

	////
	// Copy token use case: current token of the cli is copied to the kubeconfig
	////
	t.Run("copy_cli_token", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createCluster("get-kubeconfig-copy-cli-token", true),
			fetchClusterKubeconfigMetadata(true),
		),
		Cmd: "scw k8s kubeconfig get {{ ." + clusterMetaKey + ".ID }} auth-method=copy-cli-token",
		OverrideEnv: map[string]string{
			scw.ScwAccessKeyEnv: "", // Hide keys in test env
			scw.ScwSecretKeyEnv: "", // Hide keys in test env
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				secretKey, ok := ctx.Client.GetSecretKey()
				assert.True(t, ok)

				// replace token inside the golden by fake token
				replacements := []core.GoldenReplacement{
					{
						Pattern: regexp.MustCompile(
							secretKey[0:8] + "-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
						),
						Replacement: "11111111-1111-1111-1111-111111111111",
					},
					{
						Pattern:     regexp.MustCompile(`cli-token-[0-9A-Fa-f]{8}`),
						Replacement: "cli-token-11111111",
					},
				}

				// protect golden
				checkGolden := core.TestCheckGoldenAndReplacePatterns(replacements...)
				checkGolden(t, ctx)

				// then protect result
				result := ctx.Result.(string)
				var matchFailed []string
				for _, replacement := range replacements {
					if !replacement.Pattern.MatchString(result) {
						if !replacement.OptionalMatch {
							matchFailed = append(matchFailed, replacement.Pattern.String())
						}

						continue
					}
					result = replacement.Pattern.ReplaceAllString(result, replacement.Replacement)
				}

				if len(matchFailed) > 0 {
					t.Fatalf("failed to match regex in results: %#q", matchFailed)
				}

				var cliKubeconfig api.Config
				err := yaml.Unmarshal([]byte(result), &cliKubeconfig)
				require.NoError(t, err)

				require.Lenf(t, cliKubeconfig.Clusters, 1,
					"expected 1 cluster, got %d", len(cliKubeconfig.Clusters))

				require.Lenf(t, cliKubeconfig.AuthInfos, 1,
					"expected 1 user, got %d", len(cliKubeconfig.AuthInfos))

				authInfos := []api.NamedAuthInfo{
					{
						Name: "cli-token-11111111",
						AuthInfo: api.AuthInfo{
							Token: "11111111-1111-1111-1111-111111111111",
						},
					},
				}
				assert.Equal(t, authInfos, cliKubeconfig.AuthInfos)

				assert.NotEmptyf(t, cliKubeconfig.CurrentContext,
					"expected current context not empty")

				validNamedContext := api.NamedContext{
					Name: cliKubeconfig.CurrentContext,
					Context: api.Context{
						Cluster:  cliKubeconfig.Clusters[0].Name,
						AuthInfo: cliKubeconfig.AuthInfos[0].Name,
					},
				}
				assert.Equal(t, validNamedContext, cliKubeconfig.Contexts[0])
			},
		),
		AfterFunc: deleteCluster(),
	}))
}
