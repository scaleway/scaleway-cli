package applesilicon

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/ssh"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
)

func serverSSHCommand() *core.Command {
	return &core.Command{
		Short:     `SSH into a server`,
		Long:      `Connect to distant server via the SSH protocol.`,
		Namespace: "apple-silicon",
		Verb:      "ssh",
		Resource:  "server",
		Groups:    []string{"utility"},
		ArgsType:  reflect.TypeFor[ssh.Request](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      "Server ID to SSH into",
				Required:   true,
				Positional: true,
			},
			{
				Name:    "username",
				Short:   "Username used for the SSH connection",
				Default: core.DefaultValueSetter("m1"),
			},
			{
				Name:    "port",
				Short:   "Port used for the SSH connection",
				Default: core.DefaultValueSetter("22"),
			},
			{
				Name:  "command",
				Short: "Command to execute on the remote server",
			},
			core.ZoneArgSpec(),
		},
		Run:            serverSSHRun,
		ExcludeFromMCP: true,
	}
}

func serverSSHRun(ctx context.Context, argsI any) (i any, e error) {
	args := argsI.(*ssh.Request)

	client := core.ExtractClient(ctx)
	asAPI := applesilicon.NewAPI(client)
	serverResp, err := asAPI.GetServer(&applesilicon.GetServerRequest{
		Zone:     args.Zone,
		ServerID: args.ServerID,
	})
	if err != nil {
		return nil, err
	}

	if serverResp.Status != applesilicon.ServerStatusReady {
		return nil, &core.CliError{
			Err:     errors.New("server is not ready"),
			Details: fmt.Sprintf("Server %s currently in %s", serverResp.Name, serverResp.Status),
		}
	}

	return ssh.Connect(ctx, *args, serverResp.IP.String())
}
