package k8s

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type k8sPoolSetLabelRequest struct {
	PoolID string
	Key    string
	Value  string
	Region scw.Region
}

func k8sPoolSetLabelCommand() *core.Command {
	return &core.Command{
		Short:     `Add or edit a taint to a Pool`,
		Long:      `Apply a label to all nodes of the pool which will be periodically reconciled by scaleway.`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "set-label",
		Groups:    []string{"label"},
		ArgsType:  reflect.TypeOf(k8sPoolSetLabelRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			request := argsI.(*k8sPoolSetLabelRequest)

			api := k8s.NewAPI(core.ExtractClient(ctx))
			pool, err := api.GetPool(&k8s.GetPoolRequest{
				PoolID: request.PoolID,
				Region: request.Region,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			if val, ok := pool.Labels[request.Key]; ok && request.Value == val {
				return pool, nil
			}

			pool.Labels[request.Key] = request.Value

			return api.SetPoolLabels(&k8s.SetPoolLabelsRequest{
				PoolID: request.PoolID,
				Labels: pool.Labels,
				Region: request.Region,
			}, scw.WithContext(ctx))
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pool-id",
				Short:      `ID of the pool.`,
				Required:   true,
				Positional: true,
			},
			{
				Name:     "key",
				Short:    `Key of the label.`,
				Required: true,
			},
			{
				Name:     "value",
				Short:    `Value of the label.`,
				Required: true,
			},
			core.RegionArgSpec(),
		},
		Examples: []*core.Example{
			{
				Short:    "Apply a label to a specific pool",
				ArgsJSON: `{"pool_id": "11111111-1111-1111-1111-111111111111", "key": "foo", "value": "bar"}`,
			},
			{
				Short:    "Apply a full label to a specific pool",
				ArgsJSON: `{"pool_id": "11111111-1111-1111-1111-111111111111", "key": "organization.example/gpu", "value": "true"}`,
			},
		},
	}
}

type k8sPoolRemoveLabelRequest struct {
	PoolID string
	Key    string
	Region scw.Region
}

func k8sPoolRemoveLabelCommand() *core.Command {
	return &core.Command{
		Short:     `Remove a label from a Pool`,
		Long:      `Remove a label from all nodes of the pool (only apply to labels which was set through scaleway api).`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "remove-label",
		Groups:    []string{"label"},
		ArgsType:  reflect.TypeOf(k8sPoolRemoveLabelRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			request := argsI.(*k8sPoolRemoveLabelRequest)

			api := k8s.NewAPI(core.ExtractClient(ctx))
			pool, err := api.GetPool(&k8s.GetPoolRequest{
				PoolID: request.PoolID,
				Region: request.Region,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			if _, ok := pool.Labels[request.Key]; !ok {
				return pool, nil
			}

			delete(pool.Labels, request.Key)

			return api.SetPoolLabels(&k8s.SetPoolLabelsRequest{
				PoolID: request.PoolID,
				Labels: pool.Labels,
				Region: request.Region,
			}, scw.WithContext(ctx))
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pool-id",
				Short:      `ID of the pool.`,
				Required:   true,
				Positional: true,
			},
			{
				Name:     "key",
				Short:    `Key of the label.`,
				Required: true,
			},
			core.RegionArgSpec(),
		},
		Examples: []*core.Example{
			{
				Short:    "Remove a label of a specific pool",
				ArgsJSON: `{"pool_id": "11111111-1111-1111-1111-111111111111", "key": "foo"}`,
			},
		},
	}
}
