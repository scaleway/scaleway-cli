package object_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/object/v1"
	"github.com/stretchr/testify/assert"
)

const (
	testBucketNameActionCors = "-cors"
	testCorsPath             = "testdata/test-cors.json"
	testCorsContent          = `{
  "CORSRules": [
    {
      "AllowedOrigins": ["http://MY_DOMAIN_NAME", "http://www.MY_DOMAIN_NAME"],
      "AllowedHeaders": ["*"],
      "AllowedMethods": ["GET", "HEAD", "POST", "PUT", "DELETE"],
      "MaxAgeSeconds": 3000,
      "ExposeHeaders": ["Etag"]
    }
  ]
}`
)

func Test_BucketCorsCreate(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionCors)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketName),
			createFile(testCorsPath, testCorsContent),
		),
		Cmd: fmt.Sprintf(
			"scw object bucket-cors create %s cors-configuration=%s",
			bucketName,
			testCorsPath,
		),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				bucket := ctx.Result.(*object.BucketResponse).BucketInfo
				assert.Equal(t, bucketName, bucket.ID)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteBucket(bucketName),
			deleteFile(testCorsPath),
		),
	}))
}

func Test_BucketCorsDelete(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionCors)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketName),
			createFile(testCorsPath, testCorsContent),
			createCors(bucketName, testCorsPath),
		),
		Cmd: "scw object bucket-cors delete " + bucketName,
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteBucket(bucketName),
			deleteFile(testCorsPath),
		),
	}))
}

func Test_BucketCorsGet(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionCors)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketName),
			createFile(testCorsPath, testCorsContent),
			createCors(bucketName, testCorsPath),
		),
		Cmd: "scw object bucket-cors get " + bucketName,
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteBucket(bucketName),
			deleteFile(testCorsPath),
		),
	}))
}

func createFile(path, content string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		err := os.WriteFile(path, []byte(content), 0o644)
		if err != nil {
			return fmt.Errorf("error writing file %s: %w", path, err)
		}

		return nil
	}
}

func createCors(bucketName, configurationPath string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd("Cors", fmt.Sprintf(
		"scw object bucket-cors create %s cors-configuration=%s",
		bucketName,
		configurationPath,
	))
}
