package cockpit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	cockpitActionTimeout = 3 * time.Minute

	cockpitStatusMarshalSpecs = human.EnumMarshalSpecs{
		cockpit.CockpitStatusCreating: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "creating"},
		cockpit.CockpitStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
		cockpit.CockpitStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "deleting"},
		cockpit.CockpitStatusUpdating: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "updating"},
		cockpit.CockpitStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
	}
)

func cockpitWaitCommand() *core.Command {
	type cockpitWaitRequest struct {
		ProjectID string
	}

	return &core.Command{
		Short:     `Wait for a cockpit to reach a stable state (installation)`,
		Long:      `Wait for a cockpit to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the cockpit.`,
		Namespace: "cockpit",
		Resource:  "cockpit",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(cockpitWaitRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := cockpit.NewAPI(core.ExtractClient(ctx))
			logger.Debugf("starting to wait for cockpit to reach a stable delivery status")
			targetCockpit, err := api.WaitForCockpit(&cockpit.WaitForCockpitRequest{
				ProjectID:     argsI.(*cockpitWaitRequest).ProjectID,
				Timeout:       scw.TimeDurationPtr(cockpitActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
			if err != nil {
				return nil, err
			}

			if targetCockpit.Status != cockpit.CockpitStatusReady {
				return nil, &core.CliError{
					Err:     fmt.Errorf("cockpit did not reach a stable delivery status"),
					Details: fmt.Sprintf("cockpit %s is in %s status", targetCockpit.ProjectID, targetCockpit.Status),
				}
			}

			return targetCockpit, nil
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `The ID of the project the cockpit is attached to`,
				Required:   true,
				Positional: true,
			},
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a cockpit to reach a stable state",
				ArgsJSON: `{"project_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func cockpitCockpitActivateBuilder(command *core.Command) *core.Command {
	command.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		req := argsI.(*cockpit.ActivateCockpitRequest)

		client := core.ExtractClient(ctx)
		api := cockpit.NewAPI(client)
		return api.WaitForCockpit(&cockpit.WaitForCockpitRequest{
			ProjectID:     req.ProjectID,
			Timeout:       scw.TimeDurationPtr(cockpitActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}
	return command
}

func cockpitCockpitDeactivateBuilder(command *core.Command) *core.Command {
	command.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		req := argsI.(*cockpit.DeactivateCockpitRequest)

		client := core.ExtractClient(ctx)
		api := cockpit.NewAPI(client)
		_, err := api.WaitForCockpit(&cockpit.WaitForCockpitRequest{
			ProjectID:     req.ProjectID,
			Timeout:       scw.TimeDurationPtr(cockpitActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			// if we get a 404 here, it means the resource was successfully deleted
			notFoundError := &scw.ResourceNotFoundError{}
			responseError := &scw.ResponseError{}
			if errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound || errors.As(err, &notFoundError) {
				return fmt.Sprintf("Cockpit %s successfully deleted.", req.ProjectID), nil
			}

			return nil, err
		}

		return nil, fmt.Errorf("cockpit %s is still active", req.ProjectID)
	}
	return command
}
