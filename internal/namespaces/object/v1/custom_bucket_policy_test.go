package object_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/object/v1"
	"github.com/stretchr/testify/assert"
)

const (
	testBucketNameActionPolicy = "-policy"
	testPolicyPath             = "testdata/test-policy.json"
	testPolicyContent          = `{
  "Version": "2023-04-17",
  "Id": "MyBucketPolicy",
  "Statement": [
    {
      "Sid": "my-id",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "*",
      "Resource": "%[1]s"
    },
    {
      "Sid": "Scaleway secure statement",
      "Effect": "Allow",
      "Principal": {
        "SCW": "user_id:11111111-1111-1111-1111-111111111111"
      },
      "Action": "*",
      "Resource": [
        "%[1]s",
        "%[1]s/*"
      ]
    }
  ]
}`
)

func Test_BucketPolicyCreate(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionPolicy)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketName),
			createFile(testPolicyPath, fmt.Sprintf(testPolicyContent, bucketName)),
		),
		Cmd: fmt.Sprintf(
			"scw object bucket-policy create %s policy=%s",
			bucketName,
			testPolicyPath,
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
			deleteFile(testPolicyPath),
		),
	}))
}

func Test_BucketPolicyDelete(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionPolicy)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketName),
			createFile(testPolicyPath, fmt.Sprintf(testPolicyContent, bucketName)),
			createPolicy(bucketName, testPolicyPath),
		),
		Cmd: "scw object bucket-policy delete " + bucketName,
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteBucket(bucketName),
			deleteFile(testPolicyPath),
		),
	}))
}

func Test_BucketPolicyGet(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionPolicy)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketName),
			createFile(testPolicyPath, fmt.Sprintf(testPolicyContent, bucketName)),
			createPolicy(bucketName, testPolicyPath),
		),
		Cmd: "scw object bucket-policy get " + bucketName,
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteBucket(bucketName),
			deleteFile(testPolicyPath),
		),
	}))
}

func createPolicy(bucketName, policyPath string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd("Policy", fmt.Sprintf(
		"scw object bucket-policy create %s policy=%s",
		bucketName,
		policyPath,
	))
}
