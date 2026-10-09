package kafka

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/secrets"
	kafka "github.com/scaleway/scaleway-sdk-go/api/kafka/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type createClusterRequestCustom struct {
	*kafka.CreateClusterRequest
	FastConnect bool
}

func clusterCreateBuilder(c *core.Command) *core.Command {
	c.ArgsType = reflect.TypeFor[createClusterRequestCustom]()

	c.ArgSpecs.AddBefore("user-name", &core.ArgSpec{
		Name:       "fast-connect",
		Short:      secrets.FastConnectArgSpecShort,
		Required:   false,
		Deprecated: false,
		Positional: false,
		Default:    core.DefaultValueSetter("false"),
	})

	c.Run = func(ctx context.Context, args any) (any, error) {
		customRequest := args.(*createClusterRequestCustom)
		request := customRequest.CreateClusterRequest

		if customRequest.FastConnect && (request.UserName == nil || request.Password == nil) {
			return nil, errors.New("fast-connect requires both user-name and password")
		}

		client := core.ExtractClient(ctx)
		api := kafka.NewAPI(client)

		cluster, err := api.CreateCluster(request, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		if customRequest.FastConnect {
			path := secrets.Path("kafka", string(cluster.Region), cluster.ID, *request.UserName)
			if err := secrets.Persist(ctx, cluster.Region, path, secrets.Credentials{
				Username: *request.UserName,
				Password: *request.Password,
			}); err != nil {
				return nil, fmt.Errorf(
					"cluster %s was created but its credentials could not be persisted: %w",
					cluster.ID,
					err,
				)
			}
		}

		return cluster, nil
	}

	return c
}
