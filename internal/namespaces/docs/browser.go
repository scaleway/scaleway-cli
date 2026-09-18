//go:build !wasm

package docs

import (
	"context"
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/scaleway/scaleway-cli/v2/core"
)

type browserLevel int

const (
	levelNamespace browserLevel = iota
	levelResource
	levelVerb
	levelCommand
)

type browserModel struct {
	commands  *core.Commands
	ctx       context.Context //nolint:containedctx
	hierarchy *CommandHierarchy
	level     browserLevel

	namespaceCursor int
	resourceCursor  int
	verbCursor      int

	selectedNamespace string
	selectedResource  string

	commandPage string

	namespaces []string
	resources  []string
	verbs      []string

	// filtered lists (after fuzzy filtering)
	filteredNamespaces []string
	filteredResources  []string
	filteredVerbs      []string

	filter string

	width int
}

func runBrowser(ctx context.Context, commands *core.Commands) (any, error) {
	hierarchy := BuildHierarchy(commands)

	namespaces := make([]string, 0, len(hierarchy.Namespaces))
	for ns := range hierarchy.Namespaces {
		namespaces = append(namespaces, ns)
	}
	sort.Strings(namespaces)

	m := &browserModel{
		commands:   commands,
		ctx:        ctx,
		hierarchy:  hierarchy,
		level:      levelNamespace,
		namespaces: namespaces,
		width:      80,
	}
	m.applyFilter()

	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("error running docs browser: %w", err)
	}

	return &core.SuccessResult{
		Message: "Documentation browser closed",
	}, nil
}

func (m *browserModel) Init() tea.Cmd {
	return nil
}

func (m *browserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.filter != "" {
				m.filter = ""
				m.applyFilter()
			} else if m.level > levelNamespace {
				m.level--
				if m.level == levelNamespace {
					m.resources = nil
					m.verbs = nil
				}
				if m.level == levelResource {
					m.verbs = nil
				}
			}
		case "backspace":
			if m.filter != "" {
				m.filter = m.filter[:len(m.filter)-1]
				m.applyFilter()
			}
		case "up", "k":
			m.moveCursor(-1)
		case "down", "j":
			m.moveCursor(1)
		case "enter", " ":
			return m, m.handleSelect()
		default:
			if m.level != levelCommand {
				ch := msg.String()
				if len(ch) == 1 && ch[0] >= 32 && ch[0] < 127 {
					m.filter += ch
					m.applyFilter()
				}
			}
		}
	}

	return m, nil
}

//nolint:funcorder
func (m *browserModel) applyFilter() {
	switch m.level {
	case levelNamespace:
		m.filteredNamespaces = fuzzyFilter(m.namespaces, m.filter)
		m.namespaceCursor = clamp(m.namespaceCursor, len(m.filteredNamespaces)-1)
	case levelResource:
		m.filteredResources = fuzzyFilter(m.resources, m.filter)
		m.resourceCursor = clamp(m.resourceCursor, len(m.filteredResources)-1)
	case levelVerb:
		m.filteredVerbs = fuzzyFilter(m.verbs, m.filter)
		m.verbCursor = clamp(m.verbCursor, len(m.filteredVerbs)-1)
	}
}

//nolint:funcorder
func (m *browserModel) moveCursor(dir int) {
	switch m.level {
	case levelNamespace:
		m.namespaceCursor = clamp(m.namespaceCursor+dir, len(m.filteredNamespaces)-1)
	case levelResource:
		m.resourceCursor = clamp(m.resourceCursor+dir, len(m.filteredResources)-1)
	case levelVerb:
		m.verbCursor = clamp(m.verbCursor+dir, len(m.filteredVerbs)-1)
	}
}

