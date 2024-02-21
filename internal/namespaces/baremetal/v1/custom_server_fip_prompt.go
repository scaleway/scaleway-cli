package baremetal

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
)

var ipTypeOption = []string{"IPv4", "IPv6"}

func promptIPFlexibleServer(ctx context.Context, req *serverAddFlexibleIPRequest) (*serverAddFlexibleIPRequest, error) {
	if !interactive.IsInteractive {
		return nil, &core.CliError{
			Err:  fmt.Errorf("failed to create and attach a new flexible IP"),
			Hint: "Missing argument 'ip-type'",
		}
	}
	var err error
	quitMessage := terminal.Style("Type Ctrl+c or q to quit", color.FgHiBlue, color.Bold)
	fmt.Println(quitMessage)
	req.IPType, err = promptChooseIPType(ctx)
	if err != nil {
		return nil, err
	}
	req.Description, err = promptAddDescriptionFip(ctx)
	if err != nil {
		return nil, err
	}
	req.Tags, err = promptAddTags(ctx)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func promptAddTags(ctx context.Context) ([]string, error) {
	_, _ = interactive.Println()
	promptConfig := interactive.PromptStringConfig{
		Ctx:    ctx,
		Prompt: "Enter all the tags you want to associate to your flexible IP separate by a space",
	}
	res, err := interactive.PromptStringWithConfig(&promptConfig)
	_, _ = interactive.Println()
	if err != nil {
		return nil, err
	}
	return strings.Split(res, " "), nil
}

func promptAddDescriptionFip(ctx context.Context) (string, error) {
	_, _ = interactive.Println()
	promptConfig := interactive.PromptStringConfig{
		Ctx:    ctx,
		Prompt: "Enter a description for you flexible ip (max 255 words)",
	}
	res, err := interactive.PromptStringWithConfig(&promptConfig)
	if err != nil {
		return "", err
	}
	return res, nil
}

func promptChooseIPType(ctx context.Context) (string, error) {
	prompt := interactive.ListPrompt{
		Prompt:       "Choose your internet protocol",
		Choices:      ipTypeOption,
		DefaultIndex: 0,
	}
	ipTypeIndex, err := prompt.Execute(ctx)
	if err != nil {
		return "", err
	}
	return ipTypeOption[ipTypeIndex], nil
}
