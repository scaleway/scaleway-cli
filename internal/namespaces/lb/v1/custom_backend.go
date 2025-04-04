package lb

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	backendServerStatsHealthCheckStatusMarshalSpecs = human.EnumMarshalSpecs{
		lb.BackendServerStatsHealthCheckStatusPassed: &human.EnumMarshalSpec{
			Attribute: color.FgGreen,
			Value:     "passed",
		},
		lb.BackendServerStatsHealthCheckStatusFailed: &human.EnumMarshalSpec{
			Attribute: color.FgRed,
			Value:     "failed",
		},
		lb.BackendServerStatsHealthCheckStatusUnknown: &human.EnumMarshalSpec{
			Attribute: color.Faint,
			Value:     "unknown",
		},
		lb.BackendServerStatsHealthCheckStatusNeutral: &human.EnumMarshalSpec{
			Attribute: color.Faint,
			Value:     "neutral",
		},
		lb.BackendServerStatsHealthCheckStatusCondpass: &human.EnumMarshalSpec{
			Attribute: color.FgBlue,
			Value:     "condition passed",
		},
	}
	backendServerStatsServerStateMarshalSpecs = human.EnumMarshalSpecs{
		lb.BackendServerStatsServerStateStopped: &human.EnumMarshalSpec{
			Attribute: color.FgRed,
			Value:     "stopped",
		},
		lb.BackendServerStatsServerStateStarting: &human.EnumMarshalSpec{
			Attribute: color.FgBlue,
			Value:     "starting",
		},
		lb.BackendServerStatsServerStateRunning: &human.EnumMarshalSpec{
			Attribute: color.FgGreen,
			Value:     "running",
		},
		lb.BackendServerStatsServerStateStopping: &human.EnumMarshalSpec{
			Attribute: color.FgBlue,
			Value:     "stopping",
		},
	}
)

func lbBackendMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp lb.Backend
	backend := tmp(i.(lb.Backend))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "HealthCheck",
		},
		{
			FieldName: "Pool",
		},
		{
			FieldName: "LB",
		},
	}

	str, err := human.Marshal(backend, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func backendGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptBackend()

	return c
}

