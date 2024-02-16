package mnq

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	mnq "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type natsContext struct {
	Description string `json:"description"`
	URL         string `json:"url"`

	// CredentialsPath is a path to file containing credentials
	CredentialsPath string `json:"creds"`
}

func natsContextFrom(account *mnq.NatsAccount, credsPath string) []byte {
	ctx := &natsContext{
		Description:     "Nats context created by Scaleway CLI",
		URL:             account.Endpoint,
		CredentialsPath: credsPath,
	}
	b, _ := json.Marshal(ctx)

	return b
}

type CreateContextRequest struct {
	NatsAccountID   string
	ContextName     string
	CredentialsName string
	Region          scw.Region
}

func aliasCreateContextCommand() *core.Command {
	return &core.Command{
		Short:     "Create a new context for natscli",
		Namespace: "mnq",
		Resource:  "nats",
		Verb:      "create-context",
		Groups:    []string{"workflow"},
		Long: `This command help you configure your nats cli
Contexts should are stored in $HOME/.config/nats/context
Credentials and context file are saved in your nats context folder with 0600 permissions`,
		Examples: []*core.Example{
			{
				Short: "Create a custom alias 'isl' for 'instance server list'", // TODO
				Raw:   `scw alias create isl command="instance server list""`,
			},
			{
				Short: "Add an alias to a verb",
				Raw:   `scw alias create c command=create`,
			},
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "nats-account-id", // TODO should be optional and an interactive list should show accounts
				Required:   true,
				Positional: true,
				Short:      "Alias name",
			},
			{
				Name:  "name",
				Short: "Name of the saved context, defaults to account name",
			},
			{
				Name:    "credentials-name",
				Short:   "Name of the created credentials",
				Default: core.RandomValueGenerator("mnq"),
			},
			core.RegionArgSpec((*mnq.NatsAPI)(nil).Regions()...),
		},
		ArgsType: reflect.TypeOf(CreateContextRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*CreateContextRequest)
			api := mnq.NewNatsAPI(core.ExtractClient(ctx))

			natsAccount, err := api.GetNatsAccount(&mnq.NatsAPIGetNatsAccountRequest{
				Region:        args.Region,
				NatsAccountID: args.NatsAccountID,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to get nats account: %w", err)
			}

			creds, err := api.CreateNatsCredentials(&mnq.NatsAPICreateNatsCredentialsRequest{
				Region:        args.Region,
				NatsAccountID: args.NatsAccountID,
				Name:          args.CredentialsName,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to create nats credentials: %w", err)
			}

			homeDir := core.ExtractEnv(ctx, "HOME") // TODO: XDG_HOME, allow override
			natsContextDir := filepath.Join(homeDir, ".config", "nats", "context")

			credsPath := filepath.Join(natsContextDir, creds.Name+".creds")
			err = os.WriteFile(credsPath, []byte(creds.Credentials.Content), 0600)
			if err != nil {
				return nil, &core.CliError{
					Err:     err,
					Message: fmt.Sprintf("Failed to write credentials into file %q", credsPath),
					Details: fmt.Sprintf("You may want to delete created credentials %q", creds.Name),
					Code:    1,
				}
			}

			contextPath := filepath.Join(natsContextDir, natsAccount.Name+".json")
			err = os.WriteFile(contextPath, natsContextFrom(natsAccount, credsPath), 0600)
			if err != nil {
				return nil, &core.CliError{
					Err:     err,
					Message: fmt.Sprintf("Failed to write context into file %q", contextPath),
					Details: fmt.Sprintf("You may want to delete created credentials %q", creds.Name),
					Code:    1,
				}
			}

			return &core.SuccessResult{
				Message:  "Nats context successfully created",
				Details:  fmt.Sprintf("Select context using `nats context select %s`", natsAccount.Name),
				Resource: contextPath,
			}, nil
		},
	}
}
