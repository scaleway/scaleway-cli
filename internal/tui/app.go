package tui

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/scaleway/scaleway-cli/v2/core"
	accountv3 "github.com/scaleway/scaleway-sdk-go/api/account/v3"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// refreshInterval is the auto-refresh period.
// ponytail: Scaleway has no watch stream, so we poll the API less
// aggressively than k9s does the informer cache (2s).
const refreshInterval = 10 * time.Second

// tickMsg triggers an auto-refresh of the current screen.
type tickMsg struct{}

// pushScreenMsg pushes a screen on the navigation stack (k9s: Stack.Push).
type pushScreenMsg struct{ screen Screen }

// popScreenMsg pops the top of the navigation stack (k9s: Stack.Pop).
type popScreenMsg struct{}

// flashMsg sets the status bar flash message.
type flashMsg struct{ msg string }

// copyMsg asks the app to copy a value to the clipboard.
type copyMsg struct {
	value string
	label string
}

// Screen is the k9s model.Component equivalent: anything the app can
// show in the content area.
type Screen interface {
	tea.Model
	Name() string
	Hints() []string
	Size(width, height int)
}

// App is the root model: the k9s view.App+ui.App collapsed into one
// bubbletea model. It owns the k9s layout:
//
//	<header>  project/org/zone/region · shortcuts · logo
//	<content> current screen (navigation stack)
//	<prompt>  `:` / `/` bar (when active)
//	<status>  rows · flash
type App struct {
	client  *scw.Client
	meta    *core.Meta // owns the client; profile switches reload it
	loc     *Locality
	profile string
	project string
	// project/org display names, loaded once from the API.
	projectName string
	orgName     string
	orgAlias    string
	resources   map[string]*Resource
	commands    map[string]string // user `:` commands (cli.yaml: tui.commands)
	keys        *Keymap

	stack  []Screen
	prompt *Prompt
	flash  string

	width  int
	height int
}

// namesMsg carries the project/organization names loaded on startup.
type namesMsg struct {
	projectName string
	orgName     string
	orgAlias    string
}

// NewApp creates the TUI root with the initial screen (the instance
// server list). meta owns the API client (profile switches reload it
// through it), keys is the effective keymap (nil: defaults), commands
// the user `:` commands from cli.yaml, loc the session zone/region
// (shared with the resource fetches).
func NewApp(
	ctx context.Context,
	client *scw.Client,
	meta *core.Meta,
	resources map[string]*Resource,
	keys *Keymap,
	commands map[string]string,
	loc *Locality,
) *App {
	project, _ := client.GetDefaultProjectID()
	profile := scw.DefaultProfileName
	if meta != nil {
		profile = core.ExtractProfileName(core.InjectMeta(ctx, meta))
	}

	return &App{
		client:    client,
		meta:      meta,
		profile:   profile,
		loc:       loc,
		project:   project,
		resources: resources,
		commands:  commands,
		keys:      keys,
		stack:     []Screen{NewBrowser(defaultResource(resources), keys)},
	}
}

// defaultResource is the initial screen; falls back to the first
// resource if the instance server list is not registered.
func defaultResource(resources map[string]*Resource) *Resource {
	if res := resources["instance server"]; res != nil {
		return res
	}
	names := make([]string, 0, len(resources))
	for name := range resources {
		names = append(names, name)
	}
	sort.Strings(names)
	if len(names) == 0 {
		panic("tui: no resources registered")
	}

	return resources[names[0]]
}

// Init implements tea.Model.
func (a *App) Init() tea.Cmd {
	tick := tea.Tick(refreshInterval, func(time.Time) tea.Msg { return tickMsg{} })

	return tea.Batch(a.current().Init(), tick, a.loadNames())
}

