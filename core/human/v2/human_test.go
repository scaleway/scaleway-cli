package human_test

import (
	"reflect"
	"testing"

	human "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/stretchr/testify/require"
)

func TestLayout_resolveSections(t *testing.T) {
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
	human.SetLayout(reflect.TypeFor[RegistryDummyStruct](), human.Layout{
		Sections: []human.Section{
			{Path: "RegistrySection", Title: "Custom"},
		},
	})
}

func TestRender(t *testing.T) {
	for _, tc := range []struct {
		name    string
		data    any
		layout  human.Layout
		want    string
		wantErr string // substring of the expected error, empty to expect success
	}{
		{
			name: "simple struct",
			data: []DummyStruct{
				{Foo: "foo", Bar: 1},
				{Foo: "bar", Bar: 2},
			},
			want: `Foo  Bar
foo  1
bar  2
`,
		},
		{
			name: "slice of pointer to struct",
			data: []*FirstLevelDummyStruct{
				{FirstLevelFieldA: "a", FirstLevelFieldB: "b"},
				{FirstLevelFieldA: "c", FirstLevelFieldB: "d"},
			},
			want: `FirstLevelFieldA  FirstLevelFieldB
a                 b
c                 d
`,
		},
		{
			name: "slice with nil struct item",
			data: []*FirstLevelDummyStruct{
				nil,
				{FirstLevelFieldA: "c", FirstLevelFieldB: "d"},
			},
			want: `FirstLevelFieldA  FirstLevelFieldB
-                 -
c                 d
`,
		},
		{
			name: "empty string cell formatted",
			data: []DummyStruct{
				{Foo: "", Bar: 1},
			},
			want: `Foo  Bar
-    1
`,
		},
		{
			name: "struct default layout",
			data: DummyStruct{Foo: "foo", Bar: 1},
			// no explicit section: only the default pathless section.
			layout: human.Layout{},
			want: `Foo  foo
Bar  1`,
		},
		{
			name:   "pointer to struct",
			data:   &DummyStruct{Foo: "foo", Bar: 1},
			layout: human.Layout{},
			want: `Foo  foo
Bar  1`,
		},
		{
			name:   "nil pointer",
			data:   (*DummyStruct)(nil),
			layout: human.Layout{},
			want:   `-`,
		},
		{
			name:   "nil",
			data:   nil,
			layout: human.Layout{},
			want:   `-`,
		},
		{
			name:   "multi-level pointer to struct",
			data:   new(new(&DummyStruct{Foo: "foo", Bar: 1})),
			layout: human.Layout{},
			want: `Foo  foo
Bar  1`,
		},
		{
			name:   "formatter value",
			data:   FormatterStruct{Foo: "ignored"},
			layout: human.Layout{},
			want:   "cell-formatted",
		},
		{
			name:   "scalar",
			data:   42,
			layout: human.Layout{},
			want:   "42",
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
			want: `RootLevelFieldA  AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA
RootLevelFieldB  aaaaaaaa-aaaaaa-aaaaaa

First Level:
FirstLevelFieldA  BBBBBBBB-BBBB-BBBB-BBBB-BBBBBBBBBBBB
FirstLevelFieldB  BBBBB`,
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
			// The omittable section is empty so it is skipped, the required one
			// keeps its header-only table.
			want: `Required List:
FirstLevelFieldA  FirstLevelFieldB`,
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
			want: `Omittable List:
FirstLevelFieldA  FirstLevelFieldB
x                 y

Required List:
FirstLevelFieldA  FirstLevelFieldB`,
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
			want: `Allowed:
[poweron backup]`,
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
			want: `Fields:
a    one
b    two`,
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
			want: `Fields:
a.0  one
b.0  y
b.1  x`,
		},
		{
			name: "flattened struct fields",
			data: RootDummyStruct{
				RootA: "aaa",
				Sub:   NestedDummyStruct{X: "one", Y: "two"},
				RootB: "bbb",
			},
			layout: human.Layout{},
			want: `RootA  aaa
Sub.X  one
Sub.Y  two
RootB  bbb`,
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
			// The skipped omittable section must not leave an orphan blank line
			// between the two rendered ones.
			want: `RootA  aaa

C:
bbb`,
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
			want: `RootA  aaa
Sub.X  one
RootB  bbb

Y:
two`,
		},
		{
			name: "nested formatter struct rendered inline",
			data: NestedFormatterDummyStruct{
				RootA: "aaa",
				Sub:   FormatterStruct{Foo: "ignored"},
			},
			layout: human.Layout{},
			want: `RootA  aaa
Sub    cell-formatted`,
		},
		{
			name: "nested pointer-receiver formatter rendered inline",
			data: NestedPointerFormatterDummyStruct{
				RootA: "aaa",
				Sub:   &pointerRecvFormatter{Foo: "x"},
			},
			layout: human.Layout{},
			want: `RootA  aaa
Sub    pointer-formatted`,
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
			want: `RootA  aaa
RootB  bbb

Sub:
X    one
Y    two`,
		},
		{
			name: "flattened slice fields",
			data: RootSliceDummyStruct{
				RootA:   "aaa",
				Strings: []string{"s1", "s2"},
			},
			layout: human.Layout{},
			want: `RootA      aaa
Strings.0  s1
Strings.1  s2`,
		},
		{
			name: "flattened map fields",
			data: RootMapDummyStruct{
				RootA:  "aaa",
				Labels: map[string]string{"b": "two", "a": "one"},
			},
			layout: human.Layout{},
			want: `RootA     aaa
Labels.a  one
Labels.b  two`,
		},
		{
			name: "nil composite field dropped",
			data: RootSliceDummyStruct{
				RootA: "aaa",
			},
			layout: human.Layout{},
			want:   `RootA  aaa`,
		},
		{
			name: "registry layout",
			data: RegistryDummyStruct{
				RegistryField:   "value",
				RegistrySection: FirstLevelDummyStruct{FirstLevelFieldA: "aaa", FirstLevelFieldB: "bbb"},
			},
			// No layout given: the one registered for the type applies.
			layout: human.Layout{},
			want: `RegistryField  value

Custom:
FirstLevelFieldA  aaa
FirstLevelFieldB  bbb`,
		},
		{
			name: "flat table with fields",
			data: []DummyStruct{
				{Foo: "f", Bar: 1},
				{Foo: "b", Bar: 2},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Fields: []string{"Bar"}},
				},
			},
			want: `Bar
1
2
`,
		},
		{
			name: "table column with dotted path",
			data: []RootDummyStruct{
				{RootA: "a", Sub: NestedDummyStruct{X: "x", Y: "y"}},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Fields: []string{"RootA", "Sub.X"}},
				},
			},
			want: `RootA  Sub.X
a      x
`,
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
					{Fields: []string{"RootA", "Sub.X"}},
				},
			},
			want: `RootA  aaa
Sub.X  one`,
		},
		{
			name: "selected composite subtree flattens",
			data: RootSliceDummyStruct{
				RootA:   "a",
				Strings: []string{"s1", "s2"},
			},
			layout: human.Layout{
				Sections: []human.Section{
					{Fields: []string{"Strings"}},
				},
			},
			want: `Strings.0  s1
Strings.1  s2`,
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
			want: `RegistryField  value

Explicit:
FirstLevelFieldA  aaa
FirstLevelFieldB  bbb`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := human.Render(tc.data, tc.layout)

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
