package instance

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	block "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

const (
	serverActionTimeout = 10 * time.Minute
)

//
// Marshalers
//

// serverStateMarshalSpecs allows to override the displayed instance.ServerState.
var (
	serverStateMarshalSpecs = human.EnumMarshalSpecs{
		instance.ServerStateRunning: &human.EnumMarshalSpec{Attribute: color.FgGreen},
		instance.ServerStateStopped: &human.EnumMarshalSpec{
			Attribute: color.Faint,
			Value:     "archived",
		},
		instance.ServerStateStoppedInPlace: &human.EnumMarshalSpec{Attribute: color.Faint},
		instance.ServerStateStarting:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
		instance.ServerStateStopping:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
		instance.ServerStateLocked:         &human.EnumMarshalSpec{Attribute: color.FgRed},
	}
)

// serverLocationMarshalerFunc marshals a instance.ServerLocation.
func serverLocationMarshalerFunc(i any, _ *human.MarshalOpt) (string, error) {
	location := i.(instance.ServerLocation)
	zone, err := scw.ParseZone(location.ZoneID)
	if err != nil {
		return "", err
	}

	return string(zone), nil
}

// serversMarshalerFunc marshals a Server.
func serversMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	// humanServerInList is the custom Server type used for list view.
	type humanServerInList struct {
		ID                string
		Name              string
		Type              string
		State             instance.ServerState
		Zone              scw.Zone
		PublicIP          net.IP
		PrivateIP         *string
		Tags              []string
		ImageName         string
		RoutedIPEnabled   *bool
		PlacementGroup    *instance.PlacementGroup
		ModificationDate  *time.Time
		CreationDate      *time.Time
		Volumes           int
		Protected         bool
		SecurityGroupName string
		SecurityGroupID   string
		StateDetail       string
		Arch              instance.Arch
		ImageID           string
	}

	servers := i.([]*instance.Server)
	humanServers := make([]*humanServerInList, 0)
	for _, server := range servers {
		publicIPAddress := net.IP(nil)
		if server.PublicIP != nil {
			publicIPAddress = server.PublicIP.Address
		}
		serverImageID := ""
		serverImageName := ""
		if server.Image != nil {
			serverImageID = server.Image.ID
			serverImageName = server.Image.Name
		}
		humanServers = append(humanServers, &humanServerInList{
			ID:                server.ID,
			Name:              server.Name,
			Type:              server.CommercialType,
			State:             server.State,
			Zone:              server.Zone,
			PublicIP:          publicIPAddress,
			PrivateIP:         server.PrivateIP,
			Tags:              server.Tags,
			ImageName:         serverImageName,
			RoutedIPEnabled:   server.RoutedIPEnabled, //nolint: staticcheck // Field is deprecated but still supported
			PlacementGroup:    server.PlacementGroup,
			ModificationDate:  server.ModificationDate,
			CreationDate:      server.CreationDate,
			Volumes:           len(server.Volumes),
			Protected:         server.Protected,
			SecurityGroupName: server.SecurityGroup.Name,
			SecurityGroupID:   server.SecurityGroup.ID,
			StateDetail:       server.StateDetail,
			Arch:              server.Arch,
			ImageID:           serverImageID,
		})
	}

	return human.Marshal(humanServers, opt)
}

type customVolume struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Size             scw.Size   `json:"size"`
	VolumeType       string     `json:"volume_type"`
	IOPS             string     `json:"iops"`
	State            string     `json:"state"`
	CreationDate     *time.Time `json:"creation_date"`
	ModificationDate *time.Time `json:"modification_date"`
	Boot             bool       `json:"boot"`
	Zone             string     `json:"zone"`
}

// orderVolumes return an ordered slice based on the volume map key "0", "1", "2",...
func orderVolumes(v map[string]*customVolume) []*customVolume {
	indexes := []string(nil)
	for index := range v {
		indexes = append(indexes, index)
	}
	sort.Strings(indexes)

	orderedVolumes := make([]*customVolume, 0, len(indexes))
	for _, index := range indexes {
		orderedVolumes = append(orderedVolumes, v[index])
	}

	return orderedVolumes
}

type ServerWithWarningsResponse struct {
	*instance.Server
	Warnings []string
}

// serversMarshalerFunc marshals a BootscriptID.
func bootscriptMarshalerFunc(i any, _ *human.MarshalOpt) (string, error) {
	bootscript := i.(instance.Bootscript)

	return bootscript.Title, nil
}

//
// Builders
//

