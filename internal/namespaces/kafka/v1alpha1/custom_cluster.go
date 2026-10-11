package kafka

import (
	"context"
	"errors"
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

	c.ArgSpecs.AddBefore("user-name", secrets.FastConnectArgSpec())

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
			if err := secrets.PersistCredentials(
				ctx,
				cluster.Region,
				"kafka",
				string(cluster.Region),
				cluster.ID,
				*request.UserName,
				*request.Password,
				"cluster",
			); err != nil {
				return nil, err
			}
		}

		return cluster, nil
	}

	return c
}
