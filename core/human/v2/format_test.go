package human_test

import (
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	human "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

type FormatterStruct struct {
	Foo string
}

func (FormatterStruct) FormatHuman() (string, error) {
	return "cell-formatted", nil
}

type StringerStruct struct{}

func (StringerStruct) String() string {
	return "stringed"
}

type ErrorStruct struct{}

func (ErrorStruct) Error() string {
	return "something broke"
}

func TestFormat(t *testing.T) {
	date := time.Date(1990, 11, 17, 20, 20, 0, 0, time.UTC)
	humanDate := humanize.Time(date)

	for _, tc := range []struct {
		name    string
		data    any
		want    string
		wantErr string // substring of the expected error, empty to expect success
	}{
		{
			name: "formatter value",
			data: FormatterStruct{Foo: "ignored"},
			want: "cell-formatted",
		},
		{
			name: "formatter pointer",
			data: &FormatterStruct{Foo: "ignored"},
			want: "cell-formatted",
		},
		{
			name: "formatter multi-pointer",
			data: new(new(&FormatterStruct{Foo: "ignored"})),
			want: "cell-formatted",
		},
		{
			name: "string",
			data: "hello",
			want: "hello",
		},
		{
			name: "integer",
			data: 42,
			want: "42",
		},
		{
			name: "boolean",
			data: true,
			want: "true",
		},
		{
			name: "empty string",
			data: "",
			want: "-",
		},
		{
			name: "nil",
			data: nil,
			want: "-",
		},
		{
			name: "stringer value",
			data: StringerStruct{},
			want: "stringed",
		},
		{
			name: "stringer pointer",
			data: &StringerStruct{},
			want: "stringed",
		},
		{
			name: "error",
			data: ErrorStruct{},
			want: "something broke",
		},
		{
			name: "pointer to integer",
			data: new(int(42)),
			want: "42",
		},
		{
			name: "nil pointer",
			data: (*int)(nil),
			want: "-",
		},
		{
			name: "float",
			data: 3.14,
			want: "3.14",
		},
		{
			name: "time",
			data: date,
			want: humanDate,
		},
		{
			name: "nil time pointer",
			data: (*time.Time)(nil),
			want: "-",
		},
		{
			name: "size",
			data: scw.Size(13200),
			want: "13 kB",
		},
		{
			name: "size IEC notation",
			data: scw.Size(52428800),
			want: "50 MiB",
		},
		{
			name: "nil size pointer",
			data: (*scw.Size)(nil),
			want: "-",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := human.Format(tc.data)

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
