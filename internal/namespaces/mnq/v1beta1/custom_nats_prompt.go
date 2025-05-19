package mnq

import (
	"context"
	"errors"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	mnq "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
)

func promptNatsAccounts(
	ctx context.Context,
	natsAccounts []*mnq.NatsAccount,
	totalCount uint64,
) (*mnq.NatsAccount, error) {
	if totalCount == 0 {
		return nil, errors.New(
			"no nats account found, please create a NATS account with 'scw mnq nats create-account'",
		)
	}

	if !interactive.IsInteractive {
		return nil, errors.New(
			"failed to create NATS context: Multiple NATS accounts found. Please provide an account ID explicitly as the command is not running in interactive mode",
		)
	}
	if totalCount == 1 {
		return natsAccounts[0], nil
	}

	defaultIndex := 0
	natsAccountsName := make([]string, len(natsAccounts))
	for i := range natsAccounts {
		natsAccountsName[i] = fmt.Sprintf("%s %s", natsAccounts[i].Name, natsAccounts[i].Region)
	}
	prompt := interactive.ListPrompt{
		Prompt:       "Choose your nats account",
		Choices:      natsAccountsName,
		DefaultIndex: defaultIndex,
	}
	_, _ = interactive.Println()
	index, err := prompt.Execute(ctx)
	if err != nil {
		return nil, err
	}

	return natsAccounts[index], nil
}

func promptOverWriteFile(ctx context.Context, filePath string) (bool, error) {
	if !interactive.IsInteractive {
		return false, errors.New("file Exist")
	}

	config := interactive.PromptBoolConfig{
		Ctx:          ctx,
		Prompt:       "The file " + filePath + " already exists. Do you want to overwrite it?",
		DefaultValue: true,
	}
	overWrite, _ := interactive.PromptBoolWithConfig(&config)

	return overWrite, nil
}
