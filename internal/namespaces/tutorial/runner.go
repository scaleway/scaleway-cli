//go:build !wasm

package tutorial

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
)

func RunTutorial(ctx context.Context, tutorial Tutorial, startStep int) error {
	if !interactive.IsInteractive {
		return runTutorialNonInteractive(ctx, tutorial)
	}

	progress, err := LoadProgress(ctx)
	if err != nil {
		return fmt.Errorf("failed to load progress: %w", err)
	}

	err = CheckPrerequisites(ctx, tutorial)
	if err != nil {
		return err
	}

	now := time.Now()
	p, hasProgress := GetProgress(progress, tutorial.ID)
	if !hasProgress || startStep == 0 {
		p = TutorialProgress{
			TutorialID:     tutorial.ID,
			CurrentStep:    startStep,
			StartedAt:      now,
			LastAccessedAt: now,
			CLIVersion:     core.ExtractBuildInfo(ctx).Version.String(),
		}
		SetProgress(progress, p)
	} else {
		p.LastAccessedAt = now
		SetProgress(progress, p)
	}

	err = SaveProgress(ctx, progress)
	if err != nil {
		interactive.Printf("Warning: failed to save progress: %v\n", err)
	}

	for i := startStep; i < len(tutorial.Steps); i++ {
		step := tutorial.Steps[i]

		interactive.Println()
		interactive.Println(
			terminal.Style(
				fmt.Sprintf("Step %d/%d: %s", i+1, len(tutorial.Steps), step.Title),
				color.FgCyan,
				color.Bold,
			),
		)
		interactive.Println()
		interactive.PrintlnWithoutIndent(step.Concept)
		interactive.Println()

		if step.Command != "" {
			interactive.Println(terminal.Style("Command:", color.FgYellow))
			interactive.Println(terminal.Style("  "+step.Command, color.FgGreen))
			interactive.Println()

			if step.ExpectedOutput != "" {
				interactive.Println("Expected output: " + step.ExpectedOutput)
				interactive.Println()
			}

			runCmd, err := interactive.PromptBool(ctx, "Run this command?", true)
			if err != nil {
				return err
			}

			if runCmd {
				err := executeTutorialCommand(ctx, step)
				if err != nil {
					interactive.Println(terminal.Style("Error: "+err.Error(), color.FgRed))
					interactive.Println()

					if step.ErrorHint != "" {
						interactive.Println("Hint: " + step.ErrorHint)
						interactive.Println()
					}

					choice, err := promptFailureChoice(ctx)
					if err != nil {
						return err
					}

					switch choice {
					case "retry":
						i--

						continue
					case "skip":
						// continue to next step
					case "quit":
						p.CurrentStep = i
						p.LastAccessedAt = time.Now()
						SetProgress(progress, p)
						err = SaveProgress(ctx, progress)
						if err != nil {
							interactive.Printf("Warning: failed to save progress: %v\n", err)
						}
						interactive.Println(
							"Tutorial progress saved. Run 'scw tutorial resume' to continue.",
						)

						return nil
					}
				}
			} else {
				interactive.Println("You can run it manually: " + step.Command)
				interactive.Println()
			}
		} else {
			interactive.Print("Press Enter to continue...")
			_, _ = interactive.Readline(ctx, &interactive.ReadlineConfig{})
		}

		p.CurrentStep = i + 1
		p.LastAccessedAt = time.Now()
		SetProgress(progress, p)
		err = SaveProgress(ctx, progress)
		if err != nil {
			interactive.Printf("Warning: failed to save progress: %v\n", err)
		}

		if i < len(tutorial.Steps)-1 {
			continueTut, err := interactive.PromptBool(ctx, "Continue?", true)
			if err != nil {
				return err
			}
			if !continueTut {
				p.CurrentStep = i + 1
				SetProgress(progress, p)
				err = SaveProgress(ctx, progress)
				if err != nil {
					interactive.Printf("Warning: failed to save progress: %v\n", err)
				}
				interactive.Println(
					"Tutorial progress saved. Run 'scw tutorial resume' to continue.",
				)

				return nil
			}
		}
	}

	p.Completed = true
	p.LastAccessedAt = time.Now()
	SetProgress(progress, p)
	err = SaveProgress(ctx, progress)
	if err != nil {
		interactive.Printf("Warning: failed to save progress: %v\n", err)
	}

	interactive.Println()
	interactive.Println(terminal.Style("Tutorial completed!", color.FgGreen, color.Bold))

	return nil
}

func executeTutorialCommand(ctx context.Context, step TutorialStep) error {
	parts := strings.Fields(step.Command)
	if len(parts) == 0 {
		return errors.New("empty command")
	}

	binaryName := core.ExtractBinaryName(ctx)
	if binaryName == "" {
		binaryName = "scw"
	}

	if parts[0] == "scw" || parts[0] == binaryName {
		parts = parts[1:]
	}

	cmdPath, err := exec.LookPath(binaryName)
	if err != nil {
		cmdPath = os.Args[0]
	}

	cmd := exec.Command(cmdPath, parts...)
	exitCode, err := core.ExecCmd(ctx, cmd)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("command exited with code %d", exitCode)
	}

	return nil
}

func promptFailureChoice(ctx context.Context) (string, error) {
	lp := &interactive.ListPrompt{
		Prompt:  "What would you like to do?",
		Choices: []string{"retry", "skip", "quit"},
	}
	idx, err := lp.Execute(ctx)
	if err != nil {
		return "quit", err
	}

	return []string{"retry", "skip", "quit"}[idx], nil
}

func CheckPrerequisites(ctx context.Context, tutorial Tutorial) error {
	for _, prereq := range tutorial.Prerequisites {
		if prereq == "credentials-configured" {
			client := core.ExtractClient(ctx)
			if client == nil {
				runInit, err := interactive.PromptBool(
					ctx,
					"You need to run 'scw init' first. Would you like to do that now?",
					false,
				)
				if err != nil {
					return err
				}
				if runInit {
					return errors.New("please run 'scw init' and then restart the tutorial")
				}

				return errors.New("credentials must be configured to run this tutorial")
			}
		}
	}

	return nil
}

func CheckVersionMismatch(
	ctx context.Context,
	progress ProgressMap,
	tutorialID string,
) (bool, error) {
	p, hasProgress := GetProgress(progress, tutorialID)
	if !hasProgress {
		return false, nil
	}

	currentVersion := core.ExtractBuildInfo(ctx).Version.String()
	if p.CLIVersion != currentVersion {
		restart, err := interactive.PromptBool(
			ctx,
			fmt.Sprintf(
				"This tutorial was started with CLI version %s but you are now running version %s. Restart from the beginning?",
				p.CLIVersion,
				currentVersion,
			),
			true,
		)
		if err != nil {
			return false, err
		}

		return restart, nil
	}

	return false, nil
}

func runTutorialNonInteractive(_ context.Context, tutorial Tutorial) error {
	for i, step := range tutorial.Steps {
		fmt.Printf("Step %d/%d: %s\n\n", i+1, len(tutorial.Steps), step.Title)
		fmt.Println(step.Concept)
		fmt.Println()

		if step.Command != "" {
			fmt.Println("Command: " + step.Command)
			fmt.Println()
			if step.ExpectedOutput != "" {
				fmt.Println("Expected output: " + step.ExpectedOutput)
				fmt.Println()
			}
		}

		fmt.Println("---")
	}

	fmt.Println("Tutorial completed.")

	return nil
}
