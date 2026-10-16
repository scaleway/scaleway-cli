package datawarehouse

import (
	"context"
	"errors"
	"os/exec"
	"reflect"
	"strconv"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	datawarehouse "github.com/scaleway/scaleway-sdk-go/api/datawarehouse/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// dwAdminUsername is the fixed admin user of a Data Warehouse deployment.
const dwAdminUsername = "scwadmin"

var connectSeeAlsos = []*core.SeeAlso{
	{
		Command: "scw datawarehouse deployment connect",
		Short:   "Connect to a Data Warehouse deployment using locally installed clickhousectl",
	},
}

type createDeploymentRequestCustom struct {
	*datawarehouse.CreateDeploymentRequest
	CliConnect bool
}

func deploymentCreateBuilder(c *core.Command) *core.Command {
	c.ArgsType = reflect.TypeFor[createDeploymentRequestCustom]()

	c.ArgSpecs.AddBefore("password", connect.CliConnectArgSpec())

	c.SeeAlsos = []*core.SeeAlso{
		{
			Command: "scw datawarehouse deployment connect",
			Short:   "Connect to a Data Warehouse deployment using locally installed clickhousectl",
		},
		{
			Command: "scw datawarehouse user update",
			Short:   "Change the password of a user of the deployment",
		},
	}

	c.Run = func(ctx context.Context, args any) (any, error) {
		customRequest := args.(*createDeploymentRequestCustom)
		request := customRequest.CreateDeploymentRequest

		generatedPassword, err := passwordgenerator.GenerateIfNeeded(&request.Password)
		if err != nil {
			return nil, err
		}

		client := core.ExtractClient(ctx)
		api := datawarehouse.NewAPI(client)

		deployment, err := api.CreateDeployment(request, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		if generatedPassword != "" {
			interactive.Printf("Your generated password is %s \n\n", generatedPassword)
		}

		if customRequest.CliConnect {
			res := connect.Resource{
				Product:  connect.Datawarehouse,
				Location: deployment.Region.String(),
				ID:       deployment.ID,
				Username: dwAdminUsername,
			}
			if err := res.SaveSecret(
				ctx,
				deployment.Region,
				dwAdminUsername,
				request.Password,
			); err != nil {
				return nil, err
			}
		}

		return deployment, nil
	}

	return c
}

type deleteDeploymentRequestCustom struct {
	*datawarehouse.DeleteDeploymentRequest
	CliConnect bool
}

func deploymentDeleteBuilder(c *core.Command) *core.Command {
	c.ArgsType = reflect.TypeFor[deleteDeploymentRequestCustom]()
	c.ArgSpecs.AddBefore("region", connect.CliConnectDeleteArgSpec())

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		args := argsI.(*deleteDeploymentRequestCustom)

		client := core.ExtractClient(ctx)
		api := datawarehouse.NewAPI(client)
		deployment, err := api.DeleteDeployment(args.DeleteDeploymentRequest, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		if args.CliConnect {
			resource := connect.Resource{
				Product:  connect.Datawarehouse,
				Location: args.Region.String(),
				ID:       args.DeploymentID,
			}
			if err := resource.DeleteSecret(ctx, args.Region); err != nil {
				return deployment, err
			}
		}

		return deployment, nil
	}

	return c
}

type deploymentConnectArgs struct {
	Region         scw.Region
	PrivateNetwork bool
	DeploymentID   string
	Template       string
}

func checkClickhousectlInstalled() error {
	cmd := exec.Command("clickhousectl", "--version") //nolint:gosec
	if err := cmd.Run(); err != nil {
		return errors.New(
			"clickhousectl is not installed. Please install clickhousectl to use this command",
		)
	}

	return nil
}

func getPublicEndpoint(endpoints []*datawarehouse.Endpoint) (*datawarehouse.Endpoint, error) {
	for _, e := range endpoints {
		if e.Public != nil {
			return e, nil
		}
	}

	return nil, errors.New("public endpoint not found")
}

func getPrivateEndpoint(endpoints []*datawarehouse.Endpoint) (*datawarehouse.Endpoint, error) {
	for _, e := range endpoints {
		if e.PrivateNetwork != nil {
			return e, nil
		}
	}

	return nil, errors.New("private endpoint not found")
}

func getTCPServicePort(endpoint *datawarehouse.Endpoint) (uint32, error) {
	for _, s := range endpoint.Services {
		if s.Protocol == datawarehouse.EndpointServiceProtocolTCP {
			return s.Port, nil
		}
	}

	return 0, errors.New("no TCP service found on the endpoint")
}

func deploymentConnectCommand() *core.Command {
	return &core.Command{
		Namespace: "datawarehouse",
		Resource:  "deployment",
		Verb:      "connect",
		Short:     "Connect to a Data Warehouse deployment using locally installed clickhousectl",
		Long: `Connect to a Data Warehouse deployment using locally installed clickhousectl.
The scwadmin user is used and credentials are fetched from the secret created by
"datawarehouse deployment create cli-connect=true" when available.

Use --template to run an arbitrary tool: the Go text/template is rendered with
{{.Username}}, {{.Password}} (from the persisted secret), {{.Host}}, {{.Port}} and {{.Database}},
and the result is executed as the connect command.`,
		ArgsType: reflect.TypeFor[deploymentConnectArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "private-network",
				Short:    `Connect by the private network endpoint attached.`,
				Required: false,
				Default:  core.DefaultValueSetter("false"),
			},
			{
				Name:       "deployment-id",
				Short:      `UUID of the deployment`,
				Required:   true,
				Positional: true,
			},
			connect.TemplateArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*deploymentConnectArgs)

			if args.Template == "" {
				if err := checkClickhousectlInstalled(); err != nil {
					return nil, err
				}
			}

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)

			deployment, err := api.GetDeployment(&datawarehouse.GetDeploymentRequest{
				Region:       args.Region,
				DeploymentID: args.DeploymentID,
			})
			if err != nil {
				return nil, err
			}

			if len(deployment.Endpoints) == 0 {
				return nil, errors.New("no endpoint is associated with your deployment")
			}

			var endpoint *datawarehouse.Endpoint
			switch {
			case args.PrivateNetwork:
				endpoint, err = getPrivateEndpoint(deployment.Endpoints)
			default:
				endpoint, err = getPublicEndpoint(deployment.Endpoints)
			}
			if err != nil {
				return nil, err
			}

			port, err := getTCPServicePort(endpoint)
			if err != nil {
				return nil, err
			}

			// The Data Warehouse admin user is always "scwadmin"
			password, err := connect.Resource{
				Product:  connect.Datawarehouse,
				Location: deployment.Region.String(),
				ID:       deployment.ID,
				Username: dwAdminUsername,
			}.FetchOrPrompt(ctx, deployment.Region, dwAdminUsername)
			if err != nil {
				return nil, err
			}

			// Fail fast on an unreachable endpoint instead of letting the client hang on the dial
			if args.Template == "" {
				if err := core.CheckReachable(endpoint.DNSRecord, port); err != nil {
					return nil, err
				}
			}

			// ponytail: the public endpoint presents a Let's Encrypt chain
			// (CN=*.dtwh.<region>.scw.cloud) verified against the system trust store by
			// --secure alone, so no CA file is needed, as in the Scaleway-recommended command.
			// The GetDeploymentCertificate API returns a different self-signed cert for
			// private-network endpoints; pin it there if a private endpoint ever needs it.
			cmdArgs := []string{
				"clickhousectl",
				"local", "client",
				"--host", endpoint.DNSRecord,
				"--port", strconv.FormatUint(uint64(port), 10),
				"--",
				"--secure",
				"--user", dwAdminUsername,
				"--password", password,
				"--database", "default",
			}

			templateArgs, err := connect.RenderCommand(args.Template, map[string]any{
				"Username": dwAdminUsername,
				"Password": password,
				"Host":     endpoint.DNSRecord,
				"Port":     strconv.FormatUint(uint64(port), 10),
				"Database": "default",
			})
			if err != nil {
				return nil, err
			}
			if templateArgs != nil {
				cmdArgs = templateArgs
			}

			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec
			core.ExtractLogger(ctx).
				Debugf("executing %s against %s\n", cmdArgs[0], endpoint.DNSRecord)

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
				Short: "Connect to a Data Warehouse deployment",
				Raw:   `scw datawarehouse deployment connect 11111111-1111-1111-1111-111111111111`,
			},
			{
				Short: "Connect with a custom tool using a template",
				Raw:   `scw datawarehouse deployment connect 11111111-1111-1111-1111-111111111111 --template "my-ch-tool --host {{.Host}} --port {{.Port}}"`,
			},
		},
		ExcludeFromMCP: true,
	}
}