func backendCreateBuilder(c *core.Command) *core.Command {
	type lbBackendCreateRequestCustom struct {
		lb.ZonedAPICreateBackendRequest
		InstanceServerID          []string
		BaremetalServerID         []string
		UseInstanceServerPublicIP bool
		InstanceServerTag         []string
		BaremetalServerTag        []string
	}

	c.ArgsType = reflect.TypeOf(lbBackendCreateRequestCustom{})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "instance-server-id.{index}",
		Short: "UIID of the instance server.",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "instance-server-tag.{index}",
		Short: "Tag of the instance server.",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "use-instance-server-public-ip",
		Short: "Use public IP address of the instance instead of the private one",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "baremetal-server-id.{index}",
		Short: "UIID of the baremetal server.",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "baremetal-server-tag.{index}",
		Short: "Tag of the baremetal server.",
	})

	c.Run = func(ctx context.Context, argsI interface{}) (interface{}, error) {
		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		instanceAPI := instance.NewAPI(core.ExtractClient(ctx))
		baremetalAPI := baremetal.NewAPI(core.ExtractClient(ctx))

		tmpRequest := argsI.(*lbBackendCreateRequestCustom)

		request := &lb.ZonedAPICreateBackendRequest{
			Zone:                     tmpRequest.Zone,
			LBID:                     tmpRequest.LBID,
			Name:                     tmpRequest.Name,
			ForwardProtocol:          tmpRequest.ForwardProtocol,
			ForwardPort:              tmpRequest.ForwardPort,
			ForwardPortAlgorithm:     tmpRequest.ForwardPortAlgorithm,
			StickySessions:           tmpRequest.StickySessions,
			StickySessionsCookieName: tmpRequest.StickySessionsCookieName,
			HealthCheck:              tmpRequest.HealthCheck,
			ServerIP:                 tmpRequest.ServerIP,
			TimeoutServer:            tmpRequest.TimeoutServer,
			TimeoutConnect:           tmpRequest.TimeoutConnect,
			TimeoutTunnel:            tmpRequest.TimeoutTunnel,
			OnMarkedDownAction:       tmpRequest.OnMarkedDownAction,
			ProxyProtocol:            tmpRequest.ProxyProtocol,
			FailoverHost:             tmpRequest.FailoverHost,
			SslBridging:              tmpRequest.SslBridging,
			IgnoreSslServerVerify:    tmpRequest.IgnoreSslServerVerify,
		}

		// IP/ID management
		if len(tmpRequest.InstanceServerID) != 0 {
			var serverIPs []string

			for _, instanceID := range tmpRequest.InstanceServerID {
				server, err := instanceAPI.GetServer(&instance.GetServerRequest{
					Zone:     tmpRequest.Zone,
					ServerID: instanceID,
				})
				if err != nil {
					return nil, err
				}

				if tmpRequest.UseInstanceServerPublicIP {
					if server.Server.PublicIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf(
								"server %s (%s) does not have a public ip",
								server.Server.ID,
								server.Server.Name,
							),
						}
					}
					serverIPs = append(serverIPs, server.Server.PublicIP.Address.String())
				} else {
					if server.Server.PrivateIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf("server %s (%s) does not have a private ip", server.Server.ID, server.Server.Name),
							Hint:    "Private ip are assigned when the server boots, start yours with: scw instance server start " + server.Server.ID,
						}
					}
					serverIPs = append(serverIPs, *server.Server.PrivateIP)
				}

				request.ServerIP = append(request.ServerIP, serverIPs...)
			}
		}

		if len(tmpRequest.BaremetalServerID) != 0 {
			for _, baremetalID := range tmpRequest.BaremetalServerID {
				server, err := baremetalAPI.GetServer(&baremetal.GetServerRequest{
					Zone:     tmpRequest.Zone,
					ServerID: baremetalID,
				})
				if err != nil {
					return nil, err
				}

				for _, ip := range server.IPs {
					request.ServerIP = append(request.ServerIP, ip.Address.String())
				}
			}
		}

		// IP/Tag management
		if len(tmpRequest.InstanceServerTag) != 0 {
			var serverIPs []string

			listServersResponse, err := instanceAPI.ListServers(&instance.ListServersRequest{
				Zone: tmpRequest.Zone,
				Tags: tmpRequest.InstanceServerTag,
			})
			if err != nil {
				return nil, err
			}

			if len(listServersResponse.Servers) == 0 {
				return nil, &core.CliError{
					Err: fmt.Errorf(
						"there is no server with tag(s) '%v'",
						strings.Trim(strings.Join(tmpRequest.InstanceServerTag, " - "), "[]"),
					),
				}
			}

			for _, server := range listServersResponse.Servers {
				if tmpRequest.UseInstanceServerPublicIP {
					if server.PublicIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf(
								"server %s (%s) does not have a public ip",
								server.ID,
								server.Name,
							),
						}
					}
					serverIPs = append(serverIPs, server.PublicIP.Address.String())
				} else {
					if server.PrivateIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf("server %s (%s) does not have a private ip", server.ID, server.Name),
							Hint:    "Private ip are assigned when the server boots, start yours with: scw instance server start " + server.ID,
						}
					}
					serverIPs = append(serverIPs, *server.PrivateIP)
				}
			}
			request.ServerIP = append(request.ServerIP, serverIPs...)
		}

		if len(tmpRequest.BaremetalServerTag) != 0 {
			listServersResponse, err := baremetalAPI.ListServers(&baremetal.ListServersRequest{
				Zone: tmpRequest.Zone,
				Tags: tmpRequest.BaremetalServerTag,
			})
			if err != nil {
				return nil, err
			}

			if len(listServersResponse.Servers) == 0 {
				return nil, &core.CliError{
					Err: fmt.Errorf(
						"there is no server with tag(s) '%v'",
						strings.Trim(strings.Join(tmpRequest.InstanceServerTag, " - "), "[]"),
					),
				}
			}

			for i, server := range listServersResponse.Servers {
				request.ServerIP = append(request.ServerIP, server.IPs[i].Address.String())
			}
		}

		return api.CreateBackend(request)
	}

	c.Interceptor = interceptBackend()

	return c
}

func backendUpdateBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptBackend()

	return c
}

func backendDeleteBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptBackend()

	return c
}

