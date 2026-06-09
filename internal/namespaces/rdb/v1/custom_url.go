package rdb

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strconv"
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

type rdbConnectionParams struct {
	InstanceID     string
	Region         scw.Region
	User           string
	Db             string //nolint: stylecheck
	PrivateNetwork bool
}

// ConnectionInfo holds resolved connection parameters for an RDB instance.
type ConnectionInfo struct {
	EngineFamily     engineFamily
	Host             string
	Port             uint32
	User             string
	Database         string
	PrivateNetworkID string
}

func (info *ConnectionInfo) hostPort() string {
	return net.JoinHostPort(info.Host, strconv.Itoa(int(info.Port)))
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
	args := argsI.(*rdbGetURLArgs)

	client := core.ExtractClient(ctx)
	api := rdb.NewAPI(client)

	info, err := resolveRDBConnectionInfo(ctx, api, &rdbConnectionParams{
		InstanceID: args.InstanceID,
		Region:     args.Region,
		User:       args.User,
		Db:         args.Db,
	})
	if err != nil {
		return nil, err
	}

	u := connectionInfoToURL(info, args.Db != "")

	return u.String(), nil
}

func resolveRDBConnectionInfo(
	ctx context.Context,
	api *rdb.API,
	params *rdbConnectionParams,
) (*ConnectionInfo, error) {
	instance, err := api.GetInstance(&rdb.GetInstanceRequest{
		Region:     params.Region,
		InstanceID: params.InstanceID,
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get instance %q", params.InstanceID)
	}

	engineFamily, err := detectEngineFamily(instance)
	if err != nil {
		return nil, err
	}

	users, err := api.ListUsers(&rdb.ListUsersRequest{
		Region:     params.Region,
		InstanceID: params.InstanceID,
		Name:       &params.User,
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, fmt.Errorf("failed to list users for instance %q", params.InstanceID)
	}
	if users.TotalCount != 1 {
		return nil, fmt.Errorf(
			"expected 1 user with the name %q, got %d",
			params.User,
			users.TotalCount,
		)
	}

	endpoint, err := resolveRDBEndpoint(instance.Endpoints, params.PrivateNetwork)
	if err != nil {
		return nil, err
	}

	database := params.Db
	if database == "" {
		database = "rdb"
	}

	info := &ConnectionInfo{
		EngineFamily: engineFamily,
		Host:         endpoint.IP.String(),
		Port:         endpoint.Port,
		User:         users.Users[0].Name,
		Database:     database,
	}
	if endpoint.PrivateNetwork != nil {
		info.PrivateNetworkID = endpoint.PrivateNetwork.PrivateNetworkID
	}

	return info, nil
}

func resolveRDBEndpoint(endpoints []*rdb.Endpoint, preferPrivate bool) (*rdb.Endpoint, error) {
	if preferPrivate {
		return getPrivateEndpoint(endpoints)
	}

	var privateEndpoint *rdb.Endpoint
	var publicEndpoint *rdb.Endpoint
	for _, endpoint := range endpoints {
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

	return endpoint, nil
}

func connectionInfoToURL(info *ConnectionInfo, includeDatabase bool) *url.URL {
	u := &url.URL{}

	switch info.EngineFamily {
	case PostgreSQL:
		u.Scheme = "postgresql"
	case MySQL:
		u.Scheme = "mysql"
	default:
		u.Scheme = strings.ToLower(string(info.EngineFamily))
	}

	u.User = url.User(info.User)
	u.Host = info.hostPort()

	if includeDatabase {
		u = u.JoinPath(info.Database)
	}

	return u
}
