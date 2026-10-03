package mongodb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Endpoint create/delete are present in the SDK but missing from the generated CLI for v1.
func endpointDeleteCommand() *core.Command {
	return &core.Command{
		Short:     `Delete a Database Instance endpoint`,
		Long:      `Delete the endpoint of a Database Instance. You must specify the endpoint_id parameter of the endpoint you want to delete. Note that you might need to update any environment configurations that point to the deleted endpoint.`,
		Namespace: "mongodb",
		Resource:  "endpoint",
		Verb:      "delete",
		ArgsType:  reflect.TypeFor[mongodb.DeleteEndpointRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "endpoint-id",
				Short:      `UUID of the Endpoint to delete`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (any, error) {
			request := args.(*mongodb.DeleteEndpointRequest)
			api := mongodb.NewAPI(core.ExtractClient(ctx))

			err := api.DeleteEndpoint(request, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Resource: "endpoint",
				Verb:     "delete",
			}, nil
		},
	}
}

func endpointCreateCommand() *core.Command {
	return &core.Command{
		Short:     `Create a new Instance endpoint`,
		Long:      `Create a new endpoint for a MongoDB® Database Instance. You can add public_network or private_network specifications to the body of the request.`,
		Namespace: "mongodb",
		Resource:  "endpoint",
		Verb:      "create",
		ArgsType:  reflect.TypeFor[mongodb.CreateEndpointRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-id",
				Short:      `UUID of the Database Instance`,
				Required:   true,
				Positional: true,
			},
			{
				Name: "endpoint.public-network",
			},
			{
				Name:  "endpoint.private-network.private-network-id",
				Short: `UUID of the Private Network`,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (any, error) {
			request := args.(*mongodb.CreateEndpointRequest)
			api := mongodb.NewAPI(core.ExtractClient(ctx))

			return api.CreateEndpoint(request, scw.WithContext(ctx))
		},
	}
}
