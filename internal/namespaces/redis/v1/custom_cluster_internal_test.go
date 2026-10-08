package redis

import (
	"testing"

	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createRequestZone(t *testing.T) {
	// During shell completion the embedded *redis.CreateClusterRequest is
	// nil (reflect.New leaves it unallocated); this must not panic.
	t.Run("Nil embedded request does not panic", func(t *testing.T) {
		require.Equal(t, scw.Zone(""), createRequestZone(&redisCreateClusterRequestCustom{}))
	})

	t.Run("Returns zone when embedded request is set", func(t *testing.T) {
		req := &redisCreateClusterRequestCustom{
			CreateClusterRequest: &redis.CreateClusterRequest{Zone: scw.ZoneFrPar1},
		}
		assert.Equal(t, scw.ZoneFrPar1, createRequestZone(req))
	})
}
