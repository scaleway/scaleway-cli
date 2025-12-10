package redis_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/redis/v1"
	redisSDK "github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/stretchr/testify/assert"
)

func TestRedisVersionSettingsCommand(t *testing.T) {
	t.Run("List settings for Redis version", core.Test(&core.TestConfig{
		Commands: redis.GetCommands(),
		Cmd:      "scw redis version list-settings version=7.2.11",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				if ctx.Result != nil {
					settings := ctx.Result.([]*redisSDK.AvailableClusterSetting)
					assert.NotEmpty(t, settings, "should return at least one setting")
					// Verify that settings have the expected structure
					for _, setting := range settings {
						assert.NotEmpty(t, setting.Name, "setting should have a name")
					}
				}
			},
		),
	}))
}
