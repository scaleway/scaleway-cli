package core_test

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/platform/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func testAutocompleteGetCommands() *core.Commands {
	return core.NewCommands(
		&core.Command{
			Namespace: "test",
			Resource:  "flower",
			Verb:      "create",
			ArgsType:  reflect.TypeOf(struct{}{}),
			ArgSpecs: core.ArgSpecs{
				{
					Name: "name",
				},
				{
					Name:       "species",
					EnumValues: []string{"rose", "violet", "petunia", "virginia bluebell"},
				},
				{
					Name: "size",
					AutoCompleteFunc: func(_ context.Context, prefix string, _ any) core.AutocompleteSuggestions {
						value := regexp.MustCompile("[a-z]").ReplaceAllString(prefix, "")

						return core.AutocompleteSuggestions{{Value: value, Description: "the size"}}
					},
					EnumValues: []string{"S", "M", "L", "XL", "XXL"},
				},
				{
					Name:       "colours.{index}",
					EnumValues: []string{"blue", "red", "pink"},
				},
				{
					Name:       "leaves.{key}.size",
					EnumValues: []string{"S", "M", "L", "XL", "XXL"},
				},
			},
			WaitFunc: func(_ context.Context, _, _ any) (any, error) {
				return nil, nil
			},
		},
		&core.Command{
			Namespace: "test",
			Resource:  "flower",
			Verb:      "delete",
			ArgsType: reflect.TypeOf(struct {
				WithLeaves bool
			}{}),
			ArgSpecs: core.ArgSpecs{
				{
					Name:       "name",
					EnumValues: []string{"hibiscus", "anemone"},
					Positional: true,
				},
				{
					Name: "with-leaves",
				},
			},
		},
		&core.Command{
			Namespace:  "test",
			Resource:   "flower",
			Verb:       "deprecated",
			ArgsType:   reflect.TypeOf(struct{}{}),
			Deprecated: true,
			Short:      "this command is deprecated",
			Long:       "This command is deprecated and should not show up in autocomplete.",
		},
	)
}

type autoCompleteTestCase struct {
	Suggestions         core.AutocompleteSuggestions
	WordToCompleteIndex int
	Words               []string
}

func runAutocompleteTest(ctx context.Context, tc *autoCompleteTestCase) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		words := tc.Words
		if len(words) == 0 {
			name := strings.ReplaceAll(t.Name(), "TestAutocomplete/", "")
			name = strings.ReplaceAll(name, "_", " ")
			// Test can contain a sharp if duplicated
			// MyTest/scw_-flag_#01
			sharpIndex := strings.Index(name, "#")
			if sharpIndex != -1 {
				name = name[:sharpIndex]
			}
			words = strings.Split(name, " ")
		}

		wordToCompleteIndex := len(words) - 1
		if tc.WordToCompleteIndex != 0 {
			wordToCompleteIndex = tc.WordToCompleteIndex
		}
		leftWords := words[:wordToCompleteIndex]
		wordToComplete := words[wordToCompleteIndex]
		rightWord := words[wordToCompleteIndex+1:]

		result := core.AutoComplete(ctx, leftWords, wordToComplete, rightWord)
		assert.Equal(t, tc.Suggestions, result.Suggestions)
	}
}

