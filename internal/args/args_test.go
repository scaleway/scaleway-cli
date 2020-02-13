package args

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/scw"
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
