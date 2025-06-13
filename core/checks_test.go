package core_test

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func TestCheckAPIKey(t *testing.T) {
	testCommands := core.NewCommands(
		&core.Command{
			Namespace: "test",
			ArgSpecs:  core.ArgSpecs{},
			ArgsType:  reflect.TypeOf(testType{}),
			Run: func(ctx context.Context, _ any) (i any, e error) {
				// Test command reload the client so the profile used is the edited one
				return "", core.ReloadClient(ctx)
			},
		})
	metadataKey := "ApiKey"

	t.Run("basic", core.Test(&core.TestConfig{
		Commands:   testCommands,
		TmpHomeDir: true,
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			api := iam.NewAPI(ctx.Client)
			accessKey, exists := ctx.Client.GetAccessKey()
			if !exists {
				return errors.New("missing access-key")
			}

			apiKey, err := api.GetAPIKey(&iam.GetAPIKeyRequest{
				AccessKey: accessKey,
			})
			if err != nil {
				return err
			}
			expiresAt := time.Now().Add(time.Hour)
			apiKey, err = api.CreateAPIKey(&iam.CreateAPIKeyRequest{
				ApplicationID: apiKey.ApplicationID,
				UserID:        apiKey.UserID,
				ExpiresAt:     &expiresAt,
				Description:   "test-cli-TestCheckAPIKey",
			})
			if err != nil {
				return err
			}
			if !*core.UpdateCassettes {
				apiKey.AccessKey = "SCWXXXXXXXXXXXXXXXXX"
			}

			ctx.Meta[metadataKey] = apiKey
			cfg := &scw.Config{
				Profile: scw.Profile{
					AccessKey:             &apiKey.AccessKey,
					SecretKey:             apiKey.SecretKey,
					DefaultProjectID:      &apiKey.DefaultProjectID,
					DefaultOrganizationID: &apiKey.DefaultProjectID,
				},
			}
			configPath := filepath.Join(ctx.OverrideEnv["HOME"], ".config", "scw", "config.yaml")

			return cfg.SaveTo(configPath)
		},
		Cmd: "scw test",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.True(t, strings.HasPrefix(ctx.LogBuffer, "Current api key expires in"))
			},
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			return iam.NewAPI(ctx.Client).DeleteAPIKey(&iam.DeleteAPIKeyRequest{
				AccessKey: ctx.Meta[metadataKey].(*iam.APIKey).AccessKey,
			})
		},
	}))
}