func TestAutocomplete(t *testing.T) {
	ctx := core.InjectMeta(t.Context(), &core.Meta{
		Commands: testAutocompleteGetCommands(),
	})

	type testCase = autoCompleteTestCase

	run := func(tc *testCase) func(*testing.T) {
		return runAutocompleteTest(ctx, tc)
	}

	t.Run("scw ", run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "test"}}}))
	t.Run("scw te", run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "test"}}}))
	t.Run("scw test", run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "test"}}}))
	t.Run(
		"scw  flower create name=plop",
		run(&testCase{
			WordToCompleteIndex: 1,
			Suggestions:         core.AutocompleteSuggestions{{Value: "test"}},
		}),
	)
	t.Run(
		"scw te flower create name=plop",
		run(&testCase{
			WordToCompleteIndex: 1,
			Suggestions:         core.AutocompleteSuggestions{{Value: "test"}},
		}),
	)
	t.Run(
		"scw test ",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "flower"}}}),
	)
	t.Run(
		"scw test fl",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "flower"}}}),
	)
	t.Run(
		"scw test flower ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{{Value: "create"}, {Value: "delete"}},
			},
		),
	)
	t.Run(
		"scw test flower cr",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "create"}}}),
	)
	t.Run(
		"scw test flower d",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "delete"}}}),
	)
	t.Run(
		"scw test flower create ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "colours.0="},
					{Value: "leaves."},
					{Value: "name="},
					{Value: "size="},
					{Value: "species="},
				},
			},
		),
	)
	t.Run(
		"scw test flower create n",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "name="}}}),
	)
	t.Run(
		"scw test flower create name",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "name="}}}),
	)
	t.Run("scw test flower create name=", run(&testCase{Suggestions: nil}))
	t.Run(
		"scw test flower create name n",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "name="}}}),
	)
	t.Run("scw test flower create name=p", run(&testCase{Suggestions: nil}))
	t.Run(
		"scw test flower create name=p ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "colours.0="},
					{Value: "leaves."},
					{Value: "size="},
					{Value: "species="},
				},
			},
		),
	)
	t.Run("scw test flower create name=plop n", run(&testCase{Suggestions: nil}))
	t.Run(
		"scw test flower create n name=plop",
		run(&testCase{WordToCompleteIndex: 4, Suggestions: nil}),
	)
	t.Run(
		"scw test flower create s",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{{Value: "size="}, {Value: "species="}},
			},
		),
	)
	t.Run(
		"scw test flower create species=",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "species=petunia"},
					{Value: "species=rose"},
					{Value: "species=violet"},
					{Value: "species=virginia bluebell"},
				},
			},
		),
	)
	t.Run(
		"scw test flower create species=v",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "species=violet"},
					{Value: "species=virginia bluebell"},
				},
			},
		),
	)
	t.Run(
		"scw test flower create size=a1b2c",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "size=12", Description: "the size"},
				},
			},
		),
	)
	t.Run(
		"scw test flower create colo",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "colours.0="}}}),
	)
	t.Run(
		"scw test flower create colours.0",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "colours.0="}}}),
	)
	t.Run(
		"scw test flower create colours.0=",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "colours.0=blue"},
					{Value: "colours.0=pink"},
					{Value: "colours.0=red"},
				},
			},
		),
	)
	t.Run(
		"scw test flower create colours.0=r",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "colours.0=red"}}}),
	)
	t.Run(
		"scw test flower create colo colours.1=red",
		run(
			&testCase{
				WordToCompleteIndex: 4,
				Suggestions:         core.AutocompleteSuggestions{{Value: "colours.0="}},
			},
		),
	)
	t.Run(
		"scw test flower create colo colours.0=blue colours.1=red",
		run(
			&testCase{
				WordToCompleteIndex: 4,
				Suggestions:         core.AutocompleteSuggestions{{Value: "colours.2="}},
			},
		),
	)
	t.Run(
		"scw test flower create colours.0=blue colours.1=r",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "colours.1=red"}}}),
	)
	t.Run(
		"scw test flower create leaves.",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "leaves."}}}),
	)
	t.Run(
		"scw test flower create leaves.0",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "leaves.0.size="}}}),
	)
	t.Run(
		"scw test flower create leaves.0.",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "leaves.0.size="}}}),
	)
	t.Run(
		"scw test flower create leaves.0.size=M",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "leaves.0.size=M"}}}),
	)
	t.Run(
		"scw test flower create leaves.0.size=M leaves",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "leaves."}}}),
	)
	t.Run(
		"scw test flower create leaves.0.size=M leaves leaves.1.size=M",
		run(
			&testCase{
				WordToCompleteIndex: 5,
				Suggestions:         core.AutocompleteSuggestions{{Value: "leaves."}},
			},
		),
	)
	t.Run(
		"scw test flower delete ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "anemone"},
					{Value: "hibiscus"},
					{Value: "with-leaves="},
				},
			},
		),
	)
	t.Run(
		"scw test flower delete w",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "with-leaves="}}}),
	)
	t.Run(
		"scw test flower delete h",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "hibiscus"}}}),
	)
	t.Run(
		"scw test flower delete with-leaves=true ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{{Value: "anemone"}, {Value: "hibiscus"}},
			},
		),
	) // invalid notation
	t.Run("scw test flower delete hibiscus n", run(&testCase{Suggestions: nil}))
	t.Run(
		"scw test flower delete hibiscus w",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "with-leaves="}}}),
	)
	t.Run(
		"scw test flower delete hibiscus with-leaves=true",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "with-leaves=true"}}}),
	)
	t.Run(
		"scw test flower delete hibiscus with-leaves=true ",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "anemone"}}}),
	)
	t.Run(
		"scw test flower delete hibiscus with-leaves=",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "with-leaves=false"},
					{Value: "with-leaves=true"},
				},
			},
		),
	)
	t.Run(
		"scw test flower delete hibiscus with-leaves=tr",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "with-leaves=true"}}}),
	)
	t.Run("scw test flower delete hibiscus with-leaves=yes", run(&testCase{Suggestions: nil}))
	t.Run(
		"scw test flower create leaves.0.size=",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "leaves.0.size=L"},
					{Value: "leaves.0.size=M"},
					{Value: "leaves.0.size=S"},
					{Value: "leaves.0.size=XL"},
					{Value: "leaves.0.size=XXL"},
				},
			},
		),
	)
	t.Run(
		"scw -",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "--config"},
					{Value: "--debug"},
					{Value: "--help"},
					{Value: "--output"},
					{Value: "--profile"},
					{Value: "-D"},
					{Value: "-c"},
					{Value: "-h"},
					{Value: "-o"},
					{Value: "-p"},
				},
			},
		),
	)
	t.Run(
		"scw test -o j",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "json"}}}),
	)
	t.Run(
		"scw test flower -o ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: string(core.PrinterTypeHuman)},
					{Value: string(core.PrinterTypeJSON)},
					{Value: string(core.PrinterTypeTemplate)},
					{Value: string(core.PrinterTypeYAML)},
				},
			},
		),
	)
	t.Run(
		"scw test flower -o json create -",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "--config"},
					{Value: "--debug"},
					{Value: "--help"},
					{Value: "--output"},
					{Value: "--profile"},
					{Value: "--wait"},
					{Value: "-D"},
					{Value: "-c"},
					{Value: "-h"},
					{Value: "-p"},
					{Value: "-w"},
				},
			},
		),
	)
	t.Run(
		"scw test flower create name=p -o j",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "json"}}}),
	)
	t.Run(
		"scw test flower create name=p -o json ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "colours.0="},
					{Value: "leaves."},
					{Value: "size="},
					{Value: "species="},
				},
			},
		),
	)
	t.Run(
		"scw test flower create name=p -o=json ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "colours.0="},
					{Value: "leaves."},
					{Value: "size="},
					{Value: "species="},
				},
			},
		),
	)
	t.Run(
		"scw test flower create name=p -o=jso",
		run(&testCase{Suggestions: nil}),
	) // TODO: make this work
	t.Run(
		"scw test flower create name=p -o",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "-o"}}}),
	)
	t.Run(
		"scw test -o json flower create ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "colours.0="},
					{Value: "leaves."},
					{Value: "name="},
					{Value: "size="},
					{Value: "species="},
				},
			},
		),
	)
	t.Run(
		"scw test flower create name=p --profile xxxx ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "colours.0="},
					{Value: "leaves."},
					{Value: "size="},
					{Value: "species="},
				},
			},
		),
	)
	t.Run(
		"scw test --profile xxxx flower create name=p ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "colours.0="},
					{Value: "leaves."},
					{Value: "size="},
					{Value: "species="},
				},
			},
		),
	)
	t.Run("scw test flower create name=p --profile xxxx", run(&testCase{Suggestions: nil}))

	t.Run(
		"scw test flower -o json delete -",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "--config"},
					{Value: "--debug"},
					{Value: "--help"},
					{Value: "--output"},
					{Value: "--profile"},
					{Value: "-D"},
					{Value: "-c"},
					{Value: "-h"},
					{Value: "-p"},
				},
			},
		),
	)
	t.Run(
		"scw test flower delete -o ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: string(core.PrinterTypeHuman)},
					{Value: string(core.PrinterTypeJSON)},
					{Value: string(core.PrinterTypeTemplate)},
					{Value: string(core.PrinterTypeYAML)},
				},
			},
		),
	)
	t.Run(
		"scw test flower delete -o j",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "json"}}}),
	)
	t.Run(
		"scw test flower delete -o json ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "anemone"},
					{Value: "hibiscus"},
					{Value: "with-leaves="},
				},
			},
		),
	)
	t.Run(
		"scw test flower delete -o=json ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "anemone"},
					{Value: "hibiscus"},
					{Value: "with-leaves="},
				},
			},
		),
	)
	t.Run(
		"scw test flower delete -o json hibiscus w",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "with-leaves="}}}),
	)
	t.Run(
		"scw test flower delete -o=json hibiscus w",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "with-leaves="}}}),
	)
}

