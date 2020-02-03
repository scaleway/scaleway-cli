package human

import (
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/dustin/go-humanize"
)

type Struct struct {
	String      string
	Int         int
	Bool        bool
	Strings     []string
	Time        time.Time
	Struct      *Struct
	Nil         *Struct
	Structs     []*Struct
	Map         map[string]string
	Stringer    Stringer
	StringerPtr *Stringer
}

type Address struct {
	Street string
	City   string
}

type Acquaintance struct {
	Name string
	Link string
}

type Human struct {
	Name          string
	Age           int
	Address       *Address
	Acquaintances []*Acquaintance
}

type Stringer struct{}

func (s Stringer) String() string {
	return "a stringer"
}

func TestMarshal(t *testing.T) {
	type testCase struct {
		data   interface{}
		opt    *MarshalOpt
		result string
		err    error
	}

	run := func(tc *testCase) func(*testing.T) {
		return func(t *testing.T) {
			result, err := Marshal(tc.data, tc.opt)

			// Format expected to allow indentation when writing test
			expected := tc.result
			expected = strings.Replace(expected, "\t", "", -1)
			expected = strings.Trim(expected, "\n")

			if tc.result != "" {
				assert.NoError(t, err)
				assert.Equal(t, expected, result)
			} else {
				assert.Equal(t, err, err)
			}
		}
	}

	date := time.Date(1990, 11, 17, 20, 20, 0, 0, time.UTC)
	humanDate := humanize.Time(date)
	t.Run("struct", run(&testCase{
		data: &Struct{
			String:  "This is a string",
			Int:     42,
			Bool:    true,
			Strings: []string{"s1", "s2"},
			Time:    date,
			Struct:  &Struct{},
			Nil:     nil,
			Structs: []*Struct{
				{
					String: "Nested string",
				},
			},
			Map: map[string]string{
				"key1": "v1",
				"key2": "v2",
			},
			Stringer:    Stringer{},
			StringerPtr: &Stringer{},
		},
		opt: nil,
		result: `
			string              This is a string
			int                 42
			bool                true
			strings.0           s1
			strings.1           s2
			time                ` + humanDate + `
			struct.string       
			struct.int          0
			struct.bool         false
			struct.time         a long while ago
			struct.stringer     a stringer
			structs.0.string    Nested string
			structs.0.int       0
			structs.0.bool      false
			structs.0.time      a long while ago
			structs.0.stringer  a stringer
			map.key1            v1
			map.key2            v2
			stringer            a stringer
			stringer-ptr        a stringer
		`,
		err: nil,
	}))

	t.Run("struct", run(&testCase{
		data: &Human{
			Name:    "Sherlock Holmes",
			Age:     42,
			Address: &Address{Street: "221b Baker St", City: "London"},
			Acquaintances: []*Acquaintance{
				{Name: "Dr watson", Link: "Assistant"},
				{Name: "Mrs. Hudson", Link: "Landlady"},
			},
		},
		opt: &MarshalOpt{
			Title: "Personal Information",
			Sections: []*MarshalSection{
				{FieldName: "address"},
				{Title: "Relationship", FieldName: "acquaintances"},
			},
		},
		result: `
			Personal Information:
			name  Sherlock Holmes
			age   42
			
			Address:
			street  221b Baker St
			city    London
			
			Relationship:
			NAME         LINK
			Dr watson    Assistant
			Mrs. Hudson  Landlady
		`,
		err: nil,
	}))
}
