package rdb

import (
	"context"
	"reflect"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func TestRegionFromAutocompleteRequest(t *testing.T) {
	t.Parallel()

	if got := regionFromAutocompleteRequest(&rdbCreateInstanceRequestCustom{}); got != "" {
		t.Fatalf("nil embedded request: got %q", got)
	}

	if got := regionFromAutocompleteRequest((*rdbCreateInstanceRequestCustom)(nil)); got != "" {
		t.Fatalf("nil custom request: got %q", got)
	}

	if got := regionFromAutocompleteRequest((*rdbSDK.CreateInstanceRequest)(nil)); got != "" {
		t.Fatalf("nil create request: got %q", got)
	}

	if got := regionFromAutocompleteRequest((*rdbSDK.UpgradeInstanceRequest)(nil)); got != "" {
		t.Fatalf("nil upgrade request: got %q", got)
	}

	want := scw.RegionNlAms
	got := regionFromAutocompleteRequest(&rdbCreateInstanceRequestCustom{
		CreateInstanceRequest: &rdbSDK.CreateInstanceRequest{Region: want},
	})
	if got != want {
		t.Fatalf("custom create: got %q, want %q", got, want)
	}

	got = regionFromAutocompleteRequest(&rdbSDK.UpgradeInstanceRequest{Region: want})
	if got != want {
		t.Fatalf("upgrade: got %q, want %q", got, want)
	}
}

func TestAutocompleteNilSafeCacheIteration(t *testing.T) {
	t.Parallel()

	prevEngines := completeListEngineCache
	prevNodes := completeListNodeTypeCache
	t.Cleanup(func() {
		completeListEngineCache = prevEngines
		completeListNodeTypeCache = prevNodes
	})

	completeListEngineCache = &rdbSDK.ListDatabaseEnginesResponse{
		Engines: []*rdbSDK.DatabaseEngine{
			nil,
			{
				Name: "PostgreSQL",
				Versions: []*rdbSDK.EngineVersion{
					nil,
					{Name: "PostgreSQL-17", Disabled: false},
					{Name: "PostgreSQL-10", Disabled: true},
				},
			},
		},
	}
	completeListNodeTypeCache = &rdbSDK.ListNodeTypesResponse{
		NodeTypes: []*rdbSDK.NodeType{
			nil,
			{Name: "db-dev-s"},
		},
	}

	// Empty custom create request used to panic on engine= tab completion.
	req := reflect.New(reflect.TypeFor[rdbCreateInstanceRequestCustom]()).Interface()

	engines := autoCompleteDatabaseEngines(context.Background(), "PostgreSQL", req)
	assertSuggestionValues(t, engines, "PostgreSQL-17")

	nodes := autoCompleteNodeType(context.Background(), "db-dev", req)
	assertSuggestionValues(t, nodes, "db-dev-s")
}

func assertSuggestionValues(t *testing.T, got core.AutocompleteSuggestions, want ...string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len=%d (%v), want %d (%v)", len(got), suggestionValues(got), len(want), want)
	}
	for i := range want {
		if got[i].Value != want[i] {
			t.Fatalf("got %v, want %v", suggestionValues(got), want)
		}
	}
}

func suggestionValues(suggestions core.AutocompleteSuggestions) []string {
	out := make([]string, len(suggestions))
	for i, suggestion := range suggestions {
		out[i] = suggestion.Value
	}

	return out
}
