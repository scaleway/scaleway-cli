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

	// DisableTelemetry, if set to true this will disable telemetry report no matter what the config send_telemetry is set to.
	// This is useful when running test to avoid sending meaningless telemetries.
	DisableTelemetry bool
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
		command:   nil, // command is later injected by cobra_utils.go/cobraRun()
	}

	// Send Matomo telemetry when exiting the bootstrap
	start := time.Now()
	defer func() {
		// skip telemetry report when at least one of the following criteria matches:
		// - telemetry is explicitly disable in bootstrap config
		// - no command was executed
		// - telemetry is disabled on the ran command
		// - telemetry is disabled from the config (user must consent)
		if config.DisableTelemetry ||
			meta.command == nil ||
			meta.command.DisableTelemetry ||
			matomo.IsTelemetryDisabled() {
			logger.Debugf("skipping telemetry report")
			return
		}
		matomoErr := matomo.SendCommandTelemetry(&matomo.SendCommandTelemetryRequest{
			Command:       meta.command.getPath(),
			Version:       config.BuildInfo.Version.String(),
			ExecutionTime: time.Since(start),
		})
		if matomoErr != nil {
			logger.Debugf("error during telemetry reporting: %s", matomoErr)
		} else {
			logger.Debugf("telemetry successfully sent")
		}
	}()

	// Check CLI new version when exiting the bootstrap
	defer func() { // if we plan to remove defer, do not forget logger is not set until cobra pre init func
		config.BuildInfo.checkVersion()
	}()

	// cobraBuilder will build a Cobra root command from a list of Command
	builder := cobraBuilder{
		commands: config.Commands.command,
		meta:     meta,
	}

	rootCmd := builder.build()

	rootCmd.PersistentFlags().StringVarP(&meta.ProfileFlag, "profile", "p", "", "The config profile to use")
	rootCmd.PersistentFlags().VarP(&meta.PrinterTypeFlag, "output", "o", "Output format: json or human")
	rootCmd.PersistentFlags().BoolVarP(&meta.DebugModeFlag, "debug", "D", false, "Enable debug mode")

	rootCmd.SetArgs(config.Args[1:])
	err = rootCmd.Execute()

	if err != nil {
		if _, ok := err.(*interactive.InterruptError); ok {
			return 130, nil, err
		}
		printErr := meta.Printer.Print(err, nil)
		if printErr != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		return 1, nil, err
	}
	return 0, meta.result, nil
}
