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
		ArgsType:  reflect.TypeFor[instanceSSHServerRequest](),
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
		Run:            instanceServerSSHRun,
		ExcludeFromMCP: true,
	}
}

// SSHServerCmd returns the command to SSH into the given server, reusing
// the `scw instance server ssh` checks (server running, public IP). The
// caller runs it on the terminal: the CLI via core.ExecCmd, the TUI
// suspended via tea.ExecProcess.
func SSHServerCmd(
	ctx context.Context,
	client *scw.Client,
	zone scw.Zone,
	serverID string,
	username string,
	port uint64,
	command string,
) (*exec.Cmd, error) {
	apiInstance := instance.NewAPI(client)
	serverResp, err := apiInstance.GetServer(&instance.GetServerRequest{
		Zone:     zone,
		ServerID: serverID,
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
		"-p", strconv.FormatUint(port, 10),
		"-l", username,
		"-t",
	}
	if command != "" {
		sshArgs = append(sshArgs, command)
	}

	return exec.Command("ssh", sshArgs...), nil
}

func instanceServerSSHRun(ctx context.Context, argsI any) (i any, e error) {
	args := argsI.(*instanceSSHServerRequest)

	sshCmd, err := SSHServerCmd(
		ctx,
		core.ExtractClient(ctx),
		args.Zone,
		args.ServerID,
		args.Username,
		args.Port,
		args.Command,
	)
	if err != nil {
		return nil, err
	}

	exitCode, err := core.ExecCmd(ctx, sshCmd)
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		return nil, &core.CliError{Empty: true, Code: exitCode}
	}

	return &core.SuccessResult{Empty: true}, nil
}
