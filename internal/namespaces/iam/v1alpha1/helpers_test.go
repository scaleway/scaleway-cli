package iam_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func addSSHKey(metaKey string, key string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		ctx.Meta[metaKey] = ctx.ExecuteCmd([]string{
			"scw", "iam", "ssh-key", "create", "public-key=" + key,
		})

		return nil
	}
}

func newSliceResulter[T any](
	t *testing.T,
	sliceTypeName string,
	extractor func([]T) any,
) func(any) any {
	t.Helper()

	return func(result any) any {
		slice, ok := result.([]T)
		if !ok {
			t.Fatalf("expected []%s, got %T", sliceTypeName, result)
		}
		if slice == nil {
			t.Fatalf("%s slice is nil", sliceTypeName)
		}
		if len(slice) == 0 {
			t.Fatalf("no %s found", sliceTypeName)
		}

		return extractor(slice)
	}
}