func TestAutocompleteArgs(t *testing.T) {
	commands := testAutocompleteGetCommands()
	commands.Add(&core.Command{
		Namespace: "test",
		Resource:  "flower",
		Verb:      "get",
		ArgsType: reflect.TypeOf(struct {
			Name         string
			MaterialName string
		}{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "flower-id",
				Positional: true,
			},
			{
				Name: "material-name",
			},
		},
	})
	commands.Add(&core.Command{
		Namespace: "test",
		Resource:  "flower",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(struct{}{}),
		ArgSpecs:  core.ArgSpecs{},
		Run: func(_ context.Context, _ any) (any, error) {
			return []*struct {
				ID   string
				Name string
			}{
				{
					ID:   "id1",
					Name: "flower1",
				},
				{
					ID:   "id2",
					Name: "flower2",
				},
			}, nil
		},
	})
	commands.Add(&core.Command{
		Namespace: "test",
		Resource:  "material",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(struct{}{}),
		ArgSpecs:  core.ArgSpecs{},
		Run: func(_ context.Context, _ any) (any, error) {
			return []*struct {
				Name string
			}{
				{
					Name: "material1",
				},
				{
					Name: "material2",
				},
			}, nil
		},
	})
	ctx := core.InjectMeta(t.Context(), &core.Meta{
		Commands: commands,
		BetaMode: true,
	})

	type testCase = autoCompleteTestCase

	run := func(tc *testCase) func(*testing.T) {
		return runAutocompleteTest(ctx, tc)
	}

	t.Run(
		"scw test flower get ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "id1", Description: "flower1"},
					{Value: "id2", Description: "flower2"},
					{Value: "material-name="},
				},
			},
		),
	)
	t.Run(
		"scw test flower get material-name=",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "material-name=material1"},
					{Value: "material-name=material2"},
				},
			},
		),
	)
	t.Run(
		"scw test flower get material-name=mat ",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "id1", Description: "flower1"},
					{Value: "id2", Description: "flower2"},
				},
			},
		),
	)
	t.Run(
		"scw test flower create name=",
		run(&testCase{Suggestions: core.AutocompleteSuggestions(nil)}),
	)
}

