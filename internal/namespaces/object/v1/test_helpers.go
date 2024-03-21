package object

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/stretchr/testify/assert"
)

func randomNameWithPrefix(prefix string) string {
	randomInt, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s-%d", prefix, randomInt)
}

func createBucket(name string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd("Bucket", fmt.Sprintf("scw object bucket create %s", name))
}

func deleteBucket(name string) core.AfterFunc {
	return core.ExecAfterCmd(fmt.Sprintf("scw object bucket delete %s", name))
}

func checkACL(t *testing.T, actual []customS3ACLGrant, expected string, ownerID *string) {
	grantsMap := make(map[types.Permission]string)
	for _, actualACL := range actual {
		grantsMap[actualACL.Permission] = *actualACL.Grantee
	}

	switch expected {
	case "private":
		assert.Equal(t, len(grantsMap), 1)
		assert.Equal(t, *ownerID, grantsMap["FULL_CONTROL"])
	case "public-read":
		assert.Equal(t, len(grantsMap), 2)
		assert.Equal(t, *ownerID, grantsMap["FULL_CONTROL"])
		assert.Equal(t, "AllUsers", grantsMap["READ"])
	case "public-read-write":
		assert.Equal(t, len(grantsMap), 3)
		assert.Equal(t, *ownerID, grantsMap["FULL_CONTROL"])
		assert.Equal(t, "AllUsers", grantsMap["READ"])
		assert.Equal(t, "AllUsers", grantsMap["WRITE"])
	case "authenticated-read":
		assert.Equal(t, len(grantsMap), 2)
		assert.Equal(t, *ownerID, grantsMap["FULL_CONTROL"])
		assert.Equal(t, "AuthenticatedUsers", grantsMap["READ"])
	}
}
