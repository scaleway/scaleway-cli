package object_test

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/object/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

const (
	testBucketNameActionCreate = "-create"
	testBucketNameActionDelete = "-delete"
	testBucketNameActionGet    = "-get"
	testBucketNameActionUpdate = "-update"
)

func Test_BucketCreate(t *testing.T) {
	bucketName1 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionCreate)
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		Cmd:      "scw object bucket create " + bucketName1,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				bucket := ctx.Result.(*object.BucketResponse).BucketInfo
				assert.Equal(t, bucketName1, bucket.ID)
				assert.False(t, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag(nil), bucket.Tags)
				checkACL(t, "private", bucket.ACL, bucket.Owner)
			},
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteBucket(bucketName1),
	}))

	bucketName2 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionCreate)

	t.Run("With tags", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		Cmd: fmt.Sprintf(
			"scw object bucket create %s tags.0=\"key1=value1\" tags.1=\"key2=value2\"",
			bucketName2,
		),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				bucket := ctx.Result.(*object.BucketResponse).BucketInfo
				assert.Equal(t, bucketName2, bucket.ID)
				assert.False(t, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag{
					{
						Key:   scw.StringPtr("key1"),
						Value: scw.StringPtr("value1"),
					},
					{
						Key:   scw.StringPtr("key2"),
						Value: scw.StringPtr("value2"),
					},
				}, bucket.Tags)
				checkACL(t, "private", bucket.ACL, bucket.Owner)
			},
		),
		AfterFunc: deleteBucket(bucketName2),
	}))

	bucketName3 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionCreate)

	t.Run("With versioning", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		Cmd:      fmt.Sprintf("scw object bucket create %s enable-versioning=true", bucketName3),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				bucket := ctx.Result.(*object.BucketResponse).BucketInfo
				assert.Equal(t, bucketName3, bucket.ID)
				assert.True(t, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag(nil), bucket.Tags)
				checkACL(t, "private", bucket.ACL, bucket.Owner)
			},
		),
		AfterFunc: deleteBucket(bucketName3),
	}))

	bucketName4 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionCreate)

	t.Run("With ACL", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		Cmd:      fmt.Sprintf("scw object bucket create %s acl=authenticated-read", bucketName4),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				bucket := ctx.Result.(*object.BucketResponse).BucketInfo
				assert.Equal(t, bucketName4, bucket.ID)
				assert.False(t, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag(nil), bucket.Tags)
				checkACL(t, "authenticated-read", bucket.ACL, bucket.Owner)
			},
		),
		AfterFunc: deleteBucket(bucketName4),
	}))
}

func Test_BucketDelete(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionDelete)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   object.GetCommands(),
		BeforeFunc: createBucket(bucketName),
		Cmd:        "scw object bucket delete " + bucketName,
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_BucketGet(t *testing.T) {
	bucketName1 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionGet)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   object.GetCommands(),
		BeforeFunc: createBucket(bucketName1),
		Cmd:        "scw object bucket get " + bucketName1,
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				bucket := ctx.Result.(object.BucketGetResult)
				assert.Equal(t, bucketName1, bucket.ID)
				assert.False(t, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag(nil), bucket.Tags)
				checkACL(t, "private", bucket.ACL, bucket.Owner)
			},
		),
		AfterFunc: deleteBucket(bucketName1),
	}))

	bucketName2 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionGet)

	t.Run("With size", core.Test(&core.TestConfig{
		Commands:   object.GetCommands(),
		BeforeFunc: createBucket(bucketName2),
		Cmd:        fmt.Sprintf("scw object bucket get %s with-size=true", bucketName2),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				bucket := ctx.Result.(object.BucketGetResult)
				assert.Equal(t, bucketName2, bucket.ID)
				assert.False(t, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag(nil), bucket.Tags)
				checkACL(t, "private", bucket.ACL, bucket.Owner)
				assert.Equal(t, int64(0), *bucket.NbObjects)
				assert.Equal(t, int64(0), *bucket.NbParts)
				assert.Equal(t, scw.SizePtr(0), bucket.Size)
			},
		),
		AfterFunc: deleteBucket(bucketName2),
	}))
}