func TestAutocompleteRequiredListArgs(t *testing.T) {
	type listDatabaseArgs struct {
		InstanceID string
		Name       string
	}

	commands := core.NewCommands(
		&core.Command{
			Namespace: "test",
			Resource:  "backup",
			Verb:      "create",
			ArgsType: reflect.TypeOf(struct {
				InstanceID   string
				DatabaseName string
				Name         string
			}{}),
			ArgSpecs: core.ArgSpecs{
				{
					Name:     "instance-id",
					Required: true,
				},
				{
					Name:     "database-name",
					Required: true,
				},
				{
					Name: "name",
				},
			},
		},
		&core.Command{
			Namespace: "test",
			Resource:  "database",
			Verb:      "list",
			ArgsType:  reflect.TypeFor[listDatabaseArgs](),
			ArgSpecs: core.ArgSpecs{
				{
					Name:     "instance-id",
					Required: true,
				},
				{
					Name: "name",
				},
			},
			Run: func(_ context.Context, argsI any) (any, error) {
				listArgs := argsI.(*listDatabaseArgs)
				if listArgs.InstanceID == "" {
					return nil, errors.New("instance-id is required")
				}
				if listArgs.Name != "" {
					return []*struct {
						Name string
					}{}, nil
				}

				return []*struct {
					Name string
				}{
					{Name: "db1"},
					{Name: "db2"},
				}, nil
			},
		},
	)
	ctx := core.InjectMeta(t.Context(), &core.Meta{
		Commands: commands,
	})

	type testCase = autoCompleteTestCase

	run := func(tc *testCase) func(*testing.T) {
		return runAutocompleteTest(ctx, tc)
	}

	t.Run(
		"scw test backup create database-name=",
		run(&testCase{Suggestions: core.AutocompleteSuggestions(nil)}),
	)
	t.Run(
		"scw test backup create instance-id=xxx database-name=",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "database-name=db1"},
					{Value: "database-name=db2"},
				},
			},
		),
	)
	t.Run(
		"scw test backup create instance-id=xxx name=ignored database-name=",
		run(
			&testCase{
				Suggestions: core.AutocompleteSuggestions{
					{Value: "database-name=db1"},
					{Value: "database-name=db2"},
				},
			},
		),
	)
}

