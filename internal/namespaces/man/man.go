// Package man implements the `scw man` namespace: man page generation and
// installation for all built-in scw commands.
package man

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/manpage"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		manRootCommand(),
		installCommand(),
	)
}

func manRootCommand() *core.Command {
	return &core.Command{
		Groups: []string{"utility"},
		Short:  "Man pages related commands",
		Long: interactive.RemoveIndent(`
			This namespace allows you to generate and install POSIX man pages for all scw commands.
		`),
		Namespace: "man",
	}
}

func installLong() string {
	return interactive.RemoveIndent(`
		Install the man pages of all scw commands in a given directory.

		The target directory is created if it does not exist. Pages are installed in its ` + "`man1`" + `
		subdirectory, so the target should be a man root (a directory that contains a ` + "`man1`" + `
		directory), for example ` + "`/usr/local/share/man`" + ` or ` + "`~/my-man`" + `.

		On reinstall, only the pages owned by this CLI (scw.1 and scw-*.1) are managed: pages that no
		longer correspond to a command are removed, while any other file in the target is left untouched.
	`)
}

func installCommand() *core.Command {
	type installArgs struct {
		// Target is the directory where the man pages will be installed.
		Target string
	}

	return &core.Command{
		Groups:               []string{"utility"},
		Short:                "Install the man pages of all scw commands in a given directory",
		Long:                 installLong(),
		Namespace:            "man",
		Resource:             "install",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeFor[installArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "target",
				Short:      "Directory where the man pages will be installed. It is created if it does not exist. Pages are installed in its `man1` subdirectory (e.g. `/usr/local/share/man` or `~/my-man`).",
				Required:   true,
				Positional: true,
			},
		},
		Examples: []*core.Example{
			{
				Short: "Install the man pages in a standard system man directory (may require elevated privileges)",
				Raw:   "scw man install /usr/local/share/man",
			},
			{
				Short: "Install the man pages in a user directory and add it to the man search path",
				Raw:   "scw man install ~/my-man\nexport MANPATH=$MANPATH:~/my-man",
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*installArgs)

			commands := core.ExtractCommands(ctx)
			if commands == nil {
				return nil, errors.New("command list not found in context")
			}

			count, err := manpage.Install(ctx, commands, args.Target)
			if err != nil {
				return nil, err
			}

			return core.RawResult(
				[]byte(
					fmt.Sprintf(
						"Installed %d man pages to %s\n",
						count,
						filepath.Join(args.Target, manpage.SectionDir),
					),
				),
			), nil
		},
	}
}
