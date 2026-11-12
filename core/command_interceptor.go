package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

// CombineCommandInterceptor will combine one or more CommandInterceptor.
// Resulting CommandInterceptor can be viewed as a chain composed by all interceptors passed as parameter.
func CombineCommandInterceptor(interceptors ...CommandInterceptor) CommandInterceptor {
	var combinedInterceptors CommandInterceptor
	for _, interceptor := range interceptors {
		// Avoid context leaking on the anonymous function below using a variable loop-scoped
		localInterceptor := interceptor
		if interceptor == nil {
			continue
		}
		if combinedInterceptors == nil {
			combinedInterceptors = interceptor

			continue
		}

		previousInterceptor := combinedInterceptors
		combinedInterceptors = func(ctx context.Context, args any, runner CommandRunner) (any, error) {
			return previousInterceptor(
				ctx,
				args,
				func(ctx context.Context, _ any) (any, error) {
					return localInterceptor(ctx, args, runner)
				},
			)
		}
	}

	return combinedInterceptors
}

// sdkStdErrorInterceptor is a command interceptor that will catch sdk standard error and return more friendly CLI error.
func sdkStdErrorInterceptor(
	ctx context.Context,
	args any,
	runner CommandRunner,
) (any, error) {
	res, err := runner(ctx, args)

	if resourceNotFoundError, ok := errors.AsType[*scw.ResourceNotFoundError](err); ok {
		return nil, &CliError{
			Message: fmt.Sprintf(
				"cannot find resource '%v' with ID '%v'",
				resourceNotFoundError.Resource,
				resourceNotFoundError.ResourceID,
			),
			Err: err,
		}
	}

	if responseError, ok := errors.AsType[*scw.ResponseError](err); ok {
		return nil, &CliError{
			Message: responseError.Message,
			Err:     responseError,
		}
	}

	if invalidArgumentsError, ok := errors.AsType[*scw.InvalidArgumentsError](err); ok {
		reasonsMap := map[string]string{
			"unknown":    "is invalid for unexpected reasons",
			"required":   "is required",
			"format":     "is wrongly formatted",
			"constraint": "does not respect constraints",
		}

		arguments := make([]string, len(invalidArgumentsError.Details))
		reasons := make([]string, len(invalidArgumentsError.Details))
		hints := make([]string, len(invalidArgumentsError.Details))
		for i, d := range invalidArgumentsError.Details {
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
	}

	if quotasExceededError, ok := errors.AsType[*scw.QuotasExceededError](err); ok {
		invalidArgs := make([]string, len(quotasExceededError.Details))
		resources := make([]string, len(quotasExceededError.Details))
		for i, d := range quotasExceededError.Details {
			invalidArgs[i] = fmt.Sprintf(
				"- %s has reached its quota (%d/%d)",
				d.Resource,
				d.Current,
				d.Quota,
			)
			resources[i] = fmt.Sprintf("'%v'", d.Resource)
		}

		return nil, &CliError{
			Message: fmt.Sprintf("quota exceeded for resources %v", strings.Join(resources, ", ")),
			Err:     err,
			Details: strings.Join(invalidArgs, "\n"),
			Hint:    "Quotas are defined by organization. You should either delete unused resources or contact support to obtain bigger quotas.",
		}
	}

	if transientStateError, ok := errors.AsType[*scw.TransientStateError](err); ok {
		return nil, &CliError{
			Message: fmt.Sprintf(
				"transient state error for resource '%v'",
				transientStateError.Resource,
			),
			Err: err,
			Details: fmt.Sprintf("resource %s with ID %s is in a transient state '%s'",
				transientStateError.Resource,
				transientStateError.ResourceID,
				transientStateError.CurrentState),
		}
	}

	if outOfStockError, ok := errors.AsType[*scw.OutOfStockError](err); ok {
		return nil, &CliError{
			Message: fmt.Sprintf("resource out of stock '%v'", outOfStockError.Resource),
			Err:     err,
			Hint:    "Try again later :-)",
		}
	}

	if resourceExpiredError, ok := errors.AsType[*scw.ResourceExpiredError](err); ok {
		var hint string
		if resourceExpiredError.Resource == "account_token" {
			hint = "Try to generate a new token here https://console.scaleway.com/iam/api-keys"
		}

		return nil, &CliError{
			Message: fmt.Sprintf(
				"resource %s with ID %s expired since %s",
				resourceExpiredError.Resource,
				resourceExpiredError.ResourceID,
				resourceExpiredError.ExpiredSince.String(),
			),
			Err:  err,
			Hint: hint,
		}
	}

	return res, err
}

// sdkStdErrorInterceptor is a command interceptor that will catch sdk standard error and return more friendly CLI error.
func sdkStdTypeInterceptor(
	ctx context.Context,
	args any,
	runner CommandRunner,
) (any, error) {
	res, err := runner(ctx, args)
	if err != nil {
		return res, err
	}
	if sdkValue, ok := res.(*scw.File); ok {
		ExtractLogger(ctx).Debug("Intercepting scw.File type, rendering as string")
		fileContent, err := io.ReadAll(sdkValue.Content)
		if err != nil {
			return nil, err
		}

		return string(fileContent), nil
	}

	return res, err
}
