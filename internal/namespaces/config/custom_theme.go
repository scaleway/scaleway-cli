package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	tint "github.com/lrstanley/bubbletint/v2"
)

// ThemeInfo is the output struct for 'scw config theme list'.
type ThemeInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Dark   bool   `json:"dark"`
	Active bool   `json:"active"`
}

func (t *ThemeInfo) MarshalHuman() (string, error) {
	if t.Active {
		return terminal.Style(fmt.Sprintf("%-25s %-25s %s", t.ID, t.Name, themeTypeLabel(t.Dark)), color.Bold), nil
	}
	return fmt.Sprintf("%-25s %-25s %s", t.ID, t.Name, themeTypeLabel(t.Dark)), nil
}

func themeTypeLabel(dark bool) string {
	if dark {
		return "dark"
	}
	return "light"
}

func configThemeRoot() *core.Command {
	return &core.Command{
		Groups:    []string{"config"},
		Short:     `Color theme management`,
		Namespace: "config",
		Resource:  "theme",
	}
}

func configThemeListCommand() *core.Command {
	type configThemeListArgs struct {
		All bool
	}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `List available color themes`,
		Namespace:            "config",
		Resource:             "theme",
		Verb:                 "list",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeFor[configThemeListArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:  "all",
				Short: "Show all available themes (340+) instead of the curated subset",
			},
		},
		Run: func(_ context.Context, argsI any) (any, error) {
			args := argsI.(*configThemeListArgs)

			var tints []*tint.Tint
			if args.All {
				tints = tint.Tints()
			} else {
				curatedSet := make(map[string]bool)
				for _, id := range terminal.CuratedThemeIDs {
					curatedSet[id] = true
				}
				for _, t := range tint.Tints() {
					if curatedSet[t.ID] {
						tints = append(tints, t)
					}
				}
			}

			currentID := ""
			if terminal.IsThemeActive() {
				currentID = tint.Current().ID
			}

			result := make([]ThemeInfo, 0, len(tints))
			for _, t := range tints {
				result = append(result, ThemeInfo{
					ID:     t.ID,
					Name:   t.DisplayName,
					Dark:   t.Dark,
					Active: t.ID == currentID,
				})
			}

			return result, nil
		},
	}
}

func configThemeSetCommand() *core.Command {
	type configThemeSetArgs struct {
		ThemeID string
	}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Set the active color theme`,
		Namespace:            "config",
		Resource:             "theme",
		Verb:                  "set",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeFor[configThemeSetArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "theme-id",
				Short:      "The ID of the theme to activate",
				Required:   true,
				Positional: true,
				AutoCompleteFunc: func(_ context.Context, _ string, _ any) core.AutocompleteSuggestions {
					var values core.AutocompleteSuggestions
					for _, t := range tint.Tints() {
						values = append(values, t.ID)
					}
					return values
				},
				ValidateFunc: func(_ *core.ArgSpec, value any) error {
					themeID := value.(string)
					if _, ok := tint.GetTint(themeID); !ok {
						return fmt.Errorf("theme %q not found. Run 'scw config theme list' to see available themes.", themeID)
					}
					return nil
				},
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*configThemeSetArgs)
			themeID := args.ThemeID

			t, ok := tint.GetTint(themeID)
			if !ok {
				return nil, fmt.Errorf("theme %q not found. Run 'scw config theme list' to see available themes.", themeID)
			}

			cliCfg := core.ExtractCliConfig(ctx)
			if cliCfg == nil {
				return nil, fmt.Errorf("failed to get CLI config")
			}
			cliCfg.Theme = themeID
			if err := cliCfg.Save(); err != nil {
				return nil, fmt.Errorf("failed to save CLI config: %w", err)
			}

			tint.SetTintID(themeID)
			terminal.SetThemeActive(true)

			return &core.SuccessResult{
				Message: fmt.Sprintf("theme set to %q (%s)", t.DisplayName, themeID),
			}, nil
		},
	}
}

func configThemeGetCommand() *core.Command {
	type configThemeGetArgs struct{}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Show the currently active color theme`,
		Namespace:            "config",
		Resource:             "theme",
		Verb:                  "get",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeFor[configThemeGetArgs](),
		Run: func(ctx context.Context, _ any) (any, error) {
			cliCfg := core.ExtractCliConfig(ctx)
			if cliCfg == nil || cliCfg.Theme == "" {
				if t, ok := tint.GetTint("scaleway"); ok {
					return &core.SuccessResult{
						Message: fmt.Sprintf("Active theme: %s (%s)", t.DisplayName, t.ID),
					}, nil
				}
				return &core.SuccessResult{
					Message: "Active theme: default (no theme selected)",
				}, nil
			}

			t, ok := tint.GetTint(cliCfg.Theme)
			if !ok {
				return &core.SuccessResult{
					Message: fmt.Sprintf("Active theme: unknown (%s)", cliCfg.Theme),
				}, nil
			}

			return &core.SuccessResult{
				Message: fmt.Sprintf("Active theme: %s (%s)", t.DisplayName, t.ID),
			}, nil
		},
	}
}