// Update implements tea.Model.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		a.width, a.height = m.Width, m.Height
		a.sizeCurrent()

		return a, nil
	case tickMsg:
		a.flash = ""
		mdl, cmd := a.current().Update(refreshMsg{})
		a.replaceCurrent(mdl)
		tick := tea.Tick(refreshInterval, func(time.Time) tea.Msg { return tickMsg{} })

		return a, tea.Batch(cmd, tick)
	case pushScreenMsg:
		a.stack = append(a.stack, m.screen)
		a.sizeCurrent()

		return a, a.current().Init()
	case popScreenMsg:
		if len(a.stack) > 1 {
			a.stack = a.stack[:len(a.stack)-1]
			a.sizeCurrent()
		}

		return a, nil
	case gotoResourceMsg:
		return a, a.gotoResource(m.res)
	case setLocalityMsg:
		if m.isZone {
			a.loc.SetZone(m.zone)
			a.flash = "zone set to " + string(m.zone)
		} else {
			a.loc.SetRegion(m.region)
			a.flash = "region set to " + string(m.region)
		}
		// Close the picker (it is on top of the stack) and re-fetch
		// the screen under it with the new locality.
		if len(a.stack) > 1 {
			a.stack = a.stack[:len(a.stack)-1]
			a.sizeCurrent()
		}
		mdl, cmd := a.current().Update(refreshMsg{})
		a.replaceCurrent(mdl)

		return a, cmd
	case setProfileMsg:
		a.switchProfile(m.profile)
		// Close the picker (it is on top of the stack) and re-fetch
		// the screen under it with the new profile.
		if len(a.stack) > 1 {
			a.stack = a.stack[:len(a.stack)-1]
			a.sizeCurrent()
		}
		mdl, cmd := a.current().Update(refreshMsg{})
		a.replaceCurrent(mdl)

		return a, tea.Batch(cmd, a.loadNames())
	case setProjectMsg:
		a.switchProject(m.projectID, m.name)
		// Close the picker (it is on top of the stack) and re-fetch
		// the screen under it with the new project.
		if len(a.stack) > 1 {
			a.stack = a.stack[:len(a.stack)-1]
			a.sizeCurrent()
		}
		mdl, cmd := a.current().Update(refreshMsg{})
		a.replaceCurrent(mdl)

		return a, tea.Batch(cmd, a.loadNames())
	case namesMsg:
		a.projectName = m.projectName
		a.orgName = m.orgName
		a.orgAlias = m.orgAlias

		return a, nil
	case flashMsg:
		a.flash = m.msg

		return a, nil
	case copyMsg:
		if err := clipboard.WriteAll(m.value); err != nil {
			a.flash = "clipboard error: " + err.Error()
		} else if m.label != "" {
			a.flash = m.label + " copied to clipboard"
		} else {
			a.flash = "copied to clipboard"
		}

		return a, nil
	case tea.KeyMsg:
		return a.updateKey(msg, m)
	}
	// Forward everything else (e.g. fetchMsg) to the current screen.
	mdl, cmd := a.current().Update(msg)
	a.replaceCurrent(mdl)

	return a, cmd
}

