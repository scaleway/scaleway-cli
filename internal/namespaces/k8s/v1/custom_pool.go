package k8s

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1/scripts"
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
	return func(ctx context.Context, _, respI interface{}) (interface{}, error) {
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
				if errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound || errors.As(err, &notFoundError) {
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
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
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
		Short: `Add an external node to a Kosmos pool`,
		// TODO: finish writing long description
		Long: `Add an external node to a Kosmos pool. 
This will connect via SSH to the node as root, download the multicloud configuration script and run it.
Keep in mind that you need SSH root access to the external node in order to be able to run this command.
Your external node needs to have bash and wget`,
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "add-external-node",
		//Groups:    []string{"workflow"}, // TODO: what group ?
		ArgsType: reflect.TypeOf(k8sPoolAddExternalNodeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "node-ip",
				Short:    `IP address of the external node`,
				Required: true,
			},
			{
				Name:     "pool-id",
				Short:    `ID of the pool the node should be added to`,
				Required: true,
			},
			{
				Name:    "username",
				Short:   "Username used for the SSH connection",
				Default: core.DefaultValueSetter("root"),
			},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			args := argsI.(*k8sPoolAddExternalNodeRequest)
			client := core.ExtractClient(ctx)
			secretKey, _ := client.GetSecretKey()

			// Generate the script with the environment variables set
			//type KosmosInstallEnvVars struct {
			//	PoolID     string
			//	PoolRegion string
			//	SecretKey  string
			//}
			//envVars := KosmosInstallEnvVars{
			//	PoolID:     args.PoolID,
			//	PoolRegion: args.Region.String(),
			//	SecretKey:  secretKey,
			//}
			//tmpl, err := template.New("kosmos-node-init").Parse(scripts.KosmosNodeInitScript)
			//if err != nil {
			//	return nil, fmt.Errorf("failed to parse template: %w", err)
			//}
			//renderedScript, err := os.Create(kosmosNodeInitScriptPath + "kosmos_node_init.sh")
			//if err != nil {
			//	return nil, fmt.Errorf("failed to create script: %w", err)
			//}
			//err = tmpl.Execute(renderedScript, envVars)
			//if err != nil {
			//	return nil, fmt.Errorf("failed to execute template: %w", err)
			//}

			sshCommonArgs := []string{
				args.NodeIP,
				"-l", args.Username,
			}

			// Copy the script to the node
			homeDir := "/root"
			if args.Username != "root" {
				homeDir = "/home/" + args.Username
			}
			//scpArgs := []string{
			//	scripts.KosmosNodeInitScript,
			//	fmt.Sprintf("%s@%s:%s/", args.Username, args.NodeIP, homeDir),
			//}
			//if err = execRemoteCommand(ctx, "scp", scpArgs); err != nil {
			//	return nil, err
			//}

			copyScriptArgs := []string{
				"cat", "<<", "EOF",
				">", homeDir + "/kosmos_node_init.sh",
				"\n",
			}
			copyScriptArgs = append(copyScriptArgs, strings.Split(scripts.KosmosNodeInitScript, " \n")...)
			copyScriptArgs = append(copyScriptArgs, "\nEOF")
			if err = execRemoteCommand(ctx, "ssh", append(sshCommonArgs, copyScriptArgs...)); err != nil {
				return nil, err
			}
			// Replace the POOL_ID and POOL_REGION variables in the script on the node
			replaceVarsArgs := []string{
				"sed", "-i",
				"-e", fmt.Sprintf("'s/<pool-region>/%s/'", args.Region.String()),
				"-e", fmt.Sprintf("'s/<pool-id>/%s/'", args.PoolID),
				"-e", fmt.Sprintf("'s/<secret-key>/$SCW_SECRET_KEY/'"),
				"kosmos_node_init.sh",
			}
			if err = execRemoteCommand(ctx, "ssh", append(sshCommonArgs, replaceVarsArgs...)); err != nil {
				return nil, err
			}

			chmodArgs := []string{
				"chmod", "+x", "kosmos_node_init.sh",
			}
			if err = execRemoteCommand(ctx, "ssh", append(sshCommonArgs, chmodArgs...)); err != nil {
				return nil, err
			}

			// Execute the script on the node
			execScriptArgs := []string{
				"",
				"SCW_SECRET_KEY=" + secretKey,
				"./kosmos_node_init.sh",
			}
			if err = execRemoteCommand(ctx, "ssh", append(sshCommonArgs, execScriptArgs...)); err != nil {
				return nil, err
			}

			return &core.SuccessResult{Empty: true}, nil
		},
	}
}
func execRemoteCommand(ctx context.Context, cmd string, args []string) error {
	remoteCmd := exec.Command(cmd, args...)

	_, _ = interactive.Println(remoteCmd)

	exitCode, err := core.ExecCmd(ctx, remoteCmd)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("%s command failed with exit code %d", cmd, exitCode)
	}
	return nil
}