func serverListBuilder(c *core.Command) *core.Command {
	type customListServersRequest struct {
		*instance.ListServersRequest
		OrganizationID *string
		ProjectID      *string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customListServersRequest{})

	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
			args := argsI.(*customListServersRequest)

			if args.ListServersRequest == nil {
				args.ListServersRequest = &instance.ListServersRequest{}
			}

			request := args.ListServersRequest
			request.Organization = args.OrganizationID
			request.Project = args.ProjectID

			return runner(ctx, request)
		},
	)

	return c
}

func serverUpdateBuilder(c *core.Command) *core.Command {
	type instanceUpdateServerRequestCustom struct {
		*instance.UpdateServerRequest
		IP               *instance.NullableStringValue
		PlacementGroupID *instance.NullableStringValue
		SecurityGroupID  *string
		VolumeIDs        *[]string
		CloudInit        string
	}

	c.ArgsType = reflect.TypeOf(instanceUpdateServerRequestCustom{})

	// Add completion functions
	c.ArgSpecs.GetByName("admin-password-encryption-ssh-key-id").AutoCompleteFunc = completeSSHKeyID

	// Rename modified arg specs.
	c.ArgSpecs.GetByName("placement-group").Name = "placement-group-id"
	c.ArgSpecs.GetByName("security-group.id").Name = "security-group-id"

	// Delete unused arg specs.
	c.ArgSpecs.DeleteByName("security-group.name")
	c.ArgSpecs.DeleteByName("volumes.{key}.name")
	c.ArgSpecs.DeleteByName("volumes.{key}.size")
	c.ArgSpecs.DeleteByName("volumes.{key}.id")
	c.ArgSpecs.DeleteByName("volumes.{key}.volume-type")

	// Add new arg specs.
	c.ArgSpecs.AddBefore("placement-group-id", &core.ArgSpec{
		Name:  "volume-ids.{index}",
		Short: "Will update ALL volume IDs at once, including the root volume of the server (use volume-ids=none to detach all volumes)",
	})
	c.ArgSpecs.AddBefore("boot-type", &core.ArgSpec{
		Name:  "ip",
		Short: `IP that should be attached to the server (use ip=none to detach)`,
	})
	c.ArgSpecs.AddBefore("boot-type", &core.ArgSpec{
		Name:        "cloud-init",
		Short:       "The cloud-init script to use",
		CanLoadFile: true,
	})

	c.Run = func(ctx context.Context, argsI any) (i any, e error) {
		customRequest := argsI.(*instanceUpdateServerRequestCustom)

		updateServerRequest := customRequest.UpdateServerRequest
		updateServerRequest.PlacementGroup = customRequest.PlacementGroupID
		if customRequest.SecurityGroupID != nil {
			updateServerRequest.SecurityGroup = &instance.SecurityGroupTemplate{
				ID: *customRequest.SecurityGroupID,
			}
		}

		attachIPRequest := (*instance.UpdateIPRequest)(nil)

		detachIP := false

		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)

		getServerResponse, err := api.GetServer(&instance.GetServerRequest{
			Zone:     updateServerRequest.Zone,
			ServerID: customRequest.ServerID,
		})
		if err != nil {
			return nil, err
		}

		switch {
		case customRequest.IP == nil:
			// ip is not set
			// do nothing

		case customRequest.IP.Null:
			// ip=none
			// detach IP from the server, only if it was attached.
			if getServerResponse.Server.PublicIP != nil {
				detachIP = true
			}

		default:
			// ip=<anything>
			// update ip
			if getServerResponse.Server.PublicIP != nil {
				detachIP = true
			}
			attachIPRequest = &instance.UpdateIPRequest{
				IP: customRequest.IP.Value,
				Server: &instance.NullableStringValue{
					Value: customRequest.ServerID,
				},
			}
		}

		// Instance API does not support detaching the existing IP and then attaching a new one to the same server
		// in 1 call only.
		// We need to do it manually in 2 calls.

		if detachIP {
			_, err = api.UpdateIP(&instance.UpdateIPRequest{
				IP: getServerResponse.Server.PublicIP.ID,
				Server: &instance.NullableStringValue{
					Null: true,
				},
			})
			if err != nil {
				return nil, err
			}
		}

		if attachIPRequest != nil {
			_, err = api.UpdateIP(attachIPRequest)
			if err != nil {
				return nil, err
			}
		}

		// Update all volume IDs at once.
		if customRequest.VolumeIDs != nil {
			volumes := make(map[string]*instance.VolumeServerTemplate)
			for i, volumeID := range *customRequest.VolumeIDs {
				index := strconv.Itoa(i)

				if volumeIsFromSBS(block.NewAPI(client), customRequest.Zone, volumeID) {
					volumes[index] = &instance.VolumeServerTemplate{
						ID:         scw.StringPtr(volumeID),
						VolumeType: instance.VolumeVolumeTypeSbsVolume,
					}
				} else {
					volumes[index] = &instance.VolumeServerTemplate{
						ID:   scw.StringPtr(volumeID),
						Name: scw.StringPtr(getServerResponse.Server.Name + "-" + index),
					}
				}
			}
			customRequest.Volumes = &volumes
		}

		// Set cloud-init
		if customRequest.CloudInit != "" {
			err := api.SetServerUserData(&instance.SetServerUserDataRequest{
				Zone:     updateServerRequest.Zone,
				ServerID: customRequest.ServerID,
				Key:      "cloud-init",
				Content:  bytes.NewBufferString(customRequest.CloudInit),
			})
			if err != nil {
				return nil, err
			}
		}

		updateServerResponse, err := api.UpdateServer(updateServerRequest)
		if err != nil {
			return nil, err
		}

		// Display warning if server-type is deprecated
		warnings := []string(nil)
		server := updateServerResponse.Server
		if server.EndOfService {
			warnings = warningServerTypeDeprecated(ctx, client, server)
		}

		return &ServerWithWarningsResponse{
			server,
			warnings,
		}, nil
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName:   "Warnings",
				Title:       "Warnings",
				HideIfEmpty: true,
			},
		},
	}

	return c
}

