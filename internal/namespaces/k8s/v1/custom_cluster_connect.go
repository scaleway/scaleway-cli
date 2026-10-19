package k8s

import (
	"context"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type clusterConnectArgs struct {
	ClusterID string
	Region    scw.Region
	Template  string
}

// clusterConnectTemplateData is the template data for the connect command:
// the base connect values plus the k8s cluster values.
type clusterConnectTemplateData struct {
	*connect.TemplateData
	ClusterID  string
	Region     scw.Region
	Kubeconfig string
}

// ClusterName returns the name of the cluster, resolved lazily.
func (d *clusterConnectTemplateData) ClusterName() string {
	cluster, err := k8s.NewAPI(d.Client()).GetCluster(&k8s.GetClusterRequest{
		ClusterID: d.ClusterID,
		Region:    d.Region,
	})
	if err != nil {
		return ""
	}

	return cluster.Name
}

func clusterConnectCommand() *core.Command {
	return &core.Command{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "connect",
		Short:     `Connect to a cluster using k9s`,
		Long: `Connect to a Kubernetes Kapsule cluster using the locally installed k9s.
The command downloads the cluster kubeconfig to a temporary file and opens k9s against it.

Use --template to run an arbitrary tool instead: the Go text/template is rendered with
{{.ClusterID}}, {{.ClusterName}}, {{.Region}}, {{.Kubeconfig}} (path to a temporary kubeconfig file, removed when the command exits),
{{.OrganizationID}}, {{.OrganizationName}}, {{.ProjectID}} and {{.ProjectName}},
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
			&core.ArgSpec{
				Name:    "template",
				Short:   "Go text/template rendered into a command line that replaces the default connect command",
				Default: core.DefaultValueSetter("k9s"),
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*clusterConnectArgs)

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
				errClose := kubeconfigFile.Close()
				if errClose != nil {
					core.ExtractLogger(ctx).Errorf("error closing kubeconfig: %v", errClose)

					return
				}
				errRemove := os.Remove(kubeconfigFile.Name())
				if errRemove != nil {
					core.ExtractLogger(ctx).Errorf("error removing kubeconfig: %v", errRemove)

					return
				}
			}()

			if _, err := kubeconfigFile.WriteString(kubeconfig); err != nil {
				return nil, err
			}

			templateData := &clusterConnectTemplateData{
				TemplateData: connect.NewTemplateData(ctx),
				ClusterID:    args.ClusterID,
				Region:       args.Region,
				Kubeconfig:   kubeconfigFile.Name(),
			}

			cmdArgs, err := connect.RenderCommand(args.Template, templateData)
			if err != nil {
				return nil, err
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

			exitCode, err := connect.ExecCmd(ctx, cmd)
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
				Short: "Connect to a cluster with k9s using the pulse view",
				Raw:   `scw k8s cluster connect 11111111-1111-1111-1111-111111111111 template="k9s -c pulses"`,
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
