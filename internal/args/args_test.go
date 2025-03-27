package args_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Basic struct {
	String    string
	Int       int
	Int16     int16
	Int32     int32
	Int64     int64
	UInt16    uint16
	UInt32    uint32
	UInt64    uint64
	Float32   float32
	Float64   float64
	StringPtr *string
	Bool      bool
}

type Slice struct {
	Strings    []string
	StringsPtr []*string
	SlicePtr   *[]string
	Basics     []Basic
}

type WellKnownTypes struct {
	Size scw.Size
	Time time.Time
}

type Nested struct {
	Basic Basic
}

type Merge1 struct {
	All       string
	Merge1    string
	MergeOnly string
}

type Merge2 struct {
	All       string
	Merge2    string
	MergeOnly string
}

type Merge struct {
	Merge1
	*Merge2

	// override nested field
	All string
}

type Map struct {
	Map map[string]string
}

type Insane struct {
	Map *map[string]**map[string]**Nested
}

type CustomStruct struct {
	value string
}

type Color string

const ColorBlue = "blue"

type Size int

const Size1 = 1

type Enum struct {
	Color Color
	Size  Size
}

type RecursiveWithMapOfRecursive struct {
	ID       int
	Name     string
	Short    string
	Elements map[string]*RecursiveWithMapOfRecursive
}

func (c *CustomStruct) UnmarshalArgs(value string) error {
	c.value = strings.ToUpper(value)

	return nil
}

func (c *CustomStruct) MarshalArgs() (string, error) {
	return strings.ToLower(c.value), nil
}

type CustomString string

func (c *CustomString) UnmarshalArgs(value string) error {
	*c = CustomString(strings.ToUpper(value))

	return nil
}

func (c *CustomString) MarshalArgs() (string, error) {
	return strings.ToLower(string(*c)), nil
}

type CustomWrapper struct {
	CustomStruct *CustomStruct
	CustomString CustomString
}

type CustomArgs struct {
	Height *height
}

type height int

func marshalHeight(src interface{}) (string, error) {
	h := src.(*height)

	return fmt.Sprintf("%dcm", *h), nil
}

func unmarshalHeight(value string, dest interface{}) error {
	h := dest.(*height)
	_, err := fmt.Sscanf(value, "%dcm", h)

	return err
}

type SamePrefixArgName struct {
	IP   string
	IPv6 string
}

func TestRawArgs_GetAll(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		a := args.RawArgs{"ssh-keys.0=foo", "ssh-keys.1=bar"}
		assert.Equal(t, []string{"foo", "bar"}, a.GetAll("ssh-keys.{index}"))
	})

	t.Run("Insane", func(t *testing.T) {
		a := args.RawArgs{
			"countries.FR.cities.paris.street.vaugirard=pouet",
			"countries.FR.cities.paris.street.besbar=tati",
			"countries.FR.cities.nice.street.promenade=anglais",
			"countries.RU.cities.moscow.street.kremelin=rouge",
		}
		assert.Equal(
			t,
			[]string{"pouet", "tati", "anglais", "rouge"},
			a.GetAll("countries.{key}.cities.{key}.street.{key}"),
		)
	})
}

func TestGetArgType(t *testing.T) {
	type TestCase struct {
		ArgType       reflect.Type
		Name          string
		ExpectedKind  reflect.Kind
		expectedError string
	}

	run := func(tc *TestCase) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			res, err := args.GetArgType(tc.ArgType, tc.Name)
			if tc.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedKind, res.Kind())
			} else {
				require.Equal(t, tc.expectedError, err.Error())
			}
		}
	}

	t.Run("Simple", run(&TestCase{
		ArgType:      reflect.TypeOf(&Basic{}),
		Name:         "string",
		ExpectedKind: reflect.String,
	}))
	t.Run("Simple int", run(&TestCase{
		ArgType:      reflect.TypeOf(&Basic{}),
		Name:         "int-64",
		ExpectedKind: reflect.Int64,
	}))
	t.Run("Ptr", run(&TestCase{
		ArgType:      reflect.TypeOf(&Basic{}),
		Name:         "string-ptr",
		ExpectedKind: reflect.String,
	}))
	t.Run("simple slice", run(&TestCase{
		ArgType:      reflect.TypeOf(&Slice{}),
		Name:         "strings.{index}",
		ExpectedKind: reflect.String,
	}))
	t.Run("simple slice ptr", run(&TestCase{
		ArgType:      reflect.TypeOf(&Slice{}),
		Name:         "slice-ptr.{index}",
		ExpectedKind: reflect.String,
	}))
	t.Run("nested simple", run(&TestCase{
		ArgType:      reflect.TypeOf(&Nested{}),
		Name:         "basic.string",
		ExpectedKind: reflect.String,
	}))
	t.Run("merge simple", run(&TestCase{
		ArgType:      reflect.TypeOf(&Merge{}),
		Name:         "merge1",
		ExpectedKind: reflect.String,
	}))
	t.Run("merge simple all", run(&TestCase{
		ArgType:      reflect.TypeOf(&Merge{}),
		Name:         "all",
		ExpectedKind: reflect.String,
	}))
}