//nolint:funcorder
func (m *browserModel) handleSelect() tea.Cmd {
	switch m.level {
	case levelNamespace:
		if len(m.filteredNamespaces) == 0 {
			return nil
		}
		m.selectedNamespace = m.filteredNamespaces[m.namespaceCursor]
		nsNode := m.hierarchy.Namespaces[m.selectedNamespace]
		if nsNode == nil {
			return nil
		}
		m.resources = make([]string, 0, len(nsNode.Children))
		for res := range nsNode.Children {
			m.resources = append(m.resources, res)
		}
		sort.Strings(m.resources)
		m.resourceCursor = 0
		m.level = levelResource
		m.filter = ""
		m.applyFilter()

		nsCmd := nsNode.Command
		if nsCmd != nil && nsCmd.Run != nil && len(nsNode.Children) == 0 {
			m.commandPage = RenderCommandPage(m.ctx, nsCmd, m.commands)
			m.level = levelCommand
		}
	case levelResource:
		if len(m.filteredResources) == 0 {
			return nil
		}
		m.selectedResource = m.filteredResources[m.resourceCursor]
		nsNode := m.hierarchy.Namespaces[m.selectedNamespace]
		if nsNode == nil {
			return nil
		}
		resNode := nsNode.Children[m.selectedResource]
		if resNode == nil {
			return nil
		}
		m.verbs = make([]string, 0, len(resNode.Verbs))
		for v := range resNode.Verbs {
			m.verbs = append(m.verbs, v)
		}
		sort.Strings(m.verbs)
		m.verbCursor = 0
		m.level = levelVerb
		m.filter = ""
		m.applyFilter()

		resCmd := resNode.Command
		if resCmd != nil && resCmd.Run != nil && len(resNode.Verbs) == 0 {
			m.commandPage = RenderCommandPage(m.ctx, resCmd, m.commands)
			m.level = levelCommand
		}
	case levelVerb:
		if len(m.filteredVerbs) == 0 {
			return nil
		}
		nsNode := m.hierarchy.Namespaces[m.selectedNamespace]
		if nsNode == nil {
			return nil
		}
		resNode := nsNode.Children[m.selectedResource]
		if resNode == nil {
			return nil
		}
		cmd := resNode.Verbs[m.filteredVerbs[m.verbCursor]]
		if cmd == nil {
			return nil
		}
		m.commandPage = RenderCommandPage(m.ctx, cmd, m.commands)
		m.level = levelCommand
	case levelCommand:
		return tea.Quit
	}

	return nil
}

func (m *browserModel) View() string {
	var sb strings.Builder

	switch m.level {
	case levelNamespace:
		sb.WriteString("Documentation Browser - Namespaces\n\n")
		sb.WriteString(
			"Type to fuzzy-filter, arrow keys (or j/k) to navigate, Enter to select, q to quit.\n\n",
		)
		if m.filter != "" {
			sb.WriteString(fmt.Sprintf("Filter: %s  (Esc to clear)\n\n", m.filter))
		}
		for i, ns := range m.filteredNamespaces {
			node := m.hierarchy.Namespaces[ns]
			short := ""
			if node != nil && node.Command != nil {
				short = node.Command.Short
			}
			if i == m.namespaceCursor {
				sb.WriteString(fmt.Sprintf("> %-20s %s\n", ns, short))
			} else {
				sb.WriteString(fmt.Sprintf("  %-20s %s\n", ns, short))
			}
		}
	case levelResource:
		sb.WriteString(
			fmt.Sprintf("Documentation Browser - %s > Resources\n\n", m.selectedNamespace),
		)
		sb.WriteString(
			"Type to fuzzy-filter, arrow keys to navigate, Enter to select, Esc to go back, q to quit.\n\n",
		)
		if m.filter != "" {
			sb.WriteString(fmt.Sprintf("Filter: %s  (Esc to clear)\n\n", m.filter))
		}
		nsNode := m.hierarchy.Namespaces[m.selectedNamespace]
		for i, res := range m.filteredResources {
			short := ""
			if nsNode != nil && nsNode.Children[res] != nil && nsNode.Children[res].Command != nil {
				short = nsNode.Children[res].Command.Short
			}
			if i == m.resourceCursor {
				sb.WriteString(fmt.Sprintf("> %-20s %s\n", res, short))
			} else {
				sb.WriteString(fmt.Sprintf("  %-20s %s\n", res, short))
			}
		}
	case levelVerb:
		sb.WriteString(
			fmt.Sprintf(
				"Documentation Browser - %s %s > Commands\n\n",
				m.selectedNamespace,
				m.selectedResource,
			),
		)
		sb.WriteString(
			"Type to fuzzy-filter, arrow keys to navigate, Enter to view, Esc to go back, q to quit.\n\n",
		)
		if m.filter != "" {
			sb.WriteString(fmt.Sprintf("Filter: %s  (Esc to clear)\n\n", m.filter))
		}
		nsNode := m.hierarchy.Namespaces[m.selectedNamespace]
		resNode := nsNode.Children[m.selectedResource]
		for i, v := range m.filteredVerbs {
			cmd := resNode.Verbs[v]
			short := ""
			if cmd != nil {
				short = cmd.Short
			}
			if i == m.verbCursor {
				sb.WriteString(fmt.Sprintf("> %-20s %s\n", v, short))
			} else {
				sb.WriteString(fmt.Sprintf("  %-20s %s\n", v, short))
			}
		}
	case levelCommand:
		sb.WriteString("Command Reference (press Enter or Esc to go back)\n\n")
		sb.WriteString(m.commandPage)
	}

	return sb.String()
}

func clamp(v, hi int) int {
	if v < 0 {
		return 0
	}
	if v > hi {
		return hi
	}

	return v
}
