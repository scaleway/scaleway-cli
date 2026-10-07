package tui

import (
	"context"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	accountv3 "github.com/scaleway/scaleway-sdk-go/api/account/v3"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// projectFetchTimeout bounds the project list call.
const projectFetchTimeout = 10 * time.Second

// projectsMsg carries the outcome of the project list fetch.
type projectsMsg struct {
	projects []*accountv3.Project
	err      error
}

// setProjectMsg asks the app to switch the session to the given project.
type setProjectMsg struct {
	projectID string
	name      string
}

// ProjectMenu lists the projects of the current organization so the
// user can pick the session's project (ctrl+e). The projects are
// fetched on Init (like the browser fetches its rows).
type ProjectMenu struct {
	*listPicker
	client *scw.Client
	ids    []string // parallel to items: the project ID of each row
	err    string
}

// NewProjectMenu builds the project picker; current is the ID of the
// project in use (the cursor starts on it once the list is loaded).
func NewProjectMenu(client *scw.Client, current string, keys *Keymap) *ProjectMenu {
	return &ProjectMenu{
		listPicker: newListPicker("Select project", nil, current, keys),
		client:     client,
	}
}

// Name implements Screen.
func (m *ProjectMenu) Name() string {
	return "select-project"
}

// Hints implements Screen.
func (m *ProjectMenu) Hints() []string {
	k := m.keys

	return []string{
		"[ type ] filter",
		"[ ↑↓ ] move",
		hint(k.Key(ActionOpen), "select project"),
		hint(k.Key(ActionBack), "cancel"),
	}
}

// Init implements Screen: fetches the organization's projects.
func (m *ProjectMenu) Init() tea.Cmd {
	orgID, _ := m.client.GetDefaultOrganizationID()
	ctx, cancel := context.WithTimeout(context.Background(), projectFetchTimeout)

	return func() tea.Msg {
		defer cancel()
		resp, err := accountv3.NewProjectAPI(m.client).ListProjects(
			&accountv3.ProjectAPIListProjectsRequest{OrganizationID: orgID},
			scw.WithContext(ctx),
		)
		var projects []*accountv3.Project
		if err == nil && resp != nil {
			projects = resp.Projects
		}

		return projectsMsg{projects: projects, err: err}
	}
}

// Update implements Screen.
func (m *ProjectMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if pm, ok := msg.(projectsMsg); ok {
		m.setProjects(pm)

		return m, nil
	}
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	return m, m.listPicker.handleKey(km, m.selectVisible)
}

// View implements Screen.
func (m *ProjectMenu) View() string {
	if m.width <= 0 {
		return ""
	}
	if len(m.items) == 0 {
		if m.err != "" {
			return dimStyle.Render("error: " + m.err)
		}

		return dimStyle.Render("loading projects…")
	}

	return strings.Join(m.listPicker.renderLines(), "\n")
}

// setProjects fills the picker from a fetch result, name-sorted, and
// moves the cursor to the project currently in use.
func (m *ProjectMenu) setProjects(pm projectsMsg) {
	if pm.err != nil {
		m.err = pm.err.Error()

		return
	}
	type proj struct {
		id   string
		name string
	}
	projs := make([]proj, 0, len(pm.projects))
	for _, p := range pm.projects {
		projs = append(projs, proj{id: p.ID, name: p.Name})
	}
	sort.Slice(projs, func(i, j int) bool {
		if projs[i].name != projs[j].name {
			return projs[i].name < projs[j].name
		}

		return projs[i].id < projs[j].id
	})
	items := make([]string, len(projs))
	m.ids = make([]string, len(projs))
	for i, p := range projs {
		if p.name != "" {
			items[i] = p.name
		} else {
			items[i] = shortID(p.id)
		}
		m.ids[i] = p.id
	}
	m.listPicker.items = items
	for i, id := range m.ids {
		if id == m.listPicker.current {
			m.listPicker.cursor = i

			break
		}
	}
}

func (m *ProjectMenu) selectVisible() tea.Cmd {
	vis := m.visible()
	if m.cursor < len(vis) {
		idx := vis[m.cursor]
		name, id := m.items[idx], m.ids[idx]

		return func() tea.Msg { return setProjectMsg{projectID: id, name: name} }
	}

	return nil
}
