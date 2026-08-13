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
			createPolicyFile(bucketName, testPolicyPath),
		),
		Cmd: fmt.Sprintf(
			"scw object bucket-policy create %s policy-path=%s",
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
		AfterFunc: deleteBucket(bucketName),
	}))
}

func Test_BucketPolicyDelete(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionPolicy)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketName),
			createPolicyFile(bucketName, testPolicyPath),
			createPolicy(bucketName, testPolicyPath),
		),
		Cmd: "scw object bucket-policy delete " + bucketName,
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteBucket(bucketName),
	}))
}

func Test_BucketPolicyGet(t *testing.T) {
	bucketName := randomNameWithPrefix(core.TestBucketNamePrefix + testBucketNameActionPolicy)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: object.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createBucket(bucketName),
			createPolicyFile(bucketName, testPolicyPath),
			createPolicy(bucketName, testPolicyPath),
		),
		Cmd: "scw object bucket-policy get " + bucketName,
		Check: core.TestCheckCombine(
			core.TestCheckS3Golden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteBucket(bucketName),
	}))
}

func createPolicyFile(bucketName, policyPath string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		// Policy content
		policy := fmt.Sprintf(testPolicyContent, bucketName)

		// Create policy file
		err := os.WriteFile(policyPath, []byte(policy), 0o644)
		if err != nil {
			fmt.Println("error writing file:", err)

			return nil
		}

		return nil
	}
}

func createPolicy(bucketName, policyPath string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd("Policy", fmt.Sprintf(
		"scw object bucket-policy create %s policy-path=%s",
		bucketName,
		policyPath,
	))
}
