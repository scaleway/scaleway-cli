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
	"sync"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	redisSDK "github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// defaultConnectTemplate renders the default redis-cli invocation: the
// endpoint, the TLS certificate (when TLS is enabled) and the cluster
// user. The password is not on the command line, where it would be visible
// in the process list: redis-cli picks it up from the REDISCLI_AUTH
// environment variable the connect command exports.
const defaultConnectTemplate = `redis-cli -h {{.Endpoint.Host}} -p {{.Endpoint.Port}}{{if .Cluster.TLSEnabled}} --tls --cacert {{.CACert}}{{end}} --user {{.Cluster.UserName}}`

type clusterConnectArgs struct {
	ClusterID      string
	Zone           scw.Zone
	Template       string
	Password       string
	PrivateNetwork bool
	VpcGatewayID   string
}

// endpoint is the connect address of a session, resolved lazily from the
// cluster's public endpoint (or the private one when private-network=true).
type endpoint struct {
	Host string
	Port uint32
}

// clusterConnectTemplateData is the template data for the connect command:
// the base connect values plus the redis cluster values.
type clusterConnectTemplateData struct {
	*connect.TemplateData
	ClusterID string
	Zone      scw.Zone
	// Password is the session password, resolved before rendering so that
	// templates can reference {{.Password}}.
	Password string

	cluster  *redisSDK.Cluster
	endpoint *endpoint
	certPath string
	// privateNetwork selects the private endpoint instead of the public one.
	privateNetwork bool

	client       *scw.Client
	clusterOnce  sync.Once
	clusterErr   error
	endpointOnce sync.Once
	endpointErr  error
	certOnce     sync.Once
	certErr      error

	// stderr receives the output of the ssh tunnel process.
	stderr io.Writer

	// vpcGatewayID, when set, routes the private endpoint through the SSH
	// bastion of the given VPC gateway. When empty, the gateway is
	// auto-detected from the gateways attached to the private network.
	vpcGatewayID     string
	privateNetworkID string
	gateway          *vpcgw.Gateway
	gatewayOnce      sync.Once
	gatewayErr       error
	tunnel           *connect.Tunnel
}

// Cluster returns the cluster, resolved lazily.
func (d *clusterConnectTemplateData) Cluster() (*redisSDK.Cluster, error) {
	d.clusterOnce.Do(func() {
		d.cluster, d.clusterErr = redisSDK.NewAPI(d.client).GetCluster(&redisSDK.GetClusterRequest{
			ClusterID: d.ClusterID,
			Zone:      d.Zone,
		})
	})

	return d.cluster, d.clusterErr
}

// Endpoint returns the connect address, resolved lazily from the cluster:
// the first IP and the port of the public endpoint, or of the attached
// private network endpoint when private-network=true. When a VPC gateway is
// selected, the private endpoint is reached through the gateway SSH bastion
// and the returned address is the local loopback port of the tunnel.
func (d *clusterConnectTemplateData) Endpoint() (*endpoint, error) {
	d.endpointOnce.Do(func() {
		cluster, err := d.Cluster()
		if err != nil {
			d.endpoint, d.endpointErr = nil, err

			return
		}

		var selected *redisSDK.Endpoint
		if d.privateNetwork {
			selected, err = getPrivateEndpoint(cluster.Endpoints)
		} else {
			selected, err = getPublicEndpoint(cluster.Endpoints)
			if err != nil {
				// No public endpoint: the private one is the only way in,
				// it will be reached through a VPC gateway.
				selected, err = getPrivateEndpoint(cluster.Endpoints)
			}
		}
		if err != nil {
			d.endpoint, d.endpointErr = nil, err

			return
		}
		if len(selected.IPs) == 0 {
			d.endpoint, d.endpointErr = nil, errors.New("endpoint has no IP addresses")
		}

		ep := &endpoint{
			Host: selected.IPs[0].String(),
			Port: selected.Port,
		}

		// A private endpoint is only reachable from inside the private
		// network: route it through a VPC gateway SSH bastion. An ssh
		// port-forward opens a local port and the connect command targets
		// it instead of the private IP.
		if selected.PrivateNetwork != nil {
			d.privateNetworkID = selected.PrivateNetwork.ID

			gw, err := d.Gateway()
			if err != nil {
				d.endpoint, d.endpointErr = nil, err

				return
			}
			tunnel, err := connect.SSHLocalForward(
				d.stderr,
				gw.IPv4.Address.String(),
				gw.BastionPort,
				ep.Host,
				ep.Port,
			)
			if err != nil {
				d.endpoint, d.endpointErr = nil, err

				return
			}
			d.tunnel = tunnel
			ep = &endpoint{
				Host: "127.0.0.1",
				Port: uint32(tunnel.LocalPort()),
			}
		}

		d.endpoint, d.endpointErr = ep, nil

		return
	})

	return d.endpoint, d.endpointErr
}

