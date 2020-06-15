package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-cli/internal/matomo"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/spf13/pflag"
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

	// Stdin stream to use. Usually os.Stdin
	Stdin io.Reader

	// If provided this client will be passed to all commands.
	// If not a client will be automatically created by the CLI using Config, Env and flags see createClient().
	Client *scw.Client

	// DisableTelemetry, if set to true this will disable telemetry report no matter what the config send_telemetry is set to.
	// This is useful when running test to avoid sending meaningless telemetries.
	DisableTelemetry bool

	// OverrideEnv overrides environment variables returned by core.ExtractEnv function.
	// This is useful for tests as it allows overriding env without relying on global state.
	OverrideEnv map[string]string

	// OverrideExec allow to override exec.Cmd.Run method. In order for this to work
	// your code must call le core.ExecCmd function to execute a given command.
	// If this function is not defined the exec.Cmd.Run function will be called directly.
	// This function is intended to be use for tests purposes.
	OverrideExec OverrideExecFunc

	// BaseContest is the base context that will be used across all function call from top to bottom.
	Ctx context.Context

	// Optional we use it if defined
	Logger *Logger
}

// Bootstrap is the main entry point. It is directly called from main.
// BootstrapConfig.Args is usually os.Args
// BootstrapConfig.Commands is a list of command available in CLI.
func Bootstrap(config *BootstrapConfig) (exitCode int, result interface{}, err error) {
	// Handles Flags
	var debug bool
	var profileName string
	var configPath string
	var printerType PrinterType

	flags := pflag.NewFlagSet(config.Args[0], pflag.ContinueOnError)
	flags.StringVarP(&profileName, "profile", "p", "", "The config profile to use")
	flags.StringVarP(&configPath, "config", "c", "", "The path to the config file")
	flags.VarP(&printerType, "output", "o", "Output format: json or human")
	flags.BoolVarP(&debug, "debug", "D", false, "Enable debug mode")
	// Ignore unknown flag
	flags.ParseErrorsWhitelist.UnknownFlags = true
	// Make sure usage is never print by the parse method. (It should only be print by cobra)
	flags.Usage = func() {}

	// We don't do any error validation as:
	// - debug is a boolean, no possible error
	// - profileName will return proper error when we try to load profile
	// - printerType will return proper error when we create the printer
	// Furthermore additional flag can be added on a per-command basis inside cobra
	// parse would fail as these flag are not known at this time.
	_ = flags.Parse(config.Args)

	// If debug flag is set enable debug mode in SDK logger
	logLevel := logger.LogLevelWarning
	if debug {
		logLevel = logger.LogLevelDebug // enable debug mode
	}

	// We force log to os.Stderr because we dont have a scoped logger feature and it create
	// concurrency situation with golden files
	log := config.Logger
	if log == nil {
		log = &Logger{
			writer: os.Stderr,
		}
	}
	log.level = logLevel
	logger.SetLogger(log)

	// The printer must be the first thing set in order to print errors
	printer, err := NewPrinter(&PrinterConfig{
		Type:   printerType,
		Stdout: config.Stdout,
		Stderr: config.Stderr,
	})
	if err != nil {
		_, _ = fmt.Fprintln(config.Stderr, err)
		return 1, nil, err
	}
	interactive.SetOutputWriter(config.Stderr) // set printer for interactive function (always stderr).

	// An authenticated client will be created later if required.
	client := config.Client
	isClientFromBootstrapConfig := true
	if client == nil {
		isClientFromBootstrapConfig = false
		client, err = createAnonymousClient(config.BuildInfo)
		if err != nil {
			printErr := printer.Print(err, nil)
			if printErr != nil {
				_, _ = fmt.Fprintln(config.Stderr, printErr)
			}
			return 1, nil, err
		}
	}

	// Meta store globally available variables like SDK client.
	// Meta is injected in a context object that will be passed to all commands.
	meta := &meta{
		ProfileFlag:  profileName,
		BinaryName:   config.Args[0],
		BuildInfo:    config.BuildInfo,
		Client:       client,
		Commands:     config.Commands,
		OverrideEnv:  config.OverrideEnv,
		OverrideExec: config.OverrideExec,

		stdout:                      config.Stdout,
		stderr:                      config.Stderr,
		stdin:                       config.Stdin,
		result:                      nil, // result is later injected by cobra_utils.go/cobraRun()
		command:                     nil, // command is later injected by cobra_utils.go/cobraRun()
		isClientFromBootstrapConfig: isClientFromBootstrapConfig,
	}
	// We make sure OverrideEnv is never nil in meta.
	if meta.OverrideEnv == nil {
		meta.OverrideEnv = map[string]string{}
	}

	// If OverrideExec was not set in the config, we set a default value.
	if meta.OverrideExec == nil {
		meta.OverrideExec = defaultOverrideExec
	}

	ctx := config.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = injectMeta(ctx, meta)

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
		commands: config.Commands.commands,
		meta:     meta,
		ctx:      ctx,
	}

	rootCmd := builder.build()

	// These flag are already handle at the beginning of this function but we keep this
	// declaration in order for them to be shown in the cobra usage documentation.
	rootCmd.PersistentFlags().StringVarP(&profileName, "profile", "p", "", "The config profile to use")
	rootCmd.PersistentFlags().StringVarP(&profileName, "config", "c", "", "The path to the config file")
	rootCmd.PersistentFlags().VarP(&printerType, "output", "o", "Output format: json or human")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "Enable debug mode")
	rootCmd.SetArgs(config.Args[1:])
	err = rootCmd.Execute()

	if err != nil {
		if _, ok := err.(*interactive.InterruptError); ok {
			return 130, nil, err
		}
		errorCode := 1
		if cliErr, ok := err.(*CliError); ok && cliErr.Code != 0 {
			errorCode = cliErr.Code
		}
		printErr := printer.Print(err, nil)
		if printErr != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		return errorCode, nil, err
	}

	if meta.command != nil {
		printErr := printer.Print(meta.result, meta.command.getHumanMarshalerOpt())
		if printErr != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	}

	return 0, meta.result, nil
}