func backendAddServersBuilder(c *core.Command) *core.Command {
	type lbBackendAddBackendServersRequestCustom struct {
		lb.ZonedAPIAddBackendServersRequest
		InstanceServerID          []string
		BaremetalServerID         []string
		UseInstanceServerPublicIP bool
		InstanceServerTag         []string
		BaremetalServerTag        []string
	}

	c.ArgsType = reflect.TypeOf(lbBackendAddBackendServersRequestCustom{})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "instance-server-id.{index}",
		Short: "UIID of the instance server.",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "instance-server-tag.{index}",
		Short: "Tag of the instance server.",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "use-instance-server-public-ip",
		Short: "Use public IP address of the instance instead of the private one",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "baremetal-server-id.{index}",
		Short: "UIID of the baremetal server.",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "baremetal-server-tag.{index}",
		Short: "Tag of the baremetal server.",
	})

	c.Run = func(ctx context.Context, argsI interface{}) (interface{}, error) {
		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		instanceAPI := instance.NewAPI(core.ExtractClient(ctx))
		baremetalAPI := baremetal.NewAPI(core.ExtractClient(ctx))

		tmpRequest := argsI.(*lbBackendAddBackendServersRequestCustom)

		request := &lb.ZonedAPIAddBackendServersRequest{
			Zone:      tmpRequest.Zone,
			BackendID: tmpRequest.BackendID,
			ServerIP:  tmpRequest.ServerIP,
		}

		// IP/ID management
		if len(tmpRequest.InstanceServerID) != 0 {
			var serverIPs []string

			for _, instanceID := range tmpRequest.InstanceServerID {
				server, err := instanceAPI.GetServer(&instance.GetServerRequest{
					Zone:     tmpRequest.Zone,
					ServerID: instanceID,
				})
				if err != nil {
					return nil, err
				}

				if tmpRequest.UseInstanceServerPublicIP {
					if server.Server.PublicIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf(
								"server %s (%s) does not have a public ip",
								server.Server.ID,
								server.Server.Name,
							),
						}
					}
					serverIPs = append(serverIPs, server.Server.PublicIP.Address.String())
				} else {
					if server.Server.PrivateIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf("server %s (%s) does not have a private ip", server.Server.ID, server.Server.Name),
							Hint:    "Private ip are assigned when the server boots, start yours with: scw instance server start " + server.Server.ID,
						}
					}
					serverIPs = append(serverIPs, *server.Server.PrivateIP)
				}

				request.ServerIP = append(request.ServerIP, serverIPs...)
			}
		}

		if len(tmpRequest.BaremetalServerID) != 0 {
			for _, baremetalID := range tmpRequest.BaremetalServerID {
				server, err := baremetalAPI.GetServer(&baremetal.GetServerRequest{
					Zone:     tmpRequest.Zone,
					ServerID: baremetalID,
				})
				if err != nil {
					return nil, err
				}

				for _, ip := range server.IPs {
					request.ServerIP = append(request.ServerIP, ip.Address.String())
				}
			}
		}

		// IP/Tag management
		if len(tmpRequest.InstanceServerTag) != 0 {
			var serverIPs []string

			listServersResponse, err := instanceAPI.ListServers(&instance.ListServersRequest{
				Zone: tmpRequest.Zone,
				Tags: tmpRequest.InstanceServerTag,
			})
			if err != nil {
				return nil, err
			}

			if len(listServersResponse.Servers) == 0 {
				return nil, &core.CliError{
					Err: fmt.Errorf(
						"there is no server with tag(s) '%v'",
						strings.Trim(strings.Join(tmpRequest.InstanceServerTag, " - "), "[]"),
					),
				}
			}

			for _, server := range listServersResponse.Servers {
				if tmpRequest.UseInstanceServerPublicIP {
					if server.PublicIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf(
								"server %s (%s) does not have a public ip",
								server.ID,
								server.Name,
							),
						}
					}
					serverIPs = append(serverIPs, server.PublicIP.Address.String())
				} else {
					if server.PrivateIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf("server %s (%s) does not have a private ip", server.ID, server.Name),
							Hint:    "Private ip are assigned when the server boots, start yours with: scw instance server start " + server.ID,
						}
					}
					serverIPs = append(serverIPs, *server.PrivateIP)
				}
			}
			request.ServerIP = append(request.ServerIP, serverIPs...)
		}

		if len(tmpRequest.BaremetalServerTag) != 0 {
			listServersResponse, err := baremetalAPI.ListServers(&baremetal.ListServersRequest{
				Zone: tmpRequest.Zone,
				Tags: tmpRequest.BaremetalServerTag,
			})
			if err != nil {
				return nil, err
			}

			if len(listServersResponse.Servers) == 0 {
				return nil, &core.CliError{
					Err: fmt.Errorf(
						"there is no server with tag(s) '%v'",
						strings.Trim(strings.Join(tmpRequest.InstanceServerTag, " - "), "[]"),
					),
				}
			}

			for i, server := range listServersResponse.Servers {
				request.ServerIP = append(request.ServerIP, server.IPs[i].Address.String())
			}
		}

		return api.AddBackendServers(request)
	}

	c.Interceptor = interceptBackend()

	return c
}

