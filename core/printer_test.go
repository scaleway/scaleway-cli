package core_test

import (
	"bytes"
	"context"
	"net/netip"
	"reflect"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func Test_CorePrinter(t *testing.T) {
	type Human struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	commands := core.NewCommands(
		&core.Command{
			Namespace: "get",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(_ context.Context, _ any) (any, error) {
				return Human{
					ID:   "111111111-111111111",
					Name: "David Copperfield",
				}, nil
			},
		},
		&core.Command{
			Namespace: "list",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(_ context.Context, _ any) (any, error) {
				return []*Human{
					{ID: "111111111-111111111", Name: "David Copperfield"},
					{ID: "222222222-222222222", Name: "Xavier Niel"},
				}, nil
			},
		},
	)

	t.Run("human-simple-without-option", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw get -o human",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("human-simple-with-options", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw get -o human=ID,Name",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("human-list-without-option", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw list -o human",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("human-list-with-options", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw list -o human=Name,ID",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("human-list-with-options-unknown-column", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw -D list -o human=Name,ID,Unknown",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("human-list-no-headers", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw list -o human -H",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("human-list-no-headers-with-options", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw list -o human=Name,ID -H",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_YamlPrinter(t *testing.T) {
	type Human struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	commands := core.NewCommands(
		&core.Command{
			Namespace: "get",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(_ context.Context, _ any) (any, error) {
				return Human{
					ID:   "111111111-111111111",
					Name: "David Copperfield",
				}, nil
			},
		},
		&core.Command{
			Namespace: "list",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(_ context.Context, _ any) (any, error) {
				return []*Human{
					{ID: "111111111-111111111", Name: "David Copperfield"},
					{ID: "222222222-222222222", Name: "Xavier Niel"},
				}, nil
			},
		},
		&core.Command{
			Namespace: "NilSlice",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(_ context.Context, _ any) (any, error) {
				return []Human{}, nil
			},
		},
	)

	t.Run("human-simple-without-option", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw get -o yaml",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("human-simple-with-options", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw get -o yaml=ID,Name",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("human-list-without-option", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw list -o yaml",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("human-list-with-options", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw list -o yaml=Name,ID",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("human-list-with-options-unknown-column", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw -D list -o yaml=Name,ID,Unknown",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("nil-slice", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw NilSlice -o yaml",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_TemplatePrinter(t *testing.T) {
	type Human struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	commands := core.NewCommands(
		&core.Command{
			Namespace: "get",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(_ context.Context, _ any) (any, error) {
				return Human{
					ID:   "111111111-111111111",
					Name: "David Copperfield",
				}, nil
			},
		},
		&core.Command{
			Namespace: "list",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(_ context.Context, _ any) (any, error) {
				return []*Human{
					{ID: "111111111-111111111", Name: "David Copperfield"},
					{ID: "222222222-222222222", Name: "Xavier Niel"},
				}, nil
			},
		},
		&core.Command{
			Namespace: "NilSlice",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(_ context.Context, _ any) (any, error) {
				return []Human{}, nil
			},
		},
	)

	t.Run("template-simple-without-option", core.Test(&core.TestConfig{
		Commands: commands,
		Args: []string{
			"scw", "get", "-o", "template",
		},
		Check: core.TestCheckGolden(),
	}))

	t.Run("template-simple-with-options", core.Test(&core.TestConfig{
		Commands: commands,
		Args: []string{
			// We escape this sequence because there is already golang template rendering on commands in core.Test
			"scw", "get", "-o", "{{`template={{ .ID }}`}}",
		},
		Check: core.TestCheckGolden(),
	}))

	t.Run("template-list-without-option", core.Test(&core.TestConfig{
		Commands: commands,
		Args: []string{
			"scw", "list", "-o", "template",
		},
		Check: core.TestCheckGolden(),
	}))

	t.Run("template-list-with-options", core.Test(&core.TestConfig{
		Commands: commands,
		Args: []string{
			// We escape this sequence because there is already golang template rendering on commands in core.Test
			"scw", "list", "-o", "{{`template={{ .Name }} <-> {{ .ID }}`}}",
		},
		Check: core.TestCheckGolden(),
	}))

	t.Run("template-list-with-options-unknown-column", core.Test(&core.TestConfig{
		Commands: commands,
		Args: []string{
			// We escape this sequence because there is already golang template rendering on commands in core.Test
			"scw", "list", "-o", "{{`template={{ .Unknown }}`}}",
		},
		Check: core.TestCheckGolden(),
	}))

	t.Run("nil-slice", core.Test(&core.TestConfig{
		Commands: commands,
		Args: []string{
			"scw", "NilSlice", "-o", "template={{ .ID }}",
		},
		Check: core.TestCheckGolden(),
	}))
}

func Test_JSONPrinter(t *testing.T) {
	type Empty struct{}
	type Endpoint struct {
		ID           string `json:"id"`
		LoadBalancer *Empty `json:"load_balancer,omitempty"`
	}

	commands := core.NewCommands(
		&core.Command{
			Namespace: "time",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(_ context.Context, _ any) (any, error) {
				return time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC), nil
			},
		},
		&core.Command{
			Namespace: "addr",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(_ context.Context, _ any) (any, error) {
				return netip.MustParseAddr("127.0.0.1"), nil
			},
		},
		&core.Command{
			Namespace: "buffer",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(_ context.Context, _ any) (any, error) {
				return &bytes.Buffer{}, nil
			},
		},
		&core.Command{
			Namespace: "endpoint",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(_ context.Context, _ any) (any, error) {
				return Endpoint{ID: "111111111", LoadBalancer: &Empty{}}, nil
			},
		},
	)

	t.Run("time", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw time -o json",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("addr", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw addr -o json",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("buffer", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw buffer -o json",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("union-marker", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd:      "scw endpoint -o json",
		Check:    core.TestCheckGolden(),
	}))
}
