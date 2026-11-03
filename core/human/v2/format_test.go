package human_test

import (
	"net"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	version "github.com/hashicorp/go-version"
	human "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

// makeIPNet parses a CIDR into a scw.IPNet for tests.
func makeIPNet(cidr string) scw.IPNet {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}

	return scw.IPNet{IPNet: *ipNet}
}

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

// Registry test types are distinct per test: the registry is global, so
// sharing a type between tests would couple their run order.

type RegistryFormatted struct {
	Foo string
}

type RegistryPointerFormatted struct {
	Foo string
}

type RegistryPrecedenceFormatted struct{}

func (RegistryPrecedenceFormatted) FormatHuman() (string, error) {
	return "interface-formatted", nil
}

type RegistryNestedStruct struct {
	Root string
	Sub  RegistryFormatted
}

func TestFormatRegistered(t *testing.T) {
	human.RegisterFormatFunc[RegistryFormatted](func(v RegistryFormatted) (string, error) {
		return "registered-" + v.Foo, nil
	})

	got, err := human.Format(RegistryFormatted{Foo: "a"})
	require.NoError(t, err)
	require.Equal(t, "registered-a", got)

	// A pointer derefs to the registered value type.
	got, err = human.Format(&RegistryFormatted{Foo: "b"})
	require.NoError(t, err)
	require.Equal(t, "registered-b", got)
}

func TestFormatRegisteredPointer(t *testing.T) {
	human.RegisterFormatFunc[*RegistryPointerFormatted](func(v *RegistryPointerFormatted) (string, error) {
		return "pointer-" + v.Foo, nil
	})

	got, err := human.Format(&RegistryPointerFormatted{Foo: "c"})
	require.NoError(t, err)
	require.Equal(t, "pointer-c", got)
}

func TestFormatRegisteredPrecedence(t *testing.T) {
	human.RegisterFormatFunc[RegistryPrecedenceFormatted](func(RegistryPrecedenceFormatted) (string, error) {
		return "registered", nil
	})

	// Like v1, a registered format function wins over Formatter.
	got, err := human.Format(RegistryPrecedenceFormatted{})
	require.NoError(t, err)
	require.Equal(t, "registered", got)
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
		{
			name: "duration",
			data: scw.Duration{Seconds: 604800},
			want: "7 days",
		},
		{
			name: "duration combined",
			data: scw.Duration{Seconds: 9000},
			want: "2 hours 30 minutes",
		},
		{
			name: "duration zero",
			data: scw.Duration{},
			want: "0 seconds",
		},
		{
			name: "duration pointer",
			data: &scw.Duration{Seconds: 30},
			want: "30 seconds",
		},
		{
			name: "sizes",
			data: []scw.Size{13200, 52428800},
			want: "13 kB, 50 MiB",
		},
		{
			name: "ip",
			data: net.ParseIP("1.2.3.4"),
			want: "1.2.3.4",
		},
		{
			name: "ips",
			data: []net.IP{net.ParseIP("1.2.3.4")},
			want: "[1.2.3.4]",
		},
		{
			name: "ipnet",
			data: makeIPNet("10.0.0.0/24"),
			want: "10.0.0.0/24",
		},
		{
			name: "nil ipnet",
			data: scw.IPNet{},
			want: "-",
		},
		{
			name: "version",
			data: version.Must(version.NewVersion("1.2.3")),
			want: "1.2.3",
		},
		{
			name: "json object",
			data: scw.JSONObject{"a": 1},
			want: `{"a":1}`,
		},
		{
			name: "bytes",
			data: []byte{0, 1},
			want: "AAE=",
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