// View implements tea.Model.
func (a *App) View() string {
	if a.width == 0 {
		return "starting…"
	}

	parts := []string{a.viewHeader(), a.current().View()}
	if a.prompt != nil && a.prompt.Active {
		parts = append(parts, a.viewPrompt())
	}
	parts = append(parts, a.viewStatus())

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (a *App) current() Screen {
	return a.stack[len(a.stack)-1]
}

func (a *App) replaceCurrent(mdl tea.Model) {
	a.stack[len(a.stack)-1] = mdl.(Screen)
}

func (a *App) updateKey(msg tea.Msg, k tea.KeyMsg) (tea.Model, tea.Cmd) {
	if a.prompt != nil && a.prompt.Active {
		return a.updatePromptKey(k)
	}
	if k.Type == tea.KeyCtrlC {
		return a, tea.Quit
	}
	// In the goto menu and the pickers every keystroke feeds the
	// widget (k9s command buffer): the global shortcuts are
	// disabled while one is open.
	switch a.current().(type) {
	case *Menu, *LocalityMenu, *ProfileMenu, *ProjectMenu:
		mdl, cmd := a.current().Update(msg)
		a.replaceCurrent(mdl)

		return a, cmd
	}
	// Global actions, bound through the keymap so cli.yaml
	// (tui.keybindings) can override them.
	switch a.keys.Action(keyOf(k)) {
	case ActionQuit:
		return a, tea.Quit
	case ActionMenu:
		a.stack = append(a.stack, NewMenu(a.resources, a.keys))
		a.sizeCurrent()

		return a, nil
	case ActionZone:
		a.stack = append(a.stack, NewLocalityMenu(true, a.loc, a.keys))
		a.sizeCurrent()

		return a, nil
	case ActionRegion:
		a.stack = append(a.stack, NewLocalityMenu(false, a.loc, a.keys))
		a.sizeCurrent()

		return a, nil
	case ActionProfile:
		a.stack = append(a.stack, NewProfileMenu(a.profileList(), a.profile, a.keys))
		a.sizeCurrent()

		return a, nil
	case ActionProject:
		a.stack = append(a.stack, NewProjectMenu(a.client, a.project, a.keys))
		a.sizeCurrent()

		return a, a.current().Init()
	case ActionCommand:
		a.prompt = newCmdPrompt()
		a.sizeCurrent()

		return a, nil
	case ActionFilter:
		switch a.current().(type) {
		case *Browser:
			a.prompt = newFilterPrompt()
		case *Details:
			a.prompt = newSearchPrompt()
		default:
			return a, nil
		}
		a.sizeCurrent()

		return a, nil
	case ActionRefresh:
		mdl, cmd := a.current().Update(refreshMsg{})
		a.replaceCurrent(mdl)

		return a, cmd
	case ActionBack:
		if len(a.stack) > 1 {
			return a, func() tea.Msg { return popScreenMsg{} }
		}
	}

	// Everything else goes to the current screen (k9s: per-view key actions).
	mdl, cmd := a.current().Update(msg)
	a.replaceCurrent(mdl)

	return a, cmd
}

func (a *App) updatePromptKey(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	p := a.prompt
	switch k.Type {
	case tea.KeyEscape:
		p.Active = false
		a.sizeCurrent()
	case tea.KeyEnter:
		cmd := a.execPrompt()
		p.Active = false
		a.sizeCurrent()

		return a, cmd
	default:
		if edited, ok := lineEdit(k.Type, p.Buffer); ok {
			if edited != p.Buffer {
				p.Buffer = edited
				a.applyFilter()
			}
		} else if len(k.String()) == 1 && len(p.Buffer) < 120 {
			// ponytail: ASCII printable only; tea returns the rune as a string.
			p.Buffer += k.String()
			a.applyFilter()
		}
	}

	return a, nil
}

// applyFilter pushes the prompt buffer to the table filter (filter mode).
func (a *App) applyFilter() {
	if a.prompt.Mode != PromptFilter {
		return
	}
	if b, ok := a.current().(*Browser); ok {
		b.Table().SetFilter(a.prompt.Buffer)
	}
}

// execPrompt runs a `:` command. Filter mode just closes on enter;
// search mode finds the attribute and copies its value.
func (a *App) execPrompt() tea.Cmd {
	switch a.prompt.Mode {
	case PromptFilter:
		return nil
	case PromptSearch:
		if d, ok := a.current().(*Details); ok {
			return d.SearchAndCopy(a.prompt.Buffer)
		}

		return nil
	}
	cmd := strings.TrimSpace(a.prompt.Buffer)
	if cmd == "" {
		return nil
	}
	name := strings.ToLower(cmd)
	// User `:` commands (cli.yaml: tui.commands) resolve one level to a
	// resource name or a built-in command name.
	if target, ok := a.commands[name]; ok {
		name = strings.ToLower(strings.TrimSpace(target))
	}
	switch name {
	case "q", "quit":
		return tea.Quit
	case "help", "h":
		a.flash = fmt.Sprintf(
			"%d resources available — press %s to browse",
			len(a.resources),
			a.keys.Key(ActionMenu),
		)

		return nil
	}
	res, ok := a.resources[name]
	if !ok {
		if alias, ok2 := resourceAliases[name]; ok2 {
			res, ok = a.resources[alias]
		}
	}
	if !ok {
		a.flash = "unknown command: " + cmd

		return nil
	}

	return a.gotoResource(res)
}

// goto replaces the stack with a browser for res (k9s gotoResource).
func (a *App) gotoResource(res *Resource) tea.Cmd {
	a.stack = []Screen{NewBrowser(res, a.keys)}
	a.flash = ""
	a.sizeCurrent()

	return a.current().Init()
}

// profileList returns the profiles defined in the scw config file: the
// default profile first, then the named ones alphabetically.
func (a *App) profileList() []string {
	names := []string{scw.DefaultProfileName}
	if a.meta != nil {
		ctx := core.InjectMeta(context.Background(), a.meta)
		if config, err := scw.LoadConfigFromPath(core.ExtractConfigPath(ctx)); err == nil {
			for name := range config.Profiles {
				if name != scw.DefaultProfileName {
					names = append(names, name)
				}
			}
			sort.Strings(names[1:])
		}
	}

	return names
}

// switchProfile re-creates the API client for the given profile
// (reusing the CLI's ReloadClient, which reads the profile from the
// config file) and resets the header data derived from the previous
// one.
func (a *App) switchProfile(name string) {
	if a.meta == nil || a.meta.Platform == nil {
		a.flash = "profile switch unavailable"

		return
	}
	a.meta.ProfileFlag = name
	if err := core.ReloadClient(core.InjectMeta(context.Background(), a.meta)); err != nil {
		a.flash = "profile error: " + err.Error()

		return
	}
	a.client = a.meta.Client
	a.profile = name
	a.project, _ = a.client.GetDefaultProjectID()
	a.projectName, a.orgName, a.orgAlias = "", "", ""
	a.flash = "profile set to " + name
}

// switchProject re-creates the API client with the given project as
// the default one (session only: the config file is not touched).
func (a *App) switchProject(projectID, name string) {
	client, err := a.clientWithProject(projectID)
	if err != nil {
		a.flash = "project error: " + err.Error()

		return
	}
	a.client = client
	if a.meta != nil {
		a.meta.Client = client
	}
	a.project = projectID
	a.projectName = ""
	a.flash = "project set to " + name
}

// clientWithProject re-creates the API client for the current profile
// with projectID as the default project. The profile is read back from
// the platform's loaded config; without one it is rebuilt from the
// current client's credentials.
func (a *App) clientWithProject(projectID string) (*scw.Client, error) {
	var opts []scw.ClientOption
	if a.meta != nil && a.meta.Platform != nil {
		if config := a.meta.Platform.ScwConfig(); config != nil {
			profile, err := config.GetProfile(a.profile)
			if err != nil {
				return nil, err
			}
			opts = append(opts, scw.WithProfile(profile))
		}
	}
	if len(opts) == 0 {
		accessKey, _ := a.client.GetAccessKey()
		secretKey, _ := a.client.GetSecretKey()
		orgID, _ := a.client.GetDefaultOrganizationID()
		region, _ := a.client.GetDefaultRegion()
		zone, _ := a.client.GetDefaultZone()
		opts = []scw.ClientOption{
			scw.WithAuth(accessKey, secretKey),
			scw.WithDefaultOrganizationID(orgID),
			scw.WithDefaultRegion(region),
			scw.WithDefaultZone(zone),
		}
	}
	if ua, ok := a.client.GetUserAgent(); ok {
		opts = append(opts, scw.WithUserAgent(ua))
	}

	return scw.NewClient(append(opts, scw.WithDefaultProjectID(projectID))...)
}

// sizeCurrent fits the current screen in the content area.
func (a *App) sizeCurrent() {
	if a.width == 0 {
		return
	}
	chrome := len(scalewayLogo) + 1 // header (logo height) + status
	if a.prompt != nil && a.prompt.Active {
		chrome++
	}
	h := max(a.height-chrome, 1)
	a.current().Size(a.width, h)
}

// headerInfoWidth is the fixed width of the header's left block
// (project/organization/zone/region).
const headerInfoWidth = 46

func (a *App) viewHeader() string {
	// Left block: profile, project, organization, zone/region (top-left).
	left := []string{
		headerField("profile", a.profile),
		headerField("project", a.projectLabel()),
		headerField("organization", a.orgLabel()),
		headerField("zone", string(a.loc.GetZone())) +
			"  " + headerField("region", a.regionString()),
	}
	// Center: the current screen's key hints, at most 3 columns per row.
	centerW := a.width - headerInfoWidth - lipgloss.Width(scalewayLogo[0]) - 2
	var hints []string
	if centerW >= 10 {
		hints = layoutHints(a.current().Hints(), 3)
	}
	rows := make([]string, len(scalewayLogo))
	for i := range rows {
		var b strings.Builder
		if i < len(left) {
			b.WriteString(left[i])
		}
		line := b.String()
		if i < len(hints) {
			line += strings.Repeat(" ", max(0, headerInfoWidth-lipgloss.Width(line))) +
				hintStyle.Render(truncate(hints[i], centerW))
		}
		line += strings.Repeat(
			" ",
			max(0, a.width-lipgloss.Width(scalewayLogo[0])-lipgloss.Width(line)),
		)
		rows[i] = line + logoStyle.Render(scalewayLogo[i])
	}

	return strings.Join(rows, "\n")
}

// layoutHints groups key hints into rows of at most cols entries,
// padding each column to the widest hint of that column (shared by
// every row, so the columns stay aligned).
func layoutHints(groups []string, cols int) []string {
	rows := make([][]string, 0, (len(groups)+cols-1)/cols)
	for start := 0; start < len(groups); start += cols {
		rows = append(rows, groups[start:min(start+cols, len(groups))])
	}
	widths := make([]int, cols)
	for _, row := range rows {
		for j, g := range row {
			widths[j] = max(widths[j], lipgloss.Width(g))
		}
	}
	out := make([]string, len(rows))
	for i, row := range rows {
		cells := make([]string, len(row))
		for j, g := range row {
			cells[j] = g + strings.Repeat(" ", widths[j]-lipgloss.Width(g))
		}
		out[i] = strings.TrimRight(strings.Join(cells, "  "), " ")
	}

	return out
}

// headerField renders one left-header entry: dim key, lightsky value.
func headerField(key, value string) string {
	return dimStyle.Render(key+":") + " " + statusStyle.Render(value)
}

// projectLabel is the project name with its short ID as alias.
func (a *App) projectLabel() string {
	if a.projectName == "" {
		return shortID(a.project)
	}
	if a.project != "" {
		return truncate(a.projectName, 26) + " (" + shortID(a.project) + ")"
	}

	return a.projectName
}

// orgLabel is the organization name (alias) with its short ID.
func (a *App) orgLabel() string {
	name := a.orgName
	if a.orgAlias != "" && a.orgAlias != name {
		name += "/" + a.orgAlias
	}
	if name == "" {
		return "-"
	}
	orgID := a.orgID()
	if orgID != "" {
		return truncate(name, 24) + " (" + shortID(orgID) + ")"
	}

	return name
}

func (a *App) orgID() string {
	if a.client == nil {
		return ""
	}
	id, _ := a.client.GetDefaultOrganizationID()

	return id
}

func (a *App) regionString() string {
	zone, region := a.loc.Get()
	if region != "" {
		return string(region)
	}
	if r, err := zone.Region(); err == nil {
		return string(r)
	}

	return "-"
}

// loadNames fetches the project and organization names once, on
// startup. Failures degrade to the short IDs.
func (a *App) loadNames() tea.Cmd {
	return func() tea.Msg {
		names := namesMsg{}
		if a.client == nil {
			return names
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if orgID, _ := a.client.GetDefaultOrganizationID(); orgID != "" {
			if org, err := iam.NewAPI(a.client).GetOrganization(
				&iam.GetOrganizationRequest{OrganizationID: orgID},
				scw.WithContext(ctx),
			); err == nil {
				names.orgName, names.orgAlias = org.Name, org.Alias
			}
		}
		if p, _ := a.client.GetDefaultProjectID(); p != "" {
			if proj, err := accountv3.NewProjectAPI(a.client).GetProject(
				&accountv3.ProjectAPIGetProjectRequest{ProjectID: p},
				scw.WithContext(ctx),
			); err == nil {
				names.projectName = proj.Name
			}
		}

		return names
	}
}

func (a *App) viewPrompt() string {
	label := ": "
	switch a.prompt.Mode {
	case PromptFilter:
		label = "/ "
	case PromptSearch:
		label = "find: "
	}
	line := promptStyle.Render(label+a.prompt.Buffer) + "█"

	return pad(line, a.width)
}

func (a *App) viewStatus() string {
	parts := []string{}
	if b, ok := a.current().(*Browser); ok {
		parts = append(parts, fmt.Sprintf("%d %s", len(b.Table().Visible()), b.res.Name))
		if !b.lastUpdate.IsZero() {
			parts = append(
				parts,
				"updated "+time.Since(b.lastUpdate).Round(time.Second).String()+" ago",
			)
		}
	}
	left := statusStyle.Render(strings.Join(parts, "  "))
	right := flashStyle.Render(a.flash)
	avail := a.width - lipgloss.Width(left) - 2
	if avail < 0 {
		return truncate(left, a.width)
	}

	return left + strings.Repeat(" ", max(0, avail-lipgloss.Width(right))) + truncate(right, avail)
}
