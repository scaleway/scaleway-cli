package core

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testAutocompleteGetCommands() *Commands {
	return NewCommands(
		&Command{
			Namespace: "test",
			Resource:  "flower",
			Verb:      "create",
			ArgSpecs: ArgSpecs{
				{
					Name: "name",
				},
				{
					Name:       "species",
					EnumValues: []string{"rose", "violet", "petunia", "virginia bluebell"},
				},
				{
					Name: "size",
					AutoCompleteFunc: func(ctx context.Context, prefix string) AutocompleteSuggestions {
						return []string{regexp.MustCompile("[a-z]").ReplaceAllString(prefix, "")}
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
			WaitFunc: func(ctx context.Context, argsI, respI interface{}) error {
				return nil
			},
		},
		&Command{
			Namespace: "test",
			Resource:  "flower",
			Verb:      "delete",
			ArgSpecs:  nil,
		},
	)
}

func TestAutocomplete(t *testing.T) {
	ctx := injectCommands(context.Background(), testAutocompleteGetCommands())

	type testCase struct {
		Suggestions         AutocompleteSuggestions
		WordToCompleteIndex int
		Words               []string
	}

	run := func(tc *testCase) func(*testing.T) {
		return func(t *testing.T) {

			words := tc.Words
			if len(words) == 0 {
				name := strings.Replace(t.Name(), "TestAutocomplete/", "", -1)
				name = strings.Replace(name, "_", " ", -1)
				words = strings.Split(name, " ")
			}

			wordToCompleteIndex := len(words) - 1
			if tc.WordToCompleteIndex != 0 {
				wordToCompleteIndex = tc.WordToCompleteIndex
			}
			leftWords := words[:wordToCompleteIndex]
			wordToComplete := words[wordToCompleteIndex]
			rightWord := words[wordToCompleteIndex+1:]

			result := AutoComplete(ctx, leftWords, wordToComplete, rightWord)
			assert.Equal(t, tc.Suggestions, result.Suggestions)
		}
	}

	t.Run("scw", run(&testCase{Suggestions: AutocompleteSuggestions{"scw"}}))
	t.Run("scw ", run(&testCase{Suggestions: AutocompleteSuggestions{"test"}}))
	t.Run("scw te", run(&testCase{Suggestions: AutocompleteSuggestions{"test"}}))
	t.Run("scw test", run(&testCase{Suggestions: AutocompleteSuggestions{"test"}}))
	t.Run("scw  flower create name=plop", run(&testCase{WordToCompleteIndex: 1, Suggestions: AutocompleteSuggestions{"test"}}))
	t.Run("scw te flower create name=plop", run(&testCase{WordToCompleteIndex: 1, Suggestions: AutocompleteSuggestions{"test"}}))
	t.Run("scw test ", run(&testCase{Suggestions: AutocompleteSuggestions{"flower"}}))
	t.Run("scw test fl", run(&testCase{Suggestions: AutocompleteSuggestions{"flower"}}))
	t.Run("scw test flower ", run(&testCase{Suggestions: AutocompleteSuggestions{"create", "delete"}}))
	t.Run("scw test flower cr", run(&testCase{Suggestions: AutocompleteSuggestions{"create"}}))
	t.Run("scw test flower d", run(&testCase{Suggestions: AutocompleteSuggestions{"delete"}}))
	t.Run("scw test flower create ", run(&testCase{Suggestions: AutocompleteSuggestions{"colours.0=", "leaves.0.size=", "name=", "size=", "species="}}))
	t.Run("scw test flower create n", run(&testCase{Suggestions: AutocompleteSuggestions{"name="}}))
	t.Run("scw test flower create name", run(&testCase{Suggestions: AutocompleteSuggestions{"name="}}))
	t.Run("scw test flower create name=", run(&testCase{Suggestions: nil}))
	t.Run("scw test flower create name n", run(&testCase{Suggestions: nil})) // happens when name is a boolean
	t.Run("scw test flower create name=p", run(&testCase{Suggestions: nil}))
	t.Run("scw test flower create name=p ", run(&testCase{Suggestions: AutocompleteSuggestions{"colours.0=", "leaves.0.size=", "size=", "species="}}))
	t.Run("scw test flower create name=plop n", run(&testCase{Suggestions: nil}))
	t.Run("scw test flower create n name=plop", run(&testCase{WordToCompleteIndex: 4, Suggestions: nil}))
	t.Run("scw test flower create s", run(&testCase{Suggestions: AutocompleteSuggestions{"size=", "species="}}))
	t.Run("scw test flower create species=", run(&testCase{Suggestions: AutocompleteSuggestions{"species=petunia", "species=rose", "species=violet", "species=virginia bluebell"}}))
	t.Run("scw test flower create species=v", run(&testCase{Suggestions: AutocompleteSuggestions{"species=violet", "species=virginia bluebell"}}))
	t.Run("scw test flower create size=a1b2c", run(&testCase{Suggestions: AutocompleteSuggestions{"size=12"}}))
	t.Run("scw test flower create colo", run(&testCase{Suggestions: AutocompleteSuggestions{"colours.0="}}))
	t.Run("scw test flower create colours.0", run(&testCase{Suggestions: AutocompleteSuggestions{"colours.0="}}))
	t.Run("scw test flower create colours.0=", run(&testCase{Suggestions: AutocompleteSuggestions{"colours.0=blue", "colours.0=pink", "colours.0=red"}}))
	t.Run("scw test flower create colours.0=r", run(&testCase{Suggestions: AutocompleteSuggestions{"colours.0=red"}}))
	t.Run("scw test flower create colo colours.1=red", run(&testCase{WordToCompleteIndex: 4, Suggestions: AutocompleteSuggestions{"colours.0="}}))
	t.Run("scw test flower create colo colours.0=blue colours.1=red", run(&testCase{WordToCompleteIndex: 4, Suggestions: AutocompleteSuggestions{"colours.2="}}))
	t.Run("scw test flower create colours.0=blue colours.1=r", run(&testCase{Suggestions: AutocompleteSuggestions{"colours.1=red"}}))
	t.Run("scw test flower create leaves.", run(&testCase{Suggestions: AutocompleteSuggestions{"leaves.0.size="}}))
	t.Run("scw test flower create leaves.0", run(&testCase{Suggestions: AutocompleteSuggestions{"leaves.0.size="}}))
	t.Run("scw test flower create leaves.0.", run(&testCase{Suggestions: AutocompleteSuggestions{"leaves.0.size="}}))
	t.Run("scw test flower create leaves.0.size=M leaves", run(&testCase{Suggestions: AutocompleteSuggestions{"leaves.1.size="}}))
	t.Run("scw test flower create leaves.0.size=M leaves leaves.1.size=M", run(&testCase{WordToCompleteIndex: 5, Suggestions: AutocompleteSuggestions{"leaves.2.size="}}))
	// TODO: t.Run("scw test flower create leaves.0.size=", run(&testCase{Suggestions: AutocompleteSuggestions{"L", "M", "S", "XL", "XXL"}}))

	t.Run("scw -", run(&testCase{Suggestions: AutocompleteSuggestions{"--access-key", "--debug", "--help", "--output", "--profile", "--secret-key", "-D", "-h", "-o", "-p"}}))
	t.Run("scw test -o j", run(&testCase{Suggestions: AutocompleteSuggestions{"json"}}))
	t.Run("scw test flower -o ", run(&testCase{Suggestions: AutocompleteSuggestions{"human", "json"}}))
	t.Run("scw test flower -o json create -", run(&testCase{Suggestions: AutocompleteSuggestions{"--access-key", "--debug", "--help", "--output", "--profile", "--secret-key", "--wait", "-D", "-h", "-p", "-w"}}))
	t.Run("scw test flower create name=p -o j", run(&testCase{Suggestions: AutocompleteSuggestions{"json"}}))
	t.Run("scw test flower create name=p -o json ", run(&testCase{Suggestions: AutocompleteSuggestions{"colours.0=", "leaves.0.size=", "size=", "species="}}))
	t.Run("scw test -o json flower create ", run(&testCase{Suggestions: AutocompleteSuggestions{"colours.0=", "leaves.0.size=", "name=", "size=", "species="}}))
}

func TestWordIndex(t *testing.T) {
	type testCase struct {
		CharIndex int
		Words     []string
		WordIndex int
	}

	run := func(tc *testCase) func(*testing.T) {
		return func(t *testing.T) {
			words := tc.Words
			if len(words) == 0 {
				name := strings.ReplaceAll(t.Name(), "TestWordIndex/", "")
				name = strings.ReplaceAll(name, " ", "_")
				words = strings.Split(name, "_")
			}

			assert.Equal(t, tc.WordIndex, WordIndex(tc.CharIndex, words))
		}
	}

	t.Run("scw", run(&testCase{CharIndex: 3, WordIndex: 0}))
	t.Run("scw ", run(&testCase{CharIndex: 3, WordIndex: 0}))
	t.Run("scw ", run(&testCase{CharIndex: 4, WordIndex: 1}))
	t.Run("scw plop", run(&testCase{CharIndex: 4, WordIndex: 1}))
	t.Run("scw plop", run(&testCase{CharIndex: 8, WordIndex: 1}))
	t.Run("scw in security-groups list", run(&testCase{CharIndex: 6, WordIndex: 1}))
	t.Run("scw in security-groups list", run(&testCase{CharIndex: 27, WordIndex: 3}))
	t.Run("scw in security-groups list ", run(&testCase{CharIndex: 28, WordIndex: 4}))
}
