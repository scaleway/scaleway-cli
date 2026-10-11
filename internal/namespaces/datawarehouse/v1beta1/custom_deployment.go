package datawarehouse

import (
	"context"
	"fmt"
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

type createDeploymentRequestCustom struct {
	*datawarehouse.CreateDeploymentRequest
	GeneratePassword bool
	CliConnect       bool
}

func deploymentCreateBuilder(c *core.Command) *core.Command {
	c.ArgsType = reflect.TypeFor[createDeploymentRequestCustom]()

	c.ArgSpecs.AddBefore("password", passwordgenerator.GeneratePasswordArgSpec())
	c.ArgSpecs.AddBefore("password", connect.CliConnectArgSpec())

	c.SeeAlsos = []*core.SeeAlso{
		{
			Command: "scw datawarehouse deployment connect",
			Short:   "Connect to a Data Warehouse deployment using locally installed clickhousectl",
		},
	}

	c.Run = func(ctx context.Context, args any) (any, error) {
		customRequest := args.(*createDeploymentRequestCustom)
		request := customRequest.CreateDeploymentRequest

		generatedPassword, err := passwordgenerator.GenerateIfNeeded(
			customRequest.GeneratePassword,
			request.Password,
		)
		if err != nil {
			return nil, err
		}
		if generatedPassword != "" {
			request.Password = generatedPassword
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
			if err := connect.PersistCredentials(
				ctx,
				deployment.Region,
				"datawarehouse",
				string(deployment.Region),
				deployment.ID,
				dwAdminUsername,
				request.Password,
				"deployment",
			); err != nil {
				return nil, err
			}
		}

		return deployment, nil
	}

	return c
}

func deploymentGetBuilder(c *core.Command) *core.Command {
	c.SeeAlsos = []*core.SeeAlso{
		{
			Command: "scw datawarehouse deployment connect",
			Short:   "Connect to a Data Warehouse deployment using locally installed clickhousectl",
		},
	}

	return c
}

type deploymentConnectArgs struct {
	Region         scw.Region
	PrivateNetwork bool
	DeploymentID   string
	CliClickhouse  *string
}

const (
	connectErrorMessagePublicEndpointNotFound  = "public endpoint not found"
	connectErrorMessagePrivateEndpointNotFound = "private endpoint not found"
	connectErrorMessageEndpointNotFound        = "no endpoint is associated with your deployment"
	connectErrorMessageTCPServiceNotFound      = "no TCP service found on the endpoint"
	connectErrorMessageCliNotFound             = "clickhousectl is not installed. Please install clickhousectl to use this command"
)

func checkClickhousectlInstalled(cliClickhouse string) error {
	cmd := exec.Command(cliClickhouse, "--version") //nolint:gosec
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s", connectErrorMessageCliNotFound)
	}

	return nil
}

func getPublicEndpoint(endpoints []*datawarehouse.Endpoint) (*datawarehouse.Endpoint, error) {
	for _, e := range endpoints {
		if e.Public != nil {
			return e, nil
		}
	}

	return nil, fmt.Errorf("%s", connectErrorMessagePublicEndpointNotFound)
}

func getPrivateEndpoint(endpoints []*datawarehouse.Endpoint) (*datawarehouse.Endpoint, error) {
	for _, e := range endpoints {
		if e.PrivateNetwork != nil {
			return e, nil
		}
	}

	return nil, fmt.Errorf("%s", connectErrorMessagePrivateEndpointNotFound)
}

func getTCPServicePort(endpoint *datawarehouse.Endpoint) (uint32, error) {
	for _, s := range endpoint.Services {
		if s.Protocol == datawarehouse.EndpointServiceProtocolTCP {
			return s.Port, nil
		}
	}

	return 0, fmt.Errorf("%s", connectErrorMessageTCPServiceNotFound)
}

func deploymentConnectCommand() *core.Command {
	return &core.Command{
		Namespace: "datawarehouse",
		Resource:  "deployment",
		Verb:      "connect",
		Short:     "Connect to a Data Warehouse deployment using locally installed clickhousectl",
		Long: `Connect to a Data Warehouse deployment using locally installed clickhousectl.
The scwadmin user is used and credentials are fetched from the secret created by
"datawarehouse deployment create cli-connect=true" when available.`,
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
			{
				Name:  "cli-clickhouse",
				Short: "Command line tool to use, default to clickhousectl",
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*deploymentConnectArgs)

			cliClickhouse := "clickhousectl"
			if args.CliClickhouse != nil {
				cliClickhouse = *args.CliClickhouse
			}

			if err := checkClickhousectlInstalled(cliClickhouse); err != nil {
				return nil, err
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
				return nil, fmt.Errorf("%s", connectErrorMessageEndpointNotFound)
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
			creds, err := connect.FetchOrPrompt(
				ctx,
				deployment.Region,
				connect.Path(
					"datawarehouse",
					string(deployment.Region),
					deployment.ID,
					dwAdminUsername,
				),
			)
			if err != nil {
				return nil, err
			}
			password := creds.Password

			// Fail fast on an unreachable endpoint instead of letting the client hang on the dial
			if err := core.CheckReachable(endpoint.DNSRecord, port); err != nil {
				return nil, err
			}

			// ponytail: the public endpoint presents a Let's Encrypt chain
			// (CN=*.dtwh.<region>.scw.cloud) verified against the system trust store by
			// --secure alone, so no CA file is needed, as in the Scaleway-recommended command.
			// The GetDeploymentCertificate API returns a different self-signed cert for
			// private-network endpoints; pin it there if a private endpoint ever needs it.
			cmdArgs := []string{
				cliClickhouse,
				"local", "client",
				"--host", endpoint.DNSRecord,
				"--port", strconv.FormatUint(uint64(port), 10),
				"--",
				"--secure",
				"--user", dwAdminUsername,
				"--password", password,
				"--database", "default",
			}

			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec
			core.ExtractLogger(ctx).
				Debugf("executing %s against %s\n", cliClickhouse, endpoint.DNSRecord)

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
		},
		ExcludeFromMCP: true,
	}
}
