package mongodb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"sync"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// defaultConnectTemplate renders the default mongosh invocation.
// ConnURL has no password (mongosh prompts when --username is set without one).
// CertPath is quoted so paths with spaces survive shlex splitting.
const defaultConnectTemplate = `mongosh {{.ConnURL}} --username {{.Username}} --tlsCAFile "{{.CertPath}}"`

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
	instance, err := d.Instance()
	if err != nil {
		return nil, err
	}

	if d.privateNetwork {
		return getPrivateEndpoint(instance.Endpoints)
	}

	return getPublicEndpoint(instance.Endpoints)
}

// CertPath returns the path of the temporary TLS certificate file.
func (d *instanceConnectTemplateData) CertPath() (string, error) {
	d.certOnce.Do(func() {
		d.certPath, d.certErr = d.downloadCertificate()
	})

	return d.certPath, d.certErr
}

// ConnURL returns a mongodb+srv URI. Credentials are included only when
// Password is set (password=), so the default interactive path never puts a
// password in the child process arguments — mongosh prompts instead.
func (d *instanceConnectTemplateData) ConnURL() (string, error) {
	endpoint, err := d.Endpoint()
	if err != nil {
		return "", err
	}

	u := &url.URL{
		Scheme: "mongodb+srv",
		Host:   endpoint.DNSRecord,
	}
	if d.Password != "" {
		u.User = url.UserPassword(d.Username, d.Password)
	}

	return u.String(), nil
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

	f, err := os.CreateTemp("", "mongodb-cert-*.crt")
	if err != nil {
		return "", fmt.Errorf("failed to create certificate file: %w", err)
	}
	certPath := f.Name()

	if _, err := f.Write(certContent); err != nil {
		_ = f.Close()
		_ = os.Remove(certPath)

		return "", fmt.Errorf("failed to write certificate: %w", err)
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(certPath)

		return "", fmt.Errorf("failed to close certificate file: %w", err)
	}
	if err := os.Chmod(certPath, 0o600); err != nil {
		_ = os.Remove(certPath)

		return "", fmt.Errorf("failed to set certificate permissions: %w", err)
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

When no username is given, the first db_admin user is used (falling back to the first user). With the default template, mongosh prompts for the password unless password= is provided (in which case it is placed in ConnURL). The instance TLS certificate is downloaded to a unique temporary file (removed when the command exits).

Use template to connect with a different tool instead. Its value (default: the mongosh command below) is rendered as a Go text/template and the result is executed as the connect command.
Available variables:
{{.InstanceID}}, {{.Instance}} (the instance object), {{.Endpoint}} ({{.Endpoint.DNSRecord}}), {{.Region}}, {{.Username}}, {{.Password}}, {{.ConnURL}}, {{.CertPath}}, {{.Organization}} and {{.Project}}.`,
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
				Short: `Name of the user to connect with (defaults to the first db_admin, else the first user)`,
			},
			{
				Name:  "password",
				Short: `Password for ConnURL. If omitted, the default template lets mongosh prompt for it`,
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
				users, err := api.ListUsers(&mongodb.ListUsersRequest{
					Region:     args.Region,
					InstanceID: args.InstanceID,
				}, scw.WithAllPages())
				if err != nil {
					return nil, err
				}
				username, err = pickAdminUser(users.Users)
				if err != nil {
					return nil, err
				}
			}

			templateData := &instanceConnectTemplateData{
				TemplateData:   connect.NewTemplateData(ctx),
				InstanceID:     args.InstanceID,
				Region:         args.Region,
				Username:       username,
				Password:       args.Password,
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
