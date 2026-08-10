package core_test

import (
	"context"
	"reflect"
	"testing"

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
