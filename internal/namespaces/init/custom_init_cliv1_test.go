package init

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type cliV1 struct {
	OrganizationId string `json:"organization"`
	Token          string `json:"token"`
	Version        string `json:"version"`
}

func setUpConfV1(ctx *core.BeforeFuncCtx, confV1 *cliV1) error {
	res, err := json.Marshal(*confV1)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(ctx.OverrideEnv["HOME"], ".scwrc"), res, 0600)
	if err != nil {
		return err
	}
	return nil
}

func Test_InitCLIv1(t *testing.T) {
	secretKey := dummyUUID
	// if you are recording, you must place a valid token in the environment variable SCW_TEST_SECRET_KEY
	if os.Getenv("SCW_TEST_SECRET_KEY") != "" {
		secretKey = os.Getenv("SCW_TEST_SECRET_KEY")
	}
	t.Run("CLIv1Config", func(t *testing.T) {
		dummySecretV1 := "44444444-4444-4444-4444-444444444444"
		dummyConfigV1 := &cliV1{
			OrganizationId: "33333333-3333-3333-3333-333333333333",
			Token:          dummySecretV1,
			Version:        "v1.20",
		}

		t.Run("Import", func(t *testing.T) {
			core.Test(&core.TestConfig{
				Commands: GetCommands(),
				BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
					return setUpConfV1(ctx, dummyConfigV1)
				},
				Cmd: "scw init secret-key=" + secretKey + " organization-id=11111111-1111-1111-1111-111111111111 send-telemetry=true install-autocomplete=false remove-v1-config=false with-ssh-key=false",
				Check: core.TestCheckCombine(
					core.TestCheckExitCode(0),
					core.TestCheckGolden(),
					checkConfig(func(t *testing.T, config *scw.Config) {
						assert.Equal(t, dummySecretV1, *config.SecretKey)
					}),
				),
				TmpHomeDir:          true,
				PromptResponseMocks: nil,
			})(t)
		})

		t.Run("ImportAndDelete", func(t *testing.T) {
			core.Test(&core.TestConfig{
				Commands: GetCommands(),
				BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
					return setUpConfV1(ctx, dummyConfigV1)
				},
				Cmd: "scw init secret-key=" + secretKey + " organization-id=11111111-1111-1111-1111-111111111111 send-telemetry=true install-autocomplete=false remove-v1-config=true with-ssh-key=false",
				Check: core.TestCheckCombine(
					core.TestCheckExitCode(0),
					core.TestCheckGolden(),
					checkConfig(func(t *testing.T, config *scw.Config) {
						assert.Equal(t, dummySecretV1, *config.SecretKey)
					}),
					func(t *testing.T, ctx *core.CheckFuncCtx) {
						// We check that the CLIv1 configuration file is deleted
					},
				),
				TmpHomeDir:          true,
				PromptResponseMocks: nil,
			})(t)

		})

		t.Run("ImportAndNoDelete", func(t *testing.T) {
			core.Test(&core.TestConfig{
				Commands: GetCommands(),
				BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
					return setUpConfV1(ctx, dummyConfigV1)
				},
				Cmd: "scw init secret-key=" + secretKey + " organization-id=11111111-1111-1111-1111-111111111111 send-telemetry=true install-autocomplete=false remove-v1-config=false with-ssh-key=false",
				Check: core.TestCheckCombine(
					core.TestCheckExitCode(0),
					core.TestCheckGolden(),
					checkConfig(func(t *testing.T, config *scw.Config) {
						assert.Equal(t, dummySecretV1, *config.SecretKey)
					}),
					func(t *testing.T, ctx *core.CheckFuncCtx) {
						// We check that the CLIv1 configuration file is not deleted
					},
				),
				TmpHomeDir:          true,
				PromptResponseMocks: nil,
			})(t)
		})
	})

}
