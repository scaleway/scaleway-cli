package mongodb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// defaultConnectTemplate renders the default mongosh invocation.
// Credentials are in ConnURL (mongosh has no redis-cli-style auth env var).
const defaultConnectTemplate = `mongosh {{.ConnURL}} --tlsCAFile {{.CertPath}}`

type instanceConnectArgs struct {
	Region         scw.Region
	PrivateNetwork bool
	InstanceID     string
	Username       string
	Password       string
	Template       string
}

func getPublicEndpoint(endpoints []*mongodb.Endpoint) (*mongodb.Endpoint, error) {
	for _, e := range endpoints {
		if e.PublicNetwork != nil {
			return e, nil
		}
	}

	return nil, errors.New("public endpoint not found")
}

func getPrivateEndpoint(endpoints []*mongodb.Endpoint) (*mongodb.Endpoint, error) {
	for _, e := range endpoints {
		if e.PrivateNetwork != nil {
			return e, nil
		}
	}

	return nil, errors.New("private endpoint not found")
}

// pickAdminUser returns the name of the first user holding the db_admin role,
// falling back to the first user in the list.
// ponytail: db_admin is used as the "admin" marker (the initial user created at
// provision time). If an instance has several db_admin users the first one wins;
// upgrade to a dedicated flag/role if that proves ambiguous.
func pickAdminUser(users []*mongodb.User) (string, error) {
	fallback := ""
	for _, user := range users {
		if fallback == "" {
			fallback = user.Name
		}
		for _, role := range user.Roles {
			if role.Role == mongodb.UserRoleRoleDbAdmin {
				return user.Name, nil
			}
		}
	}

	if fallback == "" {
		return "", errors.New("no user found on instance")
	}

	return fallback, nil
}

func findAdminUser(api *mongodb.API, instanceID string, region scw.Region) (string, error) {
	users, err := api.ListUsers(&mongodb.ListUsersRequest{
		Region:     region,
		InstanceID: instanceID,
	}, scw.WithAllPages())
	if err != nil {
		return "", err
	}

	return pickAdminUser(users.Users)
}

// instanceConnectTemplateData is the template data for the connect command.
type instanceConnectTemplateData struct {
	*connect.TemplateData
	InstanceID string
	Region     scw.Region
	Username   string
	Password   string

	privateNetwork bool
	client         *scw.Client

	instance     *mongodb.Instance
	instanceOnce sync.Once
	instanceErr  error

	endpoint     *mongodb.Endpoint
	endpointOnce sync.Once
	endpointErr  error

	certPath string
	certOnce sync.Once
	certErr  error
}

func (d *instanceConnectTemplateData) Instance() (*mongodb.Instance, error) {
	d.instanceOnce.Do(func() {
		d.instance, d.instanceErr = mongodb.NewAPI(d.client).
			GetInstance(&mongodb.GetInstanceRequest{
				Region:     d.Region,
				InstanceID: d.InstanceID,
			})
	})

	return d.instance, d.instanceErr
}

func (d *instanceConnectTemplateData) Endpoint() (*mongodb.Endpoint, error) {
	d.endpointOnce.Do(func() {
		instance, err := d.Instance()
		if err != nil {
			d.endpoint, d.endpointErr = nil, err

			return
		}
		if len(instance.Endpoints) == 0 {
			d.endpoint, d.endpointErr = nil, errors.New(
				"no endpoint is associated with your instance",
			)

			return
		}

		if d.privateNetwork {
			d.endpoint, d.endpointErr = getPrivateEndpoint(instance.Endpoints)
		} else {
			d.endpoint, d.endpointErr = getPublicEndpoint(instance.Endpoints)
		}
	})

	return d.endpoint, d.endpointErr
}

// DNSRecord is the endpoint DNS name used in mongodb+srv URIs.
func (d *instanceConnectTemplateData) DNSRecord() (string, error) {
	endpoint, err := d.Endpoint()
	if err != nil {
		return "", err
	}

	return endpoint.DNSRecord, nil
}

// CertPath returns the path of the temporary TLS certificate file.
func (d *instanceConnectTemplateData) CertPath() (string, error) {
	d.certOnce.Do(func() {
		d.certPath, d.certErr = d.downloadCertificate()
	})

	return d.certPath, d.certErr
}

// ConnURL returns a mongodb+srv URI with username and password.
func (d *instanceConnectTemplateData) ConnURL() (string, error) {
	dnsRecord, err := d.DNSRecord()
	if err != nil {
		return "", err
	}

	return (&url.URL{
		Scheme: "mongodb+srv",
		Host:   dnsRecord,
		User:   url.UserPassword(d.Username, d.Password),
	}).String(), nil
}

