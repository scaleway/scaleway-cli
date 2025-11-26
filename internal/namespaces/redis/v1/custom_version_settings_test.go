package redis_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	redis "github.com/scaleway/scaleway-cli/v2/internal/namespaces/redis/v1"
	redisSDK "github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisVersionSettingsCommand(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/redis/v1/zones/fr-par-1/cluster-versions", r.URL.Path)
		q := r.URL.Query()
		assert.Equal(t, "7.2.11", q.Get("version"))

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(&redisSDK.ListClusterVersionsResponse{
			Versions: []*redisSDK.ClusterVersion{
				{
					Version: "7.2.11",
					AvailableSettings: []*redisSDK.AvailableClusterSetting{
						{
							Name:         "maxclients",
							DefaultValue: scw.StringPtr("10000"),
							Description:  "Maximum number of connected clients",
						},
					},
					Zone: scw.ZoneFrPar1,
				},
			},
		})
		if err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := scw.NewClient(
		scw.WithAPIURL(server.URL),
		scw.WithDefaultZone(scw.ZoneFrPar1),
	)
	require.NoError(t, err)

	ctx := core.InjectMeta(context.Background(), &core.Meta{
		Client: client,
	})

	cmds := redis.GetCommands()
	cmd := cmds.MustFind("redis", "version", "settings")

	args := reflect.New(cmd.ArgsType).Interface()
	argsValue := reflect.ValueOf(args).Elem()
	argsValue.FieldByName("Version").SetString("7.2.11")
	argsValue.FieldByName("Zone").Set(reflect.ValueOf(scw.ZoneFrPar1))

	result, err := cmd.Run(ctx, args)
	require.NoError(t, err)

	settings, ok := result.([]*redisSDK.AvailableClusterSetting)
	require.True(t, ok)
	require.Len(t, settings, 1)
	require.Equal(t, "maxclients", settings[0].Name)
	require.Equal(t, "10000", *settings[0].DefaultValue)
}