func Test_BucketList(t *testing.T) {
	bucketName1 := randomNameWithPrefix("cli-test-bucket-list")
	bucketName2 := randomNameWithPrefix("cli-test-bucket-list")

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketName1),
			createBucket(bucketName2),
		),
		Cmd: "scw object bucket list",
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteBucket(bucketName1),
			deleteBucket(bucketName2),
		),
	}))
}

func Test_BucketUpdate(t *testing.T) {
	bucketName1 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionUpdate)

	t.Run("Tags", core.Test(&core.TestConfig{
		Commands:   object.GetCommands(),
		BeforeFunc: createBucket(bucketName1),
		Cmd:        fmt.Sprintf("scw object bucket update %s tags.0=\"key1=value1\"", bucketName1),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				bucket := ctx.Result.(*object.BucketResponse).BucketInfo
				assert.Equal(t, bucketName1, bucket.ID)
				assert.False(t, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag{
					{
						Key:   scw.StringPtr("key1"),
						Value: scw.StringPtr("value1"),
					},
				}, bucket.Tags)
				checkACL(t, "private", bucket.ACL, bucket.Owner)
			},
		),
		AfterFunc: deleteBucket(bucketName1),
	}))

	bucketName2 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionUpdate)

	t.Run("All options", core.Test(&core.TestConfig{
		Commands:   object.GetCommands(),
		BeforeFunc: createBucket(bucketName2),
		Cmd: fmt.Sprintf(
			"scw object bucket update %s enable-versioning=true acl=public-read-write tags.0=\"key1=value1\" tags.1=\"key2=value2\"",
			bucketName2,
		),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				bucket := ctx.Result.(*object.BucketResponse).BucketInfo
				assert.Equal(t, bucketName2, bucket.ID)
				assert.True(t, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag{
					{
						Key:   scw.StringPtr("key1"),
						Value: scw.StringPtr("value1"),
					},
					{
						Key:   scw.StringPtr("key2"),
						Value: scw.StringPtr("value2"),
					},
				}, bucket.Tags)
				checkACL(t, "public-read-write", bucket.ACL, bucket.Owner)
			},
		),
		AfterFunc: deleteBucket(bucketName2),
	}))
}

func randomNameWithPrefix(prefix string) string {
	randomInt, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt))
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s-%d", prefix, randomInt)
}

func createBucket(name string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd("Bucket", "scw object bucket create "+name)
}

func deleteBucket(name string) core.AfterFunc {
	return core.ExecAfterCmd("scw object bucket delete " + name)
}

func checkACL(t *testing.T, expected string, actual []object.CustomS3ACLGrant, owner string) {
	t.Helper()
	grantsMap := make(map[types.Permission]string)
	for _, actualACL := range actual {
		if actualACL.Grantee == nil {
			t.Fatalf("ACL Grantee should not be nil")
		}
		grantsMap[actualACL.Permission] = *actualACL.Grantee
	}

	switch expected {
	case "private":
		assert.Len(t, grantsMap, 1)
		assert.Equal(t, owner, grantsMap["FULL_CONTROL"])
	case "public-read":
		assert.Len(t, grantsMap, 2)
		assert.Equal(t, owner, grantsMap["FULL_CONTROL"])
		assert.Equal(t, "AllUsers", grantsMap["READ"])
	case "public-read-write":
		assert.Len(t, grantsMap, 3)
		assert.Equal(t, owner, grantsMap["FULL_CONTROL"])
		assert.Equal(t, "AllUsers", grantsMap["READ"])
		assert.Equal(t, "AllUsers", grantsMap["WRITE"])
	case "authenticated-read":
		assert.Len(t, grantsMap, 2)
		assert.Equal(t, owner, grantsMap["FULL_CONTROL"])
		assert.Equal(t, "AuthenticatedUsers", grantsMap["READ"])
	}
}
