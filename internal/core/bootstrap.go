package core

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-cli/internal/matomo"
	"github.com/scaleway/scaleway-cli/internal/printer"
	"github.com/scaleway/scaleway-sdk-go/logger"
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
// BootstrapConfig.Args is usually os.Args
// BootstrapConfig.Commands is a list of command available in CLI.
func Bootstrap(config *BootstrapConfig) (exitCode int, result interface{}, err error) {
	// The global printer must be the first thing set in order to print errors
	globalPrinter, err := printer.New(printer.Human, config.Stdout, config.Stderr)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1, nil, err
	}

	// Meta store globally available variables like SDK client.
	// Meta is injected in a context object that will be passed to all commands.
	meta := &meta{
		BuildInfo: config.BuildInfo,
		stdout:    config.Stdout,
		stderr:    config.Stderr,
		Client:    config.Client,
		Commands:  config.Commands,
		Printer:   globalPrinter,
		result:    nil, // result is later injected by cobra_utils.go/cobraRun()
	}

	// Send Matomo report when exiting the bootstrap
	if (matomo.ForceTracking || config.BuildInfo.isRelease()) && matomo.IsTelemetryEnabled() {
		start := time.Now()
		defer func() {
			matomoErr := matomo.TrackCommand(&matomo.TrackCommandRequest{
				RunCommand:    meta.runCommand.getPath(),
				Version:       config.BuildInfo.Version,
				ExecutionTime: time.Since(start),
			})
			if matomoErr != nil {
				logger.Warningf("Error during telemetry reporting: %s", matomoErr)
			}
		}()
	}

	// cobraBuilder will build a Cobra root command from a list of Command
	builder := cobraBuilder{
		commands: config.Commands.command,
		meta:     meta,
	}

	rootCmd := builder.build()

	rootCmd.PersistentFlags().StringVarP(&meta.AccessKeyFlag, "access-key", "", "", "Scaleway access key")
	rootCmd.PersistentFlags().StringVarP(&meta.SecretKeyFlag, "secret-key", "", "", "Scaleway secret key")
	rootCmd.PersistentFlags().StringVarP(&meta.ProfileFlag, "profile", "p", "", "The config profile to use")
	rootCmd.PersistentFlags().VarP(&meta.PrinterTypeFlag, "output", "o", "Output format: json or human")
	rootCmd.PersistentFlags().BoolVarP(&meta.DebugModeFlag, "debug", "D", false, "Enable debug mode")

	rootCmd.SetArgs(config.Args[1:])
	err = rootCmd.Execute()

	if err != nil {
		if _, ok := err.(*interactive.InterruptError); ok {
			return 130, meta.result, err
		}
		err = meta.Printer.Print(err, nil)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		return 1, meta.result, err
	}
	return 0, meta.result, nil
}
