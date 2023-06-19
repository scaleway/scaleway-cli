package instance

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/sshconfig"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type sshConfigServer struct {
	Name              string
	Address           string
	User              string
	PrivateNetworksID []string
}

func (s sshConfigServer) InPrivateNetwork(id string) bool {
	for _, pnID := range s.PrivateNetworksID {
		if pnID == id {
			return true
		}
	}

	return false
}

type sshConfigRequest struct {
	Zone      scw.Zone
	ProjectID *string
}

func sshConfigInstallCommand() *core.Command {
	return &core.Command{
		Namespace: "instance",
		Resource:  "ssh",
		Verb:      "install-config",
		Short: `Install a ssh config with all your servers as host
It generate hosts for instance servers, baremetal, apple-silicon and bastions`,
		Long:     "Path of the config will be $HOME/.ssh/scaleway.config",
		ArgsType: reflect.TypeOf(sshConfigRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.ZoneArgSpec(((*instance.API)(nil)).Zones()...),
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*sshConfigRequest)
			homeDir := core.ExtractUserHomeDir(ctx)

			// Start server list with instances
			servers, err := sshConfigListServers(ctx, args)
			if err != nil {
				return nil, err
			}

			// Add baremetal servers
			baremetalServers, err := sshConfigListBaremetalServers(ctx, args)
			if err != nil {
				return nil, err
			}
			servers = append(servers, baremetalServers...)

			// Add Apple-Silicon servers
			siliconServers, err := sshConfigListAppleSiliconServers(ctx, args)
			if err != nil {
				return nil, err
			}
			servers = append(servers, siliconServers...)

			// Fill hosts with servers
			hosts := make([]sshconfig.Host, 0, len(servers))
			for _, server := range servers {
				if server.Address == "" {
					continue
				}
				hosts = append(hosts, sshconfig.SimpleHost{
					Name:    server.Name,
					Address: server.Address,
				})
			}

			// Add Bastions to hosts
			bastionHosts, err := sshConfigBastionHosts(ctx, args, servers)
			if err != nil {
				return nil, err
			}
			hosts = append(hosts, bastionHosts...)

			err = sshconfig.Save(homeDir, hosts)
			if err != nil {
				return nil, fmt.Errorf("failed to save config file: %w", err)
			}

			configFilePath := sshconfig.ConfigFilePath(homeDir)
			includePrompt := fmt.Sprintf(`Generated config file needs to be included in your default ssh config (%s)
Do you want the include statement to be added at the beginning of your file ?`, sshconfig.DefaultConfigFilePath(homeDir))

			// Generated config needs an include statement in default config
			included, err := sshconfig.ConfigIsIncluded(homeDir)
			if err != nil {
				if err == sshconfig.ErrFileNotFound {
					includePrompt += "\nFile was not found, it will be created"
				} else {
					logger.Warningf("Failed to check default config file, skipping include prompt\n")
					return &core.SuccessResult{
						Message: "Config file was generated to " + configFilePath,
					}, nil
				}
			}

			// Generated config is already included
			if included {
				return &core.SuccessResult{
					Message: "Config file was generated to " + configFilePath,
				}, nil
			}

			shouldIncludeConfig, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
				Ctx:          ctx,
				Prompt:       includePrompt,
				DefaultValue: true,
			})
			if err != nil {
				logger.Warningf("Failed to prompt, skipping include\n")
				return &core.SuccessResult{
					Message: "Config file was generated to " + configFilePath,
				}, nil
			}

			if shouldIncludeConfig {
				err := sshconfig.IncludeConfigFile(homeDir)
				if err != nil {
					return nil, fmt.Errorf("failed to add include statement: %w", err)
				}
			}

			return &core.SuccessResult{
				Message: "Config file was generated to " + configFilePath,
			}, nil
		},
		Groups: []string{"workflow"},
	}
}

func sshConfigListServers(ctx context.Context, args *sshConfigRequest) ([]sshConfigServer, error) {
	instanceAPI := instance.NewAPI(core.ExtractClient(ctx))

	reqOpts := []scw.RequestOption{scw.WithAllPages()}
	if args.Zone == scw.Zone(core.AllLocalities) {
		reqOpts = append(reqOpts, scw.WithZones(instanceAPI.Zones()...))
	}

	listServers, err := instanceAPI.ListServers(&instance.ListServersRequest{
		Zone:    args.Zone,
		Project: args.ProjectID,
	}, reqOpts...)
	if err != nil {
		return nil, err
	}

	servers := make([]sshConfigServer, len(listServers.Servers))

	for i, server := range listServers.Servers {
		pnIDs := make([]string, len(server.PrivateNics))
		for j, nic := range server.PrivateNics {
			pnIDs[j] = nic.PrivateNetworkID
		}

		serverAddress := ""
		if server.PublicIP != nil {
			serverAddress = server.PublicIP.Address.String()
		}

		servers[i] = sshConfigServer{
			Name:              server.Name,
			Address:           serverAddress,
			PrivateNetworksID: pnIDs,
		}
	}

	return servers, nil
}

