package docs

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
)

type DocsArgs struct {
	Search   string
	PathArgs []string
}

func GetCommands() *core.Commands {
	return core.NewCommands(
		docsRoot(),
	)
}

func docsRoot() *core.Command {
	return &core.Command{
		Short:                "Browse the CLI documentation interactively",
		Long:                 "Open an interactive documentation browser that lets you navigate all CLI commands organized by namespace, resource, and verb. Use arrow keys to browse, or pass a namespace/command path as argument to jump directly to a specific page.",
		Namespace:            "docs",
		AllowAnonymousClient: true,
		Groups:               []string{"utility"},
		ArgsType:             reflect.TypeFor[args.RawArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:  "search",
				Short: "Search for commands matching a keyword",
			},
		},
		ValidateFunc: func(ctx context.Context, cmd *core.Command, cmdArgs any, rawArgs args.RawArgs) error {
			return nil
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			rawArgs := *argsI.(*args.RawArgs)
			search, _ := rawArgs.Get("search")
			pathArgs := rawArgs.GetPositionalArgs()

			commands := core.ExtractCommands(ctx)
			if commands == nil {
				return nil, &core.CliError{
					Err: errors.New("failed to extract commands from context"),
				}
			}

			if search != "" {
				results := SearchCommands(commands, search)
				if len(results) == 0 {
					return "No commands found matching \"" + search + "\".", nil
				}

				return formatSearchResults(results), nil
			}

			if len(pathArgs) > 0 {
				return renderDirectPath(ctx, commands, pathArgs)
			}

			if interactive.IsInteractive {
				return runBrowser(ctx, commands)
			}

			return RenderAllNamespaces(commands), nil
		},
	}
}

func renderDirectPath(
	ctx context.Context,
	commands *core.Commands,
	pathArgs []string,
) (any, error) {
	path := strings.Join(pathArgs, ".")
	cmd := commands.Find(pathArgs...)
	if cmd == nil {
		if len(pathArgs) == 1 {
			return nil, &core.CliError{
				Err: fmt.Errorf(
					"Namespace %q not found. Run \"scw docs\" to see available namespaces.",
					pathArgs[0],
				),
			}
		}

		return nil, &core.CliError{
			Err: fmt.Errorf(
				"Command %q not found. Run \"scw docs\" to see available commands.",
				path,
			),
		}
	}

	if cmd.Namespace != "" && cmd.Resource == "" && cmd.Verb == "" {
		return RenderNamespacePage(commands, cmd.Namespace), nil
	}

	if cmd.Namespace != "" && cmd.Resource != "" && cmd.Verb == "" {
		return RenderResourcePage(commands, cmd.Namespace, cmd.Resource), nil
	}

	return RenderCommandPage(ctx, cmd, commands), nil
}

func formatSearchResults(results []SearchResult) string {
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(r.Path + " - " + r.Short + "\n")
	}

	return strings.TrimSuffix(sb.String(), "\n")
}
