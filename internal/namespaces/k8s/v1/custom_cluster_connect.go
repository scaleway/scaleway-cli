package k8s

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const errorMessageK9sNotFound = "k9s is not installed. Please install k9s to use this command (https://k9scli.io)"

type clusterConnectArgs struct {
	ClusterID string
	Region    scw.Region
	Template  string
}

func clusterConnectCommand() *core.Command {
	return &core.Command{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "connect",
		Short:     `Connect to a cluster using k9s`,
		Long: `Connect to a Kubernetes Kapsule cluster using the locally installed k9s.
The command checks that k9s is installed, downloads the cluster kubeconfig to a temporary file and opens k9s against it.

Use --template to run an arbitrary tool: the Go text/template is rendered with
{{.ClusterID}}, {{.Region}} and {{.Kubeconfig}} (path to a temporary kubeconfig file),
and the result is executed as the connect command.`,
		ArgsType: reflect.TypeFor[clusterConnectArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      "Cluster ID to connect to",
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(),
			connect.TemplateArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*clusterConnectArgs)

			if args.Template == "" {
				exitCode, err := core.ExecCmd(ctx, exec.Command("k9s", "version"))
				if err != nil || exitCode != 0 {
					return nil, errors.New(errorMessageK9sNotFound)
				}
			}

			kubeconfigRes, err := k8sKubeconfigGetRun(ctx, &k8sKubeconfigGetRequest{
				ClusterID:  args.ClusterID,
				Region:     args.Region,
				AuthMethod: authMethods(defaultAuthMethod),
			})
			if err != nil {
				return nil, err
			}

			kubeconfig := kubeconfigRes.(string)

			// The kubeconfig contains secrets, keep it 0600 and clean it up on exit
			kubeconfigFile, err := os.CreateTemp("", "scw-k8s-kubeconfig-*.yaml")
			if err != nil {
				return nil, err
			}

			defer func() {
				kubeconfigFile.Close()
				os.Remove(kubeconfigFile.Name())
			}()

			if _, err := kubeconfigFile.WriteString(kubeconfig); err != nil {
				return nil, err
			}

			cmdArgs, err := connect.RenderCommand(args.Template, map[string]any{
				"ClusterID":  args.ClusterID,
				"Region":     args.Region,
				"Kubeconfig": kubeconfigFile.Name(),
			})
			if err != nil {
				return nil, err
			}

			if cmdArgs == nil {
				cmdArgs = []string{"k9s"}
			}

			// Drop any user KUBECONFIG so the tool only sees this cluster
			env := make([]string, 0, len(os.Environ())+1)
			for _, e := range os.Environ() {
				if !strings.HasPrefix(e, "KUBECONFIG=") {
					env = append(env, e)
				}
			}

			env = append(env, "KUBECONFIG="+kubeconfigFile.Name())

			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec
			cmd.Env = env

			exitCode, err := core.ExecCmd(ctx, cmd)
			if err != nil {
				return nil, err
			}
			if exitCode != 0 {
				return nil, &core.CliError{Empty: true, Code: exitCode}
			}

			return &core.SuccessResult{Empty: true}, nil
		},
		Examples: []*core.Example{
			{
				Short: "Connect to a cluster with k9s",
				Raw:   `scw k8s cluster connect 11111111-1111-1111-1111-111111111111`,
			},
			{
				Short:    "Connect to a cluster in a specific region with k9s",
				ArgsJSON: `{"cluster_id": "11111111-1111-1111-1111-111111111111", "region": "fr-par"}`,
			},
			{
				Short: "Connect to a cluster with an arbitrary tool via template",
				Raw:   `scw k8s cluster connect 11111111-1111-1111-1111-111111111111 template="kubectl --kubeconfig {{.Kubeconfig}} get pods"`,
			},
		},
		ExcludeFromMCP: true,
	}
}