func sshConfigListBaremetalServers(ctx context.Context, args *sshConfigRequest) ([]sshConfigServer, error) {
	baremetalAPI := baremetal.NewAPI(core.ExtractClient(ctx))
	baremetalPNAPI := baremetal.NewPrivateNetworkAPI(core.ExtractClient(ctx))

	reqOpts := []scw.RequestOption{scw.WithAllPages()}
	if args.Zone == scw.Zone(core.AllLocalities) {
		reqOpts = append(reqOpts, scw.WithZones(baremetalAPI.Zones()...))
	}

	listServers, err := baremetalAPI.ListServers(&baremetal.ListServersRequest{
		Zone:      args.Zone,
		ProjectID: args.ProjectID,
	}, reqOpts...)
	if err != nil {
		// TODO: check permissions and print warning
		return nil, err
	}
	listPNs, err := baremetalPNAPI.ListServerPrivateNetworks(&baremetal.PrivateNetworkAPIListServerPrivateNetworksRequest{
		Zone: args.Zone,
	}, reqOpts...)
	if err != nil {
		// TODO: check permissions and print warning
		return nil, err
	}

	servers := make([]sshConfigServer, len(listServers.Servers))

	for i, server := range listServers.Servers {
		pnIDs := []string(nil)
		for _, pn := range listPNs.ServerPrivateNetworks {
			if pn.ServerID == server.ID {
				pnIDs = append(pnIDs, pn.PrivateNetworkID)
			}
		}

		address := ""
		if len(server.IPs) > 0 {
			address = server.IPs[0].Address.String()
		}

		servers[i] = sshConfigServer{
			Name:              server.Name,
			Address:           address,
			PrivateNetworksID: pnIDs,
		}
	}

	return servers, nil
}

func sshConfigListAppleSiliconServers(ctx context.Context, args *sshConfigRequest) ([]sshConfigServer, error) {
	siliconAPI := applesilicon.NewAPI(core.ExtractClient(ctx))

	reqOpts := []scw.RequestOption{scw.WithAllPages()}
	if args.Zone == scw.Zone(core.AllLocalities) {
		reqOpts = append(reqOpts, scw.WithZones(siliconAPI.Zones()...))
	}

	listServers, err := siliconAPI.ListServers(&applesilicon.ListServersRequest{
		Zone:      args.Zone,
		ProjectID: args.ProjectID,
	}, reqOpts...)
	if err != nil {
		return nil, err
	}

	servers := make([]sshConfigServer, len(listServers.Servers))

	for i, server := range listServers.Servers {
		servers[i] = sshConfigServer{
			Name:    server.Name,
			Address: server.IP.String(),
		}
	}

	return servers, nil
}

func sshConfigBastionHosts(ctx context.Context, args *sshConfigRequest, servers []sshConfigServer) ([]sshconfig.Host, error) {
	gwAPI := vpcgw.NewAPI(core.ExtractClient(ctx))

	reqOpts := []scw.RequestOption{scw.WithAllPages()}
	if args.Zone == scw.Zone(core.AllLocalities) {
		reqOpts = append(reqOpts, scw.WithZones(gwAPI.Zones()...))
	}

	listGateways, err := gwAPI.ListGateways(&vpcgw.ListGatewaysRequest{
		Zone: args.Zone,
	}, reqOpts...)
	if err != nil {
		// TODO: check permissions and print warning
		return nil, err
	}

	hosts := []sshconfig.Host(nil)

	for _, gateway := range listGateways.Gateways {
		if !gateway.BastionEnabled {
			continue
		}
		for _, network := range gateway.GatewayNetworks {
			bastionHost := sshconfig.BastionHost{
				Name:    network.DHCP.DNSLocalName,
				Address: gateway.IP.Address.String(),
				Port:    gateway.BastionPort,
			}

			for _, server := range servers {
				if server.InPrivateNetwork(network.PrivateNetworkID) {
					bastionHost.Hosts = append(bastionHost.Hosts, sshconfig.SimpleHost{
						Name:    server.Name,
						Address: server.Address,
						User:    server.User,
					})
				}
			}

			hosts = append(hosts, bastionHost)
		}
	}

	return hosts, nil
}
