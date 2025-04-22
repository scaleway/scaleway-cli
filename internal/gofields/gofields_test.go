package gofields_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/gofields"
	"github.com/stretchr/testify/assert"
)

type Friends struct {
	Name string
}

type Animal struct {
	Species string
}

type Pet struct {
	Animal
	Name string
}

type Address struct {
	Zip string

	// Private field
	digicode string
}

type User struct {
	Name    string
	Address *Address
	Friends []*Friends
	Pets    map[string]*Pet
}

func TestGetValue(t *testing.T) {
	type TestCase struct {
		Data     interface{}
		Path     string
		Expected interface{}
	}

	run := func(tc *TestCase) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			actual, err := gofields.GetValue(tc.Data, tc.Path)
			if err != nil {
				assert.Equal(t, tc.Expected, err.Error())
			} else {
				assert.Equal(t, tc.Expected, actual)
			}
		}
	}

	john := &User{
		Name: "John",
		Address: &Address{
			Zip:      "75008",
			digicode: "0000",
		},
		Friends: []*Friends{
			{
				Name: "Alice",
			},
		},
		Pets: map[string]*Pet{
			"rex": {
				Animal: Animal{
					Species: "dog",
				},
				Name: "Rex",
			},
		},
	}

	mark := &User{}

	t.Run("Simple", run(&TestCase{
		Data:     john,
		Path:     "Name",
		Expected: "John",
	}))

	t.Run("Nested", run(&TestCase{
		Data:     john,
		Path:     "Address.Zip",
		Expected: "75008",
	}))
	t.Run("Slice", run(&TestCase{
		Data:     john,
		Path:     "Friends.0.Name",
		Expected: "Alice",
	}))
	t.Run("Map", run(&TestCase{
		Data:     john,
		Path:     "Pets.rex.Name",
		Expected: "Rex",
	}))
	t.Run("UnknownField", run(&TestCase{
		Data:     john,
		Path:     "Unknown",
		Expected: "field Unknown does not exist in $",
	}))
	t.Run("UnknownNestedField", run(&TestCase{
		Data:     john,
		Path:     "Address.Unknown",
		Expected: "field Unknown does not exist in $.Address",
	}))
	t.Run("InvalidSliceIndex", run(&TestCase{
		Data:     john,
		Path:     "Friends.NotANumber",
		Expected: "trying to access array $.Friends but NotANumber is not a numerical index",
	}))
	t.Run("OutOfRangeSliceIndex", run(&TestCase{
		Data:     john,
		Path:     "Friends.13",
		Expected: "trying to access array $.Friends but 13 is out of range",
	}))
	t.Run("UnknownMapKey", run(&TestCase{
		Data:     john,
		Path:     "Pets.unknown",
		Expected: "trying to access map $.Pets but unknown key does not exist",
	}))
	t.Run("Nil", run(&TestCase{
		Data:     nil,
		Path:     "Pets",
		Expected: "field $ is nil",
	}))
	t.Run("NestedNil", run(&TestCase{
		Data:     mark,
		Path:     "Address.Zip",
		Expected: "field $.Address is nil",
	}))
	t.Run("AnonymousField", run(&TestCase{
		Data:     john,
		Path:     "Pets.rex.Species",
		Expected: "dog",
	}))
	t.Run("PrivateField", run(&TestCase{
		Data:     john,
		Path:     "Address.digicode",
		Expected: "field digicode is private in $.Address",
	}))
}

func TestGetType(t *testing.T) {
	type TestCase struct {
		Data     reflect.Type
		Path     string
		Expected interface{}
	}

	run := func(tc *TestCase) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			actual, err := gofields.GetType(tc.Data, tc.Path)
			if err != nil {
				assert.Equal(t, tc.Expected, err.Error())
			} else {
				assert.Equal(t, tc.Expected, actual)
			}
		}
	}

	t.Run("Simple", run(&TestCase{
		Data:     reflect.TypeOf(&User{}),
		Path:     "Name",
		Expected: reflect.TypeOf(""),
	}))

	t.Run("Nested", run(&TestCase{
		Data:     reflect.TypeOf(&User{}),
		Path:     "Address.Zip",
		Expected: reflect.TypeOf(""),
	}))
	t.Run("Slice", run(&TestCase{
		Data:     reflect.TypeOf(&User{}),
		Path:     "Friends.0.Name",
		Expected: reflect.TypeOf(""),
	}))
	t.Run("Map", run(&TestCase{
		Data:     reflect.TypeOf(&User{}),
		Path:     "Pets.rex.Name",
		Expected: reflect.TypeOf(""),
	}))
	t.Run("UnknownField", run(&TestCase{
		Data:     reflect.TypeOf(&User{}),
		Path:     "Unknown",
		Expected: "field Unknown does not exist in $",
	}))
	t.Run("UnknownNestedField", run(&TestCase{
		Data:     reflect.TypeOf(&User{}),
		Path:     "Address.Unknown",
		Expected: "field Unknown does not exist in $.Address",
	}))
	t.Run("InvalidSliceIndex", run(&TestCase{
		Data:     reflect.TypeOf(&User{}),
		Path:     "Friends.NotANumber",
		Expected: "trying to access array $.Friends but NotANumber is not a numerical index",
	}))
	t.Run("OutOfRangeSliceIndex", run(&TestCase{
		Data:     reflect.TypeOf(&User{}),
		Path:     "Friends.13.Name",
		Expected: reflect.TypeOf(""),
	}))
	t.Run("UnknownMapKey", run(&TestCase{
		Data:     reflect.TypeOf(&User{}),
		Path:     "Pets.unknown.Name",
		Expected: reflect.TypeOf(""),
	}))
	t.Run("AnonymousField", run(&TestCase{
		Data:     reflect.TypeOf(&User{}),
		Path:     "Pets.rex.Species",
		Expected: reflect.TypeOf(""),
	}))
	t.Run("PrivateField", run(&TestCase{
		Data:     reflect.TypeOf(&User{}),
		Path:     "Address.digicode",
		Expected: "field digicode is private in $.Address",
	}))
}

func TestListFields(t *testing.T) {
	type TestCase struct {
		Data     reflect.Type
		Expected []string
	}

	run := func(tc *TestCase) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			actual := gofields.ListFields(tc.Data)
			assert.Equal(t, tc.Expected, actual)
		}
	}

	t.Run("Simple", run(&TestCase{
		Data: reflect.TypeOf(&User{}),
		Expected: []string{
			"Name",
			"Address.Zip",
			"Friends.<index>.Name",
			"Pets.<key>.Species",
			"Pets.<key>.Name",
		},
	}))
}

func TestListFieldsWithFilter(t *testing.T) {
	type TestCase struct {
		Data     reflect.Type
		Expected []string
		Filter   gofields.ListFieldFilter
	}

	run := func(tc *TestCase) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			actual := gofields.ListFieldsWithFilter(tc.Data, tc.Filter)
			assert.Equal(t, tc.Expected, actual)
		}
	}

	t.Run("Simple", run(&TestCase{
		Data:     reflect.TypeOf(&User{}),
		Expected: []string{"Address.Zip"},
		Filter: func(_ reflect.Type, s string) bool {
			return strings.Contains(s, "Zip")
		},
	}))
}
