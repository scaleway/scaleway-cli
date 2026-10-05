package tutorial

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
)

type TutorialArgs struct {
	TutorialID string
}

func GetCommands() *core.Commands {
	return core.NewCommands(
		tutorialRoot(),
		tutorialResume(),
		tutorialList(),
	)
}

func tutorialRoot() *core.Command {
	return &core.Command{
		Short:                "Follow guided tutorials to learn the CLI",
		Long:                 "Browse and start interactive tutorials that walk you through common Scaleway CLI tasks step by step. Each tutorial explains concepts, shows commands, and lets you run them with confirmation.",
		Namespace:            "tutorial",
		AllowAnonymousClient: true,
		Groups:               []string{"utility"},
		ArgsType:             reflect.TypeFor[TutorialArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tutorial-id",
				Short:      "Tutorial to start",
				Positional: true,
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*TutorialArgs)

			if args.TutorialID != "" {
				tut, ok := GetTutorial(args.TutorialID)
				if !ok {
					return nil, &core.CliError{
						Err:  fmt.Errorf("Tutorial %q not found", args.TutorialID),
						Hint: "Run \"scw tutorial list\" to see available tutorials.",
					}
				}

				progress, err := LoadProgress(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to load progress: %w", err)
				}

				startStep := 0
				if existing, hasProgress := GetProgress(
					progress,
					tut.ID,
				); hasProgress &&
					!existing.Completed {
					restart, err := CheckVersionMismatch(ctx, progress, tut.ID)
					if err != nil {
						return nil, err
					}
					if restart {
						ResetProgress(progress, tut.ID)
						startStep = 0
					} else {
						startStep = existing.CurrentStep
					}
				}

				err = RunTutorial(ctx, tut, startStep)
				if err != nil {
					return nil, err
				}

				return &core.SuccessResult{
					Message: "Tutorial finished",
				}, nil
			}

			if interactive.IsInteractive {
				tutorials := ListTutorials()
				labels := make([]string, len(tutorials))
				for i, t := range tutorials {
					label := t.Title
					if t.Recommended {
						label += " (recommended)"
					}
					label += " - " + t.Description
					labels[i] = label
				}

				lp := &interactive.ListPrompt{
					Prompt:  "Available tutorials (select one to start):",
					Choices: labels,
				}
				idx, err := lp.Execute(ctx)
				if err != nil {
					return nil, err
				}

				selected := tutorials[idx]
				progress, err := LoadProgress(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to load progress: %w", err)
				}

				startStep := 0
				if existing, hasProgress := GetProgress(
					progress,
					selected.ID,
				); hasProgress &&
					!existing.Completed {
					restart, err := CheckVersionMismatch(ctx, progress, selected.ID)
					if err != nil {
						return nil, err
					}
					if restart {
						ResetProgress(progress, selected.ID)
						startStep = 0
					} else {
						startStep = existing.CurrentStep
					}
				}

				err = RunTutorial(ctx, selected, startStep)
				if err != nil {
					return nil, err
				}

				return &core.SuccessResult{
					Message: "Tutorial finished",
				}, nil
			}

			tutorials := ListTutorials()
			var sb strings.Builder
			sb.WriteString("Available tutorials:\n\n")
			for _, t := range tutorials {
				marker := "  "
				if t.Recommended {
					marker = "* "
				}
				sb.WriteString(
					fmt.Sprintf(
						"%s%s - %s (%s)\n",
						marker,
						t.Title,
						t.Description,
						t.EstimatedDuration,
					),
				)
			}
			sb.WriteString("\nRun 'scw tutorial <title>' to start a tutorial.")

			return sb.String(), nil
		},
	}
}

func tutorialResume() *core.Command {
	return &core.Command{
		Short:                "Resume an interrupted tutorial",
		Namespace:            "tutorial",
		Resource:             "resume",
		AllowAnonymousClient: true,
		Groups:               []string{"utility"},
		ArgsType:             reflect.TypeOf(struct{}{}),
		ArgSpecs:             core.ArgSpecs{},
		Run: func(ctx context.Context, _ any) (any, error) {
			progress, err := LoadProgress(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to load progress: %w", err)
			}

			var lastTutorial *TutorialProgress
			var lastTutorialID string
			for id, p := range progress {
				if p.Completed {
					continue
				}
				if lastTutorial == nil || p.LastAccessedAt.After(lastTutorial.LastAccessedAt) {
					progCopy := p
					lastTutorial = &progCopy
					lastTutorialID = id
				}
			}

			if lastTutorial == nil {
				return nil, &core.CliError{
					Err:  errors.New("No tutorial in progress"),
					Hint: "Run \"scw tutorial\" to start one.",
				}
			}

			tut, ok := GetTutorial(lastTutorialID)
			if !ok {
				return nil, &core.CliError{
					Err:  fmt.Errorf("Tutorial %q not found", lastTutorialID),
					Hint: "This tutorial may have been removed. Run \"scw tutorial list\" to see available tutorials.",
				}
			}

			restart, err := CheckVersionMismatch(ctx, progress, tut.ID)
			if err != nil {
				return nil, err
			}

			startStep := lastTutorial.CurrentStep
			if restart {
				ResetProgress(progress, tut.ID)
				startStep = 0
			}

			err = RunTutorial(ctx, tut, startStep)
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Message: "Tutorial resumed and finished",
			}, nil
		},
	}
}

func tutorialList() *core.Command {
	return &core.Command{
		Short:                "List available tutorials",
		Namespace:            "tutorial",
		Resource:             "list",
		AllowAnonymousClient: true,
		Groups:               []string{"utility"},
		ArgsType:             reflect.TypeOf(struct{}{}),
		ArgSpecs:             core.ArgSpecs{},
		Run: func(ctx context.Context, _ any) (any, error) {
			tutorials := ListTutorials()
			var sb strings.Builder
			sb.WriteString("Available tutorials:\n\n")
			for _, t := range tutorials {
				marker := "  "
				if t.Recommended {
					marker = "* "
				}
				sb.WriteString(
					fmt.Sprintf(
						"%s%s - %s (%s, %s)\n",
						marker,
						t.Title,
						t.Description,
						t.Difficulty,
						t.EstimatedDuration,
					),
				)
			}
			sb.WriteString("\n* = recommended for new users")
			sb.WriteString("\nRun 'scw tutorial <id>' to start a tutorial.")

			return sb.String(), nil
		},
	}
}
