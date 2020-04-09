package instance

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/fatih/color"

	"github.com/hashicorp/go-multierror"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
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
		instance.ServerStateRunning:        &human.EnumMarshalSpec{Attribute: color.FgGreen},
		instance.ServerStateStopped:        &human.EnumMarshalSpec{Attribute: color.Faint, Value: "archived"},
		instance.ServerStateStoppedInPlace: &human.EnumMarshalSpec{Attribute: color.Faint},
		instance.ServerStateStarting:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
		instance.ServerStateStopping:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
		instance.ServerStateLocked:         &human.EnumMarshalSpec{Attribute: color.FgRed},
	}
)

// serverLocationMarshalerFunc marshals a instance.ServerLocation.
func serverLocationMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	location := i.(instance.ServerLocation)
	zone, err := scw.ParseZone(location.ZoneID)
	if err != nil {
		return "", err
	}
	return string(zone), nil
}

// serversMarshalerFunc marshals a Server.
func serversMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
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
		PlacementGroup    *instance.PlacementGroup
		ModificationDate  time.Time
		CreationDate      time.Time
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

func getServerResponseMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	serverResponse := i.(instance.GetServerResponse)

	// Sections
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "server",
			Title:     "Server",
		},
		{
			FieldName: "server.image",
			Title:     "Server Image",
		}, {
			FieldName: "server.allowed-actions",
		}, {
			FieldName: "volumes",
			Title:     "Volumes",
		},
	}

	customServer := &struct {
		Server  *instance.Server
		Volumes []*instance.Volume
	}{
		serverResponse.Server,
		orderVolumes(serverResponse.Server.Volumes),
	}

	str, err := human.Marshal(customServer, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

// orderVolumes return an ordered slice based on the volume map key "0", "1", "2",...
func orderVolumes(v map[string]*instance.Volume) []*instance.Volume {
	indexes := []string(nil)
	for index := range v {
		indexes = append(indexes, index)
	}
	sort.Strings(indexes)
	var orderedVolumes []*instance.Volume
	for _, index := range indexes {
		orderedVolumes = append(orderedVolumes, v[index])
	}
	return orderedVolumes
}

// serversMarshalerFunc marshals a BootscriptID.
func bootscriptMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
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
	}

	c.ArgSpecs.GetByName(oldOrganizationFieldName).Name = newOrganizationFieldName

	c.ArgsType = reflect.TypeOf(customListServersRequest{})

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		args := argsI.(*customListServersRequest)

		if args.ListServersRequest == nil {
			args.ListServersRequest = &instance.ListServersRequest{}
		}

		request := args.ListServersRequest
		request.Organization = args.OrganizationID

		return runner(ctx, request)
	})
	return c
}

func serverUpdateBuilder(c *core.Command) *core.Command {
	type instanceUpdateServerRequestCustom struct {
		*instance.UpdateServerRequest
		IP               *instance.NullableStringValue
		PlacementGroupID *instance.NullableStringValue
		SecurityGroupID  *string
		VolumeIDs        *[]string
	}

	c.ArgsType = reflect.TypeOf(instanceUpdateServerRequestCustom{})

	// Rename modified arg specs.
	c.ArgSpecs.GetByName("placement-group").Name = "placement-group-id"
	c.ArgSpecs.GetByName("security-group.id").Name = "security-group-id"

	// Delete unused arg specs.
	c.ArgSpecs.DeleteByName("security-group.name")
	c.ArgSpecs.DeleteByName("volumes.{key}.name")
	c.ArgSpecs.DeleteByName("volumes.{key}.size")
	c.ArgSpecs.DeleteByName("volumes.{key}.id")
	c.ArgSpecs.DeleteByName("volumes.{key}.volume-type")
	c.ArgSpecs.DeleteByName("volumes.{key}.organization")

	// Add new arg specs.
	c.ArgSpecs.AddBefore("placement-group-id", &core.ArgSpec{
		Name:  "volume-ids.{index}",
		Short: "Will update ALL volume IDs at once, including the root volume of the server (use volume-ids=none to detach all volumes)",
	})
	c.ArgSpecs.AddBefore("boot-type", &core.ArgSpec{
		Name:  "ip",
		Short: `IP that should be attached to the server (use ip=none to detach)`,
	})

	c.Run = func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
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
			volumes := make(map[string]*instance.VolumeTemplate)
			for i, volumeID := range *customRequest.VolumeIDs {
				index := strconv.Itoa(i)
				volumes[index] = &instance.VolumeTemplate{ID: volumeID, Name: getServerResponse.Server.Name + "-" + index}
			}
			customRequest.Volumes = &volumes
		}

		updateServerResponse, err := api.UpdateServer(updateServerRequest)
		if err != nil {
			return nil, err
		}

		return updateServerResponse, nil
	}

	return c
}

//
// Commands
//

