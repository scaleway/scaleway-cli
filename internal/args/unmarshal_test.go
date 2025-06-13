package args_test

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	args.TestForceNow = scw.TimePtr(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
}

func TestUnmarshalStruct(t *testing.T) {
	type TestCase struct {
		args     []string
		error    string
		expected any
		data     any
	}

	stringPtr := "test"
	slicePtr := []string{"0", "1", "2"}

	run := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			if testCase.data == nil {
				testCase.data = reflect.New(reflect.TypeOf(testCase.expected).Elem()).Interface()
			}
			err := args.UnmarshalStruct(testCase.args, testCase.data)

			if testCase.error == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.expected, testCase.data)
			} else {
				assert.Equal(t, testCase.error, err.Error())
			}
		}
	}

	args.RegisterUnmarshalFunc((*height)(nil), unmarshalHeight)

	t.Run("basic", run(TestCase{
		args: []string{
			"string=test",
			"int=42",
			"int16=16",
			"int32=32",
			"int64=64",
			"u-int16=16",
			"u-int32=32",
			"u-int64=64",
			"float32=3.2",
			"float64=6.4",
			"string-ptr=test",
			"bool=true",
		},
		expected: &Basic{
			String:    "test",
			Int:       42,
			Int16:     16,
			Int32:     32,
			Int64:     64,
			UInt16:    16,
			UInt32:    32,
			UInt64:    64,
			Float32:   3.2,
			Float64:   6.4,
			StringPtr: &stringPtr,
			Bool:      true,
		},
	}))

	t.Run("false-bool", run(TestCase{
		args: []string{
			"bool=false",
		},
		expected: &Basic{
			Bool: false,
		},
	}))

	t.Run("data-must-be-a-pointer", run(TestCase{
		args:  []string{},
		data:  Basic{},
		error: "data must be a pointer to a struct",
	}))

	t.Run("invalid-arg-name", run(TestCase{
		args: []string{
			"testCase=12",
		},
		expected: &Basic{},
		error:    "cannot unmarshal arg 'testCase=12': arg name must only contain lowercase letters, numbers or dashes",
	}))

	t.Run("field-do-not-exist", run(TestCase{
		args: []string{
			"unknown-field=12",
		},
		expected: &Basic{},
		error:    "cannot unmarshal arg 'unknown-field=12': unknown argument",
	}))

	t.Run("invalid-bool", run(TestCase{
		args: []string{
			"bool=invalid",
		},
		expected: &Basic{},
		error:    "cannot unmarshal arg 'bool=invalid': *bool is not unmarshalable: invalid boolean value",
	}))

	t.Run("missing-slice-index", run(TestCase{
		args: []string{
			"strings.1=2",
		},
		expected: &Slice{},
		error:    "cannot unmarshal arg 'strings.1=2': missing index 0, all indices prior to 1 must be set as well",
	}))

	t.Run("missing-slice-indices", run(TestCase{
		args: []string{
			"strings.5=2",
		},
		expected: &Slice{},
		error:    "cannot unmarshal arg 'strings.5=2': missing indices, 0,1,2,3,4 all indices prior to 5 must be set as well",
	}))

	t.Run("missing-slice-indices-overflow", run(TestCase{
		args: []string{
			"strings.99999=2",
		},
		expected: &Slice{},
		error:    "cannot unmarshal arg 'strings.99999=2': missing indices, 0,1,2,3,4,5,6,7,8,9,... all indices prior to 99999 must be set as well",
	}))

	t.Run("duplicate-slice-index", run(TestCase{
		args: []string{
			"basics.0.string=2",
			"basics.0.string=2",
		},
		expected: &Slice{},
		error:    "cannot unmarshal arg 'basics.0.string=2': duplicate argument",
	}))

	t.Run("slice-with-negative-index", run(TestCase{
		args: []string{
			"strings.0=2",
			"strings.-1=2",
		},
		expected: &Slice{},
		error:    "cannot unmarshal arg 'strings.-1=2': invalid index '-1' is not a positive integer",
	}))

	t.Run("nested-slice-with-invalid-index", run(TestCase{
		args: []string{
			"basics.string=test",
		},
		expected: &Slice{},
		error:    "cannot unmarshal arg 'basics.string=test': invalid index 'string' is not a positive integer",
	}))

	t.Run("basic-slice", run(TestCase{
		args: []string{
			"strings.0=1",
			"strings.1=2",
			"strings.2=3",
			"strings.3=test",
			"strings-ptr.0=test",
			"strings-ptr.1=test",
			"slice-ptr.0=0",
			"slice-ptr.1=1",
			"slice-ptr.2=2",
			"basics.0.string=test",
			"basics.0.int=42",
			"basics.1.string=test",
			"basics.1.int=42",
		},
		expected: &Slice{
			Strings:    []string{"1", "2", "3", "test"},
			StringsPtr: []*string{&stringPtr, &stringPtr},
			SlicePtr:   &slicePtr,
			Basics: []Basic{
				{
					String: "test",
					Int:    42,
				},
				{
					String: "test",
					Int:    42,
				},
			},
		},
	}))

	t.Run("empty-slice", run(TestCase{
		args: []string{
			"slice-ptr=none",
		},
		expected: &Slice{
			Strings:    []string(nil),
			SlicePtr:   scw.StringsPtr([]string{}),
			StringsPtr: []*string(nil),
		},
	}))

	t.Run("none-on-non-pointer-slice", run(TestCase{
		args: []string{
			"strings=none",
		},
		error: "cannot unmarshal arg 'strings=none': missing index on the array",
		data:  &Slice{},
	}))

	t.Run("simple-parent-child-conflict", run(TestCase{
		args: []string{
			"slice-ptr=none",
			"slice-ptr.0=none",
		},
		error: "arguments 'slice-ptr' and 'slice-ptr.0' cannot be used simultaneously",
		data:  &Slice{},
	}))

	t.Run("simple-child-parent-conflict", run(TestCase{
		args: []string{
			"slice-ptr.0=none",
			"slice-ptr=none",
		},
		error: "arguments 'slice-ptr.0' and 'slice-ptr' cannot be used simultaneously",
		data:  &Slice{},
	}))

	t.Run("well-known-types", run(TestCase{
		args: []string{
			"size=20gb",
		},
		expected: &WellKnownTypes{
			Size: 20 * scw.GB,
		},
	}))

	t.Run("Absolute date", run(TestCase{
		args: []string{
			"time=2006-01-02T15:04:05Z",
		},
		expected: &WellKnownTypes{
			Time: time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
		},
	}))

	t.Run("Relative date positive", run(TestCase{
		args: []string{
			"time=+1m1s",
		},
		expected: &WellKnownTypes{
			Time: time.Date(1970, 1, 1, 0, 1, 1, 0, time.UTC),
		},
	}))

	t.Run("Relative date negative", run(TestCase{
		args: []string{
			"time=-1m1s",
		},
		expected: &WellKnownTypes{
			Time: time.Date(1969, 12, 31, 23, 58, 59, 0, time.UTC),
		},
	}))

	t.Run("Unknown relative date markers", run(TestCase{
		data: &time.Time{},
		args: []string{
			"time=-1R",
		},
		error: `cannot unmarshal arg 'time=-1R': cannot set nested field for unmarshalable type time.Time`,
	}))

	t.Run("nested-basic", run(TestCase{
		args: []string{
			"basic.string=test",
		},
		expected: &Nested{
			Basic: Basic{
				String: "test",
			},
		},
	}))

	t.Run("map-basic", run(TestCase{
		args: []string{
			"map.key1=value1",
			"map.key2=value2",
		},
		expected: &Map{
			Map: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}))

	t.Run("custom", run(TestCase{
		args: []string{
			"custom-struct=test",
			"custom-string=test",
		},
		expected: &CustomWrapper{
			CustomStruct: &CustomStruct{
				value: "TEST",
			},
			CustomString: CustomString("TEST"),
		},
	}))

	t.Run("insane", run(TestCase{
		args: []string{
			"map.key1.key2.basic.string=test",
		},
		expected: func() any {
			n1 := &Nested{Basic: Basic{String: "test"}}
			m1 := &map[string]**Nested{"key2": &n1}
			m2 := map[string]**map[string]**Nested{"key1": &m1}

			return &Insane{Map: &m2}
		}(),
	}))

	t.Run("data-is-a-map", run(TestCase{
		args: []string{
			"key1=v1",
			"key2=v2",
		},
		expected: &map[string]string{
			"key1": "v1",
			"key2": "v2",
		},
		data: &map[string]string{},
	}))

	t.Run("IP", func(_ *testing.T) {
		ip := net.IPv4(1, 2, 3, 4)
		run(TestCase{
			args: []string{
				"ip=1.2.3.4",
			},
			expected: ip,
			data:     net.IPv4,
		})
	})

	t.Run("data-is-an-enum", run(TestCase{
		args: []string{
			"color=blue",
			"size=1",
		},
		expected: &Enum{Color: ColorBlue, Size: Size1},
	}))

	t.Run("data-is-raw-args", run(TestCase{
		args: []string{
			"pro.access_key",
			"access_key",
		},
		expected: &args.RawArgs{
			"pro.access_key",
			"access_key",
		},
	}))

	h := height(14)
	t.Run("height-set", run(TestCase{
		args: []string{
			"height=14cm",
		},
		data: &CustomArgs{},
		expected: &CustomArgs{
			Height: &h,
		},
	}))

	t.Run("height-not-set", run(TestCase{
		args:     []string{},
		data:     &CustomArgs{},
		expected: &CustomArgs{},
	}))

	t.Run("duplicate-keys-simple", run(TestCase{
		args: []string{
			"custom-struct=test",
			"custom-struct=test2",
			"custom-string=test",
		},
		data:  &CustomWrapper{},
		error: "cannot unmarshal arg 'custom-struct=test2': duplicate argument",
	}))

	t.Run("duplicate-keys-insane", run(TestCase{
		args: []string{
			"map.key1.key2.basic.string=test",
			"map.key1.key2.basic.string=test2",
		},
		data:  &Insane{},
		error: "cannot unmarshal arg 'map.key1.key2.basic.string=test2': duplicate argument",
	}))

	t.Run("anonymous-nested-field", run(TestCase{
		args: []string{
			"all=all",
			"merge1=1",
			"merge2=2",
			"merge-only=2",
		},
		expected: &Merge{
			Merge1: Merge1{
				All:       "",
				Merge1:    "1",
				MergeOnly: "",
			},
			Merge2: &Merge2{
				All:       "",
				Merge2:    "2",
				MergeOnly: "2",
			},
			All: "all",
		},
	}))

	t.Run("recursive-with-map-of-recursive-with-one-field-set", run(TestCase{
		args: []string{
			"name=coucou",
			"elements.0.name=bob",
			"elements.0.elements.plop.name=world",
		},
		expected: &RecursiveWithMapOfRecursive{
			Name: "coucou",
			Elements: map[string]*RecursiveWithMapOfRecursive{
				"0": {
					Name: "bob",
					Elements: map[string]*RecursiveWithMapOfRecursive{
						"plop": {
							Name: "world",
						},
					},
				},
			},
		},
	}))

	t.Run("recursive-with-map-of-recursive-with-multiple-fields-set", run(TestCase{
		args: []string{
			"name=coucou",
			"elements.0.id=1453",
			"elements.0.name=bob",
			"elements.0.elements.plop.name=world",
			"elements.0.elements.plop.short=long",
		},
		expected: &RecursiveWithMapOfRecursive{
			Name: "coucou",
			Elements: map[string]*RecursiveWithMapOfRecursive{
				"0": {
					ID:   1453,
					Name: "bob",
					Elements: map[string]*RecursiveWithMapOfRecursive{
						"plop": {
							Name:  "world",
							Short: "long",
						},
					},
				},
			},
		},
	}))

	t.Run("common-prefix-args", run(TestCase{
		args: []string{
			"ip=ip",
			"ipv6=ipv6",
		},
		expected: &SamePrefixArgName{
			IP:   "ip",
			IPv6: "ipv6",
		},
	}))

	t.Run("bool-without-equal", run(TestCase{
		args: []string{
			"bool",
		},
		data:  &Basic{},
		error: "cannot unmarshal arg 'bool': *bool is not unmarshalable: invalid boolean value",
	}))

	t.Run("bool-without-value", run(TestCase{
		args: []string{
			"bool=",
		},
		data:  &Basic{},
		error: "cannot unmarshal arg 'bool': *bool is not unmarshalable: invalid boolean value",
	}))

	t.Run("string-without-equal", run(TestCase{
		args: []string{
			"string",
		},
		expected: &Basic{},
	}))

	t.Run("string-without-value", run(TestCase{
		args: []string{
			"string=",
		},
		expected: &Basic{},
	}))

	t.Run("strings-without-equal", run(TestCase{
		args: []string{
			"strings",
		},
		data:  &Slice{},
		error: "cannot unmarshal arg 'strings': missing index on the array",
	}))

	t.Run("strings-without-value", run(TestCase{
		args: []string{
			"strings=",
		},
		data:  &Slice{},
		error: "cannot unmarshal arg 'strings': missing index on the array",
	}))
}

