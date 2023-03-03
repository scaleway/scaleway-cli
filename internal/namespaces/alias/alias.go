package alias

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/shlex"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
)

func GetCommands() *core.Commands {
	human.RegisterMarshalerFunc(map[string][]string(nil), func(i interface{}, opt *human.MarshalOpt) (string, error) {
		aliasMap := i.(map[string][]string)
		type humanAlias struct {
			Alias   string
			Command string
		}
		aliases := make([]humanAlias, 0, len(aliasMap))

		for key, value := range aliasMap {
			aliases = append(aliases, humanAlias{
				Alias:   key,
				Command: strings.Join(value, " "),
			})
		}
		return human.Marshal(aliases, opt)
	})

	return core.NewCommands(
		aliasRootCommand(),
		aliasCreateCommand(),
		aliasListCommand(),
		aliasDeleteCommand(),
	)
}

func aliasRootCommand() *core.Command {
	return &core.Command{
		Short: "Alias related commands",
		Long: `This namespace allows you to manage your aliases
You can use multiple aliases in one command
aliases in your commands are evaluated and you get completion
  with: isl = instance server list
    "scw isl <TAB>" will complete as "scw instance server list <TAB>"
    "scw <TAB>" will complete "isl"
`,
		Examples: []*core.Example{
			{
				Short: "Create a custom alias 'isl' for 'instance server list'",
				Raw:   `scw alias create isl command="instance server list"`,
			},
			{
				Short: "Create an alias for a verb",
				Raw:   `scw alias create c command=create`,
			},
		},
		Namespace: "alias",
	}
}

type CreateRequest struct {
	Alias   string `json:"alias"`
	Command string `json:"command"`
}

func aliasCreateCommand() *core.Command {
	return &core.Command{
		Short:     "Create a new alias for a command",
		Namespace: "alias",
		Resource:  "create",
		Long:      `This command help you create aliases and save it to your config`,
		Examples: []*core.Example{
			{
				Short: "Create a custom alias 'isl' for 'instance server list'",
				Raw:   `scw alias create isl command="instance server list""`,
			},
			{
				Short: "Add an alias to a verb",
				Raw:   `scw alias create c command=create`,
			},
		},
		AllowAnonymousClient: true,
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "alias",
				Required:   true,
				Positional: true,
				Short:      "Alias name",
			},
			{
				Name:       "command",
				OneOfGroup: "command",
				Short:      "Command to create an alias for",
			},
		},
		ArgsType: reflect.TypeOf(CreateRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*CreateRequest)
			cfg := core.ExtractCliConfig(ctx)

			response := struct {
				Alias string `json:"alias"`
			}{}
			command, err := shlex.Split(args.Command)
			if err != nil {
				return nil, fmt.Errorf("failed to parse command: %w", err)
			}
			replaced := cfg.Alias.AddAlias(args.Alias, command)
			if replaced {
				response.Alias = "replaced"
			} else {
				response.Alias = "created"
			}

			err = cfg.Save()
			if err != nil {
				return nil, fmt.Errorf("failed to save aliases: %w", err)
			}

			return response, nil
		},
	}
}

type ListRequest struct {
	Alias string `json:"alias"`
}

func aliasListCommand() *core.Command {
	return &core.Command{
		Short:                "List aliases and their commands",
		Namespace:            "alias",
		Resource:             "list",
		AllowAnonymousClient: true,
		ArgSpecs:             core.ArgSpecs{},
		ArgsType:             reflect.TypeOf(ListRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			aliasCfg := core.ExtractAliases(ctx)

			return aliasCfg.Aliases, nil
		},
	}
}

type DeleteRequest struct {
	Alias    string `json:"alias"`
	Resource string `json:"resource"`
}

func aliasDeleteCommand() *core.Command {
	return &core.Command{
		Short:                "Delete an alias",
		Namespace:            "alias",
		Resource:             "delete",
		AllowAnonymousClient: true,
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "alias",
				Positional: true,
				Short:      "alias name",
			},
		},
		ArgsType: reflect.TypeOf(DeleteRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*DeleteRequest)
			cfg := core.ExtractCliConfig(ctx)

			deleted := cfg.Alias.DeleteAlias(args.Alias)

			err := cfg.Save()
			if err != nil {
				return nil, fmt.Errorf("failed to save aliases: %w", err)
			}

			response := struct {
				Alias string `json:"alias"`
			}{}
			if deleted {
				response.Alias = "Deleted"
			} else {
				response.Alias = "Not found"
			}
			return response, err
		},
	}
}