func configThemePreviewCommand() *core.Command {
	type configThemePreviewArgs struct {
		ThemeID string
	}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Preview a color theme`,
		Namespace:            "config",
		Resource:             "theme",
		Verb:                  "preview",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeFor[configThemePreviewArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "theme-id",
				Short:      "The ID of the theme to preview",
				Required:   true,
				Positional: true,
				AutoCompleteFunc: func(_ context.Context, _ string, _ any) core.AutocompleteSuggestions {
					var values core.AutocompleteSuggestions
					for _, t := range tint.Tints() {
						values = append(values, t.ID)
					}
					return values
				},
				ValidateFunc: func(_ *core.ArgSpec, value any) error {
					themeID := value.(string)
					if _, ok := tint.GetTint(themeID); !ok {
						return fmt.Errorf("theme %q not found. Run 'scw config theme list' to see available themes.", themeID)
					}
					return nil
				},
			},
		},
		Run: func(_ context.Context, argsI any) (any, error) {
			args := argsI.(*configThemePreviewArgs)
			themeID := args.ThemeID

			t, ok := tint.GetTint(themeID)
			if !ok {
				return nil, fmt.Errorf("theme %q not found. Run 'scw config theme list' to see available themes.", themeID)
			}

			var wasActive bool
			var prevID string
			if terminal.IsThemeActive() {
				wasActive = true
				prevID = tint.Current().ID
			}

			tint.SetTintID(themeID)
			terminal.SetThemeActive(true)

			lines := []string{
				fmt.Sprintf("Theme: %s (%s)\n", t.DisplayName, t.ID),
				terminal.Style("✅ Success    Operation completed successfully", color.FgGreen),
				terminal.Style("❌ Error      Something went wrong", color.FgRed),
				terminal.Style("⚠  Warning    This action is irreversible", color.FgYellow),
				terminal.Style("ℹ  Info       Server is running", color.FgCyan),
				fmt.Sprintf("%s / %s", terminal.Style("true", color.FgGreen), terminal.Style("false", color.FgRed)),
				terminal.Style("Header     Server Info", color.Bold),
				fmt.Sprintf("%s / %s / %s",
					terminal.Style("running", color.FgGreen),
					terminal.Style("stopped", color.FgYellow),
					terminal.Style("error", color.FgRed),
				),
			}

			if wasActive {
				tint.SetTintID(prevID)
				terminal.SetThemeActive(true)
			} else {
				terminal.SetThemeActive(false)
			}

			return core.RawResult(append([]byte(strings.Join(lines, "\n")), '\n')), nil
		},
	}
}

func configThemeImportCommand() *core.Command {
	type configThemeImportArgs struct {
		File string
	}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Import a custom color theme from a JSON file`,
		Namespace:            "config",
		Resource:             "theme",
		Verb:                  "import",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeFor[configThemeImportArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "file",
				Short:      "Path to the JSON theme file",
				Required:   true,
				Positional: true,
				CanLoadFile: true,
			},
		},
		Run: func(_ context.Context, argsI any) (any, error) {
			args := argsI.(*configThemeImportArgs)

			data, err := os.ReadFile(args.File)
			if err != nil {
				return nil, fmt.Errorf("failed to read theme file: %w", err)
			}

			var customTint tint.Tint
			if err := json.Unmarshal(data, &customTint); err != nil {
				return nil, fmt.Errorf("failed to parse theme file: %w", err)
			}

			if customTint.ID == "" {
				return nil, fmt.Errorf("theme file must contain a non-empty \"id\" field")
			}
			if customTint.DisplayName == "" {
				return nil, fmt.Errorf("theme file must contain a non-empty \"display_name\" field")
			}

			if _, exists := tint.GetTint(customTint.ID); exists {
				return nil, fmt.Errorf("theme ID %q already exists. Choose a different ID.", customTint.ID)
			}

			tint.Register(&customTint)

			return &core.SuccessResult{
				Message: fmt.Sprintf("custom theme %q (%s) imported successfully. Run 'scw config theme set %s' to activate it.", customTint.DisplayName, customTint.ID, customTint.ID),
			}, nil
		},
	}
}