// customServer is a copy of instance.Server without the fields that are deprecated or duplicated in another section.
// It is used by the `instance server get` command.
type customServer struct {
	ID                              string                         `json:"id"`
	Name                            string                         `json:"name"`
	Organization                    string                         `json:"organization"`
	Project                         string                         `json:"project"`
	AllowedActions                  []instance.ServerAction        `json:"allowed_actions"`
	Tags                            []string                       `json:"tags"`
	CommercialType                  string                         `json:"commercial_type"`
	CreationDate                    *time.Time                     `json:"creation_date"`
	DynamicIPRequired               bool                           `json:"dynamic_ip_required"`
	RoutedIPEnabled                 *bool                          `json:"routed_ip_enabled"`
	Hostname                        string                         `json:"hostname"`
	Image                           *instance.Image                `json:"image"`
	Protected                       bool                           `json:"protected"`
	PrivateIP                       *string                        `json:"private_ip"`
	PublicIPs                       []*instance.ServerIP           `json:"public_ips"`
	MacAddress                      string                         `json:"mac_address"`
	ModificationDate                *time.Time                     `json:"modification_date"`
	State                           instance.ServerState           `json:"state"`
	StateDetail                     string                         `json:"state_detail"`
	IPv6                            *instance.ServerIPv6           `json:"ipv6"`
	BootType                        instance.BootType              `json:"boot_type"`
	SecurityGroup                   *instance.SecurityGroupSummary `json:"security_group"`
	Arch                            instance.Arch                  `json:"arch"`
	PlacementGroup                  *instance.PlacementGroup       `json:"placement_group"`
	Zone                            scw.Zone                       `json:"zone"`
	AdminPasswordEncryptionSSHKeyID *string                        `json:"admin_password_encryption_ssh_key_id"`
	AdminPasswordEncryptedValue     *string                        `json:"admin_password_encrypted_value"`
	Filesystems                     []*instance.ServerFilesystem   `json:"filesystems"`
	EndOfService                    bool                           `json:"end_of_service"`
}

func customServerFromInstanceServer(server *instance.Server) *customServer {
	return &customServer{
		ID:                              server.ID,
		Name:                            server.Name,
		Organization:                    server.Organization,
		Project:                         server.Project,
		AllowedActions:                  server.AllowedActions,
		Tags:                            server.Tags,
		CommercialType:                  server.CommercialType,
		CreationDate:                    server.CreationDate,
		DynamicIPRequired:               server.DynamicIPRequired,
		RoutedIPEnabled:                 server.RoutedIPEnabled,
		Hostname:                        server.Hostname,
		Image:                           server.Image,
		Protected:                       server.Protected,
		PrivateIP:                       server.PrivateIP,
		PublicIPs:                       server.PublicIPs,
		MacAddress:                      server.MacAddress,
		ModificationDate:                server.ModificationDate,
		State:                           server.State,
		StateDetail:                     server.StateDetail,
		BootType:                        server.BootType,
		SecurityGroup:                   server.SecurityGroup,
		Arch:                            server.Arch,
		PlacementGroup:                  server.PlacementGroup,
		Zone:                            server.Zone,
		AdminPasswordEncryptionSSHKeyID: server.AdminPasswordEncryptionSSHKeyID,
		AdminPasswordEncryptedValue:     server.AdminPasswordEncryptedValue,
		Filesystems:                     server.Filesystems,
		EndOfService:                    server.EndOfService,
	}
}

