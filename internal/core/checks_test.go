package core

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func TestCheckAPIKey(t *testing.T) {
	testCommands := NewCommands(
		&Command{
			Namespace: "test",
			ArgSpecs:  ArgSpecs{},
			ArgsType:  reflect.TypeOf(testType{}),
			Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
				// Test command reload the client so the profile used is the edited one
				return "", ReloadClient(ctx)
			},
		})
	metadataKey := "ApiKey"

	t.Run("basic", Test(&TestConfig{
		Commands:   testCommands,
		TmpHomeDir: true,
		BeforeFunc: func(ctx *BeforeFuncCtx) error {
			api := iam.NewAPI(ctx.Client)
			accessKey, exists := ctx.Client.GetAccessKey()
			if !exists {
				return fmt.Errorf("missing access-key")
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
			if !*UpdateCassettes {
				apiKey.AccessKey = "SCWXXXXXXXXXXXXXXXXX"
			}

			ctx.Meta[metadataKey] = apiKey
			cfg := &scw.Config{
				Profile: scw.Profile{
					AccessKey: &apiKey.AccessKey,
					SecretKey: apiKey.SecretKey,
				},
			}
			configPath := filepath.Join(ctx.OverrideEnv["HOME"], ".config", "scw", "config.yaml")

			return cfg.SaveTo(configPath)
		},
		Cmd: "scw test",
		Check: TestCheckCombine(
			TestCheckExitCode(0),
			func(t *testing.T, ctx *CheckFuncCtx) {
				assert.True(t, strings.HasPrefix(ctx.LogBuffer, "Current api key expires in"))
			},
		),
		AfterFunc: func(ctx *AfterFuncCtx) error {
			return iam.NewAPI(ctx.Client).DeleteAPIKey(&iam.DeleteAPIKeyRequest{
				AccessKey: ctx.Meta[metadataKey].(*iam.APIKey).AccessKey,
			})
		},
	}))
}
