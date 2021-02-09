package applesilicon

import (
	"context"
	"fmt"
	"os/exec"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type serverSSHConnectRequest struct {
	Zone     scw.Zone
	ServerID string
	Username string
	Port     uint
	Command  string
}

func serverSSHCommand() *core.Command {
	return &core.Command{
		Short:     `SSH into a server`,
		Long:      `Connect to distant server via the SSH protocol.`,
		Namespace: "apple-silicon",
		Verb:      "ssh",
		Resource:  "server",
		ArgsType:  reflect.TypeOf(serverSSHConnectRequest{}),
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
		Run: serverSSHRun,
	}
}

func serverSSHRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	args := argsI.(*serverSSHConnectRequest)

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
			Err:     fmt.Errorf("server is not ready"),
			Details: fmt.Sprintf("Server %s currently in %s", serverResp.Name, serverResp.Status),
		}
	}

	sshArgs := []string{
		serverResp.IP.String(),
		"-p", fmt.Sprintf("%d", args.Port),
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
