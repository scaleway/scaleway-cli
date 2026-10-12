package tui

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type fakeUpdateRequest struct {
	Zone       scw.Zone `json:"-"`
	ServerID   string   `json:"-"`
	Name       *string  `json:"name"`
	Tags       []string `json:"tags"`
	IsDisabled bool     `json:"is_disabled"`
}

func fakeUpdateCommand() *core.Command {
	return &core.Command{
		Namespace: "instance", Resource: "server", Verb: "update",
		ArgsType: reflect.TypeFor[fakeUpdateRequest](),
		ArgSpecs: core.ArgSpecs{
			{Name: "server-id", Short: "Server ID", Required: true, Positional: true},
			{Name: "name", Short: "Name of the server"},
			{Name: "tags.{index}", Short: "Tags of the server"},
			{Name: "is-disabled", Short: "Disable the server"},
		},
		Run: func(ctx context.Context, argsI any) (any, error) { return nil, nil },
	}
}

func TestEditDocument(t *testing.T) {
	type fakeServer struct {
		ID   string
		Name string
		Tags []string
	}
	res := &Resource{Name: "instance server", UpdateCmd: fakeUpdateCommand()}
	resource := &fakeServer{Name: "web-1", Tags: []string{"a", "b"}}

	content, tmpl, err := editDocument(res, NewLocality("fr-par-1"), "abc", resource)
	if err != nil {
		t.Fatal(err)
	}
	s := string(content)
	// The current values are present.
	if !strings.Contains(s, "name: web-1") || !strings.Contains(s, "is_disabled: false") {
		t.Fatalf("missing values:\n%s", s)
	}
	// Each field has its doc as a comment above it (nested args index
	// their base key).
	for _, pair := range []struct{ doc, field string }{
		{"# Name of the server", "name: web-1"},
		{"# Tags of the server", "tags:"},
		{"# Disable the server", "is_disabled: false"},
	} {
		docIdx := strings.Index(s, pair.doc+"\n")
		fieldIdx := strings.Index(s, pair.field)
		if docIdx < 0 {
			t.Fatalf("missing doc %q:\n%s", pair.doc, s)
		}
		if fieldIdx < docIdx+len(pair.doc)+1 {
			t.Fatalf("doc %q must sit above %q:\n%s", pair.doc, pair.field, s)
		}
	}
	// Path parameters (json:"-") are not in the file.
	if strings.Contains(s, "zone:") || strings.Contains(s, "server_id:") {
		t.Fatalf("path parameters must not be in the file:\n%s", s)
	}
	// The tmpl keeps the id and zone for the post-edit restore.
	req := tmpl.(*fakeUpdateRequest)
	if req.ServerID != "abc" || req.Zone != "fr-par-1" {
		t.Fatalf("tmpl = %+v, want id and zone filled", req)
	}
}

func TestFinalUpdateRequest(t *testing.T) {
	cmd := fakeUpdateCommand()
	tmpl := &fakeUpdateRequest{ServerID: "abc", Zone: "fr-par-1"}
	// ponytail: a bare `y`/`n` would be parsed as a YAML 1.1 boolean
	// (ghodss/yaml behavior, same as the CLI `edit` commands) — values
	// like that must be quoted in the file.
	edited := []byte("name: renamed\ntags:\n    - x\n    - z\nis_disabled: true\n")

	req, err := finalUpdateRequest(cmd, tmpl, edited)
	if err != nil {
		t.Fatal(err)
	}
	r := req.(*fakeUpdateRequest)
	if r.Name == nil || *r.Name != "renamed" {
		t.Fatalf("name = %v, want renamed", r.Name)
	}
	if len(r.Tags) != 2 || r.Tags[0] != "x" || r.Tags[1] != "z" {
		t.Fatalf("tags = %v, want [x z]", r.Tags)
	}
	if !r.IsDisabled {
		t.Fatal("is_disabled = false, want true")
	}
	// Path params (json:"-") are restored from the tmpl.
	if r.ServerID != "abc" || r.Zone != "fr-par-1" {
		t.Fatalf("path params = %q/%q, want restored from the tmpl", r.ServerID, r.Zone)
	}
}

func TestFinalUpdateRequestInvalidYAML(t *testing.T) {
	_, err := finalUpdateRequest(
		fakeUpdateCommand(),
		&fakeUpdateRequest{},
		[]byte("name: [unclosed"),
	)
	if err == nil {
		t.Fatal("expected an error for invalid YAML")
	}
}

func TestTopLevelYAMLKey(t *testing.T) {
	if k, ok := topLevelYAMLKey("name: foo"); !ok || k != "name" {
		t.Fatalf("topLevelYAMLKey('name: foo') = %q, %v", k, ok)
	}
	if k, ok := topLevelYAMLKey("tags:"); !ok || k != "tags" {
		t.Fatalf("topLevelYAMLKey('tags:') = %q, %v", k, ok)
	}
	for _, line := range []string{"  - a", "    - b", "  name: nested", ""} {
		if _, ok := topLevelYAMLKey(line); ok {
			t.Fatalf("topLevelYAMLKey(%q) must not be a top-level key", line)
		}
	}
}

func TestResourcesFromCommandsWiresEdit(t *testing.T) {
	mk := func(ns, res, verb string) *core.Command {
		return &core.Command{
			Namespace: ns, Resource: res, Verb: verb,
			Run: func(ctx context.Context, argsI any) (any, error) { return nil, nil },
		}
	}
	list := mk("instance", "server", "list")
	get := mk("instance", "server", "get")
	update := fakeUpdateCommand()
	listOnly := mk("lb", "pool", "list")

	res := ResourcesFromCommands(
		[]*core.Command{list, get, update, listOnly},
		nil,
		NewLocality("fr-par-1"),
	)
	s := res["instance server"]
	if s == nil {
		t.Fatalf("missing 'instance server' in %v", res)
	}
	if s.GetCmd != get || s.UpdateCmd != update {
		t.Fatal("get/update commands not wired")
	}
	if !s.Editable() {
		t.Fatal("'instance server' must be editable")
	}
	if p := res["lb pool"]; p == nil || p.Editable() {
		t.Fatal("'lb pool' (no get/update) must not be editable")
	}
}
