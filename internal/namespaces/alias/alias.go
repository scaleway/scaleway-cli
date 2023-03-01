package alias

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/alias"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
)

func GetCommands() *core.Commands {
	human.RegisterMarshalerFunc(alias.Config{}, func(i interface{}, opt *human.MarshalOpt) (string, error) {
		cfg := i.(alias.Config)
		// To avoid recursion of human.Marshal we create a dummy type
		type tmp alias.Config
		cfgTmp := tmp(cfg)

		// Sections
		opt.Sections = []*human.MarshalSection{
			{
				FieldName: "Aliases",
			},
			{
				FieldName: "ResourceAliases",
			},
		}

		str, err := human.Marshal(cfgTmp, opt)
		if err != nil {
			return "", err
		}

		return str, nil
	})

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
Two alias types exist
Raw aliases:
	A command can contains multiple aliases
	aliases in your commands are expanded and you get completion
	with: isl = instance server list
	"scw isl <TAB>" will complete as "scw instance server list <TAB>"
	but "scw <TAB>" will not complete "isl"

Resource aliases:
	These are alternate names for namespaces, resources or verbs
	resource aliases are provided by the completion
	with: instance server create = c
	"scw instance server <TAB>" will get "c" in completion
`,
		Examples: []*core.Example{
			{
				Short: "Create a custom alias 'isl' for 'instance server list'",
				Raw:   "scw alias create isl command.0=instance command.1=server command.2=list",
			},
			{
				Short: "Add an alias to a verb",
				Raw:   `scw alias create c resource=instance.server.create`,
			},
		},
		Namespace: "alias",
	}
}

type CreateRequest struct {
	Alias    string   `json:"alias"`
	Command  []string `json:"command"`
	Resource string   `json:"resource"`
}

func aliasCreateCommand() *core.Command {
	return &core.Command{
		Short:     "Create a new alias for a command",
		Namespace: "alias",
		Resource:  "create",
		Long: `This command help you create aliases and save it to your config
use command argument to create a raw alias
use resource argument to add a resource alias
`,
		Examples: []*core.Example{
			{
				Short: "Create a custom alias 'isl' for 'instance server list'",
				Raw:   "scw alias create isl command.0=instance command.1=server command.2=list",
			},
			{
				Short: "Add an alias to a verb",
				Raw:   `scw alias create c resource=instance.server.create`,
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
				Name:       "command.{index}",
				OneOfGroup: "command",
				Short:      "Command to expand to for a raw alias",
			},
			{
				Name:       "resource",
				OneOfGroup: "command",
				Short:      "resource to add alias to",
			},
		},
		ArgsType: reflect.TypeOf(CreateRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*CreateRequest)
			cfg := core.ExtractCliConfig(ctx)

			var replaced bool
			response := struct {
				Alias string `json:"alias"`
			}{}

			if args.Resource != "" {
				// Resource alias
				resourcePath := alias.SplitResourcePath(args.Resource)

				cmd := core.ExtractCommands(ctx).Find(resourcePath...)
				if cmd == nil {
					return nil, fmt.Errorf("resource not found: %s", args.Resource)
				}

				cfg.Alias.AddResourceAlias(resourcePath, args.Alias)
			} else {
				// Raw alias
				replaced = cfg.Alias.AddAlias(args.Alias, args.Command)
			}

			if replaced {
				response.Alias = "replaced"
			} else {
				response.Alias = "created"
			}
			err := cfg.Save()
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

			return aliasCfg, nil
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
			{
				Name:  "resource",
				Short: "resource path is for resource aliases",
			},
		},
		ArgsType: reflect.TypeOf(DeleteRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*DeleteRequest)
			cfg := core.ExtractCliConfig(ctx)

			var deleted bool

			if args.Resource != "" {
				// Resources Alias
				resource := alias.SplitResourcePath(args.Resource)

				cmd := core.ExtractCommands(ctx).Find(resource...)
				if cmd == nil {
					return nil, fmt.Errorf("resource not found: %s", args.Resource)
				}

				deleted = cfg.Alias.DeleteResourceAlias(resource, args.Alias)
			} else {
				deleted = cfg.Alias.DeleteAlias(args.Alias)
			}

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
