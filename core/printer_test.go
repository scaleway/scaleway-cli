package core_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/stretchr/testify/assert"
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

func Test_WarnUnknownTemplateFields(t *testing.T) {
	type Human struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	cmd := &core.Command{
		Namespace:    "get",
		ResponseType: reflect.TypeFor[*Human](),
	}

	tests := map[string]struct {
		outputFlag    string
		expectedWarns bool
	}{
		"no-response-type":       {outputFlag: "template={{ .Foo }}", expectedWarns: false},
		"non-template-output":    {outputFlag: "json", expectedWarns: false},
		"valid-field":            {outputFlag: "template={{ .ID }}", expectedWarns: false},
		"valid-field-chain":      {outputFlag: "template={{ .ID }}-{{ .Name }}", expectedWarns: false},
		"unknown-field":          {outputFlag: "template={{ .Foo }}", expectedWarns: true},
		"slice-response-element": {outputFlag: "template={{ .ID }} {{ .Name }}", expectedWarns: false},
		"empty-template":         {outputFlag: "template=", expectedWarns: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c := cmd
			if name == "no-response-type" {
				c = &core.Command{Namespace: "get"}
			}
			if name == "slice-response-element" {
				c = &core.Command{Namespace: "list", ResponseType: reflect.TypeFor[[]*Human]()}
			}

			var stderr strings.Builder
			core.WarnUnknownTemplateFields(c, tc.outputFlag, &stderr)

			if tc.expectedWarns {
				assert.Contains(t, stderr.String(), "not found in")
				assert.Contains(t, stderr.String(), "Foo")
				assert.Contains(t, stderr.String(), "ID, Name")
			} else {
				assert.Empty(t, stderr.String())
			}
		})
	}
}
