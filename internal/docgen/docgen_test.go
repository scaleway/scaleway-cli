package docgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_escapeTableCell(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no pipe",
			in:   "Server tags",
			want: "Server tags",
		},
		{
			name: "list of accepted values",
			in:   "Filter by IP type (server | rpnv2_subnet)",
			want: `Filter by IP type (server \| rpnv2_subnet)`,
		},
		{
			name: "every pipe is escaped",
			in:   "(new | ipv4 | ipv6 | both)",
			want: `(new \| ipv4 \| ipv6 \| both)`,
		},
		{
			name: "empty",
			in:   "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, escapeTableCell(tt.in))
		})
	}
}
