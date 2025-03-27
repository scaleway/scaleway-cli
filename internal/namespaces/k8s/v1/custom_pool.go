package k8s

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	poolActionTimeout = 10 * time.Minute
)

//
// Marshalers
//

// poolStatusMarshalSpecs marshals a k8s.PoolStatus.
var (
	poolStatusMarshalSpecs = human.EnumMarshalSpecs{
		k8s.PoolStatusScaling:   &human.EnumMarshalSpec{Attribute: color.FgBlue},
		k8s.PoolStatusReady:     &human.EnumMarshalSpec{Attribute: color.FgGreen},
		k8s.PoolStatusLocked:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		k8s.PoolStatusUpgrading: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		k8s.PoolStatusWarning:   &human.EnumMarshalSpec{Attribute: color.FgHiYellow},
	}
)

const (
	poolActionCreate = iota
	poolActionUpdate
	poolActionUpgrade
	poolActionDelete
)

func poolCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForPoolFunc(poolActionCreate)

	c.ArgSpecs.GetByName("size").Default = core.DefaultValueSetter("1")
	c.ArgSpecs.GetByName("node-type").Default = core.DefaultValueSetter("DEV1-M")

	return c
}

func poolDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForPoolFunc(poolActionDelete)

	return c
}

func poolUpgradeBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForPoolFunc(poolActionUpgrade)

	return c
}

func poolUpdateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForPoolFunc(poolActionUpdate)

	return c
}

func waitForPoolFunc(action int) core.WaitFunc {
	return func(ctx context.Context, _, respI interface{}) (interface{}, error) {
		pool, err := k8s.NewAPI(core.ExtractClient(ctx)).WaitForPool(&k8s.WaitForPoolRequest{
			Region:        respI.(*k8s.Pool).Region,
			PoolID:        respI.(*k8s.Pool).ID,
			Timeout:       scw.TimeDurationPtr(poolActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		switch action {
		case poolActionCreate:
			return pool, err
		case poolActionUpdate:
			return pool, err
		case poolActionUpgrade:
			return pool, err
		case poolActionDelete:
			if err != nil {
				// if we get a 404 here, it means the resource was successfully deleted
				notFoundError := &scw.ResourceNotFoundError{}
				responseError := &scw.ResponseError{}
				if errors.As(err, &responseError) &&
					responseError.StatusCode == http.StatusNotFound ||
					errors.As(err, &notFoundError) {
					return fmt.Sprintf("Pool %s successfully deleted.", respI.(*k8s.Pool).ID), nil
				}
			}
		}

		return nil, err
	}
}

func k8sPoolWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for a pool to reach a stable state`,
		Long:      `Wait for a pool to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the node.`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(k8s.WaitForPoolRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := k8s.NewAPI(core.ExtractClient(ctx))

			return api.WaitForPool(&k8s.WaitForPoolRequest{
				Region:        argsI.(*k8s.WaitForPoolRequest).Region,
				PoolID:        argsI.(*k8s.WaitForPoolRequest).PoolID,
				Timeout:       argsI.(*k8s.WaitForPoolRequest).Timeout,
				RetryInterval: core.DefaultRetryInterval,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pool-id",
				Short:      `ID of the pool.`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(),
			core.WaitTimeoutArgSpec(poolActionTimeout),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a pool to reach a stable state",
				ArgsJSON: `{"pool_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