func serverAttachVolumeCommand() *core.Command {
	return &core.Command{
		Short:     `Attach a volume to a server`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "attach-volume",
		ArgsType:  reflect.TypeOf(instance.AttachVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "server-id",
				Short:    `ID of the server`,
				Required: true,
			},
			{
				Name:     "volume-id",
				Short:    `ID of the volume to attach`,
				Required: true,
			},
			core.ZoneArgSpec(),
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			request := argsI.(*instance.AttachVolumeRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.AttachVolume(request)
		},
		Examples: []*core.Example{
			{
				Short:   "Attach a volume to a server",
				Request: `{"server_id": "11111111-1111-1111-1111-111111111111","volume_id": "22222222-1111-5555-2222-666666111111"}`,
			},
		},
	}
}

func serverDetachVolumeCommand() *core.Command {
	return &core.Command{
		Short:     `Detach a volume from its server`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "detach-volume",
		ArgsType:  reflect.TypeOf(instance.DetachVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "volume-id",
				Short:    `ID of the volume to detach`,
				Required: true,
			},
			core.ZoneArgSpec(),
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			request := argsI.(*instance.DetachVolumeRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.DetachVolume(request)
		},
		Examples: []*core.Example{
			{
				Short:   "Detach a volume from its server",
				Request: `{"volume_id": "22222222-1111-5555-2222-666666111111"}`,
			},
		},
	}
}

type instanceActionRequest struct {
	Zone     scw.Zone
	ServerID string
}

var serverActionArgSpecs = core.ArgSpecs{
	{
		Name:       "server-id",
		Short:      `ID of the server affected by the action.`,
		Required:   true,
		Positional: true,
	},
	core.ZoneArgSpec(),
}