func backendRemoveServersBuilder(c *core.Command) *core.Command {
	type lbBackendRemoveBackendServersRequestCustom struct {
		lb.ZonedAPIRemoveBackendServersRequest
		InstanceServerID          []string
		BaremetalServerID         []string
		UseInstanceServerPublicIP bool
		InstanceServerTag         []string
		BaremetalServerTag        []string
	}

	c.ArgsType = reflect.TypeOf(lbBackendRemoveBackendServersRequestCustom{})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "instance-server-id.{index}",
		Short: "UIID of the instance server.",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "instance-server-tag.{index}",
		Short: "Tag of the instance server.",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "use-instance-server-public-ip",
		Short: "Use public IP address of the instance instead of the private one",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "baremetal-server-id.{index}",
		Short: "UIID of the baremetal server.",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "baremetal-server-tag.{index}",
		Short: "Tag of the baremetal server.",
	})

	c.Run = func(ctx context.Context, argsI interface{}) (interface{}, error) {
		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		instanceAPI := instance.NewAPI(core.ExtractClient(ctx))
		baremetalAPI := baremetal.NewAPI(core.ExtractClient(ctx))

		tmpRequest := argsI.(*lbBackendRemoveBackendServersRequestCustom)

		request := &lb.ZonedAPIRemoveBackendServersRequest{
			Zone:      tmpRequest.Zone,
			BackendID: tmpRequest.BackendID,
			ServerIP:  tmpRequest.ServerIP,
		}

		// IP/ID management
		if len(tmpRequest.InstanceServerID) != 0 {
			var serverIPs []string

			for _, instanceID := range tmpRequest.InstanceServerID {
				server, err := instanceAPI.GetServer(&instance.GetServerRequest{
					Zone:     tmpRequest.Zone,
					ServerID: instanceID,
				})
				if err != nil {
					return nil, err
				}

				if tmpRequest.UseInstanceServerPublicIP {
					if server.Server.PublicIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf(
								"server %s (%s) does not have a public ip",
								server.Server.ID,
								server.Server.Name,
							),
						}
					}
					serverIPs = append(serverIPs, server.Server.PublicIP.Address.String())
				} else {
					if server.Server.PrivateIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf("server %s (%s) does not have a private ip", server.Server.ID, server.Server.Name),
							Hint:    "Private ip are assigned when the server boots, start yours with: scw instance server start " + server.Server.ID,
						}
					}
					serverIPs = append(serverIPs, *server.Server.PrivateIP)
				}

				request.ServerIP = append(request.ServerIP, serverIPs...)
			}
		}

		if len(tmpRequest.BaremetalServerID) != 0 {
			for _, baremetalID := range tmpRequest.BaremetalServerID {
				server, err := baremetalAPI.GetServer(&baremetal.GetServerRequest{
					Zone:     tmpRequest.Zone,
					ServerID: baremetalID,
				})
				if err != nil {
					return nil, err
				}

				for _, ip := range server.IPs {
					request.ServerIP = append(request.ServerIP, ip.Address.String())
				}
			}
		}

		// IP/Tag management
		if len(tmpRequest.InstanceServerTag) != 0 {
			var serverIPs []string

			listServersResponse, err := instanceAPI.ListServers(&instance.ListServersRequest{
				Zone: tmpRequest.Zone,
				Tags: tmpRequest.InstanceServerTag,
			})
			if err != nil {
				return nil, err
			}

			if len(listServersResponse.Servers) == 0 {
				return nil, &core.CliError{
					Err: fmt.Errorf(
						"there is no server with tag(s) '%v'",
						strings.Trim(strings.Join(tmpRequest.InstanceServerTag, " - "), "[]"),
					),
				}
			}

			for _, server := range listServersResponse.Servers {
				if tmpRequest.UseInstanceServerPublicIP {
					if server.PublicIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf(
								"server %s (%s) does not have a public ip",
								server.ID,
								server.Name,
							),
						}
					}
					serverIPs = append(serverIPs, server.PublicIP.Address.String())
				} else {
					if server.PrivateIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf("server %s (%s) does not have a private ip", server.ID, server.Name),
							Hint:    "Private ip are assigned when the server boots, start yours with: scw instance server start " + server.ID,
						}
					}
					serverIPs = append(serverIPs, *server.PrivateIP)
				}
			}
			request.ServerIP = append(request.ServerIP, serverIPs...)
		}

		if len(tmpRequest.BaremetalServerTag) != 0 {
			listServersResponse, err := baremetalAPI.ListServers(&baremetal.ListServersRequest{
				Zone: tmpRequest.Zone,
				Tags: tmpRequest.BaremetalServerTag,
			})
			if err != nil {
				return nil, err
			}

			if len(listServersResponse.Servers) == 0 {
				return nil, &core.CliError{
					Err: fmt.Errorf(
						"there is no server with tag(s) '%v'",
						strings.Trim(strings.Join(tmpRequest.InstanceServerTag, " - "), "[]"),
					),
				}
			}

			for i, server := range listServersResponse.Servers {
				request.ServerIP = append(request.ServerIP, server.IPs[i].Address.String())
			}
		}

		return api.RemoveBackendServers(request)
	}

	c.Interceptor = interceptBackend()

	return c
}