func TestIsUmarshalableValue(t *testing.T) {
	type TestCase struct {
		expected bool
		data     any
	}

	run := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			value := args.IsUmarshalableValue(testCase.data)
			assert.Equal(t, testCase.expected, value)
		}
	}

	args.RegisterUnmarshalFunc((*height)(nil), unmarshalHeight)

	strPtr := "This is a pointer"
	heightPtr := height(42)
	customStringPtr := CustomString("test")

	t.Run("string", run(TestCase{
		data:     "a simple string",
		expected: true,
	}))

	t.Run("int", run(TestCase{
		data:     42,
		expected: true,
	}))

	t.Run("custom", run(TestCase{
		data:     CustomString("CUSTOM-STRING"),
		expected: true,
	}))

	t.Run("nil", run(TestCase{
		data:     nil,
		expected: false,
	}))

	t.Run("custom-func", run(TestCase{
		data:     height(42),
		expected: true,
	}))

	t.Run("a-struct", run(TestCase{
		data:     &Basic{},
		expected: false,
	}))

	t.Run("str-pointer", run(TestCase{
		data:     &strPtr,
		expected: true,
	}))
	t.Run("custom-func-pointer", run(TestCase{
		data:     &heightPtr,
		expected: true,
	}))
	t.Run("custom-pointer", run(TestCase{
		data:     &customStringPtr,
		expected: true,
	}))
	t.Run("custom-pointer", run(TestCase{
		data:     map[string]string{},
		expected: false,
	}))

	t.Run("custom-pointer", run(TestCase{
		data:     []string{},
		expected: false,
	}))
}
