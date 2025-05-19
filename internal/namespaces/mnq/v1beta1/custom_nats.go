package mnq

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	mnq "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type natsContext struct {
	Description string `json:"description"`
	URL         string `json:"url"`

	// CredentialsPath is a path to file containing credentials
	CredentialsPath string `json:"creds"`
}

type CreateContextRequest struct {
	NatsAccountID   string
	ContextName     string
	CredentialsName string
	Region          scw.Region
}

func createContextCommand() *core.Command {
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
				Short: "Create a context in your nats server",
				Raw:   `scw mnq nats create-context <nats-account-id> credentials-name=<credential-name> region=fr-par`,
			},
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:  "nats-account-id",
				Short: "ID of the NATS account",
			},
			{
				Name:  "name",
				Short: "Name of the saved context, defaults to account name",
			},
			{
				Name:  "credentials-name",
				Short: "Name of the created credentials",
			},
			core.RegionArgSpec((*mnq.NatsAPI)(nil).Regions()...),
		},
		ArgsType: reflect.TypeOf(CreateContextRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*CreateContextRequest)
			api := mnq.NewNatsAPI(core.ExtractClient(ctx))
			natsAccount, err := getNatsAccountID(ctx, args, api)
			if err != nil {
				return nil, err
			}

			var credentialsName string
			if args.CredentialsName != "" {
				credentialsName = args.CredentialsName
			} else {
				credentialsName = natsAccount.Name + core.GetRandomName("creds")
			}
			credentials, err := api.CreateNatsCredentials(&mnq.NatsAPICreateNatsCredentialsRequest{
				Region:        args.Region,
				NatsAccountID: natsAccount.ID,
				Name:          credentialsName,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}
			contextPath, err := saveNATSCredentials(ctx, credentials, natsAccount)
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Message: "Nats context successfully created",
				Details: fmt.Sprintf(
					"%s nats credentials was created\nSelect context using `nats context select %s`",
					credentials.Name,
					natsAccount.Name,
				),
				Resource: contextPath,
			}, nil
		},
	}
}

func mnqNatsGetAccountBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("nats-account-id").AutoCompleteFunc = autocompleteNatsAccountID

	return c
}

func mnqNatsListCredentialsBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("nats-account-id").AutoCompleteFunc = autocompleteNatsAccountID

	return c
}

var completeListNatsAccountIDCache *mnq.ListNatsAccountsResponse

func autocompleteNatsAccountID(
	ctx context.Context,
	prefix string,
	request any,
) core.AutocompleteSuggestions {
	region := scw.Region("")
	switch req := request.(type) {
	case *mnq.NatsAPIGetNatsAccountRequest:
		region = req.Region
	case *mnq.NatsAPIListNatsCredentialsRequest:
		region = req.Region
	}

	suggestions := core.AutocompleteSuggestions(nil)

	client := core.ExtractClient(ctx)
	api := mnq.NewNatsAPI(client)

	if completeListNatsAccountIDCache == nil {
		res, err := api.ListNatsAccounts(&mnq.NatsAPIListNatsAccountsRequest{
			Region: region,
		})
		if err != nil {
			return nil
		}
		completeListNatsAccountIDCache = res
	}

	for _, natsAccountID := range completeListNatsAccountIDCache.NatsAccounts {
		if strings.HasPrefix(natsAccountID.ID, prefix) {
			suggestions = append(suggestions, natsAccountID.ID)
		}
	}

	return suggestions
}