// Gateway returns the VPC gateway whose SSH bastion is used to reach the
// private endpoint, resolved lazily. When vpc-gateway-id is set the gateway
// is fetched and validated; otherwise it is auto-detected from the gateways
// attached to the private network of the selected endpoint.
func (d *clusterConnectTemplateData) Gateway() (*vpcgw.Gateway, error) {
	d.gatewayOnce.Do(func() {
		if d.vpcGatewayID == "" {
			if d.privateNetworkID == "" {
				d.gatewayErr = errors.New(
					"no private endpoint selected, cannot auto-detect the VPC gateway",
				)

				return
			}

			d.gateway, d.gatewayErr = connect.PickGatewayForPrivateNetwork(
				d.client,
				d.Zone,
				d.privateNetworkID,
			)

			return
		}

		gw, err := vpcgw.NewAPI(d.client).GetGateway(&vpcgw.GetGatewayRequest{
			GatewayID: d.vpcGatewayID,
			Zone:      d.Zone,
		})
		if err != nil {
			d.gatewayErr = err

			return
		}
		switch {
		case !gw.BastionEnabled:
			d.gatewayErr = fmt.Errorf("the SSH bastion is not enabled on gateway %s", gw.ID)
		case gw.IPv4 == nil:
			d.gatewayErr = errors.New("the gateway has no public IPv4 address")
		default:
			d.gateway = gw
		}
	})

	return d.gateway, d.gatewayErr
}

// CACert returns the path of the temporary file holding the cluster TLS
// certificate, downloaded on first use. It is empty when TLS is disabled;
// the connect command removes the file when it exits.
func (d *clusterConnectTemplateData) CACert() (string, error) {
	d.certOnce.Do(func() {
		d.certPath, d.certErr = d.downloadCertificate()
	})

	return d.certPath, d.certErr
}