func TestAutocompleteProfiles(t *testing.T) {
	commands := testAutocompleteGetCommands()
	ctx := core.InjectMeta(t.Context(), &core.Meta{
		Commands: commands,
		BetaMode: true,
		Platform: terminal.NewPlatform(""),
	})

	type testCase = autoCompleteTestCase

	run := func(tc *testCase) func(*testing.T) {
		return runAutocompleteTest(ctx, tc)
	}
	t.Run("scw -p ", run(&testCase{Suggestions: nil}))
	t.Run("scw test -p ", run(&testCase{Suggestions: nil}))
	t.Run("scw test flower --profile ", run(&testCase{Suggestions: nil}))

	core.InjectConfig(ctx, &scw.Config{
		Profiles: map[string]*scw.Profile{
			"p1": nil,
			"p2": nil,
		},
	})

	t.Run(
		"scw -p ",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "p1"}, {Value: "p2"}}}),
	)
	t.Run(
		"scw test -p ",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "p1"}, {Value: "p2"}}}),
	)
	t.Run(
		"scw test flower --profile ",
		run(&testCase{Suggestions: core.AutocompleteSuggestions{{Value: "p1"}, {Value: "p2"}}}),
	)
}

func TestAutocompleteDeprecatedCommand(t *testing.T) {
	ctx := core.InjectMeta(t.Context(), &core.Meta{
		Commands: testAutocompleteGetCommands(),
	})

	type testCase = autoCompleteTestCase

	run := func(tc *testCase) func(*testing.T) {
		return runAutocompleteTest(ctx, tc)
	}

	t.Run("scw test flower deprecated", run(&testCase{Suggestions: nil}))
}
