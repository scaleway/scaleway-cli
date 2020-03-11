package core

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/printer"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/spf13/cobra"
)

// cobraRun returns a cobraRun command that wrap a CommandRunner function.
func cobraRun(ctx context.Context, cmd *Command) func(*cobra.Command, []string) error {
	return func(cobraCmd *cobra.Command, rawArgs []string) error {
		var err error
		opt := cmd.getHumanMarshalerOpt()
		meta := extractMeta(ctx)
		meta.command = cmd

		// create a new Args interface{}
		// unmarshalled arguments will be store in this interface
		cmdArgs := reflect.New(cmd.ArgsType).Interface()

		// Apply default values on missing args.
		rawArgs = ApplyDefaultValues(cmd.ArgSpecs, rawArgs)

		// Check args exist valid
		argsSlice := args.SplitRawNoError(rawArgs)
		for _, arguments := range argsSlice {
			// TODO: handle args such as tags.index
			if cmd.ArgSpecs.GetByName(arguments[0]) == nil &&
				cmd.ArgSpecs.GetByName(arguments[0]+".{index}") == nil &&
				!strings.Contains(arguments[0], ".") {
				return handleUnmarshalErrors(cmd, &args.UnmarshalArgError{
					Err:     &args.UnknownArgError{},
					ArgName: arguments[0],
				})
			}
		}

		// Unmarshal args.
		// After that we are done working with rawArgs
		// and will be working with cmdArgs.
		err = args.UnmarshalStruct(rawArgs, cmdArgs)
		if err != nil {
			if unmarshalError, ok := err.(*args.UnmarshalArgError); ok {
				return handleUnmarshalErrors(cmd, unmarshalError)
			}
			return err
		}

		// PreValidate hook.
		if cmd.PreValidateFunc != nil {
			err = cmd.PreValidateFunc(ctx, cmdArgs)
			if err != nil {
				return err
			}
		}

		// Validate
		validateFunc := DefaultCommandValidateFunc()
		if cmd.ValidateFunc != nil {
			validateFunc = cmd.ValidateFunc
		}
		err = validateFunc(cmd, cmdArgs)
		if err != nil {
			return err
		}

		// execute the command
		if cmd.Run != nil {
			data, err := cmd.Run(ctx, cmdArgs)
			if err != nil {
				return err
			}
			waitFlag, err := cobraCmd.PersistentFlags().GetBool("wait")
			if err == nil && cmd.WaitFunc != nil && waitFlag {
				data, err = cmd.WaitFunc(ctx, cmdArgs, data)
				if err != nil {
					return err
				}
			}
			meta.result = data
			return meta.Printer.Print(data, opt)
		}

		return nil
	}
}

func handleUnmarshalErrors(cmd *Command, unmarshalErr *args.UnmarshalArgError) error {
	wrappedErr := errors.Unwrap(unmarshalErr)

	switch e := wrappedErr.(type) {
	case *args.CannotUnmarshalError:
		hint := ""
		if _, ok := e.Dest.(*bool); ok {
			hint = "Possible values: true, false"
		}

		return &CliError{
			Err:  fmt.Errorf("invalid value for '%s' argument: %s", unmarshalErr.ArgName, e.Err),
			Hint: hint,
		}

	case *args.UnknownArgError, *args.InvalidArgNameError:
		argNames := []string(nil)
		for _, argSpec := range cmd.ArgSpecs {
			argNames = append(argNames, argSpec.Name)
		}

		return &CliError{
			Err:  fmt.Errorf("unknown argument '%s'", unmarshalErr.ArgName),
			Hint: fmt.Sprintf("Valid arguments are: %s", strings.Join(argNames, ", ")),
		}

	default:
		return &CliError{Err: unmarshalErr}
	}
}

// cobraPreRunInitMeta returns a cobraPreRun command that will initialize meta.
func cobraPreRunInitMeta(ctx context.Context, cmd *Command) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		var err error
		meta := extractMeta(ctx)

		logLevel := logger.LogLevelWarning
		if meta.DebugModeFlag {
			logLevel = logger.LogLevelDebug // enable debug mode
		}
		logger.DefaultLogger.Init(meta.stderr, logLevel)

		meta.Printer, err = printer.New(meta.PrinterTypeFlag, meta.stdout, meta.stderr)
		if err != nil {
			return &CliError{
				Err:     fmt.Errorf("invalid output flag value '%s'", meta.PrinterTypeFlag),
				Details: "Supported output format are: human or json",
			}
		}

		// If command require a client and no client was provided in BootstrapConfig
		if !cmd.NoClient && meta.Client == nil {
			meta.Client, err = createClient(meta)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
