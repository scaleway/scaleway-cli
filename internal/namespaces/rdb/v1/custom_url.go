package rdb

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type rdbGetURLArgs struct {
	InstanceID string
	Region     scw.Region
	User       string
	Db         string //nolint: stylecheck
}

func databaseGetURLCommand() *core.Command {
	return &core.Command{
		Namespace: "rdb",
		Resource:  "database",
		Verb:      "get-url",
		Short:     "Gets the URL to connect to the Database",
		Long:      "Provides the URL to connect to a Database on an Instance as the given user",
		ArgsType:  reflect.TypeOf(rdbGetURLArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-id",
				Short:      `ID of the Database Instance`,
				Required:   true,
				Positional: true,
			},
			{
				Name:       "user",
				Short:      `User of the Database`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "db",
				Short:      `Name of the Database to connect to`,
				Required:   false,
				Positional: false,
			},
		},
		Run: generateURL,
	}
}

func generateURL(ctx context.Context, argsI any) (any, error) {
	// var u url.URL
	u := &url.URL{}
	client := core.ExtractClient(ctx)
	api := rdb.NewAPI(client)
	args := argsI.(*rdbGetURLArgs)

	// First we need to determine the engine
	instance, err := api.GetInstance(&rdb.GetInstanceRequest{
		Region:     args.Region,
		InstanceID: args.InstanceID,
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get instance %q", args.InstanceID)
	}

	switch {
	case strings.HasPrefix(instance.Engine, string(PostgreSQL)):
		u.Scheme = "postgresql"
	case strings.HasPrefix(instance.Engine, string(MySQL)):
		u.Scheme = "mysql"
	default:
		return nil, fmt.Errorf("unknown engine %q", instance.Engine)
	}

	// Then we add the username
	users, err := api.ListUsers(&rdb.ListUsersRequest{
		Region:     args.Region,
		InstanceID: args.InstanceID,
		Name:       &args.User,
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, fmt.Errorf("failed to list users for instance %q", args.InstanceID)
	}
	if users.TotalCount != 1 {
		return nil, fmt.Errorf(
			"expected 1 user with the name %q, got %d",
			args.User,
			users.TotalCount,
		)
	}
	u.User = url.User(users.Users[0].Name)

	// Then we have to determine the endpoint
	var privateEndpoint *rdb.Endpoint
	var publicEndpoint *rdb.Endpoint
	for _, endpoint := range instance.Endpoints {
		if endpoint.PrivateNetwork != nil {
			privateEndpoint = endpoint
		} else if endpoint.LoadBalancer != nil || endpoint.DirectAccess != nil {
			publicEndpoint = endpoint
		}
	}
	endpoint := publicEndpoint
	if endpoint == nil {
		endpoint = privateEndpoint
	}
	if endpoint == nil {
		return nil, errors.New("instance has no endpoint therefore no url can be returned")
	}
	u.Host = fmt.Sprintf("%s:%d", endpoint.IP.String(), endpoint.Port)

	// Finally we add the database if it was given
	if args.Db != "" {
		u = u.JoinPath(args.Db)
	}

	return u.String(), nil
}
