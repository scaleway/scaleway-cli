package redis

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	redisSDK "github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// defaultConnectTemplate renders the default redis-cli invocation: the
// endpoint, the TLS certificate (when TLS is enabled) and the cluster
// user. The password is not on the command line, where it would be visible
// in the process list: redis-cli picks it up from the REDISCLI_AUTH
// environment variable the connect command exports.
const defaultConnectTemplate = `redis-cli -h {{.Host}} -p {{.Port}}{{if .TLSEnabled}} --tls --cacert {{.CACert}}{{end}}{{if .UserName}} --user {{.UserName}}{{end}}`

type clusterConnectArgs struct {
	ClusterID      string
	Zone           scw.Zone
	Template       string
	Password       string
	PrivateNetwork bool
}

// clusterConnectTemplateData is the template data for the connect command:
// the base connect values plus the redis cluster values.
type clusterConnectTemplateData struct {
	*connect.TemplateData
	// The context is stored because the no-arg template method (Password)
	// cannot receive it as a parameter.
	ctx        context.Context //nolint:containedctx
	Host       string
	Port       uint32
	UserName   string
	TLSEnabled bool
	CACert     string
	ClusterID  string
	Zone       scw.Zone
	password   string
	client     *scw.Client
}

// Cluster returns the cluster, resolved lazily.
func (d *clusterConnectTemplateData) Cluster() (*redisSDK.Cluster, error) {
	return redisSDK.NewAPI(d.client).GetCluster(&redisSDK.GetClusterRequest{
		ClusterID: d.ClusterID,
		Zone:      d.Zone,
	})
}