func serverGetBuilder(c *core.Command) *core.Command {
	// This method is here as a proof of concept before we find the correct way to implement it at larger scale
	c.ArgSpecs.GetPositionalArg().AutoCompleteFunc = func(ctx context.Context, prefix string, request any) core.AutocompleteSuggestions {
		req := request.(*instance.GetServerRequest)
		api := instance.NewAPI(core.ExtractClient(ctx))
		resp, err := api.ListServers(&instance.ListServersRequest{
			Zone: req.Zone,
		}, scw.WithAllPages())
		if err != nil {
			return nil
		}

		suggestion := core.AutocompleteSuggestions{}
		for _, s := range resp.Servers {
			if strings.HasPrefix(s.ID, prefix) {
				suggestion = append(suggestion, s.ID)
			}
		}

		return suggestion
	}

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}
		getServerResp := rawResp.(*instance.GetServerResponse)
		server := getServerResp.Server

		client := core.ExtractClient(ctx)
		vpcAPI := vpc.NewAPI(client)

		type customNICs struct {
			ID                 string `json:"id"`
			MacAddress         string `json:"mac_address"`
			PrivateNetworkName string `json:"private_network_name"`
			PrivateNetworkID   string `json:"private_network_id"`
		}

		nics := []customNICs{}

		for _, nic := range server.PrivateNics {
			region, err := server.Zone.Region()
			if err != nil {
				return nil, err
			}
			pn, err := vpcAPI.GetPrivateNetwork(&vpc.GetPrivateNetworkRequest{
				PrivateNetworkID: nic.PrivateNetworkID,
				Region:           region,
			})
			if err != nil {
				return nil, err
			}
			nics = append(nics, customNICs{
				ID:                 nic.ID,
				PrivateNetworkID:   pn.ID,
				PrivateNetworkName: pn.Name,
				MacAddress:         nic.MacAddress,
			})
		}

		volumes := map[string]*customVolume{}
		blockAPI := block.NewAPI(client)

		for _, volume := range server.Volumes {
			customVol := &customVolume{
				ID:   volume.ID,
				Zone: volume.Zone.String(),
				Boot: volume.Boot,
			}

			blockVol, _ := blockAPI.GetVolume(&block.GetVolumeRequest{
				VolumeID: volume.ID,
				Zone:     volume.Zone,
			})
			if blockVol != nil {
				customVol.Name = blockVol.Name
				customVol.Size = blockVol.Size
				customVol.VolumeType = blockVol.Type
				customVol.State = blockVol.Status.String()
				customVol.CreationDate = blockVol.CreatedAt
				customVol.ModificationDate = blockVol.UpdatedAt
				if blockVol.Specs != nil && blockVol.Specs.PerfIops != nil {
					switch *blockVol.Specs.PerfIops {
					case 5000:
						customVol.IOPS = "5K"
					case 15000:
						customVol.IOPS = "15K"
					}
				}
			} else {
				instanceVol, err := instance.NewAPI(client).GetVolume(&instance.GetVolumeRequest{
					VolumeID: volume.ID,
					Zone:     volume.Zone,
				})
				if err != nil {
					return nil, err
				}
				customVol.Name = instanceVol.Volume.Name
				customVol.Size = instanceVol.Volume.Size
				customVol.VolumeType = instanceVol.Volume.VolumeType.String()
				customVol.State = instanceVol.Volume.State.String()
				customVol.CreationDate = instanceVol.Volume.CreationDate
				customVol.ModificationDate = instanceVol.Volume.ModificationDate
			}

			volumes[volume.ID] = customVol
		}

		// Display warning if server-type is deprecated
		warnings := []string(nil)
		if server.EndOfService {
			warnings = warningServerTypeDeprecated(ctx, client, server)
		}

		return &struct {
			*customServer
			Volumes     []*customVolume
			PrivateNics []customNICs `json:"private_nics"`
			Warnings    []string     `json:"warnings"`
		}{
			customServerFromInstanceServer(server),
			orderVolumes(volumes),
			nics,
			warnings,
		}, nil
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName: "Image",
				Title:     "Server Image",
			},
			{
				FieldName: "AllowedActions",
				Title:     "Allowed Actions",
			},
			{
				FieldName: "Volumes",
				Title:     "Volumes",
			},
			{
				Title:     "Public IPs",
				FieldName: "PublicIPs",
			},
			{
				FieldName: "PrivateNics",
				Title:     "Private NICs",
			},
			{
				FieldName:   "Warnings",
				Title:       "Warnings",
				HideIfEmpty: true,
			},
			{
				FieldName: "Filesystems",
				Title:     "Server Filesystems",
			},
		},
	}

	return c
}

