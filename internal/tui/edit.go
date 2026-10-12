package tui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ghodss/yaml"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-cli/v2/internal/config"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
)

// editRequestMsg asks the app to edit a row of a resource (the app
// owns the get/update commands and the session locality).
type editRequestMsg struct {
	res *Resource
	id  string
}

// editReadyMsg carries the temp file to open in the editor (or its
// error). tmpl is the update request filled with the id and localities
// (the json:"-" fields, not present in the file) to restore on save.
type editReadyMsg struct {
	file string
	res  *Resource
	tmpl any
	err  error
}

// editDoneMsg reports the exit of the editor with the edited content.
type editDoneMsg struct {
	file    string
	res     *Resource
	tmpl    any
	content []byte
	err     error
}

// editUpdatedMsg reports the outcome of the update call.
type editUpdatedMsg struct {
	res *Resource
	err error
}

// editorPathAndArgs resolves the editor from $VISUAL then $EDITOR (the
// same resolution as the CLI `edit` commands).
func editorPathAndArgs(fileName string) (string, []string) {
	editorAndArguments := strings.Fields(config.GetDefaultEditor())
	cmdArgs := []string{fileName}
	if len(editorAndArguments) > 1 {
		cmdArgs = append(editorAndArguments[1:], cmdArgs...)
	}

	return editorAndArguments[0], cmdArgs
}

// newRequestArgs builds the request of cmd for the session locality,
// with the positional id argument set (same mapping as the CLI:
// strcase.ToPublicGoName on the arg spec name).
func newRequestArgs(cmd *core.Command, loc *Locality, id string) any {
	if cmd.ArgsType == nil {
		return nil
	}
	req := reflect.New(cmd.ArgsType).Interface()
	zone, region := loc.Get()
	fillLocalities(req, zone, region)
	if id != "" {
		if pos := cmd.ArgSpecs.GetPositionalArg(); pos != nil {
			_ = args.UnmarshalStruct([]string{pos.Name + "=" + id}, req)
		}
	}

	return req
}

// editDocument builds the file shown in the editor: the update request
// of the resource marshaled in YAML, each field commented with the doc
// of the matching `update` command argument. tmpl is the same request
// filled with the id and localities (json:"-" fields, not present in
// the file): it is restored on the edited content before the update.
func editDocument(res *Resource, loc *Locality, id string, resource any) ([]byte, any, error) {
	tmpl := newRequestArgs(res.UpdateCmd, loc, id)
	resourceV := reflect.Indirect(reflect.ValueOf(resource))
	if !resourceV.IsValid() {
		return nil, nil, fmt.Errorf("empty response for %s %s", res.Name, id)
	}
	// ponytail: the resource->request mapping copies same-named fields
	// only (editor.ValueMapper), so a custom update command embedding
	// the SDK request maps only its own fields. Same ceiling as the
	// CLI `edit` commands; add a tag-aware mapper if it bites.
	editor.ValueMapper(reflect.ValueOf(tmpl).Elem(), resourceV)
	data, err := yaml.Marshal(tmpl)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal %s: %w", res.Name, err)
	}

	return withDocComments(data, res.Name, id, res.UpdateCmd), tmpl, nil
}

