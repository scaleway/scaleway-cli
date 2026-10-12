package rdb

import (
	"context"
	"errors"
	"os/exec"

	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// ConnectCmd returns the command to connect to a database instance with
// the client matching its engine (psql for PostgreSQL, mysql for MySQL):
// the core of the `scw rdb instance connect` command. The caller runs it
// on the terminal (the CLI via core.ExecCmd, the TUI suspended via
// tea.ExecProcess). The password is not part of the command: the client
// prompts for it on the terminal (Scaleway does not return user
// passwords).
func ConnectCmd(
	ctx context.Context,
	client *scw.Client,
	region scw.Region,
	instanceID string,
	username string,
	privateNetwork bool,
) (*exec.Cmd, EngineFamily, error) {
	args := &instanceConnectArgs{Username: username}
	api := rdbSDK.NewAPI(client)
	instance, err := api.GetInstance(&rdbSDK.GetInstanceRequest{
		Region:     region,
		InstanceID: instanceID,
	})
	if err != nil {
		return nil, Unknown, err
	}
	family, err := detectEngineFamily(instance)
	if err != nil {
		return nil, Unknown, err
	}
	if len(instance.Endpoints) == 0 {
		return nil, family, errors.New(errorMessageEndpointNotFound)
	}
	var endpoint *rdbSDK.Endpoint
	if privateNetwork {
		endpoint, err = getPrivateEndpoint(instance.Endpoints)
	} else {
		endpoint, err = getPublicEndpoint(instance.Endpoints)
	}
	if err != nil {
		return nil, family, err
	}
	cmdArgs, err := createConnectCommandLineArgs(endpoint, family, args)
	if err != nil {
		return nil, family, err
	}

	return exec.Command(cmdArgs[0], cmdArgs[1:]...), family, nil //nolint:gosec
}

// ConnectUser returns the user to connect an instance with: the first
// admin user, falling back to the first user, then "admin".
func ConnectUser(
	ctx context.Context,
	client *scw.Client,
	region scw.Region,
	instanceID string,
) (string, error) {
	resp, err := rdbSDK.NewAPI(client).ListUsers(
		&rdbSDK.ListUsersRequest{Region: region, InstanceID: instanceID},
		scw.WithContext(ctx),
	)
	if err != nil {
		return "", err
	}
	for _, u := range resp.Users {
		if u.IsAdmin {
			return u.Name, nil
		}
	}
	for _, u := range resp.Users {
		if u.Name != "" {
			return u.Name, nil
		}
	}

	return "admin", nil
}
