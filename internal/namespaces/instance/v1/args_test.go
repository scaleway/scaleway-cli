package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/stretchr/testify/assert"
)

type NullableStringValueRequest struct {
	Reverse *instance.NullableStringValue
}

func TestNullableStringValueUnmarshal(t *testing.T) {
	type testCase struct {
		args           []string
		data           any
		expectedStruct any
		expectedError  error
	}

	run := func(tc *testCase) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			err := args.UnmarshalStruct(tc.args, tc.data)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStruct, tc.data)
		}
	}

	t.Run("empty-request", run(&testCase{
		args:           []string(nil),
		data:           &NullableStringValueRequest{},
		expectedStruct: &NullableStringValueRequest{},
		expectedError:  nil,
	}))

	t.Run("request-with-reverse", run(&testCase{
		args: []string{
			"reverse=test",
		},
		data: &NullableStringValueRequest{},
		expectedStruct: &NullableStringValueRequest{
			Reverse: &instance.NullableStringValue{
				Null:  false,
				Value: "test",
			},
		},
		expectedError: nil,
	}))

	t.Run("request-with-empty-reverse", run(&testCase{
		args: []string{
			"reverse=",
		},
		data: &NullableStringValueRequest{},
		expectedStruct: &NullableStringValueRequest{
			Reverse: &instance.NullableStringValue{
				Null: true,
			},
		},
		expectedError: nil,
	}))
}

func TestNullableStringValueMarshal(t *testing.T) {
	type testCase struct {
		data           any
		expectedStruct any
		expectedError  error
	}

	run := func(tc *testCase) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			args, err := args.MarshalStruct(tc.data)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStruct, args)
		}
	}

	t.Run("empty-request", run(&testCase{
		data:           &NullableStringValueRequest{},
		expectedStruct: []string(nil),
		expectedError:  nil,
	}))

	t.Run("request-with-reverse", run(&testCase{
		data: &NullableStringValueRequest{
			Reverse: &instance.NullableStringValue{
				Null:  false,
				Value: "test",
			},
		},
		expectedStruct: []string{
			"reverse=test",
		},
		expectedError: nil,
	}))

	t.Run("request-with-empty-reverse", run(&testCase{
		data: &NullableStringValueRequest{
			Reverse: &instance.NullableStringValue{
				Null:  true,
				Value: "",
			},
		},
		expectedStruct: []string{
			"reverse=",
		},
		expectedError: nil,
	}))
}
