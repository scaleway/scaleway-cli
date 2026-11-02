package human_test

import (
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
			name: "struct default layout",
			data: DummyStruct{Foo: "foo", Bar: 1},
			// no explicit section: only the default pathless section.
			layout: human.Layout{},
			want: `Foo  foo
Bar  1`,
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
