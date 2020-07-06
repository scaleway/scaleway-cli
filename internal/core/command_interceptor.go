package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

// combineCommandInterceptor will combine one or more CommandInterceptor.
// Resulting CommandInterceptor can be viewed as a chain composed by all interceptors passed as parameter.
func combineCommandInterceptor(interceptors ...CommandInterceptor) CommandInterceptor {
	var combinedInterceptors CommandInterceptor
	for _, interceptor := range interceptors {
		localInterceptor := interceptor
		if interceptor == nil {
			continue
		}
		if combinedInterceptors == nil {
			combinedInterceptors = interceptor
			continue
		}

		previousInterceptor := combinedInterceptors
		combinedInterceptors = func(ctx context.Context, args interface{}, runner CommandRunner) (interface{}, error) {
			return previousInterceptor(ctx, args, func(ctx context.Context, arg interface{}) (interface{}, error) {
				return localInterceptor(ctx, args, runner)
			})
		}
	}
	return combinedInterceptors
}

// sdkStdErrorInterceptor is a command interceptor that will catch sdk standard error and return more friendly CLI error.
func sdkStdErrorInterceptor(ctx context.Context, args interface{}, runner CommandRunner) (interface{}, error) {
	res, err := runner(ctx, args)
	switch sdkError := err.(type) {
	case *scw.ResourceNotFoundError:
		return nil, &CliError{
			Message: fmt.Sprintf("cannot find resource '%v' with ID '%v'", sdkError.Resource, sdkError.ResourceID),
			Err:     err,
		}
	case *scw.ResponseError:
		return nil, &CliError{
			Message: sdkError.Message,
			Err:     sdkError,
		}
	case *scw.InvalidArgumentsError:
		reasonsMap := map[string]string{
			"unknown":    "is invalid for unexpected reasons",
			"required":   "is required",
			"format":     "is wrongly formatted",
			"constraint": "does not respect constraints",
		}

		arguments := make([]string, len(sdkError.Details))
		reasons := make([]string, len(sdkError.Details))
		hints := make([]string, len(sdkError.Details))
		for i, d := range sdkError.Details {
			arguments[i] = "'" + d.ArgumentName + "'"
			reasons[i] = "- '" + d.ArgumentName + "' " + reasonsMap[d.Reason]
			hints[i] = d.HelpMessage
		}

		return nil, &CliError{
			Message: fmt.Sprintf("invalid arguments %v", strings.Join(arguments, ", ")),
			Err:     err,
			Details: strings.Join(reasons, "\n"),
			Hint:    strings.Join(hints, "\n"),
		}

	case *scw.QuotasExceededError:
		invalidArgs := make([]string, len(sdkError.Details))
		resources := make([]string, len(sdkError.Details))
		for i, d := range sdkError.Details {
			invalidArgs[i] = fmt.Sprintf("- %s has reached its quota (%d/%d)", d.Resource, d.Current, d.Current)
			resources[i] = fmt.Sprintf("'%v'", d.Resource)
		}

		return nil, &CliError{
			Message: fmt.Sprintf("quota exceeded for resources %v", strings.Join(resources, ", ")),
			Err:     err,
			Details: strings.Join(invalidArgs, "\n"),
			Hint:    "Quotas are defined by organization. You should either delete unused resources or contact support to obtain bigger quotas.",
		}
	case *scw.TransientStateError:
		return nil, &CliError{
			Message: fmt.Sprintf("transient state error for resource '%v'", sdkError.Resource),
			Err:     err,
			Details: fmt.Sprintf("resource %s with ID %s is in a transient state '%s'",
				sdkError.Resource,
				sdkError.ResourceID,
				sdkError.CurrentState),
		}
	case *scw.OutOfStockError:
		return nil, &CliError{
			Message: fmt.Sprintf("resource out of stock '%v'", sdkError.Resource),
			Err:     err,
			Hint:    "Try again later :-)",
		}
	case *scw.ResourceExpiredError:
		var hint string
		switch resourceName := sdkError.Resource; resourceName {
		case "account_token":
			hint = "Try to generate a new token here https://console.scaleway.com/account/credentials"
		}

		return nil, &CliError{
			Message: fmt.Sprintf("resource %s with ID %s expired since %s", sdkError.Resource, sdkError.ResourceID, sdkError.ExpiredSince.String()),
			Err:     err,
			Hint:    hint,
		}
	}

	return res, err
}

// sdkStdErrorInterceptor is a command interceptor that will catch sdk standard error and return more friendly CLI error.
func sdkStdTypeInterceptor(ctx context.Context, args interface{}, runner CommandRunner) (interface{}, error) {
	res, err := runner(ctx, args)
	if err != nil {
		return res, err
	}
	switch sdkValue := res.(type) {
	case *scw.File:
		fileContent, err := ioutil.ReadAll(sdkValue.Content)
		if err != nil {
			return nil, err
		}
		return string(fileContent), nil
	}

	return res, err
}