func serverAttachIPCommand() *core.Command {
	type customIPAttachRequest struct {
		OrganizationID *string
		ProjectID      *string
		// Server: UUID of the server you want to attach the IP to
		ServerID string   `json:"server,omitempty"`
		IP       string   `json:"-"`
		Zone     scw.Zone `json:"zone"`
	}

	return &core.Command{
		Short:     `Attach an IP to a server`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "attach-ip",
		ArgsType:  reflect.TypeOf(customIPAttachRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Positional: true,
			},
			{
				Name:     "ip",
				Short:    `UUID of the IP to attach or its UUID`,
				Required: true,
			},
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			api := instance.NewAPI(core.ExtractClient(ctx))
			args := argsI.(*customIPAttachRequest)

			var ipID string
			switch {
			case validation.IsUUID(args.IP):
				ipID = args.IP
			case net.ParseIP(args.IP) != nil:
				// Find the corresponding flexible IP UUID.
				logger.Debugf("finding public IP UUID from address: %s", args.IP)
				res, err := api.GetIP(&instance.GetIPRequest{
					Zone: args.Zone,
					IP:   args.IP,
				})
				if err != nil { // FIXME: isNotFoundError
					return nil, fmt.Errorf("%s does not belong to you", args.IP)
				}
				ipID = res.IP.ID
			default:
				return nil, fmt.Errorf(
					`invalid IP "%s", should be either an IP address ID or a reserved flexible IP address`,
					args.IP,
				)
			}

			_, err = api.UpdateIP(&instance.UpdateIPRequest{
				IP: ipID,
				Server: &instance.NullableStringValue{
					Value: args.ServerID,
				},
				Zone: args.Zone,
			})
			if err != nil {
				return nil, err
			}

			return api.GetServer(&instance.GetServerRequest{ServerID: args.ServerID})
		},
		Examples: []*core.Example{
			{
				Short:    "Attach an IP to a server",
				ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111","ip": "11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Attach an IP to a server",
				ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111","ip": "1.2.3.4"}`,
			},
		},
	}
}

func serverDetachIPCommand() *core.Command {
	type customIPDetachRequest struct {
		OrganizationID *string
		ProjectID      *string
		Zone           scw.Zone `json:"zone"`
		ServerID       string
	}

	return &core.Command{
		Short:     `Detach an IP from a server`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "detach-ip",
		ArgsType:  reflect.TypeOf(customIPDetachRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `UUID of the server.`,
				Required:   true,
				Positional: true,
			},
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			args := argsI.(*customIPDetachRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			serverResponse, err := api.GetServer(
				&instance.GetServerRequest{ServerID: args.ServerID},
			)
			if err != nil {
				return nil, err
			}

			if server := serverResponse.Server; server != nil {
				if ip := server.PublicIP; ip != nil {
					_, err := api.UpdateIP(&instance.UpdateIPRequest{
						Zone: args.Zone,
						// We detach an ip by specifying no serverResponse
						Server: &instance.NullableStringValue{
							Null: true,
						},
						IP: ip.ID,
					})
					if err != nil {
						return nil, err
					}

					return api.GetServer(&instance.GetServerRequest{ServerID: args.ServerID})
				}

				return nil, errors.New("no public ip found")
			}

			return nil, errors.New("no server found")
		},
		Examples: []*core.Example{
			{
				Short:    "Detach IP from a given server",
				ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

type serverWaitRequest struct {
	Zone     scw.Zone
	ServerID string
	Timeout  time.Duration
}

func serverWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for server to reach a stable state`,
		Long:      `Wait for server to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the server.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(serverWaitRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			args := argsI.(*serverWaitRequest)

			return instance.NewAPI(core.ExtractClient(ctx)).
				WaitForServer(&instance.WaitForServerRequest{
					Zone:          args.Zone,
					ServerID:      args.ServerID,
					Timeout:       scw.TimeDurationPtr(args.Timeout),
					RetryInterval: core.DefaultRetryInterval,
				})
		},
		ArgSpecs: core.ArgSpecs{
			core.WaitTimeoutArgSpec(serverActionTimeout),
			{
				Name:       "server-id",
				Short:      `ID of the server affected by the action.`,
				Required:   true,
				Positional: true,
			},
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a server to reach a stable state",
				ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
