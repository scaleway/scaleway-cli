package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/human"
	sdk "github.com/scaleway/scaleway-sdk-go/scw"
)

func init() {
	human.RegisterMarshalerFunc((*sdk.ResponseError)(nil), sdkResponseErrorHumanMarshalerFunc())
	human.RegisterMarshalerFunc((*sdk.InvalidArgumentsError)(nil), sdkInvalidArgumentsErrorHumanMarshalerFunc())
	human.RegisterMarshalerFunc((*sdk.QuotasExceededError)(nil), sdkQuotasExceededErrorHumanMarshalerFunc())
	human.RegisterMarshalerFunc((*sdk.TransientStateError)(nil), sdkTransientStateErrorHumanMarshalerFunc())
	human.RegisterMarshalerFunc((*sdk.ResourceNotFoundError)(nil), sdkResourceNotFoundErrorHumanMarshalerFunc())
	human.RegisterMarshalerFunc((*sdk.OutOfStockError)(nil), sdkOutOfStockErrorHumanMarshalerFunc())
	human.RegisterMarshalerFunc((*sdk.ResourceExpiredError)(nil), sdkResourceExpiredHumanMarshallFunc())
}

// CliError is an all-in-one error structure that can be used in commands to return useful errors to the user.
// CliError implements JSON and human marshaler for a smooth experience.
type CliError struct {
	Err     error
	Details string
	Hint    string
}

func (s *CliError) Error() string {
	return s.Err.Error()
}

func (s *CliError) MarshalHuman() (string, error) {
	sections := []string(nil)
	if s.Err != nil {
		str, err := human.Marshal(s.Err, nil)
		if err != nil {
			return "", err
		}
		sections = append(sections, str)
	}

	if s.Details != "" {
		str, err := human.Marshal(s.Details, &human.MarshalOpt{Title: "Details"})
		if err != nil {
			return "", err
		}
		sections = append(sections, str)
	}

	if s.Hint != "" {
		str, err := human.Marshal(s.Hint, &human.MarshalOpt{Title: "Hint"})
		if err != nil {
			return "", err
		}
		sections = append(sections, str)
	}

	return strings.Join(sections, "\n\n"), nil
}

func (s *CliError) MarshalJSON() ([]byte, error) {
	type tmpRes struct {
		Error   string `json:"error"`
		Details string `json:"details"`
		Hint    string `json:"hint"`
	}
	return json.Marshal(&tmpRes{
		Error:   s.Err.Error(),
		Details: s.Details,
		Hint:    s.Hint,
	})
}

func sdkResponseErrorHumanMarshalerFunc() human.MarshalerFunc {
	return func(i interface{}, opt *human.MarshalOpt) (string, error) {
		responseError := i.(*sdk.ResponseError)

		return human.Marshal(&CliError{
			Err: fmt.Errorf(responseError.Message),
		}, opt)
	}
}

func sdkInvalidArgumentsErrorHumanMarshalerFunc() human.MarshalerFunc {
	return func(i interface{}, opt *human.MarshalOpt) (string, error) {
		invalidArgumentsError := i.(*sdk.InvalidArgumentsError)
		reasonsMap := map[string]string{
			"unknown":    "is invalid for unexpected reason",
			"required":   "is required",
			"format":     "is wrongly formatted",
			"constraint": "does not respect constraint",
		}

		arguments := make([]string, len(invalidArgumentsError.Details))
		reasons := make([]string, len(invalidArgumentsError.Details))
		hints := make([]string, len(invalidArgumentsError.Details))
		for i, d := range invalidArgumentsError.Details {
			arguments[i] = "'" + d.ArgumentName + "'"
			reasons[i] = "- " + d.ArgumentName + " " + reasonsMap[d.Reason]
			hints[i] = d.HelpMessage
		}

		return human.Marshal(&CliError{
			Err:     fmt.Errorf("invalid arguments %v", strings.Join(arguments, ", ")),
			Details: strings.Join(reasons, "\n"),
			Hint:    strings.Join(hints, "\n"),
		}, opt)
	}
}

func sdkQuotasExceededErrorHumanMarshalerFunc() human.MarshalerFunc {
	return func(i interface{}, opt *human.MarshalOpt) (string, error) {
		quotasExceededError := i.(*sdk.QuotasExceededError)

		invalidArgs := make([]string, len(quotasExceededError.Details))
		resources := make([]string, len(quotasExceededError.Details))
		for i, d := range quotasExceededError.Details {
			invalidArgs[i] = fmt.Sprintf("- %s has reached its quota (%d/%d)", d.Resource, d.Current, d.Current)
			resources[i] = fmt.Sprintf("'%v'", d.Resource)
		}

		return human.Marshal(&CliError{
			Err:     fmt.Errorf("quota exceeded for resources %v", strings.Join(resources, ", ")),
			Details: strings.Join(invalidArgs, "\n"),
			Hint:    "Quotas are defined by organization. You should either delete unused resources or contact support to obtain bigger quotas.",
		}, opt)
	}
}

func sdkTransientStateErrorHumanMarshalerFunc() human.MarshalerFunc {
	return func(i interface{}, opt *human.MarshalOpt) (string, error) {
		transientStateError := i.(*sdk.TransientStateError)

		return human.Marshal(&CliError{
			Err: fmt.Errorf("transient state error for resource '%v'", transientStateError.Resource),
			Details: fmt.Sprintf("resource %s with ID %s is in a transient state '%s'",
				transientStateError.Resource,
				transientStateError.ResourceID,
				transientStateError.CurrentState),
		}, opt)
	}
}

func sdkResourceNotFoundErrorHumanMarshalerFunc() human.MarshalerFunc {
	return func(i interface{}, opt *human.MarshalOpt) (string, error) {
		resourceNotFoundError := i.(*sdk.ResourceNotFoundError)

		return human.Marshal(&CliError{
			Err: fmt.Errorf("cannot find resource '%v' with ID '%v'", resourceNotFoundError.Resource, resourceNotFoundError.ResourceID),
		}, opt)
	}
}

func sdkOutOfStockErrorHumanMarshalerFunc() human.MarshalerFunc {
	return func(i interface{}, opt *human.MarshalOpt) (string, error) {
		outOfStockError := i.(*sdk.OutOfStockError)

		return human.Marshal(&CliError{
			Err:  fmt.Errorf("resource out of stock '%v'", outOfStockError.Resource),
			Hint: "Try again later :-)",
		}, opt)
	}
}

func sdkResourceExpiredHumanMarshallFunc() human.MarshalerFunc {
	return func(i interface{}, opt *human.MarshalOpt) (string, error) {
		resourceExpiredError := i.(*sdk.ResourceExpiredError)

		var hint string
		switch resourceName := resourceExpiredError.Resource; resourceName {
		case "account_token":
			hint = "Try to generate a new token here https://console.scaleway.com/account/credentials"
		}

		return human.Marshal(&CliError{
			Err:  fmt.Errorf("resource %s with ID %s expired since %s", resourceExpiredError.Resource, resourceExpiredError.ResourceID, resourceExpiredError.ExpiredSince.String()),
			Hint: hint,
		}, opt)
	}
}