// withDocComments prefixes each top-level YAML field with a comment
// holding the doc of the matching update argument.
func withDocComments(data []byte, name, id string, cmd *core.Command) []byte {
	// YAML keys are the json tags (underscores); arg names are the
	// same with dashes. Nested args ("tags.{index}") index their base
	// key so slice/map fields get a doc too.
	docs := map[string]string{}
	for _, spec := range cmd.ArgSpecs {
		if spec == nil || spec.Short == "" {
			continue
		}
		if _, ok := docs[spec.Name]; !ok {
			docs[spec.Name] = spec.Short
		}
		if base, _, ok := strings.Cut(spec.Name, "."); ok {
			if _, exists := docs[base]; !exists {
				docs[base] = spec.Short
			}
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# %s %s\n", name, id)
	b.WriteString("# Edit the fields below, then save and quit to apply the update.\n")
	b.WriteString("# Quit without saving to cancel.\n")
	for line := range strings.SplitSeq(strings.TrimSuffix(string(data), "\n"), "\n") {
		if key, isField := topLevelYAMLKey(line); isField {
			if doc, ok := docs[strings.ReplaceAll(key, "_", "-")]; ok {
				b.WriteString("# " + doc + "\n")
			}
		}
		b.WriteString(line)
		b.WriteString("\n")
	}

	return []byte(b.String())
}

// topLevelYAMLKey returns the key of a top-level "key:" YAML line
// (lines with any indentation are nested content, not fields).
func topLevelYAMLKey(line string) (string, bool) {
	key, _, ok := strings.Cut(line, ":")
	if !ok || key == "" {
		return "", false
	}
	for _, r := range key {
		if !(r == '_' ||
			(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9')) {
			return "", false
		}
	}

	return key, true
}

// finalUpdateRequest unmarshals the edited YAML into a fresh update
// request and restores the path parameters (id, zone, region) that
// carry json:"-" and are not present in the file.
func finalUpdateRequest(cmd *core.Command, tmpl any, data []byte) (any, error) {
	req := reflect.New(cmd.ArgsType).Interface()
	if err := yaml.Unmarshal(data, req); err != nil {
		return nil, fmt.Errorf("failed to parse the edited file: %w", err)
	}
	editor.ValueMapper(
		reflect.ValueOf(req).Elem(),
		reflect.Indirect(reflect.ValueOf(tmpl)),
		editor.MapWithTag("-"),
	)

	return req, nil
}

// writeTempFile writes the edit content to a temp file for the editor.
func writeTempFile(content []byte) (string, error) {
	f, err := os.CreateTemp(os.TempDir(), "scw-tui-edit-*.yml")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())

		return "", fmt.Errorf("failed to write temp file %q: %w", f.Name(), err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file %q: %w", f.Name(), err)
	}

	return f.Name(), nil
}

// handleEdit routes the 'e' key edit chain: the request (fetch +
// temp file), the editor process and the update call.
func (a *App) handleEdit(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case editRequestMsg:
		if !m.res.Editable() {
			a.flash = "edit not available for " + m.res.Name

			return nil
		}

		return a.prepareEdit(m.res, m.id)
	case editReadyMsg:
		if m.err != nil {
			a.flash = "edit: " + m.err.Error()

			return nil
		}
		path, cmdArgs := editorPathAndArgs(m.file)

		// Suspend the TUI and run the editor on the terminal, like
		// the connect flow; resume (and report) when it exits.
		return tea.ExecProcess(exec.Command(path, cmdArgs...), func(err error) tea.Msg {
			content, readErr := os.ReadFile(m.file)
			if readErr != nil {
				err = readErr
			}

			return editDoneMsg{
				file: m.file, res: m.res, tmpl: m.tmpl,
				content: content, err: err,
			}
		})
	case editDoneMsg:
		if m.err != nil {
			a.flash = "edit: " + m.err.Error()

			return nil
		}
		req, err := finalUpdateRequest(m.res.UpdateCmd, m.tmpl, m.content)
		if err != nil {
			// Keep the file so the user can fix the YAML.
			a.flash = "edit: " + err.Error() + " — file kept: " + m.file

			return nil
		}
		_ = os.Remove(m.file)
		res := m.res
		a.flash = "updating " + res.Name + "…"

		return func() tea.Msg {
			ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
			defer cancel()
			_, err := runCommand(res.UpdateCmd, a.meta, a.loc, ctx, req)

			return editUpdatedMsg{res: res, err: err}
		}
	case editUpdatedMsg:
		if m.err != nil {
			a.flash = "update " + m.res.Name + ": " + m.err.Error()
		} else {
			a.flash = "updated " + m.res.Name
		}
		// The resource changed: re-fetch the current screen.
		mdl, cmd := a.current().Update(refreshMsg{})
		a.replaceCurrent(mdl)

		return cmd
	}

	return nil
}

// prepareEdit fetches the resource and writes the temp file to edit.
func (a *App) prepareEdit(res *Resource, id string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
		defer cancel()
		resource, err := runCommand(
			res.GetCmd, a.meta, a.loc, ctx, newRequestArgs(res.GetCmd, a.loc, id),
		)
		if err != nil {
			return editReadyMsg{err: err}
		}
		content, tmpl, err := editDocument(res, a.loc, id, resource)
		if err != nil {
			return editReadyMsg{err: err}
		}
		file, err := writeTempFile(content)
		if err != nil {
			return editReadyMsg{err: err}
		}

		return editReadyMsg{file: file, res: res, tmpl: tmpl}
	}
}
