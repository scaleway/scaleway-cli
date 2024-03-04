package object

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const testBucketNameActionCreate = "-create"
const testBucketNameActionDelete = "-delete"
const testBucketNameActionGet = "-get"
const testBucketNameActionUpdate = "-update"

func Test_BucketCreate(t *testing.T) {
	bucketName1 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionCreate)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      fmt.Sprintf("scw object bucket create %s", bucketName1),
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				bucket := ctx.Result.(*objectBucketResponse).BucketInfo
				assert.Equal(t, bucketName1, bucket.ID)
				assert.Equal(t, false, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag(nil), bucket.Tags)
				checkACL(t, bucket.ACL, "private", bucket.Owner)
			},
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteBucket(bucketName1),
	}))

	bucketName2 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionCreate)

	t.Run("With tags", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      fmt.Sprintf("scw object bucket create %s tags.0=\"key1=value1\" tags.1=\"key2=value2\"", bucketName2),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				bucket := ctx.Result.(*objectBucketResponse).BucketInfo
				assert.Equal(t, bucketName2, bucket.ID)
				assert.Equal(t, false, bucket.EnableVersioning)
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
				checkACL(t, bucket.ACL, "private", bucket.Owner)
			},
		),
		AfterFunc: deleteBucket(bucketName2),
	}))

	bucketName3 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionCreate)

	t.Run("With versioning", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      fmt.Sprintf("scw object bucket create %s enable-versioning=true", bucketName3),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				bucket := ctx.Result.(*objectBucketResponse).BucketInfo
				assert.Equal(t, bucketName3, bucket.ID)
				assert.Equal(t, true, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag(nil), bucket.Tags)
				checkACL(t, bucket.ACL, "private", bucket.Owner)
			},
		),
		AfterFunc: deleteBucket(bucketName3),
	}))

	bucketName4 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionCreate)

	t.Run("With ACL", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      fmt.Sprintf("scw object bucket create %s acl=authenticated-read", bucketName4),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				bucket := ctx.Result.(*objectBucketResponse).BucketInfo
				assert.Equal(t, bucketName4, bucket.ID)
				assert.Equal(t, false, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag(nil), bucket.Tags)
				checkACL(t, bucket.ACL, "authenticated-read", bucket.Owner)
			},
		),
		AfterFunc: deleteBucket(bucketName4),
	}))
}

func Test_BucketDelete(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionDelete)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createBucket(bucketName),
		Cmd:        fmt.Sprintf("scw object bucket delete %s", bucketName),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_BucketGet(t *testing.T) {
	bucketName1 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionGet)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createBucket(bucketName1),
		Cmd:        fmt.Sprintf("scw object bucket get %s", bucketName1),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				bucket := ctx.Result.(objectBucketGetResult)
				assert.Equal(t, bucketName1, bucket.ID)
				assert.Equal(t, false, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag(nil), bucket.Tags)
				checkACL(t, bucket.ACL, "private", bucket.Owner)
			},
		),
		AfterFunc: deleteBucket(bucketName1),
	}))

	bucketName2 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionGet)

	t.Run("With size", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createBucket(bucketName2),
		Cmd:        fmt.Sprintf("scw object bucket get %s with-size=true", bucketName2),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				bucket := ctx.Result.(objectBucketGetResult)
				assert.Equal(t, bucketName2, bucket.ID)
				assert.Equal(t, false, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag(nil), bucket.Tags)
				checkACL(t, bucket.ACL, "private", bucket.Owner)
				assert.Equal(t, int64(0), *bucket.NbObjects)
				assert.Equal(t, int64(0), *bucket.NbParts)
				assert.Equal(t, scw.SizePtr(0), bucket.Size)
			},
		),
		AfterFunc: deleteBucket(bucketName2),
	}))
}

func Test_BucketUpdate(t *testing.T) {
	bucketName1 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionUpdate)

	t.Run("Tags", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createBucket(bucketName1),
		Cmd:        fmt.Sprintf("scw object bucket update %s tags.0=\"key1=value1\"", bucketName1),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				bucket := ctx.Result.(*objectBucketResponse).BucketInfo
				assert.Equal(t, bucketName1, bucket.ID)
				assert.Equal(t, false, bucket.EnableVersioning)
				assert.Equal(t, []types.Tag{
					{
						Key:   scw.StringPtr("key1"),
						Value: scw.StringPtr("value1"),
					},
				}, bucket.Tags)
				checkACL(t, bucket.ACL, "private", bucket.Owner)
			},
		),
		AfterFunc: deleteBucket(bucketName1),
	}))

	bucketName2 := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionUpdate)

	t.Run("All options", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createBucket(bucketName2),
		Cmd:        fmt.Sprintf("scw object bucket update %s enable-versioning=true acl=public-read-write tags.0=\"key1=value1\" tags.1=\"key2=value2\"", bucketName2),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				bucket := ctx.Result.(*objectBucketResponse).BucketInfo
				assert.Equal(t, bucketName2, bucket.ID)
				assert.Equal(t, true, bucket.EnableVersioning)
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
				checkACL(t, bucket.ACL, "public-read-write", bucket.Owner)
			},
		),
		AfterFunc: deleteBucket(bucketName2),
	}))
}
