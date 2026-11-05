package humanv2_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	human "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/stretchr/testify/require"
)

func TestLayout_ResolveSections(t *testing.T) {
	for _, tc := range []struct {
		name   string
		layout human.Layout
		want   []human.Section
	}{
		{
			name:   "empty layout",
			layout: human.Layout{},
			want: []human.Section{
				{Path: "", Title: "", Omittable: false},
			},
		},
		{
			name: "non empty layout",
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "Foo"},
				},
			},
			want: []human.Section{
				{Path: "", Title: "", Omittable: false},
				{Path: "Foo", Title: "", Omittable: false},
			},
		},
		{
			name: "overridden default",
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "", Title: "Default"},
					{Path: "Foo"},
				},
			},
			want: []human.Section{
				{Path: "", Title: "Default", Omittable: false},
				{Path: "Foo", Title: "", Omittable: false},
			},
		},
		{
			name: "out-of-order default",
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "Foo"},
					{Path: ""},
				},
			},
			want: []human.Section{
				{Path: "Foo", Title: "", Omittable: false},
				{Path: "", Title: "", Omittable: false},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.layout.ResolveSections()
			require.Equal(t, tc.want, got)
		})
	}
}

type DummyStruct struct {
	Foo string
	Bar int
}

type UnexportedFieldStruct struct {
	Exported string
	hidden   string
}

type UnexportedFieldError struct {
	Message string
	hidden  string
}

func (e UnexportedFieldError) Error() string {
	return e.Message
}

type RootLevelDummyStruct struct {
	RootLevelFieldA   string
	RootLevelFieldB   string
	FirstLevelSection FirstLevelDummyStruct
}

type FirstLevelDummyStruct struct {
	FirstLevelFieldA string
	FirstLevelFieldB string
}

type ListDummyStruct struct {
	OmittableList []FirstLevelDummyStruct
	RequiredList  []FirstLevelDummyStruct
}

type RegistryDummyStruct struct {
	RegistryField   string
	RegistrySection FirstLevelDummyStruct
}

type ScalarSliceDummyStruct struct {
	Fields []string
}

type NestedDummyStruct struct {
	X string
	Y string
}

type RootDummyStruct struct {
	RootA string
	Sub   NestedDummyStruct
	RootB string
}

type RootSliceDummyStruct struct {
	RootA   string
	Strings []string
}

type RootMapDummyStruct struct {
	RootA  string
	Labels map[string]string
}

type NestedFormatterDummyStruct struct {
	RootA string
	Sub   FormatterStruct
}

type pointerRecvFormatter struct {
	Foo string
}

func (s *pointerRecvFormatter) FormatHuman() (string, error) {
	return "pointer-formatted", nil
}

type NestedPointerFormatterDummyStruct struct {
	RootA string
	Sub   *pointerRecvFormatter
}

type MapDummyStruct struct {
	Fields map[string]string
}

type MapListDummyStruct struct {
	Fields map[string][]string
}

func init() {
	// Registered once for the whole table; only the two registry cases use
	// RegistryDummyStruct, so no other case is affected.
	human.RegisterLayout(reflect.TypeFor[RegistryDummyStruct](), human.Layout{
		Sections: []human.Section{
			{Path: "RegistrySection", Title: "Custom"},
		},
	})
	human.RegisterFormatFunc(func(v RegistryFormatted) (string, error) {
		return "registered-" + v.Foo, nil
	})
}