func (d *clusterConnectTemplateData) downloadCertificate() (string, error) {
	cluster, err := d.Cluster()
	if err != nil {
		return "", err
	}
	if !cluster.TLSEnabled {
		return "", nil
	}

	certResp, err := redisSDK.NewAPI(d.client).GetClusterCertificate(
		&redisSDK.GetClusterCertificateRequest{
			Zone:      d.Zone,
			ClusterID: d.ClusterID,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to get certificate: %w", err)
	}

	certContent, err := io.ReadAll(certResp.Content)
	if err != nil {
		return "", fmt.Errorf("failed to read certificate content: %w", err)
	}

	certPath := filepath.Join(os.TempDir(), fmt.Sprintf("redis-cert-%s.crt", d.ClusterID))
	if err := os.WriteFile(certPath, certContent, 0o600); err != nil {
		return "", fmt.Errorf("failed to write certificate: %w", err)
	}

	return certPath, nil
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

A private endpoint is only reachable from inside the private network: the connection is routed through the SSH bastion of a VPC gateway attached to the same private network (or advertising its route) with the bastion enabled. The command starts an SSH port-forward through the gateway bastion, connects to the local port it opens, and stops the tunnel when it exits. The local ssh configuration (keys, ~/.ssh/config) is used to authenticate to the bastion. The gateway is auto-detected from the cluster's private network; when several candidates exist, one must be given explicitly with vpc-gateway-id. When the cluster has no public endpoint, the private one is used.

Use template to connect with a different tool instead. Its value (default: the redis-cli command below) is rendered as a Go text/template and the result is executed as the connect command.
Available variables:
{{.ClusterID}}, {{.Cluster}} (the cluster object, e.g. {{.Cluster.Name}}), {{.Endpoint}} (the resolved connect address: {{.Endpoint.Host}}, {{.Endpoint.Port}}), {{.Gateway}} (the VPC gateway, when vpc-gateway-id is set), {{.Zone}}, {{.Password}}, {{.CACert}} (path to the temporary certificate), {{.Organization}} and {{.Project}}.`,
		ArgsType: reflect.TypeFor[clusterConnectArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:    "private-network",
				Short:   `Connect by the private network endpoint attached.`,
				Default: core.DefaultValueSetter("false"),
			},
			{
				Name:  "vpc-gateway-id",
				Short: `UUID of the VPC gateway whose SSH bastion is used to reach the private endpoint. When not set, the gateway is auto-detected from the cluster's private network. Implies private-network=true`,
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

			// A VPC gateway only reaches private addresses: selecting one
			// means connecting through the private endpoint.
			if args.VpcGatewayID != "" {
				args.PrivateNetwork = true
			}

			// The password is exported to the child process as REDISCLI_AUTH,
			// the standard redis-cli auth environment variable, instead of
			// being passed on the command line where it would be visible in
			// the process list. It is resolved before rendering so templates
			// can reference {{.Password}}.
			password := args.Password
			if password == "" {
				var err error
				password, err = interactive.PromptPasswordWithConfig(
					ctx,
					&interactive.PromptPasswordConfig{
						Prompt: "Password",
					},
				)
				if err != nil {
					return nil, fmt.Errorf("failed to get password: %w", err)
				}
			}

			// The cluster, the endpoint and the certificate are resolved
			// lazily: the API is only called for the values the template
			// actually references (like Organization and Project).
			templateData := &clusterConnectTemplateData{
				TemplateData:   connect.NewTemplateData(ctx),
				ClusterID:      args.ClusterID,
				Zone:           args.Zone,
				privateNetwork: args.PrivateNetwork,
				vpcGatewayID:   args.VpcGatewayID,
				Password:       password,
				client:         core.ExtractClient(ctx),
				stderr:         core.ExtractStderr(ctx),
			}
			// Remove the temporary certificate, downloaded while rendering
			// the template, when the command exits.
			defer func() {
				if templateData.certPath == "" {
					return
				}
				if err := os.Remove(templateData.certPath); err != nil {
					core.ExtractLogger(ctx).Debugf("failed to remove certificate file: %v", err)
				}
			}()
			// Stop the ssh tunnel, opened while resolving the endpoint, when
			// the command exits.
			defer func() {
				if templateData.tunnel == nil {
					return
				}
				if err := templateData.tunnel.Stop(); err != nil {
					core.ExtractLogger(ctx).Debugf("failed to stop the ssh tunnel: %v", err)
				}
			}()

			cmdArgs, err := connect.RenderCommand(args.Template, templateData)
			if err != nil {
				return nil, err
			}

			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec
			// Drop any user REDISCLI_AUTH so the tool only sees this session's password.
			connect.SetEnv(cmd, "REDISCLI_AUTH", password)

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
				Short: "Connect to a Redis cluster private endpoint through a VPC gateway SSH bastion",
				Raw:   `scw redis cluster connect 11111111-1111-1111-1111-111111111111 vpc-gateway-id=22222222-2222-2222-2222-222222222222`,
			},
			{
				Short: "Connect to a Redis cluster with a custom command",
				Raw:   `scw redis cluster connect 11111111-1111-1111-1111-111111111111 template="my-redis-tool --host {{.Endpoint.Host}}"`,
			},
		},
		ExcludeFromMCP: true,
	}
}
