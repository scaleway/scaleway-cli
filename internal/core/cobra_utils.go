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

		// Handle positional argument by catching first argument `<value>` and rewrite it to `<arg-name>=<value>`.
		if err = handlePositionalArg(cmd, rawArgs); err != nil {
			return err
		}

		// Apply default values on missing args.
		rawArgs = ApplyDefaultValues(cmd.ArgSpecs, rawArgs)

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

// handlePositionalArg will catch positional argument if command has one.
// When a positional argument is found it will mutate its value in rawArgs to match the argument unmarshaller format.
// E.g.: '[value b=true c=1]' will be mutated to '[a=value b=true c=1]'.
// It returns errors when:
// - no positional argument is found.
// - a positional argument exists, but is not positional.
// - an argument duplicates a positional argument.
func handlePositionalArg(cmd *Command, rawArgs []string) error {
	positionalArg := cmd.ArgSpecs.GetPositionalArg()

	// Command does not have a positional argument.
	if positionalArg == nil {
		return nil
	}

	// Positional argument is found condition.
	positionalArgumentFound := len(rawArgs) > 0 && !strings.Contains(rawArgs[0], "=")

	// Argument exists but is not positional.
	for i, arg := range rawArgs {
		if strings.HasPrefix(arg, positionalArg.Prefix()) {
			return &CliError{
				Err:  fmt.Errorf("a positional argument is required for this command"),
				Hint: positionalArgHint(cmd, strings.TrimLeft(arg, positionalArg.Prefix()), append(rawArgs[:i], rawArgs[i+1:]...), positionalArgumentFound),
			}
		}
	}

	// If positional argument is found, prefix it with `arg-name=`.
	if positionalArgumentFound {
		rawArgs[0] = positionalArg.Prefix() + rawArgs[0]
		return nil
	}

	// No positional argument found.
	return &CliError{
		Err:  fmt.Errorf("a positional argument is required for this command"),
		Hint: positionalArgHint(cmd, "<"+positionalArg.Name+">", rawArgs, false),
	}
}

// positionalArgHint helps formatting the positional argument error.
func positionalArgHint(cmd *Command, hintValue string, otherArgs []string, positionalArgumentFound bool) string {
	suggestedArgs := []string{}

	// If no positional argument exists, suggest one.
	if !positionalArgumentFound {
		suggestedArgs = append(suggestedArgs, hintValue)
	}

	// Suggest to use the other arguments.
	suggestedArgs = append(suggestedArgs, otherArgs...)

	suggestedCommand := append([]string{"scw", cmd.getPath()}, suggestedArgs...)
	return "Try running '" + strings.Join(suggestedCommand, " ") + "'."
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
