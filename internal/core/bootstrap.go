package core

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/scaleway/scaleway-cli/internal/printer"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type BootstrapConfig struct {
	// Args to use for the command. Usually os.Args
	Args []string

	// A list of all available commands
	Commands *Commands

	// BuildInfo contains information about cli build
	BuildInfo *BuildInfo

	// Stdout stream to use. Usually os.Stdout
	Stdout io.Writer

	// Stderr stream to use. Usually os.Stderr
	Stderr io.Writer

	// If provided this client will be passed to all commands.
	// If not a client will be automatically created by the CLI using Config, Env and flags see createClient().
	Client *scw.Client
}

// Bootstrap is the main entry point. It is directly called from main.
// args is usually os.Args and commands is a list of command available in CLI.
func Bootstrap(config *BootstrapConfig) (exitCode int) {

	printer_, err := printer.New(printer.Human, config.Stdout, config.Stderr)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// Meta store globally available variables like SDK client.
	// Meta is injected in a context object that will be pass to all commands.
	m := &meta{
		Printer:   printer_,
		BuildInfo: config.BuildInfo,
		Client:    config.Client,
		stdout:    config.Stdout,
		stderr:    config.Stderr,
	}
	ctx := injectMeta(context.Background(), m)
	ctx = injectCommands(ctx, config.Commands)

	// cobraBuilder will build a Cobra root command from a list of Command
	builder := cobraBuilder{
		ctx:      ctx,
		commands: config.Commands.command,
		meta:     m,
	}

	rootCmd := builder.build()

	rootCmd.PersistentFlags().StringVarP(&m.AccessKeyFlag, "access-key", "", "", "Scaleway access key")
	rootCmd.PersistentFlags().StringVarP(&m.SecretKeyFlag, "secret-key", "", "", "Scaleway secret key")
	rootCmd.PersistentFlags().StringVarP(&m.ProfileFlag, "profile", "p", "", "The config profile to use")
	rootCmd.PersistentFlags().VarP(&m.PrinterTypeFlag, "output", "o", "Output format: json or human")
	rootCmd.PersistentFlags().BoolVarP(&m.DebugModeFlag, "debug", "D", false, "Enable debug mode")

	rootCmd.SetArgs(config.Args[1:])
	err = rootCmd.Execute()
	if err != nil {
		err = m.Printer.Print(err, nil)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		return 1
	}
	return 0
}
