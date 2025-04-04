package instance

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"strconv"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type instanceSSHServerRequest struct {
	Zone     scw.Zone
	ServerID string
	Username string
	Port     uint64
	Command  string
}

func serverSSHCommand() *core.Command {
	return &core.Command{
		Short:     `SSH into a server`,
		Long:      `Connect to distant server via the SSH protocol.`,
		Namespace: "instance",
		Verb:      "ssh",
		Resource:  "server",
		ArgsType:  reflect.TypeOf(instanceSSHServerRequest{}),
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
				Default: core.DefaultValueSetter("root"),
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
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
		Run: instanceServerSSHRun,
	}
}

func instanceServerSSHRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	args := argsI.(*instanceSSHServerRequest)

	client := core.ExtractClient(ctx)
	apiInstance := instance.NewAPI(client)
	serverResp, err := apiInstance.GetServer(&instance.GetServerRequest{
		Zone:     args.Zone,
		ServerID: args.ServerID,
	})
	if err != nil {
		return nil, err
	}

	if serverResp.Server.State != instance.ServerStateRunning {
		return nil, &core.CliError{
			Err: errors.New("server is not running"),
			Hint: fmt.Sprintf(
				"Start the instance with: %s instance server start %s --wait",
				core.ExtractBinaryName(ctx),
				serverResp.Server.ID,
			),
		}
	}

	if serverResp.Server.PublicIP == nil {
		return nil, &core.CliError{
			Err: errors.New("server does not have a public IP to connect to"),
			Hint: fmt.Sprintf(
				"Add a public IP to the instance with: %s instance server update %s ip=<ip_id>",
				core.ExtractBinaryName(ctx),
				serverResp.Server.ID,
			),
		}
	}

	sshArgs := []string{
		serverResp.Server.PublicIP.Address.String(),
		"-p", strconv.FormatUint(args.Port, 10),
		"-l", args.Username,
		"-t",
	}
	if args.Command != "" {
		sshArgs = append(sshArgs, args.Command)
	}

	sshCmd := exec.Command("ssh", sshArgs...)

	exitCode, err := core.ExecCmd(ctx, sshCmd)
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		return nil, &core.CliError{Empty: true, Code: exitCode}
	}

	return &core.SuccessResult{Empty: true}, nil
}
