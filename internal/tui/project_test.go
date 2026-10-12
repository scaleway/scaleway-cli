package tui

import (
	"context"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/platform/terminal"
	accountv3 "github.com/scaleway/scaleway-sdk-go/api/account/v3"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// TestProjectMenuSelect verifies the picker fills from the fetch
// result (name-sorted), starts on the current project, and emits a
// setProjectMsg carrying the project ID.
func TestProjectMenuSelect(t *testing.T) {
	client, err := scw.NewClient(
		scw.WithAuth("SCWAAAAAAAAAAAAAAAAA", "11111111-aaaa-bbbb-cccc-111111111111"),
		scw.WithDefaultOrganizationID("11111111-aaaa-bbbb-cccc-111111111111"),
	)
	if err != nil {
		t.Fatal(err)
	}
	m := NewProjectMenu(client, "11111111-aaaa-bbbb-cccc-333333333333", nil)
	m.Size(100, 20)

	m.Update(projectsMsg{projects: []*accountv3.Project{
		{ID: "11111111-aaaa-bbbb-cccc-333333333333", Name: "web"},
		{ID: "22222222-aaaa-bbbb-cccc-444444444444", Name: "db"},
		{ID: "33333333-aaaa-bbbb-cccc-555555555555", Name: ""},
	}})

	// Sorted by name: the unnamed project first, then db, web.
	if len(m.items) != 3 || m.items[0] != "33333333" || m.items[1] != "db" || m.items[2] != "web" {
		t.Fatalf("items = %v, want [33333333 db web]", m.items)
	}
	// The cursor starts on the current project (web).
	if m.cursor != 2 {
		t.Fatalf("cursor = %d, want 2 (current project)", m.cursor)
	}
	msg, ok := m.selectVisible()().(setProjectMsg)
	if !ok || msg.projectID != "11111111-aaaa-bbbb-cccc-333333333333" || msg.name != "web" {
		t.Fatalf("select = %+v, want web", msg)
	}
}

// TestAppSetProject verifies the ctrl+p/ctrl+e-style flow for projects:
// ctrl+e opens the picker, and a selection re-creates the API client
// with the chosen project as default (session only).
func TestAppSetProject(t *testing.T) {
	platform := terminal.NewPlatform("scw-tui-test")
	meta := &core.Meta{
		ConfigPathFlag: profileTestConfigPath(t),
		Platform:       platform,
	}
	client, err := platform.CreateClient(nil, meta.ConfigPathFlag, scw.DefaultProfileName)
	if err != nil {
		t.Fatalf("CreateClient: %v", err)
	}
	res := &Resource{Name: "instance server", Title: "instance server"}
	a := NewApp(
		context.Background(), client, meta,
		map[string]*Resource{"instance server": res}, nil, nil, NewLocality("fr-par-1"),
	)
	a.width, a.height = 100, 30
	a.sizeCurrent()

	if a.project != profileTestProjA {
		t.Fatalf("initial project = %q, want %q", a.project, profileTestProjA)
	}

	// ctrl+e opens the project picker and returns the fetch cmd
	// (without it the picker stays "loading projects…").
	_, cmd := a.Update(tea.KeyMsg{Type: tea.KeyCtrlE})
	if cmd == nil {
		t.Fatal("ctrl+e: expected the project list fetch cmd, got nil")
	}
	if _, ok := a.current().(*ProjectMenu); !ok {
		t.Fatalf("ctrl+e: expected the project picker, got %T", a.current())
	}

	// Select another project: the client is re-created for it. (The
	// returned refresh/load-names cmds are not run: they would hit the
	// API.)
	const other = "22222222-aaaa-bbbb-cccc-555555555555"
	a.Update(setProjectMsg{projectID: other, name: "staging"})

	if a.project != other {
		t.Fatalf("project = %q, want %q", a.project, other)
	}
	if p, _ := a.client.GetDefaultProjectID(); p != other {
		t.Fatalf("client project = %q, want %q (client not re-created)", p, other)
	}
	if a.meta.Client != a.client {
		t.Fatal("meta client not updated (fetches would keep the old project)")
	}
	if org, _ := a.client.GetDefaultOrganizationID(); org != profileTestOrgA {
		t.Fatalf("client org = %q, want %q (credentials lost)", org, profileTestOrgA)
	}
	if a.flash != "project set to staging" {
		t.Fatalf("flash = %q", a.flash)
	}
	if len(a.stack) != 1 {
		t.Fatalf("stack len = %d, want 1 (picker closed)", len(a.stack))
	}
	// The header shows the new project (name arrives with the
	// load-names cmd, not run here: short ID until then).
	plain := ansiRe.ReplaceAllString(a.viewHeader(), "")
	if !strings.Contains(plain, "project: "+shortID(other)) {
		t.Fatalf("header should show the new project, got:\n%s", plain)
	}
}