func TestLayout_Resolve(t *testing.T) {
	for _, tc := range []struct {
		name    string
		data    any
		layout  human.Layout
		want    []human.Block
		wantErr string // substring of the expected error, empty to expect success
	}{
		{
			name: "simple struct",
			data: []DummyStruct{
				{Foo: "foo", Bar: 1},
				{Foo: "bar", Bar: 2},
			},
			want: []human.Block{
				{
					Header: []string{"FOO", "BAR"},
					Rows:   [][]string{{"foo", "1"}, {"bar", "2"}},
				},
			},
		},
		{
			name:   "simple struct without headers",
			data:   []DummyStruct{{Foo: "foo", Bar: 1}},
			layout: human.Layout{NoHeaders: true},
			want: []human.Block{
				{
					Header:    []string{"FOO", "BAR"},
					Rows:      [][]string{{"foo", "1"}},
					NoHeaders: true,
				},
			},
		},
		{
			name: "slice of pointer to struct",
			data: []*FirstLevelDummyStruct{
				{FirstLevelFieldA: "a", FirstLevelFieldB: "b"},
				{FirstLevelFieldA: "c", FirstLevelFieldB: "d"},
			},
			want: []human.Block{
				{
					Header: []string{"FIRST LEVEL FIELD A", "FIRST LEVEL FIELD B"},
					Rows:   [][]string{{"a", "b"}, {"c", "d"}},
				},
			},
		},
		{
			name: "slice with nil struct item",
			data: []*FirstLevelDummyStruct{
				nil,
				{FirstLevelFieldA: "c", FirstLevelFieldB: "d"},
			},
			want: []human.Block{
				{
					Header: []string{"FIRST LEVEL FIELD A", "FIRST LEVEL FIELD B"},
					Rows:   [][]string{{"-", "-"}, {"c", "d"}},
				},
			},
		},
		{
			name: "empty string cell formatted",
			data: []DummyStruct{
				{Foo: "", Bar: 1},
			},
			want: []human.Block{
				{
					Header: []string{"FOO", "BAR"},
					Rows:   [][]string{{"-", "1"}},
				},
			},
		},
		{
			name: "struct default layout",
			data: DummyStruct{Foo: "foo", Bar: 1},
			// no explicit section: only the default pathless section.
			layout: human.Layout{},
			want: []human.Block{
				{Rows: [][]string{{"Foo", "foo"}, {"Bar", "1"}}},
			},
		},
		{
			name:   "pointer to struct",
			data:   &DummyStruct{Foo: "foo", Bar: 1},
			layout: human.Layout{},
			want: []human.Block{
				{Rows: [][]string{{"Foo", "foo"}, {"Bar", "1"}}},
			},
		},
		{
			name:   "nil pointer",
			data:   (*DummyStruct)(nil),
			layout: human.Layout{},
			want:   []human.Block{{Rows: [][]string{{"-"}}}},
		},
		{
			name:   "nil",
			data:   nil,
			layout: human.Layout{},
			want:   []human.Block{{Rows: [][]string{{"-"}}}},
		},
		{
			name:   "multi-level pointer to struct",
			data:   new(new(&DummyStruct{Foo: "foo", Bar: 1})),
			layout: human.Layout{},
			want: []human.Block{
				{Rows: [][]string{{"Foo", "foo"}, {"Bar", "1"}}},
			},
		},
		{
			name:   "formatter value",
			data:   FormatterStruct{Foo: "ignored"},
			layout: human.Layout{},
			want:   []human.Block{{Rows: [][]string{{"cell-formatted"}}}},
		},
		{
			name:   "scalar",
			data:   42,
			layout: human.Layout{},
			want:   []human.Block{{Rows: [][]string{{"42"}}}},
		},
		{
			name: "struct with titled struct section",
			data: RootLevelDummyStruct{
				RootLevelFieldA: "AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA",
				RootLevelFieldB: "aaaaaaaa-aaaaaa-aaaaaa",
				FirstLevelSection: FirstLevelDummyStruct{
					FirstLevelFieldA: "BBBBBBBB-BBBB-BBBB-BBBB-BBBBBBBBBBBB",
					FirstLevelFieldB: "BBBBB",
				},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "FirstLevelSection", Title: "First Level"},
				},
			},
			want: []human.Block{
				{
					Rows: [][]string{
						{"RootLevelFieldA", "AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA"},
						{"RootLevelFieldB", "aaaaaaaa-aaaaaa-aaaaaa"},
					},
				},
				{
					Title: "First Level",
					Rows: [][]string{
						{"FirstLevelFieldA", "BBBBBBBB-BBBB-BBBB-BBBB-BBBBBBBBBBBB"},
						{"FirstLevelFieldB", "BBBBB"},
					},
				},
			},
		},
		{
			name: "omittable empty section",
			data: ListDummyStruct{},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "OmittableList", Title: "Omittable List", Omittable: true},
					{Path: "RequiredList", Title: "Required List", Omittable: false},
				},
			},
			// The omittable section is empty so it is dropped, the required one
			// keeps its header-only table.
			want: []human.Block{
				{Title: "Required List", Header: []string{"FIRST LEVEL FIELD A", "FIRST LEVEL FIELD B"}},
			},
		},
		{
			name: "omittable empty but allocated slice section",
			data: ListDummyStruct{
				OmittableList: []FirstLevelDummyStruct{},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "OmittableList", Title: "Omittable List", Omittable: true},
					{Path: "RequiredList", Title: "Required List", Omittable: false},
				},
			},
			// An empty-but-allocated slice is still empty for an omittable
			// section (encoding/json omitempty), so it is dropped like the nil one.
			want: []human.Block{
				{Title: "Required List", Header: []string{"FIRST LEVEL FIELD A", "FIRST LEVEL FIELD B"}},
			},
		},
		{
			name: "omittable non-empty section",
			data: ListDummyStruct{
				OmittableList: []FirstLevelDummyStruct{
					{FirstLevelFieldA: "x", FirstLevelFieldB: "y"},
				},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "OmittableList", Title: "Omittable List", Omittable: true},
					{Path: "RequiredList", Title: "Required List", Omittable: false},
				},
			},
			want: []human.Block{
				{
					Title:  "Omittable List",
					Header: []string{"FIRST LEVEL FIELD A", "FIRST LEVEL FIELD B"},
					Rows:   [][]string{{"x", "y"}},
				},
				{Title: "Required List", Header: []string{"FIRST LEVEL FIELD A", "FIRST LEVEL FIELD B"}},
			},
		},
		{
			name: "scalar slice section",
			data: ScalarSliceDummyStruct{
				Fields: []string{"poweron", "backup"},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "Fields", Title: "Allowed"},
				},
			},
			want: []human.Block{
				{Title: "Allowed", Rows: [][]string{{"2"}}},
			},
		},
		{
			name: "map section",
			data: MapDummyStruct{
				Fields: map[string]string{"b": "two", "a": "one"},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "Fields", Title: "Fields"},
				},
			},
			want: []human.Block{
				{Title: "Fields", Rows: [][]string{{"a", "one"}, {"b", "two"}}},
			},
		},
		{
			name: "map section with slice value",
			data: MapListDummyStruct{
				Fields: map[string][]string{"b": {"y", "x"}, "a": {"one"}},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "Fields", Title: "Fields"},
				},
			},
			want: []human.Block{
				{Title: "Fields", Rows: [][]string{{"a.0", "one"}, {"b.0", "y"}, {"b.1", "x"}}},
			},
		},
		{
			name: "flattened struct fields",
			data: RootDummyStruct{
				RootA: "aaa",
				Sub:   NestedDummyStruct{X: "one", Y: "two"},
				RootB: "bbb",
			},
			layout: human.Layout{},
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}, {"Sub.X", "one"}, {"Sub.Y", "two"}, {"RootB", "bbb"}}},
			},
		},
		{
			name: "omittable skipped between two sections",
			data: RootDummyStruct{
				RootA: "aaa",
				// Sub zero -> omitted
				RootB: "bbb",
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "Sub", Title: "B", Omittable: true},
					{Path: "RootB", Title: "C"},
				},
			},
			// The skipped omittable section leaves no block at all.
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}}},
				{Title: "C", Rows: [][]string{{"bbb"}}},
			},
		},
		{
			name: "section footer drawn below the block",
			data: RootDummyStruct{
				RootA: "aaa",
				Sub:   NestedDummyStruct{X: "one", Y: "two"},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "RootA", Title: "A", Footer: "WARNING: shown only once"},
				},
			},
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}, {"Sub.X", "one"}, {"Sub.Y", "two"}}},
				{Title: "A", Rows: [][]string{{"aaa"}}, Footer: "WARNING: shown only once"},
			},
		},
		{
			name: "section footerfunc computed from data",
			data: RootDummyStruct{
				RootA: "aaa",
				RootB: "bbb",
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "RootA", Title: "A", FooterFunc: func(d any) (string, error) {
						return "RootB is " + d.(RootDummyStruct).RootB, nil
					}},
				},
			},
			want: []human.Block{
				{Rows: [][]string{{"RootB", "bbb"}}},
				{Title: "A", Rows: [][]string{{"aaa"}}, Footer: "RootB is bbb"},
			},
		},
		{
			name: "section footerfunc empty string draws nothing",
			data: RootDummyStruct{
				RootA: "aaa",
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "RootA", Title: "A", FooterFunc: func(any) (string, error) { return "", nil }},
				},
			},
			want: []human.Block{
				{Title: "A", Rows: [][]string{{"aaa"}}},
			},
		},
		{
			name: "dotted section path",
			data: RootDummyStruct{
				RootA: "aaa",
				Sub:   NestedDummyStruct{X: "one", Y: "two"},
				RootB: "bbb",
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "Sub.Y", Title: "Y"},
				},
			},
			// Sub.Y is claimed by the section, Sub.X keeps flattening into
			// the main block.
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}, {"Sub.X", "one"}, {"RootB", "bbb"}}},
				{Title: "Y", Rows: [][]string{{"two"}}},
			},
		},
		{
			name: "section derives title from path",
			data: RootDummyStruct{
				RootA: "aaa",
				RootB: "bbb",
				Sub:   NestedDummyStruct{X: "one", Y: "two"},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "RootB"}, // no Title: derived like v1
				},
			},
			// Layout.ResolveSections prepends the main block; only the selected
			// path is claimed, the rest flattens into it.
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}, {"Sub.X", "one"}, {"Sub.Y", "two"}}},
				{Title: "Root B", Rows: [][]string{{"bbb"}}},
			},
		},
		{
			name: "section derives dotted title from path",
			data: RootDummyStruct{
				RootA: "aaa",
				RootB: "bbb",
				Sub:   NestedDummyStruct{X: "one", Y: "two"},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "Sub.X"},
				},
			},
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}, {"Sub.Y", "two"}, {"RootB", "bbb"}}},
				{Title: "Sub - X", Rows: [][]string{{"one"}}},
			},
		},
		{
			name: "nested formatter struct rendered inline",
			data: NestedFormatterDummyStruct{
				RootA: "aaa",
				Sub:   FormatterStruct{Foo: "ignored"},
			},
			layout: human.Layout{},
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}, {"Sub", "cell-formatted"}}},
			},
		},
		{
			name: "nested pointer-receiver formatter rendered inline",
			data: NestedPointerFormatterDummyStruct{
				RootA: "aaa",
				Sub:   &pointerRecvFormatter{Foo: "x"},
			},
			layout: human.Layout{},
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}, {"Sub", "pointer-formatted"}}},
			},
		},
		{
			name: "dotted section path for a whole struct",
			data: RootDummyStruct{
				RootA: "aaa",
				Sub:   NestedDummyStruct{X: "one", Y: "two"},
				RootB: "bbb",
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "Sub", Title: "Sub"},
				},
			},
			// Claiming the whole Sub struct also prunes its leaves from the main
			// block, not just the exact dotted path.
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}, {"RootB", "bbb"}}},
				{Title: "Sub", Rows: [][]string{{"X", "one"}, {"Y", "two"}}},
			},
		},
		{
			name: "flattened slice fields",
			data: RootSliceDummyStruct{
				RootA:   "aaa",
				Strings: []string{"s1", "s2"},
			},
			layout: human.Layout{},
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}, {"Strings.0", "s1"}, {"Strings.1", "s2"}}},
			},
		},
		{
			name: "flattened map fields",
			data: RootMapDummyStruct{
				RootA:  "aaa",
				Labels: map[string]string{"b": "two", "a": "one"},
			},
			layout: human.Layout{},
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}, {"Labels.a", "one"}, {"Labels.b", "two"}}},
			},
		},
		{
			name: "nil composite field dropped",
			data: RootSliceDummyStruct{
				RootA: "aaa",
			},
			layout: human.Layout{},
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}}},
			},
		},
		{
			name: "flat table with fields",
			data: []DummyStruct{
				{Foo: "f", Bar: 1},
				{Foo: "b", Bar: 2},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Fields: []human.Field{{Path: "Bar"}}},
				},
			},
			want: []human.Block{
				{Header: []string{"BAR"}, Rows: [][]string{{"1"}, {"2"}}},
			},
		},
		{
			name: "table column with explicit label",
			data: []DummyStruct{
				{Foo: "f", Bar: 1},
				{Foo: "b", Bar: 2},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Fields: []human.Field{{Path: "Bar", Label: "Size"}}},
				},
			},
			want: []human.Block{
				{Header: []string{"Size"}, Rows: [][]string{{"1"}, {"2"}}},
			},
		},
		{
			name: "table column with dotted path",
			data: []RootDummyStruct{
				{RootA: "a", Sub: NestedDummyStruct{X: "x", Y: "y"}},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Fields: []human.Field{{Path: "RootA"}, {Path: "Sub.X"}}},
				},
			},
			want: []human.Block{
				{Header: []string{"ROOT A", "SUB X"}, Rows: [][]string{{"a", "x"}}},
			},
		},
		{
			name: "table column with nil intermediate",
			data: []NestedPointerFormatterDummyStruct{
				{RootA: "a"},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Fields: []human.Field{{Path: "Sub.Foo"}}},
				},
			},
			// The path is valid, but this item's Sub is nil: the cell degrades
			// to "-" instead of failing the whole table, like v1.
			want: []human.Block{
				{Header: []string{"SUB FOO"}, Rows: [][]string{{"-"}}},
			},
		},
		{
			name: "struct lines filtered",
			data: RootDummyStruct{
				RootA: "aaa",
				Sub:   NestedDummyStruct{X: "one", Y: "two"},
				RootB: "bbb",
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Fields: []human.Field{{Path: "RootA"}, {Path: "Sub.X"}}},
				},
			},
			want: []human.Block{
				{Rows: [][]string{{"RootA", "aaa"}, {"Sub.X", "one"}}},
			},
		},
		{
			name: "selected composite subtree flattens",
			data: RootSliceDummyStruct{
				RootA:   "a",
				Strings: []string{"s1", "s2"},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Fields: []human.Field{{Path: "Strings"}}},
				},
			},
			want: []human.Block{
				{Rows: [][]string{{"Strings.0", "s1"}, {"Strings.1", "s2"}}},
			},
		},
		{
			name: "explicit layout overrides registry",
			data: RegistryDummyStruct{
				RegistryField:   "value",
				RegistrySection: FirstLevelDummyStruct{FirstLevelFieldA: "aaa", FirstLevelFieldB: "bbb"},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "RegistrySection", Title: "Explicit"},
				},
			},
			want: []human.Block{
				{Rows: [][]string{{"RegistryField", "value"}}},
				{Title: "Explicit", Rows: [][]string{{"FirstLevelFieldA", "aaa"}, {"FirstLevelFieldB", "bbb"}}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.layout.Resolve(tc.data)

			if tc.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestResolveUnknownField(t *testing.T) {
	for _, tc := range []struct {
		name   string
		data   any
		layout human.Layout
	}{
		{
			name: "table column",
			data: []DummyStruct{{Foo: "f", Bar: 1}},
			layout: human.Layout{
				Sections: []human.Section{
					{Fields: []human.Field{{Path: "Nope"}}},
				},
			},
		},
		{
			name: "struct line",
			data: DummyStruct{Foo: "f", Bar: 1},
			layout: human.Layout{
				Sections: []human.Section{
					{Fields: []human.Field{{Path: "Nope"}}},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.layout.Resolve(tc.data)
			require.Error(t, err)

			unknownFieldErr := &human.UnknownFieldError{}
			require.ErrorAs(t, err, &unknownFieldErr)
			require.Equal(t, "Nope", unknownFieldErr.FieldName)
			require.Contains(t, unknownFieldErr.ValidFields, "Foo")
			require.Contains(t, unknownFieldErr.ValidFields, "Bar")
			require.NotContains(t, unknownFieldErr.ValidFields, "Nope")
		})
	}
}

func TestBlock_Draw(t *testing.T) {
	for _, tc := range []struct {
		name  string
		block human.Block
		width int
		want  string
	}{
		{
			name: "table",
			block: human.Block{
				Header: []string{"Foo", "Bar"},
				Rows:   [][]string{{"foo", "1"}, {"bar", "2"}},
			},
			want: `Foo  Bar
foo  1
bar  2`,
		},
		{
			name: "key value lines",
			block: human.Block{
				Rows: [][]string{{"Foo", "foo"}, {"Bar", "1"}},
			},
			want: `Foo  foo
Bar  1`,
		},
		{
			name: "title and table",
			block: human.Block{
				Title:  "Required List",
				Header: []string{"FirstLevelFieldA", "FirstLevelFieldB"},
			},
			want: `Required List:
FirstLevelFieldA  FirstLevelFieldB`,
		},
		{
			name: "title and lines",
			block: human.Block{
				Title: "Fields",
				Rows:  [][]string{{"a", "one"}, {"b", "two"}},
			},
			want: `Fields:
a    one
b    two`,
		},
		{
			name: "single cell",
			block: human.Block{
				Rows: [][]string{{"-"}},
			},
			want: `-`,
		},
		{
			name: "table without headers",
			block: human.Block{
				Header:    []string{"Foo", "Bar"},
				Rows:      [][]string{{"foo", "1"}, {"bar", "2"}},
				NoHeaders: true,
			},
			want: `foo  1
bar  2`,
		},
		{
			name: "table shrinks to fit width",
			block: human.Block{
				Header: []string{"ID", "NAME", "TYPE"},
				Rows:   [][]string{{"1", "xavier", "instance"}},
			},
			// Columns beyond the second drop once the running column width
			// exceeds the given width; tabwriter pads 2-char cells to 5.
			want: `ID   NAME
1    xavier`,
			width: 12,
		},
		{
			name: "no-headers table shrinks data rows",
			block: human.Block{
				Header:    []string{"ID", "NAME", "TYPE"},
				Rows:      [][]string{{"1", "xavier", "instance"}},
				NoHeaders: true,
			},
			want:  `1    xavier`,
			width: 12,
		},
		{
			name:  "empty block draws nothing",
			block: human.Block{},
			want:  ``,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, tc.block.Draw(tc.width, 0))
		})
	}
}

func TestRender(t *testing.T) {
	for _, tc := range []struct {
		name   string
		data   any
		layout human.Layout
		want   string
	}{
		{
			name:   "nil",
			data:   nil,
			layout: human.Layout{},
			want:   "-",
		},
		{
			name:   "nil pointer",
			data:   (*DummyStruct)(nil),
			layout: human.Layout{},
			want:   "-",
		},
		{
			name:   "pointer to struct",
			data:   &DummyStruct{Foo: "foo", Bar: 1},
			layout: human.Layout{},
			want: `Foo  foo
Bar  1`,
		},
		{
			name:   "multi-level pointer to struct",
			data:   new(new(&DummyStruct{Foo: "foo", Bar: 1})),
			layout: human.Layout{},
			want: `Foo  foo
Bar  1`,
		},
		{
			name:   "scalar",
			data:   42,
			layout: human.Layout{},
			want:   "42",
		},
		{
			name: "table without headers",
			data: []DummyStruct{
				{Foo: "foo", Bar: 1},
				{Foo: "bar", Bar: 2},
			},
			layout: human.Layout{NoHeaders: true},
			want: `foo  1
bar  2`,
		},
		{
			name: "table",
			data: []DummyStruct{
				{Foo: "foo", Bar: 1},
				{Foo: "bar", Bar: 2},
			},
			layout: human.Layout{},
			want: `FOO  BAR
foo  1
bar  2`,
		},
		{
			name:   "plain struct",
			data:   DummyStruct{Foo: "foo", Bar: 1},
			layout: human.Layout{},
			want: `Foo  foo
Bar  1`,
		},
		{
			name: "registry layout fallback",
			data: RegistryDummyStruct{
				RegistryField:   "value",
				RegistrySection: FirstLevelDummyStruct{FirstLevelFieldA: "aaa", FirstLevelFieldB: "bbb"},
			},
			// Zero layout: Render falls back to the one registered for the type.
			layout: human.Layout{},
			want: `RegistryField  value

Custom:
FirstLevelFieldA  aaa
FirstLevelFieldB  bbb`,
		},
		{
			name: "joined sections",
			data: RootDummyStruct{
				RootA: "aaa",
				Sub:   NestedDummyStruct{X: "one", Y: "two"},
				RootB: "bbb",
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "Sub", Title: "Sub"},
				},
			},
			want: `RootA  aaa
RootB  bbb

Sub:
X    one
Y    two`,
		},
		{
			name: "scalar sections joined with blank line",
			data: DummyStruct{Foo: "foo", Bar: 1},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "Foo", Title: "A"},
					{Path: "Bar", Title: "B"},
				},
			},
			// v1 separates consecutive scalar sections with a blank line too.
			want: `A:
foo

B:
1`,
		},
		{
			name: "omittable skipped between two sections",
			data: RootDummyStruct{
				RootA: "aaa",
				// Sub zero -> omitted
				RootB: "bbb",
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Path: "Sub", Title: "B", Omittable: true},
					{Path: "RootB", Title: "C"},
				},
			},
			// The skipped omittable section leaves no orphan blank line.
			want: `RootA  aaa

C:
bbb`,
		},
		{
			name: "layout title before main section",
			data: DummyStruct{Foo: "foo", Bar: 1},
			layout: human.Layout{
				Title: "Policy",
			},
			want: `Policy:
Foo  foo
Bar  1`,
		},
		{
			name: "layout title wraps sections",
			data: RootDummyStruct{
				RootA: "aaa",
				Sub:   NestedDummyStruct{X: "one", Y: "two"},
				RootB: "bbb",
			},
			layout: human.Layout{
				Title: "Details",
				Sections: []human.Section{
					{Path: "Sub", Title: "Sub"},
				},
			},
			want: `Details:
RootA  aaa
RootB  bbb

Sub:
X    one
Y    two`,
		},
		{
			name: "layout title wraps leaf",
			data: "text",
			layout: human.Layout{
				Title: "Details",
			},
			want: `Details:
text`,
		},
		{
			// A registered struct escapes through Format at the top level instead
			// of flattening into its raw fields.
			name: "registered format func top-level",
			data: RegistryFormatted{Foo: "x"},
			want: "registered-x",
		},
		{
			// Nested as a struct field, it formats in place rather than
			// flattening.
			name: "registered format func nested",
			data: RegistryNestedStruct{Root: "r", Sub: RegistryFormatted{Foo: "s"}},
			want: "Root  r\nSub   registered-s",
		},
		{
			// Unexported fields are skipped by walk: reflect panics if a visitor
			// calls Interface() on them.
			name: "unexported struct field skipped",
			data: UnexportedFieldStruct{Exported: "v", hidden: "secret"},
			want: "Exported  v",
		},
		{
			// Same for error values: an error struct is a leaf rendering its
			// message, so unexported fields aren't the point here anymore.
			name: "error struct with unexported field",
			data: UnexportedFieldError{Message: "boom", hidden: "secret"},
			want: "boom",
		},
		{
			// A plain error renders its message through the error case.
			name: "error message",
			data: errors.New("boom"),
			want: "boom",
		},
		{
			// Wrapped errors keep the chain in the message.
			name: "wrapped error message",
			data: fmt.Errorf("could not open file: %w", errors.New("no such file")),
			want: "could not open file: no such file",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := human.Render(tc.data, tc.layout)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestMaxCols(t *testing.T) {
	grid := [][]string{
		{"ID", "NAME", "TYPE", "STATE", "ZONE", "PUBLIC IP"},
		{"1", "xavier", "instance", "running", "fr-par-1", "1.2.3.4"},
	}
	for _, tc := range []struct {
		width      int
		colPadding int
		want       int
	}{
		{width: 3, colPadding: 2, want: 0},
		{width: 30, colPadding: 2, want: 3},
		{width: 40, colPadding: 2, want: 4},
		{width: 51, colPadding: 2, want: 5},
		{width: 60, colPadding: 2, want: 6},
		// colPadding is part of each column's estimated width: without it the
		// same width fits every column.
		{width: 40, colPadding: 0, want: 6},
	} {
		require.Equal(t, tc.want, human.MaxCols(grid, tc.width, tc.colPadding))
	}
}