func serverStartCommand() *core.Command {
	return &core.Command{
		Short:     `Power on server`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "start",
		ArgsType:  reflect.TypeOf(instanceActionRequest{}),
		Run:       getRunServerAction(instance.ServerActionPoweron),
		WaitFunc:  waitForServerFunc(),
		ArgSpecs:  serverActionArgSpecs,
		Examples: []*core.Example{
			{
				Short:   "Start a server in the default zone with a given id",
				Request: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:   "Start a server in fr-par-1 zone with a given id",
				Request: `{"zone":"fr-par-1", "server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func serverStopCommand() *core.Command {
	return &core.Command{
		Short:     `Power off server`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "stop",
		ArgsType:  reflect.TypeOf(instanceActionRequest{}),
		Run:       getRunServerAction(instance.ServerActionPoweroff),
		WaitFunc:  waitForServerFunc(),
		ArgSpecs:  serverActionArgSpecs,
		Examples: []*core.Example{
			{
				Short:   "Stop a server in the default zone with a given id",
				Request: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:   "Stop a server in fr-par-1 zone with a given id",
				Request: `{"zone":"fr-par-1", "server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func serverStandbyCommand() *core.Command {
	return &core.Command{
		Short:     `Put server in standby mode`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "standby",
		ArgsType:  reflect.TypeOf(instanceActionRequest{}),
		Run:       getRunServerAction(instance.ServerActionStopInPlace),
		WaitFunc:  waitForServerFunc(),
		ArgSpecs:  serverActionArgSpecs,
		Examples: []*core.Example{
			{
				Short:   "Put in standby a server in the default zone with a given id",
				Request: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:   "Put in standby a server in fr-par-1 zone with a given id",
				Request: `{"zone":"fr-par-1", "server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func serverRebootCommand() *core.Command {
	return &core.Command{
		Short:     `Reboot server`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "reboot",
		ArgsType:  reflect.TypeOf(instanceActionRequest{}),
		Run:       getRunServerAction(instance.ServerActionReboot),
		WaitFunc:  waitForServerFunc(),
		ArgSpecs:  serverActionArgSpecs,
		Examples: []*core.Example{
			{
				Short:   "Reboot a server in the default zone with a given id",
				Request: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:   "Reboot a server in fr-par-1 zone with a given id",
				Request: `{"zone":"fr-par-1", "server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func serverWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for server to reach a stable state`,
		Long:      `Wait for server to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the server.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(instanceActionRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			return waitForServerFunc()(ctx, argsI, nil)
		},
		ArgSpecs: serverActionArgSpecs,
		Examples: []*core.Example{
			{
				Short:   "Wait for a server to reach a stable state",
				Request: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func waitForServerFunc() core.WaitFunc {
	return func(ctx context.Context, argsI, _ interface{}) (interface{}, error) {
		return instance.NewAPI(core.ExtractClient(ctx)).WaitForServer(&instance.WaitForServerRequest{
			Zone:     argsI.(*instanceActionRequest).Zone,
			ServerID: argsI.(*instanceActionRequest).ServerID,
			Timeout:  serverActionTimeout,
		})
	}
}

func getRunServerAction(action instance.ServerAction) core.CommandRunner {
	return func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
		args := argsI.(*instanceActionRequest)

		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)

		_, err := api.ServerAction(&instance.ServerActionRequest{
			Zone:     args.Zone,
			ServerID: args.ServerID,
			Action:   action,
		})
		return &core.SuccessResult{Message: fmt.Sprintf("%s successful for the server", action)}, err
	}
}

type customDeleteServerRequest struct {
	Zone          scw.Zone
	ServerID      string
	WithVolumes   withVolumes
	WithIP        bool
	ForceShutdown bool
}

type withVolumes string

const (
	withVolumesNone  = withVolumes("none")
	withVolumesLocal = withVolumes("local")
	withVolumesBlock = withVolumes("block")
	withVolumesRoot  = withVolumes("root")
	withVolumesAll   = withVolumes("all")
)

func serverDeleteCommand() *core.Command {
	return &core.Command{
		Short:     `Delete server`,
		Long:      `Delete a server with the given ID.`,
		Namespace: "instance",
		Verb:      "delete",
		Resource:  "server",
		ArgsType:  reflect.TypeOf(customDeleteServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Required:   true,
				Positional: true,
			},
			{
				Name:    "with-volumes",
				Short:   "Delete the volumes attached to the server",
				Default: core.DefaultValueSetter("none"),
				EnumValues: []string{
					string(withVolumesNone),
					string(withVolumesLocal),
					string(withVolumesBlock),
					string(withVolumesRoot),
					string(withVolumesAll),
				},
			},
			{
				Name:  "with-ip",
				Short: "Delete the IP attached to the server",
			},
			{
				Name:  "force-shutdown",
				Short: "Force shutdown of the instance server before deleting it",
			},
			core.ZoneArgSpec(),
		},
		Examples: []*core.Example{
			{
				Short:   "Delete a server in the default zone with a given id",
				Request: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:   "Delete a server in fr-par-1 zone with a given id",
				Request: `{"zone":"fr-par-1", "server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw instance server stop",
				Short:   "Stop a running server",
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			deleteServerArgs := argsI.(*customDeleteServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			server, err := api.GetServer(&instance.GetServerRequest{
				Zone:     deleteServerArgs.Zone,
				ServerID: deleteServerArgs.ServerID,
			})
			if err != nil {
				return nil, err
			}

			if deleteServerArgs.ForceShutdown {
				finalStateServer, err := api.WaitForServer(&instance.WaitForServerRequest{
					Zone:     deleteServerArgs.Zone,
					ServerID: deleteServerArgs.ServerID,
				})
				if err != nil {
					return nil, err
				}

				if finalStateServer.State != instance.ServerStateStopped {
					err = api.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
						Zone:     deleteServerArgs.Zone,
						ServerID: deleteServerArgs.ServerID,
						Action:   instance.ServerActionPoweroff,
					})
					if err != nil {
						return nil, err
					}
				}
			}

			err = api.DeleteServer(&instance.DeleteServerRequest{
				Zone:     deleteServerArgs.Zone,
				ServerID: deleteServerArgs.ServerID,
			})
			if err != nil {
				return nil, err
			}

			var multiErr error
			if deleteServerArgs.WithIP && server.Server.PublicIP != nil && !server.Server.PublicIP.Dynamic {
				err = api.DeleteIP(&instance.DeleteIPRequest{
					Zone: deleteServerArgs.Zone,
					IP:   server.Server.PublicIP.ID,
				})
				if err != nil {
					multiErr = multierror.Append(multiErr, err)
				} else {
					_, _ = interactive.Printf("successfully deleted ip %s\n", server.Server.PublicIP.Address.String())
				}
			}

			deletedVolumeMessages := [][2]string(nil)
			for index, volume := range server.Server.Volumes {
				switch {
				case deleteServerArgs.WithVolumes == withVolumesNone:
					break
				case deleteServerArgs.WithVolumes == withVolumesRoot && index != "0":
					continue
				case deleteServerArgs.WithVolumes == withVolumesLocal && volume.VolumeType != instance.VolumeTypeLSSD:
					continue
				case deleteServerArgs.WithVolumes == withVolumesBlock && volume.VolumeType != instance.VolumeTypeBSSD:
					continue
				}
				err = api.DeleteVolume(&instance.DeleteVolumeRequest{
					Zone:     deleteServerArgs.Zone,
					VolumeID: volume.ID,
				})
				if err != nil {
					multiErr = multierror.Append(multiErr, err)
				} else {
					humanSize, err := human.Marshal(volume.Size, nil)
					if err != nil {
						logger.Debugf("cannot marshal human size %v", volume.Size)
					}
					deletedVolumeMessages = append(deletedVolumeMessages, [2]string{
						index,
						fmt.Sprintf("successfully deleted volume %s (%s %s)", volume.Name, humanSize, volume.VolumeType),
					})
				}
			}
			if multiErr != nil {
				return nil, &core.CliError{
					Err:  multiErr,
					Hint: "Make sure these resources have been deleted or try to delete it manually.",
				}
			}

			// Sort and print deleted volume messages
			sort.Slice(deletedVolumeMessages, func(i, j int) bool {
				return deletedVolumeMessages[i][0] < deletedVolumeMessages[j][0]
			})
			for _, message := range deletedVolumeMessages {
				_, _ = interactive.Println(message[1])
			}

			return &core.SuccessResult{}, nil
		},
	}
}
