package k8s

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func coreV1TaintEffectEnums() []string {
	sdkEffects := k8s.CoreV1TaintEffect("").Values()
	enums := make([]string, len(sdkEffects))
	for i, v := range sdkEffects {
		enums[i] = v.String()
	}

	return enums
}

type k8sPoolSetTaintRequest struct {
	PoolID string
	Key    string
	Value  string
	Effect string
	Region scw.Region
}

func k8sPoolSetTaintCommand() *core.Command {
	return &core.Command{
		Short:     `Apply a taint to a Pool`,
		Long:      `Apply a taint to all nodes of the pool which will be periodically reconciled by scaleway.`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "set-taint",
		Groups:    []string{"taint"},
		ArgsType:  reflect.TypeOf(k8sPoolSetTaintRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			request := argsI.(*k8sPoolSetTaintRequest)

			api := k8s.NewAPI(core.ExtractClient(ctx))
			pool, err := api.GetPool(&k8s.GetPoolRequest{
				PoolID: request.PoolID,
				Region: request.Region,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			found := false
			for _, t := range pool.Taints {
				if t.Key == request.Key {
					t.Value = request.Value
					t.Effect = k8s.CoreV1TaintEffect(request.Effect)
					found = true

					break
				}
			}

			if !found {
				pool.Taints = append(pool.Taints, &k8s.CoreV1Taint{
					Key:    request.Key,
					Value:  request.Value,
					Effect: k8s.CoreV1TaintEffect(request.Effect),
				})
			}

			return api.SetPoolTaints(&k8s.SetPoolTaintsRequest{
				PoolID: request.PoolID,
				Taints: pool.Taints,
				Region: request.Region,
			})
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
				Short:    `Key of the taint.`,
				Required: true,
			},
			{
				Name:     "value",
				Short:    `Value of the taint.`,
				Required: true,
			},
			{
				Name:       "effect",
				Short:      `Effect of the taint.`,
				Required:   true,
				EnumValues: coreV1TaintEffectEnums(),
			},
			core.RegionArgSpec(),
		},
		Examples: []*core.Example{
			{
				Short:    "Apply a taint to a specific pool",
				ArgsJSON: `{"pool_id": "11111111-1111-1111-1111-111111111111", "key": "foo", "value": "bar", "effect": "NoSchedule"}`,
			},
			{
				Short:    "Apply a full taint to a specific pool",
				ArgsJSON: `{"pool_id": "11111111-1111-1111-1111-111111111111", "key": "organization.example/gpu", "value": "true", "effect": "NoSchedule"}`,
			},
		},
	}
}

type k8sPoolRemoveTaintRequest struct {
	PoolID string
	Key    string
	Region scw.Region
}

func k8sPoolRemoveTaintCommand() *core.Command {
	return &core.Command{
		Short:     `Remove a taint from a Pool`,
		Long:      `Remove a taint from all all nodes of the pool (only apply to taints which was set through scaleway api).`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "remove-taint",
		Groups:    []string{"taint"},
		ArgsType:  reflect.TypeOf(k8sPoolRemoveTaintRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			request := argsI.(*k8sPoolRemoveTaintRequest)

			api := k8s.NewAPI(core.ExtractClient(ctx))
			pool, err := api.GetPool(&k8s.GetPoolRequest{
				PoolID: request.PoolID,
				Region: request.Region,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			found := false
			for i, t := range pool.Taints {
				if t.Key == request.Key {
					pool.Taints = append(pool.Taints[:i], pool.Taints[i+1:]...)
					found = true

					break
				}
			}

			if !found {
				return pool, nil
			}

			return api.SetPoolTaints(&k8s.SetPoolTaintsRequest{
				PoolID: request.PoolID,
				Taints: pool.Taints,
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
				Short:    `Key of the taint.`,
				Required: true,
			},
			core.RegionArgSpec(),
		},
		Examples: []*core.Example{
			{
				Short:    "Remove a taint to a specific pool",
				ArgsJSON: `{"pool_id": "11111111-1111-1111-1111-111111111111", "key": "foo"}`,
			},
		},
	}
}

type k8sPoolSetStartupTaintRequest struct {
	PoolID string
	Key    string
	Value  string
	Effect string
	Region scw.Region
}

func k8sPoolSetStartupTaintCommand() *core.Command {
	return &core.Command{
		Short:     `Apply a startup taint to a Pool`,
		Long:      `Apply a taint at node creation but does not reconcile after.`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "set-startup-taint",
		Groups:    []string{"taint"},
		ArgsType:  reflect.TypeOf(k8sPoolSetStartupTaintRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			request := argsI.(*k8sPoolSetStartupTaintRequest)

			api := k8s.NewAPI(core.ExtractClient(ctx))
			pool, err := api.GetPool(&k8s.GetPoolRequest{
				PoolID: request.PoolID,
				Region: request.Region,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			found := false
			for _, t := range pool.StartupTaints {
				if t.Key == request.Key {
					t.Value = request.Value
					t.Effect = k8s.CoreV1TaintEffect(request.Effect)
					found = true

					break
				}
			}

			if !found {
				pool.StartupTaints = append(pool.StartupTaints, &k8s.CoreV1Taint{
					Key:    request.Key,
					Value:  request.Value,
					Effect: k8s.CoreV1TaintEffect(request.Effect),
				})
			}

			return api.SetPoolStartupTaints(&k8s.SetPoolStartupTaintsRequest{
				PoolID:        request.PoolID,
				StartupTaints: pool.StartupTaints,
				Region:        request.Region,
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
				Short:    `Key of the taint.`,
				Required: true,
			},
			{
				Name:     "value",
				Short:    `Value of the taint.`,
				Required: true,
			},
			{
				Name:       "effect",
				Short:      `Effect of the taint.`,
				Required:   true,
				EnumValues: coreV1TaintEffectEnums(),
			},
			core.RegionArgSpec(),
		},
		Examples: []*core.Example{
			{
				Short:    "Apply a startup taint to a specific pool",
				ArgsJSON: `{"pool_id": "11111111-1111-1111-1111-111111111111", "key": "foo", "value": "bar", "effect": "NoSchedule"}`,
			},
			{
				Short:    "Apply a full startup taint to a specific pool",
				ArgsJSON: `{"pool_id": "11111111-1111-1111-1111-111111111111", "key": "organization.example/gpu", "value": "true", "effect": "NoSchedule"}`,
			},
		},
	}
}

type k8sPoolRemoveStartupTaintRequest struct {
	PoolID string
	Key    string
	Region scw.Region
}

func k8sPoolRemoveStartupTaintCommand() *core.Command {
	return &core.Command{
		Short:     `Remove a startup taint from a Pool`,
		Long:      `New nodes will not have this taint at startup (does not remove taints from kubernetes side).`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "remove-startup-taint",
		Groups:    []string{"taint"},
		ArgsType:  reflect.TypeOf(k8sPoolRemoveStartupTaintRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			request := argsI.(*k8sPoolRemoveStartupTaintRequest)

			api := k8s.NewAPI(core.ExtractClient(ctx))
			pool, err := api.GetPool(&k8s.GetPoolRequest{
				PoolID: request.PoolID,
				Region: request.Region,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			found := false
			for i, t := range pool.StartupTaints {
				if t.Key == request.Key {
					pool.StartupTaints = append(pool.StartupTaints[:i], pool.StartupTaints[i+1:]...)
					found = true

					break
				}
			}

			if !found {
				return pool, nil
			}

			return api.SetPoolStartupTaints(&k8s.SetPoolStartupTaintsRequest{
				PoolID:        request.PoolID,
				StartupTaints: pool.StartupTaints,
				Region:        request.Region,
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
				Short:    `Key of the taint.`,
				Required: true,
			},
			core.RegionArgSpec(),
		},
		Examples: []*core.Example{
			{
				Short:    "Remove a startup taint of a specific pool",
				ArgsJSON: `{"pool_id": "11111111-1111-1111-1111-111111111111", "key": "foo"}`,
			},
		},
	}
}
