package lb

import (
	"context"
	"fmt"
	"net"
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
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

func addServerCommand() *core.Command {
	type importInstanceArgs struct {
		InstanceServerID  string
		BaremetalServerID string
		LBID              string
		IP                *net.IPAddr
		Protocol          *lb.Protocol
		Port              int32
		InstanceZone      scw.Zone
		UsePublicIP       bool
		Region            scw.Region
	}

	return &core.Command{
		Short:     `Import an instance as a load balancer backend`,
		Long:      `Import an instance as a load balancer backend.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "import",
		ArgsType:  reflect.TypeOf(importInstanceArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-server-id",
				Short:      "ID of the instance server.",
				OneOfGroup: "id",
			},
			{
				Name:       "baremetal-server-id",
				Short:      "ID of the baremetal server.",
				OneOfGroup: "id",
			},
			{
				Name:       "ip",
				Short:      "IP of the server you want to add.",
				OneOfGroup: "id",
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
				Name:  "use-instance-server-public-ip",
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
			req := &lb.CreateBackendRequest{}

			req.LBID = args.LBID

			req.ForwardPort = args.Port
			req.HealthCheck = &lb.HealthCheck{
				CheckMaxRetries: 5,
				TCPConfig:       &lb.HealthCheckTCPConfig{},
				Port:            args.Port,
			}

			req.ForwardProtocol = lb.ProtocolTCP
			if args.Protocol != nil {
				req.ForwardProtocol = *args.Protocol
			}

			if args.InstanceServerID != "" {
				instanceServerID := args.InstanceServerID
				zone := args.InstanceZone
				instanceAPI := instance.NewAPI(core.ExtractClient(ctx))
				server, err := instanceAPI.GetServer(&instance.GetServerRequest{
					Zone:     zone,
					ServerID: instanceServerID,
				})
				if err != nil {
					return nil, err
				}

				req.Name = server.Server.Name

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
				req.ServerIP = []string{ip}

				lbAPI := lb.NewAPI(core.ExtractClient(ctx))
				backend, err := lbAPI.CreateBackend(req)
				if err != nil {
					return nil, err
				}
				return backend, nil
			}

			if args.BaremetalServerID != "" {
				baremetalServerID := args.BaremetalServerID
				zone := args.InstanceZone
				baremetalAPI := baremetal.NewAPI(core.ExtractClient(ctx))
				server, err := baremetalAPI.GetServer(&baremetal.GetServerRequest{
					Zone:     zone,
					ServerID: baremetalServerID,
				})
				if err != nil {
					return nil, err
				}

				req.Name = server.Name

				var ips []string
				for _, ip := range server.IPs {
					ips = append(ips, ip.Address.String())
				}
				req.ServerIP = ips

				lbAPI := lb.NewAPI(core.ExtractClient(ctx))
				backend, err := lbAPI.CreateBackend(req)
				if err != nil {
					return nil, err
				}
				return backend, nil
			}

			if args.IP != nil {
				req.Name = args.IP.String()
				req.ServerIP = []string{args.IP.String()}

				lbAPI := lb.NewAPI(core.ExtractClient(ctx))
				backend, err := lbAPI.CreateBackend(req)
				if err != nil {
					return nil, err
				}
				return backend, nil
			}

			return nil, &core.CliError{
				Message: "Both instance-server-id and baremetal-server-id seems to be empty",
				Details: fmt.Sprintf("instance-server-id: %s | baremetal-server-id: %s", args.InstanceServerID, args.BaremetalServerID),
				Hint:    "Specify one of instance-server-id or baremetal-server-id",
			}
		},
	}
}