func (d *instanceConnectTemplateData) downloadCertificate() (string, error) {
	cert, err := mongodb.NewAPI(d.client).
		GetInstanceCertificate(&mongodb.GetInstanceCertificateRequest{
			Region:     d.Region,
			InstanceID: d.InstanceID,
		})
	if err != nil {
		return "", fmt.Errorf("failed to get certificate: %w", err)
	}

	certContent, err := io.ReadAll(cert.Content)
	if err != nil {
		return "", fmt.Errorf("failed to read certificate content: %w", err)
	}

	certPath := filepath.Join(os.TempDir(), fmt.Sprintf("mongodb-cert-%s.crt", d.InstanceID))
	if err := os.WriteFile(certPath, certContent, 0o600); err != nil {
		return "", fmt.Errorf("failed to write certificate: %w", err)
	}

	return certPath, nil
}

func instanceConnectCommand() *core.Command {
	return &core.Command{
		Namespace: "mongodb",
		Resource:  "instance",
		Verb:      "connect",
		Short:     `Connect to a MongoDB® Database Instance using locally installed mongosh, or a custom command via template`,
		Long: `Connect to a MongoDB® Database Instance using locally installed mongosh.

When no username is given, the first db_admin user is used (falling back to the first user). The password is prompted unless password= is provided. The instance TLS certificate is downloaded to a temporary file (removed when the command exits).

Use template to connect with a different tool instead. Its value (default: the mongosh command below) is rendered as a Go text/template and the result is executed as the connect command.
Available variables:
{{.InstanceID}}, {{.Instance}} (the instance object), {{.Endpoint}} ({{.Endpoint.DNSRecord}}), {{.DNSRecord}}, {{.Region}}, {{.Username}}, {{.Password}}, {{.ConnURL}}, {{.CertPath}}, {{.Organization}} and {{.Project}}.`,
		ArgsType: reflect.TypeFor[instanceConnectArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:    "private-network",
				Short:   `Connect by the private network endpoint attached.`,
				Default: core.DefaultValueSetter("false"),
			},
			{
				Name:       "instance-id",
				Short:      `UUID of the instance`,
				Required:   true,
				Positional: true,
			},
			{
				Name:  "username",
				Short: `Name of the user to connect with (defaults to the admin user)`,
			},
			{
				Name:  "password",
				Short: `Password to connect with. If not provided, the user is prompted for it`,
			},
			&core.ArgSpec{
				Name:    "template",
				Short:   `Command line to run instead of the default mongosh one, as a Go text/template`,
				Default: core.DefaultValueSetter(defaultConnectTemplate),
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*instanceConnectArgs)
			client := core.ExtractClient(ctx)
			api := mongodb.NewAPI(client)

			username := args.Username
			if username == "" {
				var err error
				username, err = findAdminUser(api, args.InstanceID, args.Region)
				if err != nil {
					return nil, err
				}
			}

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

			templateData := &instanceConnectTemplateData{
				TemplateData:   connect.NewTemplateData(ctx),
				InstanceID:     args.InstanceID,
				Region:         args.Region,
				Username:       username,
				Password:       password,
				privateNetwork: args.PrivateNetwork,
				client:         client,
			}
			defer func() {
				if templateData.certPath == "" {
					return
				}
				if err := os.Remove(templateData.certPath); err != nil {
					core.ExtractLogger(ctx).Debugf("failed to remove certificate file: %v", err)
				}
			}()

			// Fail fast on the default mongosh path when the +srv record is
			// not resolvable, instead of letting mongosh time out. Skip when
			// OverrideExec is set (cassette tests): DNS is not recorded and
			// the instance may already be gone at replay time.
			meta := core.ExtractMeta(ctx)
			if args.Template == defaultConnectTemplate &&
				(meta == nil || meta.OverrideExec == nil) {
				dnsRecord, err := templateData.DNSRecord()
				if err != nil {
					return nil, err
				}
				if _, _, srvErr := net.DefaultResolver.LookupSRV(
					ctx,
					"mongodb",
					"tcp",
					dnsRecord,
				); srvErr != nil {
					return nil, fmt.Errorf(
						"cannot resolve the MongoDB endpoint %s from this machine (%v); "+
							"check your network/DNS, or use private-network=true if a private endpoint is attached",
						dnsRecord,
						srvErr,
					)
				}
			}

			cmdArgs, err := connect.RenderCommand(args.Template, templateData)
			if err != nil {
				return nil, err
			}

			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec
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
				Short: "Connect to a MongoDB instance",
				Raw:   `scw mongodb instance connect 11111111-1111-1111-1111-111111111111`,
			},
			{
				Short: "Connect via private network",
				Raw:   `scw mongodb instance connect 11111111-1111-1111-1111-111111111111 private-network=true`,
			},
			{
				Short: "Connect with a custom command",
				Raw:   `scw mongodb instance connect 11111111-1111-1111-1111-111111111111 template="my-mongo-tool --uri {{.ConnURL}}"`,
			},
		},
		ExcludeFromMCP: true,
	}
}
