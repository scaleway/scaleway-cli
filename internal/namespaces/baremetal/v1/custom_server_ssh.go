package baremetal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/ssh"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
)

func serverSSHCommand() *core.Command {
	return &core.Command{
		Short:     `SSH into a server`,
		Long:      `Connect to distant server via the SSH protocol.`,
		Namespace: "baremetal",
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
				Name:  "username",
				Short: "Username used for the SSH connection (defaults to the server install user, or root if the server is not installed)",
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
	api := baremetal.NewAPI(client)
	serverResp, err := api.GetServer(&baremetal.GetServerRequest{
		Zone:     args.Zone,
		ServerID: args.ServerID,
	})
	if err != nil {
		return nil, err
	}

	if serverResp.Status != baremetal.ServerStatusReady {
		return nil, &core.CliError{
			Err: errors.New("server is not ready"),
			Hint: fmt.Sprintf(
				"Start the server with: %s baremetal server start %s --wait",
				core.ExtractBinaryName(ctx),
				serverResp.ID,
			),
		}
	}

	ip, err := serverSSHIP(ctx, serverResp)
	if err != nil {
		return nil, err
	}

	if args.Username == "" {
		if serverResp.Install != nil && serverResp.Install.User != "" {
			args.Username = serverResp.Install.User
		} else {
			args.Username = "root"
		}
	}

	return ssh.Connect(ctx, *args, ip.String())
}

// serverSSHIP returns the public IP to connect to.
func serverSSHIP(ctx context.Context, server *baremetal.Server) (net.IP, error) {
	for _, ip := range server.IPs {
		if ip.Version == baremetal.IPVersionIPv4 {
			return ip.Address, nil
		}
	}

	return nil, &core.CliError{
		Err: errors.New("server does not have a public IP to connect to"),
		Hint: fmt.Sprintf(
			"Add a public IP to the server with: %s baremetal server add-flexible-ip %s ip-type=IPv4",
			core.ExtractBinaryName(ctx),
			server.ID,
		),
	}
}
