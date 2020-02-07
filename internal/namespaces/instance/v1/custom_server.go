package instance

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/hashicorp/go-multierror"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-cli/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	serverActionTimeout = 10 * time.Minute
)

//
// Marshalers
//

// serverStateMarshalerFunc marshals a instance.ServerState.
func serverStateMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	// The Scaleway console shows "archived" for a stopped server.
	if i.(instance.ServerState) == instance.ServerStateStopped {
		return terminal.Style("archived", color.Faint), nil
	}
	return human.BindAttributesMarshalFunc(serverStateAttributes)(i, opt)
}

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
		State             instance.ServerState
		Zone              scw.Zone
		PublicIP          net.IP
		PrivateIP         *string
		ImageName         string
		Tags              []string
		ModificationDate  time.Time
		CreationDate      time.Time
		ImageID           string
		Protected         bool
		Volumes           int
		SecurityGroupID   string
		SecurityGroupName string
		StateDetail       string
		Arch              instance.Arch
		PlacementGroup    *instance.PlacementGroup
	}

	servers := i.([]*instance.Server)
	humanServers := make([]*humanServerInList, 0)
	for _, server := range servers {
		var zone scw.Zone
		if server.Location != nil {
			parsedZone, err := scw.ParseZone(server.Location.ZoneID)
			if err != nil {
				return "", err
			}
			zone = parsedZone
		}
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
			State:             server.State,
			Zone:              zone,
			ModificationDate:  server.ModificationDate,
			CreationDate:      server.CreationDate,
			ImageID:           serverImageID,
			ImageName:         serverImageName,
			Protected:         server.Protected,
			PublicIP:          publicIPAddress,
			PrivateIP:         server.PrivateIP,
			Volumes:           len(server.Volumes),
			SecurityGroupID:   server.SecurityGroup.ID,
			SecurityGroupName: server.SecurityGroup.Name,
			StateDetail:       server.StateDetail,
			Arch:              server.Arch,
			PlacementGroup:    server.PlacementGroup,
			Tags:              server.Tags,
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
// Commands
//

type instanceActionRequest struct {
	Zone     scw.Zone
	ServerID string
}

var serverActionArgSpecs = core.ArgSpecs{
	{
		Name:     "server-id",
		Short:    `ID of the server affected by the action.`,
		Required: true,
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
				Request: `{"server_id": "4fbb5119-4542-489c-9443-aa8b574bb6ad"}`,
			},
			{
				Short:   "Start a server in fr-par-1 zone with a given id",
				Request: `{"zone":"fr-par-1", "server_id": "3f38e99f-07bd-458a-9512-b1d15433a102"}`,
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
				Request: `{"server_id": "db21e0d5-f0f0-4535-8815-88b80abac8a9"}`,
			},
			{
				Short:   "Stop a server in fr-par-1 zone with a given id",
				Request: `{"zone":"fr-par-1", "server_id": "1a147062-b8f8-410c-91f7-c150f356e38f"}`,
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
				Request: `{"server_id": "d9439f78-df31-411b-b292-8ef56dd2920c"}`,
			},
			{
				Short:   "Put in standby a server in fr-par-1 zone with a given id",
				Request: `{"zone":"fr-par-1", "server_id": "d36cfc02-5d79-417f-a785-79da256c53d0"}`,
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
				Request: `{"server_id": "887274a9-dd80-41ff-aab9-764e935b77b8"}`,
			},
			{
				Short:   "Reboot a server in fr-par-1 zone with a given id",
				Request: `{"zone":"fr-par-1", "server_id": "bde358bc-298f-43b5-9093-683645fbae95"}`,
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
	DeleteIP      bool
	DeleteVolumes bool
	ForceShutdown bool
}

func serverDeleteCommand() *core.Command {
	return &core.Command{
		Short:     `Delete server`,
		Long:      `Delete a server with the given ID.`,
		Namespace: "instance",
		Verb:      "delete",
		Resource:  "server",
		ArgsType:  reflect.TypeOf(customDeleteServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(),
			{
				Name:     "server-id",
				Required: true,
			},
			{
				Name:  "delete-ip",
				Short: "Delete the IP attached to the server as well",
			},
			{
				Name:  "delete-volumes",
				Short: "Delete the volumes attached to the server as well",
			},
			{
				Name:  "force-shutdown",
				Short: "Force shutdown of the instance server before deleting it",
			},
		},
		Examples: []*core.Example{
			{
				Short:   "Delete a server in the default zone with a given id",
				Request: `{"server_id": "c5a72f0e-d768-414f-95d3-f58a07e83147"}`,
			},
			{
				Short:   "Delete a server in fr-par-1 zone with a given id",
				Request: `{"zone":"fr-par-1", "server_id": "6bc6479d-5cda-4a25-b279-67e44b606eac"}`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw instance server stop",
				Short:   "Stop a running server",
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*customDeleteServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			server, err := api.GetServer(&instance.GetServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			})
			if err != nil {
				return nil, err
			}

			if args.ForceShutdown {
				finalStateServer, err := api.WaitForServer(&instance.WaitForServerRequest{
					Zone:     args.Zone,
					ServerID: args.ServerID,
				})
				if err != nil {
					return nil, err
				}

				if finalStateServer.State != instance.ServerStateStopped {
					err = api.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
						Zone:     args.Zone,
						ServerID: args.ServerID,
						Action:   instance.ServerActionPoweroff,
					})
					if err != nil {
						return nil, err
					}
				}
			}

			err = api.DeleteServer(&instance.DeleteServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			})
			if err != nil {
				return nil, err
			}

			var multiErr error
			if args.DeleteIP && server.Server.PublicIP != nil {
				err = api.DeleteIP(&instance.DeleteIPRequest{
					Zone: args.Zone,
					IP:   server.Server.PublicIP.ID,
				})
				if err != nil {
					multiErr = multierror.Append(multiErr, err)
				}
			}

			if args.DeleteVolumes {
				for _, volume := range server.Server.Volumes {
					err = api.DeleteVolume(&instance.DeleteVolumeRequest{
						Zone:     args.Zone,
						VolumeID: volume.ID,
					})
					if err != nil {
						multiErr = multierror.Append(multiErr, err)
					}
				}
			}
			if multiErr != nil {
				return nil, &core.CliError{
					Err:  multiErr,
					Hint: "Make sure these resources have been deleted or try to delete it manually.",
				}
			}

			return &core.SuccessResult{}, nil
		},
	}
}
