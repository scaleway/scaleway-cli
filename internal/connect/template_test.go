package connect_test

import (
	"reflect"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/connect"
)

func Test_RenderCommand(t *testing.T) {
	data := map[string]any{
		"Username": "admin",
		"Password": "s3cr3t",
		"Host":     "1.2.3.4",
		"Port":     "5432",
	}

	tests := []struct {
		name    string
		tmpl    string
		want    []string
		wantErr bool
	}{
		{name: "empty returns nil", tmpl: "", want: nil},
		{name: "whitespace returns nil", tmpl: "   ", want: nil},
		{
			name: "renders fields",
			tmpl: "psql --host {{.Host}} --port {{.Port}} --username {{.Username}}",
			want: []string{"psql", "--host", "1.2.3.4", "--port", "5432", "--username", "admin"},
		},
		{name: "no placeholders", tmpl: "mytool", want: []string{"mytool"}},
		{name: "invalid template", tmpl: "{{.Host", wantErr: true},
		{name: "missing key", tmpl: "{{.Missing}}", wantErr: true},
		{name: "renders empty command", tmpl: "{{if false}}x{{end}}", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := connect.RenderCommand(tt.tmpl, data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RenderCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("RenderCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