func backendSetServersBuilder(c *core.Command) *core.Command {
	type lbBackendSetBackendServersRequestCustom struct {
		lb.ZonedAPISetBackendServersRequest
		InstanceServerID          []string
		BaremetalServerID         []string
		UseInstanceServerPublicIP bool
		InstanceServerTag         []string
		BaremetalServerTag        []string
	}

	c.ArgsType = reflect.TypeOf(lbBackendSetBackendServersRequestCustom{})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "instance-server-id.{index}",
		Short: "UIID of the instance server.",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "instance-server-tag.{index}",
		Short: "Tag of the instance server.",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "use-instance-server-public-ip",
		Short: "Use public IP address of the instance instead of the private one",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "baremetal-server-id.{index}",
		Short: "UIID of the baremetal server.",
	})

	c.ArgSpecs.AddBefore("server-ip.{index}", &core.ArgSpec{
		Name:  "baremetal-server-tag.{index}",
		Short: "Tag of the baremetal server.",
	})

	c.Run = func(ctx context.Context, argsI interface{}) (interface{}, error) {
		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		instanceAPI := instance.NewAPI(core.ExtractClient(ctx))
		baremetalAPI := baremetal.NewAPI(core.ExtractClient(ctx))

		tmpRequest := argsI.(*lbBackendSetBackendServersRequestCustom)

		request := &lb.ZonedAPISetBackendServersRequest{
			Zone:      tmpRequest.Zone,
			BackendID: tmpRequest.BackendID,
			ServerIP:  tmpRequest.ServerIP,
		}

		// IP/ID management
		if len(tmpRequest.InstanceServerID) != 0 {
			var serverIPs []string

			for _, instanceID := range tmpRequest.InstanceServerID {
				server, err := instanceAPI.GetServer(&instance.GetServerRequest{
					Zone:     tmpRequest.Zone,
					ServerID: instanceID,
				})
				if err != nil {
					return nil, err
				}

				if tmpRequest.UseInstanceServerPublicIP {
					if server.Server.PublicIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf(
								"server %s (%s) does not have a public ip",
								server.Server.ID,
								server.Server.Name,
							),
						}
					}
					serverIPs = append(serverIPs, server.Server.PublicIP.Address.String())
				} else {
					if server.Server.PrivateIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf("server %s (%s) does not have a private ip", server.Server.ID, server.Server.Name),
							Hint:    "Private ip are assigned when the server boots, start yours with: scw instance server start " + server.Server.ID,
						}
					}
					serverIPs = append(serverIPs, *server.Server.PrivateIP)
				}

				request.ServerIP = append(request.ServerIP, serverIPs...)
			}
		}

		if len(tmpRequest.BaremetalServerID) != 0 {
			for _, baremetalID := range tmpRequest.BaremetalServerID {
				server, err := baremetalAPI.GetServer(&baremetal.GetServerRequest{
					Zone:     tmpRequest.Zone,
					ServerID: baremetalID,
				})
				if err != nil {
					return nil, err
				}

				for _, ip := range server.IPs {
					request.ServerIP = append(request.ServerIP, ip.Address.String())
				}
			}
		}

		// IP/Tag management
		if len(tmpRequest.InstanceServerTag) != 0 {
			var serverIPs []string

			listServersResponse, err := instanceAPI.ListServers(&instance.ListServersRequest{
				Zone: tmpRequest.Zone,
				Tags: tmpRequest.InstanceServerTag,
			})
			if err != nil {
				return nil, err
			}

			if len(listServersResponse.Servers) == 0 {
				return nil, &core.CliError{
					Err: fmt.Errorf(
						"there is no server with tag(s) '%v'",
						strings.Trim(strings.Join(tmpRequest.InstanceServerTag, " - "), "[]"),
					),
				}
			}

			for _, server := range listServersResponse.Servers {
				if tmpRequest.UseInstanceServerPublicIP {
					if server.PublicIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf(
								"server %s (%s) does not have a public ip",
								server.ID,
								server.Name,
							),
						}
					}
					serverIPs = append(serverIPs, server.PublicIP.Address.String())
				} else {
					if server.PrivateIP == nil {
						return nil, &core.CliError{
							Message: fmt.Sprintf("server %s (%s) does not have a private ip", server.ID, server.Name),
							Hint:    "Private ip are assigned when the server boots, start yours with: scw instance server start " + server.ID,
						}
					}
					serverIPs = append(serverIPs, *server.PrivateIP)
				}
			}
			request.ServerIP = append(request.ServerIP, serverIPs...)
		}

		if len(tmpRequest.BaremetalServerTag) != 0 {
			listServersResponse, err := baremetalAPI.ListServers(&baremetal.ListServersRequest{
				Zone: tmpRequest.Zone,
				Tags: tmpRequest.BaremetalServerTag,
			})
			if err != nil {
				return nil, err
			}

			if len(listServersResponse.Servers) == 0 {
				return nil, &core.CliError{
					Err: fmt.Errorf(
						"there is no server with tag(s) '%v'",
						strings.Trim(strings.Join(tmpRequest.InstanceServerTag, " - "), "[]"),
					),
				}
			}

			for i, server := range listServersResponse.Servers {
				request.ServerIP = append(request.ServerIP, server.IPs[i].Address.String())
			}
		}

		return api.SetBackendServers(request)
	}

	c.Interceptor = interceptBackend()

	return c
}

func backendUpdateHealthcheckBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptBackend()

	return c
}

func interceptBackend() core.CommandInterceptor {
	return func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		backend, err := getBackendBeforeAction(api, argsI)
		if err != nil {
			return nil, err
		}

		res, err := runner(ctx, argsI)
		if err != nil {
			var invalidArgErr *scw.InvalidArgumentsError
			if errors.As(err, &invalidArgErr) {
				for _, detail := range invalidArgErr.Details {
					switch detail.ArgumentName {
					case "Port":
						return nil, &core.CliError{
							Err: errors.New("missing or invalid 'health-check.port' argument"),
						}
					case "CheckMaxRetries":
						return nil, &core.CliError{
							Err: errors.New(
								"missing or invalid 'health-check.check-max-retries' argument",
							),
						}
					}
				}
			}

			return nil, err
		}

		switch res.(type) {
		case *core.SuccessResult, *lb.HealthCheck:
			if len(backend.LB.Tags) != 0 && backend.LB.Tags[0] == kapsuleTag {
				return warningKapsuleTaggedMessageView(), nil
			}
		}

		return res, nil
	}
}

func getBackendBeforeAction(api *lb.ZonedAPI, argsI interface{}) (*lb.Backend, error) {
	switch args := argsI.(type) {
	case *lb.ZonedAPIDeleteBackendRequest:
		return api.GetBackend(&lb.ZonedAPIGetBackendRequest{
			Zone:      args.Zone,
			BackendID: args.BackendID,
		})
	case *lb.ZonedAPIUpdateHealthCheckRequest:
		return api.GetBackend(&lb.ZonedAPIGetBackendRequest{
			Zone:      args.Zone,
			BackendID: args.BackendID,
		})
	default:
		return nil, nil
	}
}
