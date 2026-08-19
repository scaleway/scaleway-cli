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
	testBucketNameActionLifecycle = "-lifecycle"
	testLifecyclePath             = "testdata/test-lifecycle.json"
	testLifecycleContent          = `{
  "Rules": [
    {
      "ID": "ArchiveAndPurgeLogs",
      "Status": "Enabled",
      "Filter": {
        "Prefix": "logs/"
      },
      "Transitions": [
        {
          "Days": 30,
          "StorageClass": "ONEZONE_IA"
        },
        {
          "Days": 90,
          "StorageClass": "GLACIER"
        }
      ],
      "Expiration": {
        "Days": 365
      }
    },
    {
      "ID": "CleanupIncompleteMultipartUploads",
      "Status": "Enabled",
      "Filter": {},
      "AbortIncompleteMultipartUpload": {
        "DaysAfterInitiation": 7
      }
    }
  ]
}`
)

func Test_BucketLifecycleCreate(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionLifecycle)

	bucketNameEmpty := bucketName + "-empty"
	testLifecyclePathEmpty := "testdata/test-lifecycle-empty.json"

	t.Run("Empty", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketNameEmpty),
			createFile(testLifecyclePathEmpty, "{}"),
		),
		Cmd: fmt.Sprintf(
			"scw object bucket-lifecycle create %s lifecycle-configuration=%s",
			bucketNameEmpty,
			testLifecyclePathEmpty,
		),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteBucket(bucketNameEmpty),
			deleteFile(testLifecyclePathEmpty),
		),
	}))

	bucketNameEmptyRules := bucketName + "-empty-rules"
	testLifecyclePathEmptyRules := "testdata/test-lifecycle-empty-rules.json"

	t.Run("Empty rules", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketNameEmptyRules),
			createFile(testLifecyclePathEmptyRules, `{"Rules": []}`),
		),
		Cmd: fmt.Sprintf(
			"scw object bucket-lifecycle create %s lifecycle-configuration=%s",
			bucketNameEmptyRules,
			testLifecyclePathEmptyRules,
		),
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains("Lifecycle configuration cannot be empty"),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteBucket(bucketNameEmptyRules),
			deleteFile(testLifecyclePathEmptyRules),
		),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketName),
			createFile(testLifecyclePath, testLifecycleContent),
		),
		Cmd: fmt.Sprintf(
			"scw object bucket-lifecycle create %s lifecycle-configuration=%s",
			bucketName,
			testLifecyclePath,
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
			deleteFile(testLifecyclePath),
		),
	}))
}

func Test_BucketLifecycleDelete(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionLifecycle)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketName),
			createFile(testLifecyclePath, testLifecycleContent),
			createLifecycleConfiguration(bucketName, testLifecyclePath),
		),
		Cmd: "scw object bucket-lifecycle delete " + bucketName,
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteBucket(bucketName),
			deleteFile(testLifecyclePath),
		),
	}))
}

func Test_BucketLifecycleGet(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionLifecycle)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketName),
			createFile(testLifecyclePath, testLifecycleContent),
			createLifecycleConfiguration(bucketName, testLifecyclePath),
		),
		Cmd: "scw object bucket-lifecycle get " + bucketName,
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteBucket(bucketName),
			deleteFile(testLifecyclePath),
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

func createLifecycleConfiguration(bucketName, configurationPath string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd("Lifecycle", fmt.Sprintf(
		"scw object bucket-lifecycle create %s lifecycle-configuration=%s",
		bucketName,
		configurationPath,
	))
}
