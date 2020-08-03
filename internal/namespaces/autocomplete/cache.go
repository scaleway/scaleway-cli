package autocomplete

import (
	"context"
	"fmt"
	"reflect"

	// Used for cache
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/scaleway/scaleway-cli/internal/core"
	account "github.com/scaleway/scaleway-cli/internal/namespaces/account/v2alpha1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/baremetal/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/k8s/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/lb/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/rdb/v1"
	"github.com/scaleway/scaleway-cli/internal/namespaces/registry/v1"
)

func autocompleteCacheCommand() *core.Command {
	return &core.Command{
		Short:            `Autocomplete cache commands`,
		Long:             `Autocomplete cache commands.`,
		Namespace:        "autocomplete",
		Resource:         "cache",
		DisableTelemetry: true,
		ArgsType:         reflect.TypeOf(struct{}{}),
	}
}

type autocompleteRefreshArgs struct {
	Namespace string
}

func autocompleteRefreshCommand() *core.Command {
	return &core.Command{
		Short:            `Refresh cache for a given namespace`,
		Long:             `Refresh cache for a given namespace.`,
		Namespace:        "autocomplete",
		Resource:         "cache",
		Verb:             "refresh",
		DisableTelemetry: true,
		ArgsType:         reflect.TypeOf(autocompleteRefreshArgs{}),
		ArgSpecs: []*core.ArgSpec{
			{
				Name:       "namespace",
				Short:      "Name of the namespace you want to refresh the cache of",
				Positional: true,
				Required:   true,
				EnumValues: []string{
					"all",
					"account",
					"baremetal",
					"instance",
					"k8s",
					"lb",
					"rdb",
					"registry",
				},
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*autocompleteRefreshArgs)
			switch args.Namespace {
			case "all":
				_, err := account.RefreshCache(ctx)
				if err != nil {
					return nil, err
				}
				_, err = baremetal.RefreshCache(ctx)
				if err != nil {
					return nil, err
				}
				_, err = instance.RefreshCache(ctx)
				if err != nil {
					return nil, err
				}
				_, err = k8s.RefreshCache(ctx)
				if err != nil {
					return nil, err
				}
				_, err = lb.RefreshCache(ctx)
				if err != nil {
					return nil, err
				}
				_, err = rdb.RefreshCache(ctx)
				if err != nil {
					return nil, err
				}
				_, err = registry.RefreshCache(ctx)
				if err != nil {
					return nil, err
				}
				return fmt.Sprintln("Cache refreshed for all resources"), nil
			case "account":
				return account.RefreshCache(ctx)
			case "baremetal":
				return baremetal.RefreshCache(ctx)
			case "instance":
				return instance.RefreshCache(ctx)
			case "k8s":
				return k8s.RefreshCache(ctx)
			case "lb":
				return lb.RefreshCache(ctx)
			case "rdb":
				return rdb.RefreshCache(ctx)
			case "registry":
				return registry.RefreshCache(ctx)
			default:
				return nil, fmt.Errorf("unknown cached namespace")
			}
		},
	}
}
