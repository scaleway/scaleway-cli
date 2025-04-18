package alias

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/pkg/shlex"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		aliasRootCommand(),
		aliasCreateCommand(),
		aliasListCommand(),
		aliasDeleteCommand(),
	)
}

func aliasRootCommand() *core.Command {
	return &core.Command{
		Groups: []string{"config"},
		Short:  "Alias related commands",
		Long: `This namespace allows you to manage your aliases
Aliases are stored in cli config file, Default path for this configuration file is based on the following priority order:

- $SCW_CLI_CONFIG_PATH
- $XDG_CONFIG_HOME/scw/config.yaml
- $HOME/.config/scw/config.yaml
- $USERPROFILE/.config/scw/config.yaml

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
	Alias   string `json:"alias"`
	OrderBy string `json:"order-by"`
	Command string `json:"command"`
}

type aliasListItem struct {
	Alias   string
	Command string
}

type aliasListResponse []aliasListItem

type aliasListOrderFunction func(aliases aliasListResponse) func(i, j int) bool

var aliasListOrderFunctions = map[string]aliasListOrderFunction{
	"command_asc": func(aliases aliasListResponse) func(i int, j int) bool {
		return func(i, j int) bool {
			return aliases[i].Command < aliases[j].Command
		}
	},
	"command_desc": func(aliases aliasListResponse) func(i int, j int) bool {
		return func(i, j int) bool {
			return aliases[i].Command > aliases[j].Command
		}
	},
	"alias_asc": func(aliases aliasListResponse) func(i int, j int) bool {
		return func(i, j int) bool {
			return aliases[i].Alias < aliases[j].Alias
		}
	},
	"alias_desc": func(aliases aliasListResponse) func(i int, j int) bool {
		return func(i, j int) bool {
			return aliases[i].Alias > aliases[j].Alias
		}
	},
}

func aliasListCommand() *core.Command {
	return &core.Command{
		Short:                "List aliases and their commands",
		Namespace:            "alias",
		Resource:             "list",
		AllowAnonymousClient: true,
		ArgSpecs: core.ArgSpecs{
			{
				Name:    "order-by",
				Default: core.DefaultValueSetter("command_asc"),
				EnumValues: []string{
					"command_asc",
					"command_desc",
					"alias_asc",
					"alias_desc",
				},
			},
			{
				Name:  "command",
				Short: "filter command",
			},
			{
				Name:  "alias",
				Short: "filter alias",
			},
		},
		ArgsType: reflect.TypeOf(ListRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*ListRequest)
			aliasCfg := core.ExtractAliases(ctx)
			aliases := make(aliasListResponse, 0, len(aliasCfg.Aliases))

			for key, value := range aliasCfg.Aliases {
				aliases = append(aliases, aliasListItem{
					Alias:   key,
					Command: strings.Join(value, " "),
				})
			}

			orderFunction, exists := aliasListOrderFunctions[args.OrderBy]
			if !exists {
				return nil, fmt.Errorf("order-by %s is not supported", args.OrderBy)
			}

			sort.Slice(aliases, orderFunction(aliases))

			filteredAliases := aliasListResponse{}

			for _, aliasItem := range aliases {
				if strings.Contains(aliasItem.Command, args.Command) &&
					strings.Contains(aliasItem.Alias, args.Alias) {
					filteredAliases = append(filteredAliases, aliasItem)
				}
			}

			return filteredAliases, nil
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
