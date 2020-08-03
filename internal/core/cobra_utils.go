package core

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/spf13/cobra"
)

// cobraRun returns a cobraRun command that wrap a CommandRunner function.
func cobraRun(ctx context.Context, cmd *Command) func(*cobra.Command, []string) error {
	return func(cobraCmd *cobra.Command, rawArgsStr []string) error {
		rawArgs := args.RawArgs(rawArgsStr)

		meta := extractMeta(ctx)
		meta.command = cmd

		// If command requires authentication and the client was not directly provided in the bootstrap config, we create a new client and overwrite the existing one
		if !cmd.AllowAnonymousClient && !meta.isClientFromBootstrapConfig {
			client, err := createClient(meta.httpClient, meta.BuildInfo, ExtractProfileName(ctx))
			if err != nil {
				return err
			}
			meta.Client = client
			err = validateClient(meta.Client)
			if err != nil {
				return err
			}
		}

		// If command has no Run method there is nothing to do.
		if cmd.Run == nil {
			return nil
		}

		// Apply default values on missing args.
		rawArgs = ApplyDefaultValues(ctx, cmd.ArgSpecs, rawArgs)

		positionalArgSpec := cmd.ArgSpecs.GetPositionalArg()

		// If this command has no positional argument we execute the run
		if positionalArgSpec == nil {
			data, err := run(ctx, cobraCmd, cmd, rawArgs)
			if err != nil {
				return err
			}

			meta.result = data
			return nil
		}

		positionalArgs := rawArgs.GetPositionalArgs()

		// Return an error if a positional argument was provide using `key=value` syntax.
		value, exist := rawArgs.Get(positionalArgSpec.Name)
		if exist {
			otherArgs := rawArgs.Remove(positionalArgSpec.Name)
			return &CliError{
				Err:  fmt.Errorf("a positional argument is required for this command"),
				Hint: positionalArgHint(meta.BinaryName, cmd, value, otherArgs, len(positionalArgs) > 0),
			}
		}

		// If no positional arguments were provided, return an error
		if len(positionalArgs) == 0 {
			return &CliError{
				Err:  fmt.Errorf("a positional argument is required for this command"),
				Hint: positionalArgHint(meta.BinaryName, cmd, "<"+positionalArgSpec.Name+">", rawArgs, false),
			}
		}

		results := MultiResults(nil)
		rawArgs = rawArgs.RemoveAllPositional()
		for _, positionalArg := range positionalArgs {
			rawArgsWithPositional := rawArgs.Add(positionalArgSpec.Name, positionalArg)

			result, err := run(ctx, cobraCmd, cmd, rawArgsWithPositional)
			if err != nil {
				return err
			}

			results = append(results, result)
		}
		// If only one positional parameter was provided we return the result directly instead of
		// an array of results
		if len(results) == 1 {
			meta.result = results[0]
		} else {
			meta.result = results
		}

		return nil
	}
}

func run(ctx context.Context, cobraCmd *cobra.Command, cmd *Command, rawArgs []string) (interface{}, error) {
	var err error

	// create a new Args interface{}
	// unmarshalled arguments will be store in this interface
	cmdArgs := reflect.New(cmd.ArgsType).Interface()

	// Unmarshal args.
	// After that we are done working with rawArgs
	// and will be working with cmdArgs.
	err = args.UnmarshalStruct(rawArgs, cmdArgs)
	if err != nil {
		if unmarshalError, ok := err.(*args.UnmarshalArgError); ok {
			return nil, handleUnmarshalErrors(cmd, unmarshalError)
		}
		return nil, err
	}

	// PreValidate hook.
	if cmd.PreValidateFunc != nil {
		err = cmd.PreValidateFunc(ctx, cmdArgs)
		if err != nil {
			return nil, err
		}
	}

	// Validate
	validateFunc := DefaultCommandValidateFunc()
	if cmd.ValidateFunc != nil {
		validateFunc = cmd.ValidateFunc
	}
	err = validateFunc(cmd, cmdArgs, rawArgs)
	if err != nil {
		return nil, err
	}

	// execute the command
	interceptor := combineCommandInterceptor(
		sdkStdErrorInterceptor,
		sdkStdTypeInterceptor,
		cmd.Interceptor,
	)

	data, err := interceptor(ctx, cmdArgs, func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
		return cmd.Run(ctx, argsI)
	})
	if err != nil {
		return nil, err
	}
	waitFlag, err := cobraCmd.PersistentFlags().GetBool("wait")
	if err == nil && cmd.WaitFunc != nil && waitFlag {
		data, err = cmd.WaitFunc(ctx, cmdArgs, data)
		if err != nil {
			return nil, err
		}
	}
	// update cache
	if meta := extractMeta(ctx); !meta.disableCache {
		if cmd.RefreshCacheFunc != nil {
			cmd.RefreshCacheFunc(ctx, data)
		}
	}
	return data, nil
}

// positionalArgHint formats the positional argument error hint.
func positionalArgHint(binaryName string, cmd *Command, hintValue string, otherArgs []string, positionalArgumentFound bool) string {
	suggestedArgs := []string{}

	// If no positional argument exists, suggest one.
	if !positionalArgumentFound {
		suggestedArgs = append(suggestedArgs, hintValue)
	}

	// Suggest to use the other arguments.
	suggestedArgs = append(suggestedArgs, otherArgs...)

	suggestedCommand := append([]string{cmd.GetCommandLine(binaryName)}, suggestedArgs...)
	return "Try running: " + strings.Join(suggestedCommand, " ")
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
