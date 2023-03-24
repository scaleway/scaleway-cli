package container

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	containerNamespaceActionTimeout = 5 * time.Minute

	namespaceStatusMarshalSpecs = human.EnumMarshalSpecs{
		container.NamespaceStatusCreating: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.NamespaceStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.NamespaceStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		container.NamespaceStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		container.NamespaceStatusPending:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.NamespaceStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		container.NamespaceStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint},
	}
)

func containerNamespaceWaitCommand() *core.Command {
	type containerNamespaceWaitRequest struct {
		Region scw.Region

		NamespaceID string
	}

	return &core.Command{
		Short:     `Wait for a namespace to reach a stable state (installation)`,
		Long:      `Wait for a namespace to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the namespace.`,
		Namespace: "container",
		Resource:  "namespace",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(containerNamespaceWaitRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			req := argsI.(*containerNamespaceWaitRequest)
			api := container.NewAPI(core.ExtractClient(ctx))

			logger.Debugf("starting to wait for the namespace to reach a stable delivery status")
			namespace, err := api.WaitForNamespace(&container.WaitForNamespaceRequest{
				Region:        req.Region,
				NamespaceID:   req.NamespaceID,
				Timeout:       scw.TimeDurationPtr(containerNamespaceActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
			if err != nil {
				return nil, err
			}

			if namespace.Status != container.NamespaceStatusReady {
				return nil, &core.CliError{
					Err:     fmt.Errorf("the namespace did not reach a stable delivery status"),
					Details: fmt.Sprintf("the namespace %s is in %s status", namespace.RegistryNamespaceID, namespace.Status),
				}
			}

			return namespace, nil
		},
	}
}

func containerNamespaceCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		res := respI.(*container.Namespace)

		client := core.ExtractClient(ctx)
		api := container.NewAPI(client)
		return api.WaitForNamespace(&container.WaitForNamespaceRequest{
			NamespaceID:   res.ID,
			Region:        res.Region,
			Timeout:       scw.TimeDurationPtr(containerNamespaceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func containerNamespaceDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		req := argsI.(*container.DeleteNamespaceRequest)

		client := core.ExtractClient(ctx)
		api := container.NewAPI(client)
		_, err := api.WaitForNamespace(&container.WaitForNamespaceRequest{
			NamespaceID:   req.NamespaceID,
			Region:        req.Region,
			Timeout:       scw.TimeDurationPtr(containerNamespaceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			if core.IsNotFoundError(err) {
				return nil, nil
			}

			return nil, err
		}

		return nil, nil
	}

	return c
}
