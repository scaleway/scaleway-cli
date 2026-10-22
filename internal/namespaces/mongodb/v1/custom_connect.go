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

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type instanceConnectArgs struct {
	Region         scw.Region
	PrivateNetwork bool
	InstanceID     string
	Username       string
	Template       string
}

func checkMongoshInstalled() error {
	cmd := exec.Command("mongosh", "--version") //nolint:gosec
	if err := cmd.Run(); err != nil {
		return errors.New("mongosh is not installed. Please install mongosh to use this command")
	}

	return nil
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

// findAdminUser returns the name of the first user holding the db_admin role,
// falling back to the first user in the list.
// ponytail: db_admin is used as the "admin" marker (the initial user created at
// provision time). If an instance has several db_admin users the first one wins;
// upgrade to a dedicated flag/role if that proves ambiguous.
func findAdminUser(api *mongodb.API, instanceID string, region scw.Region) (string, error) {
	users, err := api.ListUsers(&mongodb.ListUsersRequest{
		Region:     region,
		InstanceID: instanceID,
	}, scw.WithAllPages())
	if err != nil {
		return "", err
	}

	fallback := ""
	for _, user := range users.Users {
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
		return "", fmt.Errorf("no user found on instance %s", instanceID)
	}

	return fallback, nil
}

func instanceConnectCommand() *core.Command {
	return &core.Command{
		Namespace: "mongodb",
		Resource:  "instance",
		Verb:      "connect",
		Short:     "Connect to a MongoDB® Database Instance using locally installed mongosh",
		Long: `Connect to a MongoDB® Database Instance using locally installed mongosh.
When no username is given, the admin user is used. The password is prompted.

Use --template to run an arbitrary tool: the Go text/template is rendered with
{{.Username}}, {{.Password}}, {{.ConnURL}}, {{.DNSRecord}} and
{{.CertPath}}, and the result is executed as the connect command.`,
		ArgsType: reflect.TypeFor[instanceConnectArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "private-network",
				Short:    `Connect by the private network endpoint attached.`,
				Required: false,
				Default:  core.DefaultValueSetter("false"),
			},
			{
				Name:       "instance-id",
				Short:      `UUID of the instance`,
				Required:   true,
				Positional: true,
			},
			{
				Name:     "username",
				Short:    "Name of the user to connect with to the database (defaults to the admin user)",
				Required: false,
			},
			connect.TemplateArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*instanceConnectArgs)

			if args.Template == "" {
				if err := checkMongoshInstalled(); err != nil {
					return nil, err
				}
			}

			client := core.ExtractClient(ctx)
			api := mongodb.NewAPI(client)

			instance, err := api.GetInstance(&mongodb.GetInstanceRequest{
				Region:     args.Region,
				InstanceID: args.InstanceID,
			})
			if err != nil {
				return nil, err
			}

			if len(instance.Endpoints) == 0 {
				return nil, errors.New("no endpoint is associated with your instance")
			}

			var endpoint *mongodb.Endpoint
			switch {
			case args.PrivateNetwork:
				endpoint, err = getPrivateEndpoint(instance.Endpoints)
			default:
				endpoint, err = getPublicEndpoint(instance.Endpoints)
			}
			if err != nil {
				return nil, err
			}

			username := args.Username
			if username == "" {
				username, err = findAdminUser(api, instance.ID, instance.Region)
				if err != nil {
					return nil, err
				}
			}

			password, err := interactive.PromptPassword(ctx, "Password")
			if err != nil {
				return nil, fmt.Errorf("failed to get password: %w", err)
			}

			// TLS is mandatory: download the instance certificate and point mongosh at it
			cert, err := api.GetInstanceCertificate(&mongodb.GetInstanceCertificateRequest{
				Region:     args.Region,
				InstanceID: instance.ID,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to get certificate: %w", err)
			}

			certContent, err := io.ReadAll(cert.Content)
			if err != nil {
				return nil, fmt.Errorf("failed to read certificate content: %w", err)
			}

			certPath := filepath.Join(os.TempDir(), fmt.Sprintf("mongodb-cert-%s.crt", instance.ID))
			if err := os.WriteFile(certPath, certContent, 0o600); err != nil {
				return nil, fmt.Errorf("failed to write certificate: %w", err)
			}
			defer func() {
				if err := os.Remove(certPath); err != nil {
					core.ExtractLogger(ctx).Debugf("failed to remove certificate file: %v", err)
				}
			}()

			// Fail fast when the +srv record is not (yet) resolvable, instead of letting
			// mongosh time out on server selection.
			if args.Template == "" {
				if _, _, srvErr := net.DefaultResolver.LookupSRV(
					ctx,
					"mongodb",
					"tcp",
					endpoint.DNSRecord,
				); srvErr != nil {
					return nil, fmt.Errorf(
						"cannot resolve the MongoDB endpoint %s from this machine (%v); "+
							"check your network/DNS, or use --private-network if a private endpoint is attached",
						endpoint.DNSRecord,
						srvErr,
					)
				}
			}

			// mongodb+srv://{user}:{pass}@{dns-record}/ (credentials escaped by url.UserPassword)
			connURL := &url.URL{
				Scheme: "mongodb+srv",
				Host:   endpoint.DNSRecord,
				User:   url.UserPassword(username, password),
			}

			cmdArgs, err := connect.RenderCommand(args.Template, map[string]any{
				"Username":  username,
				"Password":  password,
				"ConnURL":   connURL.String(),
				"DNSRecord": endpoint.DNSRecord,
				"CertPath":  certPath,
			})
			if err != nil {
				return nil, err
			}
			if cmdArgs == nil {
				cmdArgs = []string{
					"mongosh",
					connURL.String(),
					"--tlsCAFile", certPath,
				}
			}

			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec
			core.ExtractLogger(ctx).Debugf("executing %s\n", cmdArgs[0])

			exitCode, err := core.ExecCmd(ctx, cmd)
			if err != nil {
				return nil, err
			}
			if exitCode != 0 {
				return nil, &core.CliError{Empty: true, Code: exitCode}
			}

			return &core.SuccessResult{
				Empty: true, // the program will output the success message
			}, nil
		},
		Examples: []*core.Example{
			{
				Short: "Connect to a MongoDB instance",
				Raw:   `scw mongodb instance connect 11111111-1111-1111-1111-111111111111`,
			},
			{
				Short: "Connect with a custom tool using a template",
				Raw:   `scw mongodb instance connect 11111111-1111-1111-1111-111111111111 --template "my-mongo-tool --auth {{.Username}}:{{.Password}} {{.DNSRecord}}"`,
			},
		},
		ExcludeFromMCP: true,
	}
}
