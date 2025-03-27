package baremetal

import (
	"context"
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
)

var ipTypeOption = []string{"IPv4", "IPv6"}

func promptIPFlexibleServer(
	ctx context.Context,
	req *serverAddFlexibleIPRequest,
) (*serverAddFlexibleIPRequest, error) {
	if !interactive.IsInteractive {
		return nil, &core.CliError{
			Err:  errors.New("failed to create and attach a new flexible IP"),
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

	return req, nil
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