// Password returns the password of the connect session: the one passed on
// the command line when provided, otherwise the user is prompted for it.
// The connect command always resolves it up front (it is exported to the
// child process as REDISCLI_AUTH); the result is cached so a template
// referencing {{.Password}} does not prompt again.
func (d *clusterConnectTemplateData) Password() (string, error) {
	if d.password != "" {
		return d.password, nil
	}

	password, err := interactive.PromptPasswordWithConfig(
		d.ctx,
		&interactive.PromptPasswordConfig{
			Prompt: "Password",
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to get password: %w", err)
	}

	// Cache the prompted value so a template referencing {{.Password}}
	// more than once does not prompt again.
	d.password = password

	return password, nil
}

func clusterConnectCommand() *core.Command {
	return &core.Command{
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "connect",
		Short:     `Connect to a Redis cluster using locally installed redis-cli, or a custom command via template`,
		Long: `Connect to a Redis cluster using locally installed redis-cli.

The command prompts for the password (unless password= is provided), downloads the cluster TLS certificate to a temporary file (removed when the command exits) when TLS is enabled, and connects to the cluster endpoint with redis-cli's -h/-p options.

The password is not on the command line, where it would be visible in the process list: it is exported to the connect command as the REDISCLI_AUTH environment variable (the standard redis-cli auth variable). Any REDISCLI_AUTH already set in the environment is replaced by it.

Use template to connect with a different tool instead. Its value (default: the redis-cli command below) is rendered as a Go text/template and the result is executed as the connect command.
Available variables:
{{.ClusterID}}, {{.Cluster}} (the cluster object, e.g. {{.Cluster.Name}}), {{.Zone}}, {{.Host}}, {{.Port}}, {{.Password}}, {{.UserName}}, {{.TLSEnabled}}, {{.CACert}} (path to the temporary certificate), {{.Organization}} and {{.Project}}.`,
		ArgsType: reflect.TypeFor[clusterConnectArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:    "private-network",
				Short:   `Connect by the private network endpoint attached.`,
				Default: core.DefaultValueSetter("false"),
			},
			{
				Name:       "cluster-id",
				Short:      `UUID of the cluster`,
				Required:   true,
				Positional: true,
			},
			&core.ArgSpec{
				Name:    "template",
				Short:   `Command line to run instead of the default redis-cli one, as a Go text/template`,
				Default: core.DefaultValueSetter(defaultConnectTemplate),
			},
			{
				Name:  "password",
				Short: `Password to connect with. If not provided, the user is prompted for it. Exported to the connect command as the REDISCLI_AUTH environment variable`,
			},
			core.ZoneArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*clusterConnectArgs)

			client := core.ExtractClient(ctx)
			api := redisSDK.NewAPI(client)
			cluster, err := api.GetCluster(&redisSDK.GetClusterRequest{
				Zone:      args.Zone,
				ClusterID: args.ClusterID,
			})
			if err != nil {
				return nil, err
			}

			if len(cluster.Endpoints) == 0 {
				return nil, fmt.Errorf("%s", errorMessageEndpointNotFound)
			}

			var endpoint *redisSDK.Endpoint
			switch {
			case args.PrivateNetwork:
				endpoint, err = getPrivateEndpoint(cluster.Endpoints)
				if err != nil {
					return nil, err
				}
			default:
				endpoint, err = getPublicEndpoint(cluster.Endpoints)
				if err != nil {
					return nil, err
				}
			}

			if len(endpoint.IPs) == 0 {
				return nil, errors.New("endpoint has no IP addresses")
			}

			var certPath string
			if cluster.TLSEnabled {
				certResp, err := api.GetClusterCertificate(&redisSDK.GetClusterCertificateRequest{
					Zone:      args.Zone,
					ClusterID: args.ClusterID,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to get certificate: %w", err)
				}

				certContent, err := io.ReadAll(certResp.Content)
				if err != nil {
					return nil, fmt.Errorf("failed to read certificate content: %w", err)
				}

				certPath = filepath.Join(
					os.TempDir(),
					fmt.Sprintf("redis-cert-%s.crt", args.ClusterID),
				)
				if err := os.WriteFile(certPath, certContent, 0o600); err != nil {
					return nil, fmt.Errorf("failed to write certificate: %w", err)
				}
				defer func() {
					if err := os.Remove(certPath); err != nil {
						core.ExtractLogger(ctx).Debugf("failed to remove certificate file: %v", err)
					}
				}()
			}

			templateData := &clusterConnectTemplateData{
				TemplateData: connect.NewTemplateData(ctx),
				ctx:          ctx,
				Host:         endpoint.IPs[0].String(),
				Port:         endpoint.Port,
				UserName:     cluster.UserName,
				TLSEnabled:   cluster.TLSEnabled,
				CACert:       certPath,
				ClusterID:    args.ClusterID,
				Zone:         args.Zone,
				password:     args.Password,
				client:       client,
			}

			cmdArgs, err := connect.RenderCommand(args.Template, templateData)
			if err != nil {
				return nil, err
			}

			// The password is exported as REDISCLI_AUTH, the standard redis-cli
			// auth environment variable, instead of being passed on the command
			// line where it would be visible in the process list.
			password, err := templateData.Password()
			if err != nil {
				return nil, err
			}

			// Drop any user REDISCLI_AUTH so the tool only sees this session's password
			env := make([]string, 0, len(os.Environ())+1)
			for _, e := range os.Environ() {
				if !strings.HasPrefix(e, "REDISCLI_AUTH=") {
					env = append(env, e)
				}
			}

			env = append(env, "REDISCLI_AUTH="+password)

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
				Short: "Connect to a Redis cluster",
				Raw:   `scw redis cluster connect 11111111-1111-1111-1111-111111111111`,
			},
			{
				Short: "Connect to a Redis cluster via private network",
				Raw:   `scw redis cluster connect 11111111-1111-1111-1111-111111111111 private-network=true`,
			},
			{
				Short: "Connect to a Redis cluster with a custom command",
				Raw:   `scw redis cluster connect 11111111-1111-1111-1111-111111111111 template="my-redis-tool --host {{.Host}}"`,
			},
		},
		ExcludeFromMCP: true,
	}
}
