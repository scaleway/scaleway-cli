package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// deleteRequestMsg asks the app to delete a row of a resource (the app
// owns the delete command and the session locality).
type deleteRequestMsg struct {
	res *Resource
	id  string
}

// deleteDoneMsg reports the outcome of the delete call.
type deleteDoneMsg struct {
	res *Resource
	id  string
	err error
}

// confirmState is the pending delete confirmation: one modal line above
// the status bar.
type confirmState struct {
	res *Resource
	id  string
}

// handleDelete routes the 'd' key delete chain: the request (the
// confirmation modal, or the delete call directly when
// tui.delete_confirm is false) and the delete call.
func (a *App) handleDelete(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case deleteRequestMsg:
		if !m.res.Deleteable() {
			a.flash = "delete not available for " + m.res.Name

			return nil
		}
		if !a.confirmDelete {
			a.flash = "deleting " + m.res.Name + " " + shortID(m.id) + "…"

			return a.runDelete(m.res, m.id)
		}
		a.confirm = &confirmState{res: m.res, id: m.id}
		a.sizeCurrent()

		return nil
	case deleteDoneMsg:
		if m.err != nil {
			a.flash = "delete " + m.res.Name + ": " + m.err.Error()
		} else {
			a.flash = "deleted " + m.res.Name + " " + shortID(m.id)
		}
		// The resource is gone: re-fetch the current screen.
		mdl, cmd := a.current().Update(refreshMsg{})
		a.replaceCurrent(mdl)

		return cmd
	}

	return nil
}

// updateConfirmKey handles a keystroke while the delete confirmation is
// open: y/enter delete, n/esc cancel, everything else is swallowed.
//
// ponytail: the confirm keys (y/n) are hardcoded, not keymap actions;
// add them to the keymap if overridable confirm keys are wanted.
func (a *App) updateConfirmKey(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case k.Type == tea.KeyEnter || k.String() == "y":
		res, id := a.confirm.res, a.confirm.id
		a.confirm = nil
		a.sizeCurrent()
		a.flash = "deleting " + res.Name + " " + shortID(id) + "…"

		return a, a.runDelete(res, id)
	case k.Type == tea.KeyEscape || k.String() == "n":
		a.confirm = nil
		a.sizeCurrent()
	}

	return a, nil
}

// runDelete deletes res id, built the TUI way: the positional id and
// the session localities (newRequestArgs, same as the edit flow).
func (a *App) runDelete(res *Resource, id string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
		defer cancel()
		args := newRequestArgs(res.DeleteCmd, a.loc, id)
		_, err := runCommand(res.DeleteCmd, a.meta, a.loc, ctx, args)

		return deleteDoneMsg{res: res, id: id, err: err}
	}
}

// viewConfirm renders the delete confirmation line.
func (a *App) viewConfirm() string {
	line := confirmStyle.Render(fmt.Sprintf(
		"delete %s %s?  [ y ] delete   [ n ] cancel",
		a.confirm.res.Name, shortID(a.confirm.id),
	))

	return pad(line, a.width)
}
