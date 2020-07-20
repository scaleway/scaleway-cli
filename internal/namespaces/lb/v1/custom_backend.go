package lb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	backendServerStatsHealthCheckStatusMarshalSpecs = human.EnumMarshalSpecs{
		lb.BackendServerStatsHealthCheckStatusPassed:   &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "passed"},
		lb.BackendServerStatsHealthCheckStatusFailed:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "failed"},
		lb.BackendServerStatsHealthCheckStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
		lb.BackendServerStatsHealthCheckStatusNeutral:  &human.EnumMarshalSpec{Attribute: color.Faint, Value: "neutral"},
		lb.BackendServerStatsHealthCheckStatusCondpass: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "condition passed"},
	}
	backendServerStatsServerStateMarshalSpecs = human.EnumMarshalSpecs{
		lb.BackendServerStatsServerStateStopped:  &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "stopped"},
		lb.BackendServerStatsServerStateStarting: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "starting"},
		lb.BackendServerStatsServerStateRunning:  &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "running"},
		lb.BackendServerStatsServerStateStopping: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "stopping"},
	}
)

func importInstanceCommand() *core.Command {
	type importInstanceArgs struct {
		InstanceID   string
		LBID         string
		Protocol     *lb.Protocol
		Port         int32
		InstanceZone scw.Zone
		UsePublicIP  bool
		Region       scw.Region
	}

	return &core.Command{
		Short:     `Import an instance as a load balancer backend`,
		Long:      `Import an instance as a load balancer backend.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "import-instance",
		ArgsType:  reflect.TypeOf(importInstanceArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "instance-id",
				Short:    `ID of the instance.`,
				Required: true,
			},
			{
				Name:     "lb-id",
				Short:    "ID of the load balancer you want to import the instance into",
				Required: true,
			},
			{
				Name:     "protocol",
				Short:    "Protocol used by the backend",
				Required: true,
			},
			{
				Name:  "use-public",
				Short: "Use public IP address of the instance instead of the private one",
			},
			{
				Name:     "port",
				Short:    "Port number used by the backend",
				Required: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			args := argsI.(*importInstanceArgs)
			instanceID := args.InstanceID
			zone := args.InstanceZone
			instanceAPI := instance.NewAPI(core.ExtractClient(ctx))
			server, err := instanceAPI.GetServer(&instance.GetServerRequest{
				Zone:     zone,
				ServerID: instanceID,
			})
			if err != nil {
				return nil, err
			}

			forwardPort := args.Port

			forwardProtocol := lb.ProtocolTCP
			if args.Protocol != nil {
				forwardProtocol = *args.Protocol
			}

			if server.Server.PrivateIP == nil {
				return nil, &core.CliError{
					Message: fmt.Sprintf("server %s (%s) does not have a private ip", server.Server.ID, server.Server.Name),
					Hint:    fmt.Sprintf("Private ip are assigned when the server boots, start yours with: scw instance server start %s", server.Server.ID),
				}
			}
			ip := *server.Server.PrivateIP
			if args.UsePublicIP {
				if server.Server.PublicIP == nil {
					return nil, &core.CliError{
						Message: fmt.Sprintf("server %s (%s) does not have a public ip", server.Server.ID, server.Server.Name),
					}
				}
				ip = server.Server.PublicIP.Address.String()
			}

			lbAPI := lb.NewAPI(core.ExtractClient(ctx))
			backend, err := lbAPI.CreateBackend(&lb.CreateBackendRequest{
				Name:            server.Server.Name,
				LBID:            args.LBID,
				ForwardProtocol: forwardProtocol,
				ForwardPort:     forwardPort,
				HealthCheck: &lb.HealthCheck{
					CheckMaxRetries: 5,
					TCPConfig:       &lb.HealthCheckTCPConfig{},
					Port:            args.Port,
				},
				ServerIP: []string{ip},
			})
			if err != nil {
				return nil, err
			}

			return backend, nil
		},
	}
}
