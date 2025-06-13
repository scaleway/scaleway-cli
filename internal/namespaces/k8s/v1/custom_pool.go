package k8s

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	poolActionTimeout = 10 * time.Minute
)

//
// Marshalers
//

// poolStatusMarshalSpecs marshals a k8s.PoolStatus.
var (
	poolStatusMarshalSpecs = human.EnumMarshalSpecs{
		k8s.PoolStatusScaling:   &human.EnumMarshalSpec{Attribute: color.FgBlue},
		k8s.PoolStatusReady:     &human.EnumMarshalSpec{Attribute: color.FgGreen},
		k8s.PoolStatusLocked:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		k8s.PoolStatusUpgrading: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		k8s.PoolStatusWarning:   &human.EnumMarshalSpec{Attribute: color.FgHiYellow},
	}
)

const (
	poolActionCreate = iota
	poolActionUpdate
	poolActionUpgrade
	poolActionDelete
)

func poolCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForPoolFunc(poolActionCreate)

	c.ArgSpecs.GetByName("size").Default = core.DefaultValueSetter("1")
	c.ArgSpecs.GetByName("node-type").Default = core.DefaultValueSetter("DEV1-M")

	return c
}

func poolDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForPoolFunc(poolActionDelete)

	return c
}

func poolUpgradeBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForPoolFunc(poolActionUpgrade)

	return c
}

func poolUpdateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForPoolFunc(poolActionUpdate)

	return c
}

func waitForPoolFunc(action int) core.WaitFunc {
	return func(ctx context.Context, _, respI any) (any, error) {
		pool, err := k8s.NewAPI(core.ExtractClient(ctx)).WaitForPool(&k8s.WaitForPoolRequest{
			Region:        respI.(*k8s.Pool).Region,
			PoolID:        respI.(*k8s.Pool).ID,
			Timeout:       scw.TimeDurationPtr(poolActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		switch action {
		case poolActionCreate:
			return pool, err
		case poolActionUpdate:
			return pool, err
		case poolActionUpgrade:
			return pool, err
		case poolActionDelete:
			if err != nil {
				// if we get a 404 here, it means the resource was successfully deleted
				notFoundError := &scw.ResourceNotFoundError{}
				responseError := &scw.ResponseError{}
				if errors.As(err, &responseError) &&
					responseError.StatusCode == http.StatusNotFound ||
					errors.As(err, &notFoundError) {
					return fmt.Sprintf("Pool %s successfully deleted.", respI.(*k8s.Pool).ID), nil
				}
			}
		}

		return nil, err
	}
}

func k8sPoolWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for a pool to reach a stable state`,
		Long:      `Wait for a pool to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the node.`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(k8s.WaitForPoolRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			api := k8s.NewAPI(core.ExtractClient(ctx))

			return api.WaitForPool(&k8s.WaitForPoolRequest{
				Region:        argsI.(*k8s.WaitForPoolRequest).Region,
				PoolID:        argsI.(*k8s.WaitForPoolRequest).PoolID,
				Timeout:       argsI.(*k8s.WaitForPoolRequest).Timeout,
				RetryInterval: core.DefaultRetryInterval,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pool-id",
				Short:      `ID of the pool.`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(),
			core.WaitTimeoutArgSpec(poolActionTimeout),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a pool to reach a stable state",
				ArgsJSON: `{"pool_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

type k8sPoolAddExternalNodeRequest struct {
	NodeIP   string
	PoolID   string
	Username string
	Region   scw.Region
}

func k8sPoolAddExternalNodeCommand() *core.Command {
	return &core.Command{
		Short: `Add an external node to a Kosmos Pool`,
		Long: `Add an external node to a Kosmos Pool. 
This will connect via SSH to the node, download the multicloud configuration script and run it with sudo privileges.
Keep in mind that your external node needs to have wget in order to download the script.`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "add-external-node",
		ArgsType:  reflect.TypeOf(k8sPoolAddExternalNodeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "node-ip",
				Short:    `IP address of the external node`,
				Required: true,
			},
			{
				Name:     "pool-id",
				Short:    `ID of the Pool the node should be added to`,
				Required: true,
			},
			{
				Name:    "username",
				Short:   "Username used for the SSH connection",
				Default: core.DefaultValueSetter("root"),
			},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			args := argsI.(*k8sPoolAddExternalNodeRequest)
			sshCommonArgs := []string{
				args.NodeIP,
				"-t",
				"-l", args.Username,
			}

			// Set POOL_ID and REGION in the node init script and copy it to the remote
			homeDir := "/root"
			if args.Username != "root" {
				homeDir = "/home/" + args.Username
			}
			nodeInitScript := buildNodeInitScript(args.PoolID, args.Region)
			copyScriptArgs := []string{
				"cat", "<<", "EOF",
				">", homeDir + "/init_kosmos_node.sh",
				"\n",
			}
			copyScriptArgs = append(copyScriptArgs, strings.Split(nodeInitScript, " \n")...)
			if err = execSSHCommand(ctx, append(sshCommonArgs, copyScriptArgs...), true); err != nil {
				return nil, err
			}
			chmodArgs := []string{"chmod", "+x", homeDir + "/init_kosmos_node.sh"}
			if err = execSSHCommand(ctx, append(sshCommonArgs, chmodArgs...), true); err != nil {
				return nil, err
			}

			// Execute the script with SCW_SECRET_KEY set
			client := core.ExtractClient(ctx)
			secretKey, _ := client.GetSecretKey()
			execScriptArgs := []string{
				"", // Adding a space to prevent the command from being logged in history as it shows the secret key
				"SCW_SECRET_KEY=" + secretKey,
				"./init_kosmos_node.sh",
			}
			if err = execSSHCommand(ctx, append(sshCommonArgs, execScriptArgs...), false); err != nil {
				return nil, err
			}

			return &core.SuccessResult{Empty: true}, nil
		},
	}
}

func execSSHCommand(ctx context.Context, args []string, printSeparator bool) error {
	remoteCmd := exec.Command("ssh", args...)
	_, _ = interactive.Println(remoteCmd)

	exitCode, err := core.ExecCmd(ctx, remoteCmd)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("ssh command failed with exit code %d", exitCode)
	}
	if printSeparator {
		_, _ = interactive.Println("-----")
	}

	return nil
}

func buildNodeInitScript(poolID string, region scw.Region) string {
	return fmt.Sprintf(`#!/usr/bin/env sh

set -e
wget https://scwcontainermulticloud.s3.fr-par.scw.cloud/node-agent_linux_amd64 --no-verbose
chmod +x node-agent_linux_amd64
export POOL_ID=%s  POOL_REGION=%s SCW_SECRET_KEY=\$SCW_SECRET_KEY
sudo -E ./node-agent_linux_amd64 -loglevel 0 -no-controller
EOF`, poolID, region.String())
}
