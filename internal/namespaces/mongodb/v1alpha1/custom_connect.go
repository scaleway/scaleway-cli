package mongodb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/scaleway/scaleway-cli/v2/core"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	errorMessageMongoCliNotFound        = "mongosh is not installed. Please install mongosh to use this command"
	errorMessagePublicEndpointNotFound  = "public endpoint not found"
	errorMessagePrivateEndpointNotFound = "private endpoint not found"
	errorMessageEndpointNotFound        = "any endpoint is associated on your instance"
)

type instanceConnectArgs struct {
	Region         scw.Region
	PrivateNetwork bool
	InstanceID     string
	Username       string
	Database       *string
	CliMongo       *string
	CliArgs        []string
}

func instanceConnectCommand() *core.Command {
	return &core.Command{
		Namespace: "mongodb",
		Resource:  "instance",
		Verb:      "connect",
		Short:     "Connect to an instance using locally installed mongosh",
		Long:      "Connect to an instance using locally installed mongosh. The command verifies mongosh is installed, downloads the CA certificate for TLS, and lets mongosh prompt for the password.",
		ArgsType:  reflect.TypeFor[instanceConnectArgs](),
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
				Short:    "Name of the user to connect with to the database",
				Required: true,
			},
			{
				Name:    "database",
				Short:   "Name of the database to connect to",
				Default: core.DefaultValueSetter("admin"),
			},
			{
				Name:  "cli-mongo",
				Short: "Command line tool to use, defaults to mongosh",
			},
			{
				Name:     "cli-args",
				Short:    "Additional arguments to pass to mongosh",
				Required: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*instanceConnectArgs)

			cliMongo := "mongosh"
			if args.CliMongo != nil {
				cliMongo = *args.CliMongo
			}

			if _, err := exec.LookPath(cliMongo); err != nil {
				return nil, fmt.Errorf("%s", errorMessageMongoCliNotFound)
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
				return nil, fmt.Errorf("%s", errorMessageEndpointNotFound)
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

			if len(endpoint.IPs) == 0 {
				return nil, errors.New("endpoint has no IP addresses")
			}

			certPath, err := downloadInstanceCertificate(api, args.Region, args.InstanceID)
			if err != nil {
				return nil, err
			}
			defer os.Remove(certPath)

			database := "admin"
			if args.Database != nil {
				database = *args.Database
			}

			hostPort := net.JoinHostPort(
				endpoint.IPs[0].String(),
				strconv.FormatUint(uint64(endpoint.Port), 10),
			)
			connURI := fmt.Sprintf(
				"mongodb://%s@%s/%s?tls=true&tlsCAFile=%s&authSource=admin",
				args.Username,
				hostPort,
				database,
				certPath,
			)

			cmdArgs := append([]string{cliMongo, connURI}, args.CliArgs...)

			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec
			core.ExtractLogger(ctx).Debugf("executing: %s\n", cmd.Args)

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
				Short: "Connect to an instance",
				Raw:   `scw mongodb instance connect 11111111-1111-1111-1111-111111111111 --username=admin`,
			},
			{
				Short: "Connect to an instance via private network",
				Raw:   `scw mongodb instance connect 11111111-1111-1111-1111-111111111111 --username=admin --private-network=true --database=mydb`,
			},
		},
		ExcludeFromMCP: true,
	}
}

func getPublicEndpoint(endpoints []*mongodb.Endpoint) (*mongodb.Endpoint, error) {
	for _, e := range endpoints {
		if e.Public != nil {
			return e, nil
		}
	}

	return nil, fmt.Errorf("%s", errorMessagePublicEndpointNotFound)
}

func getPrivateEndpoint(endpoints []*mongodb.Endpoint) (*mongodb.Endpoint, error) {
	for _, e := range endpoints {
		if e.PrivateNetwork != nil {
			return e, nil
		}
	}

	return nil, fmt.Errorf("%s", errorMessagePrivateEndpointNotFound)
}

// downloadInstanceCertificate fetches the CA certificate of the instance and
// writes it to a temp file so mongosh can validate the TLS connection.
func downloadInstanceCertificate(
	api *mongodb.API,
	region scw.Region,
	instanceID string,
) (string, error) {
	certResp, err := api.GetInstanceCertificate(&mongodb.GetInstanceCertificateRequest{
		Region:     region,
		InstanceID: instanceID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get certificate: %w", err)
	}

	certContent, err := io.ReadAll(certResp.Content)
	if err != nil {
		return "", fmt.Errorf("failed to read certificate content: %w", err)
	}

	certPath := filepath.Join(os.TempDir(), fmt.Sprintf("mongodb-cert-%s.crt", instanceID))
	if err := os.WriteFile(certPath, certContent, 0o600); err != nil {
		return "", fmt.Errorf("failed to write certificate: %w", err)
	}

	return certPath, nil
}
