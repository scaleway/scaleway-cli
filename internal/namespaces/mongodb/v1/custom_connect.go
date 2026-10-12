package mongodb

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strconv"

	"github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// ConnectCmd returns the mongosh command to connect to a MongoDB instance,
// against its public endpoint and one of its users. The Scaleway API does
// not return user passwords, so mongosh prompts for it on the terminal.
// The caller runs it (the TUI does, suspended via tea.ExecProcess).
func ConnectCmd(
	ctx context.Context,
	client *scw.Client,
	region scw.Region,
	instanceID string,
) (*exec.Cmd, error) {
	api := mongodb.NewAPI(client)
	inst, err := api.GetInstance(
		&mongodb.GetInstanceRequest{Region: region, InstanceID: instanceID},
		scw.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	var host string
	var port uint32
	for _, e := range inst.Endpoints {
		if e.PrivateNetwork == nil { // the public endpoint
			host, port = e.DNSRecord, e.Port

			break
		}
	}
	if host == "" {
		return nil, errors.New("instance has no public endpoint")
	}
	user, err := connectUser(ctx, api, region, instanceID)
	if err != nil {
		return nil, err
	}
	hostPort := net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))
	uri := fmt.Sprintf("mongodb://%s@%s/?authSource=admin&tls=true", user, hostPort)

	return exec.Command("mongosh", uri), nil
}

// connectUser returns the user to connect with: the first user of the
// instance, falling back to "admin".
func connectUser(
	ctx context.Context,
	api *mongodb.API,
	region scw.Region,
	instanceID string,
) (string, error) {
	resp, err := api.ListUsers(
		&mongodb.ListUsersRequest{Region: region, InstanceID: instanceID},
		scw.WithContext(ctx),
	)
	if err != nil {
		return "", err
	}
	for _, u := range resp.Users {
		if u.Name != "" {
			return u.Name, nil
		}
	}

	return "admin", nil
}
