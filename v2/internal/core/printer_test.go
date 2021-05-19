package core

import (
	"context"
	"reflect"
	"testing"
)

func Test_CorePrinter(t *testing.T) {
	type Human struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	commands := NewCommands(
		&Command{
			Namespace: "get",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
				return Human{
					ID:   "111111111-111111111",
					Name: "David Copperfield",
				}, nil
			},
		},
		&Command{
			Namespace: "list",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
				return []*Human{
					{ID: "111111111-111111111", Name: "David Copperfield"},
					{ID: "222222222-222222222", Name: "Xavier Niel"},
				}, nil
			},
		},
	)

	t.Run("human-simple-without-option", Test(&TestConfig{
		Commands: commands,
		Cmd:      "scw get -o human",
		Check:    TestCheckGolden(),
	}))

	t.Run("human-simple-with-options", Test(&TestConfig{
		Commands: commands,
		Cmd:      "scw get -o human=ID,Name",
		Check:    TestCheckGolden(),
	}))

	t.Run("human-list-without-option", Test(&TestConfig{
		Commands: commands,
		Cmd:      "scw list -o human",
		Check:    TestCheckGolden(),
	}))

	t.Run("human-list-with-options", Test(&TestConfig{
		Commands: commands,
		Cmd:      "scw list -o human=Name,ID",
		Check:    TestCheckGolden(),
	}))

	t.Run("human-list-with-options-unknown-column", Test(&TestConfig{
		Commands: commands,
		Cmd:      "scw -D list -o human=Name,ID,Unknown",
		Check:    TestCheckGolden(),
	}))
}

func Test_YamlPrinter(t *testing.T) {
	type Human struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	commands := NewCommands(
		&Command{
			Namespace: "get",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
				return Human{
					ID:   "111111111-111111111",
					Name: "David Copperfield",
				}, nil
			},
		},
		&Command{
			Namespace: "list",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
				return []*Human{
					{ID: "111111111-111111111", Name: "David Copperfield"},
					{ID: "222222222-222222222", Name: "Xavier Niel"},
				}, nil
			},
		},
		&Command{
			Namespace: "NilSlice",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
				return []Human{}, nil
			},
		},
	)

	t.Run("human-simple-without-option", Test(&TestConfig{
		Commands: commands,
		Cmd:      "scw get -o yaml",
		Check:    TestCheckGolden(),
	}))

	t.Run("human-simple-with-options", Test(&TestConfig{
		Commands: commands,
		Cmd:      "scw get -o yaml=ID,Name",
		Check:    TestCheckGolden(),
	}))

	t.Run("human-list-without-option", Test(&TestConfig{
		Commands: commands,
		Cmd:      "scw list -o yaml",
		Check:    TestCheckGolden(),
	}))

	t.Run("human-list-with-options", Test(&TestConfig{
		Commands: commands,
		Cmd:      "scw list -o yaml=Name,ID",
		Check:    TestCheckGolden(),
	}))

	t.Run("human-list-with-options-unknown-column", Test(&TestConfig{
		Commands: commands,
		Cmd:      "scw -D list -o yaml=Name,ID,Unknown",
		Check:    TestCheckGolden(),
	}))

	t.Run("nil-slice", Test(&TestConfig{
		Commands: commands,
		Cmd:      "scw NilSlice -o yaml",
		Check:    TestCheckGolden(),
	}))
}

func Test_TemplatePrinter(t *testing.T) {
	type Human struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	commands := NewCommands(
		&Command{
			Namespace: "get",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
				return Human{
					ID:   "111111111-111111111",
					Name: "David Copperfield",
				}, nil
			},
		},
		&Command{
			Namespace: "list",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
				return []*Human{
					{ID: "111111111-111111111", Name: "David Copperfield"},
					{ID: "222222222-222222222", Name: "Xavier Niel"},
				}, nil
			},
		},
		&Command{
			Namespace: "NilSlice",
			ArgsType:  reflect.TypeOf(struct{}{}),
			Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
				return []Human{}, nil
			},
		},
	)

	t.Run("template-simple-without-option", Test(&TestConfig{
		Commands: commands,
		Args: []string{
			"scw", "get", "-o", "template",
		},
		Check: TestCheckGolden(),
	}))

	t.Run("template-simple-with-options", Test(&TestConfig{
		Commands: commands,
		Args: []string{
			// We escape this sequence because there is already golang template rendering on commands in core.Test
			"scw", "get", "-o", "{{`template={{ .ID }}`}}",
		},
		Check: TestCheckGolden(),
	}))

	t.Run("template-list-without-option", Test(&TestConfig{
		Commands: commands,
		Args: []string{
			"scw", "list", "-o", "template",
		},
		Check: TestCheckGolden(),
	}))

	t.Run("template-list-with-options", Test(&TestConfig{
		Commands: commands,
		Args: []string{
			// We escape this sequence because there is already golang template rendering on commands in core.Test
			"scw", "list", "-o", "{{`template={{ .Name }} <-> {{ .ID }}`}}",
		},
		Check: TestCheckGolden(),
	}))

	t.Run("template-list-with-options-unknown-column", Test(&TestConfig{
		Commands: commands,
		Args: []string{
			// We escape this sequence because there is already golang template rendering on commands in core.Test
			"scw", "list", "-o", "{{`template={{ .Unknown }}`}}",
		},
		Check: TestCheckGolden(),
	}))

	t.Run("nil-slice", Test(&TestConfig{
		Commands: commands,
		Args: []string{
			"scw", "NilSlice", "-o", "template={{ .ID }}",
		},
		Check: TestCheckGolden(),
	}))
}
