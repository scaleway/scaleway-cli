package tui

// PromptMode is the mode of the prompt bar.
type PromptMode int

const (
	// PromptCmd is the `:` command mode (navigate between resources).
	PromptCmd PromptMode = iota
	// PromptFilter is the `/` filter mode (filter the current table).
	PromptFilter
	// PromptSearch is the `/` search mode (find an attribute in the
	// details view; enter copies its value to the clipboard).
	PromptSearch
)

// Prompt is the k9s `:` command bar (the model.FishBuff equivalent). It
// is driven directly by the App, not a standalone tea.Model.
type Prompt struct {
	Active bool
	Mode   PromptMode
	Buffer string
}

func newCmdPrompt() *Prompt {
	return &Prompt{Active: true, Mode: PromptCmd}
}

func newFilterPrompt() *Prompt {
	return &Prompt{Active: true, Mode: PromptFilter}
}

func newSearchPrompt() *Prompt {
	return &Prompt{Active: true, Mode: PromptSearch}
}
