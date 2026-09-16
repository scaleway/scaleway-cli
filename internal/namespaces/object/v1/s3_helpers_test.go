package object_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/object/v1"
)

func Test_FormatAccessKey(t *testing.T) {
	cases := []struct {
		testName         string
		accessKey        string
		argProjectID     string
		defaultProjectID string
		expectedFormat   string
	}{
		{
			testName:         "standard access key",
			accessKey:        "SCW00000000000000000",
			argProjectID:     "11111111-1111-1111-1111-111111111111",
			defaultProjectID: "22222222-2222-2222-2222-222222222222",
			expectedFormat:   "SCW00000000000000000@11111111-1111-1111-1111-111111111111",
		},

		{
			testName:         "access key with project id, should not change",
			accessKey:        "SCW00000000000000000@11111111-1111-1111-1111-111111111111",
			argProjectID:     "11111111-1111-1111-1111-111111111111",
			defaultProjectID: "22222222-2222-2222-2222-222222222222",
			expectedFormat:   "SCW00000000000000000@11111111-1111-1111-1111-111111111111",
		},

		{
			testName:         "access key with project id, should be replaced with arg",
			accessKey:        "SCW00000000000000000@22222222-2222-2222-2222-222222222222",
			argProjectID:     "11111111-1111-1111-1111-111111111111",
			defaultProjectID: "22222222-2222-2222-2222-222222222222",
			expectedFormat:   "SCW00000000000000000@11111111-1111-1111-1111-111111111111",
		},

		{
			testName:         "access key with project id, should be replaced with arg",
			accessKey:        "SCW00000000000000000@33333333-3333-3333-3333-333333333333",
			argProjectID:     "11111111-1111-1111-1111-111111111111",
			defaultProjectID: "22222222-2222-2222-2222-222222222222",
			expectedFormat:   "SCW00000000000000000@11111111-1111-1111-1111-111111111111",
		},

		{
			testName:         "access key with project id, should be replaced with default",
			accessKey:        "SCW00000000000000000@33333333-3333-3333-3333-333333333333",
			argProjectID:     "",
			defaultProjectID: "22222222-2222-2222-2222-222222222222",
			expectedFormat:   "SCW00000000000000000@22222222-2222-2222-2222-222222222222",
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			resultAccessKey := object.FormatAccessKey(
				c.accessKey,
				c.argProjectID,
				c.defaultProjectID,
			)
			if resultAccessKey != c.expectedFormat {
				t.Fatalf("expected '%s', got '%s'", c.expectedFormat, resultAccessKey)
			}
		})
	}
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

func deleteFile(filePath string) core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		return os.Remove(filePath)
	}
}
