package instance

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	block "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

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
				Default: core.DefaultValueSetter("all"),
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
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
		Examples: []*core.Example{
			{
				Short:    "Delete a server in the default zone with a given id",
				ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Delete a server in fr-par-1 zone with a given id",
				ArgsJSON: `{"zone":"fr-par-1", "server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw instance server terminate",
				Short:   "Terminate a running server",
			},
			{
				Command: "scw instance server stop",
				Short:   "Stop a running server",
			},
		},
		WaitUsage: "wait until the server and its resources are deleted",
		WaitFunc: func(ctx context.Context, _, respI any) (any, error) {
			server := respI.(*core.SuccessResult).TargetResource.(*instance.Server)
			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			notFoundErr := &scw.ResourceNotFoundError{}

			_, err := api.WaitForServer(&instance.WaitForServerRequest{
				Zone:          server.Zone,
				ServerID:      server.ID,
				Timeout:       scw.TimeDurationPtr(serverActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
			if err != nil {
				err = errors.Unwrap(err)
				if !errors.As(err, &notFoundErr) {
					return nil, err
				}
			}

			return respI, nil
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
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
					Zone:          deleteServerArgs.Zone,
					ServerID:      deleteServerArgs.ServerID,
					Timeout:       scw.TimeDurationPtr(serverActionTimeout),
					RetryInterval: core.DefaultRetryInterval,
				})
				if err != nil {
					return nil, err
				}

				if finalStateServer.State != instance.ServerStateStopped {
					err = api.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
						Zone:          deleteServerArgs.Zone,
						ServerID:      deleteServerArgs.ServerID,
						Action:        instance.ServerActionPoweroff,
						Timeout:       scw.TimeDurationPtr(serverActionTimeout),
						RetryInterval: core.DefaultRetryInterval,
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

			if deleteServerArgs.WithIP && server.Server.PublicIPs != nil {
				for _, ip := range server.Server.PublicIPs {
					if ip.Dynamic {
						continue
					}
					err = api.DeleteIP(&instance.DeleteIPRequest{
						Zone: deleteServerArgs.Zone,
						IP:   ip.ID,
					})
					if err != nil {
						return nil, err
					}
					_, _ = interactive.Printf("successfully deleted ip %s\n", ip.Address.String())
				}
			}

			deletedVolumeMessages := [][2]string(nil)
		volumeDelete:
			for index, volume := range server.Server.Volumes {
				switch {
				case deleteServerArgs.WithVolumes == withVolumesNone:
					break volumeDelete
				case deleteServerArgs.WithVolumes == withVolumesRoot && index != "0":
					continue
				case deleteServerArgs.WithVolumes == withVolumesLocal && volume.VolumeType != instance.VolumeServerVolumeTypeLSSD:
					continue
				case deleteServerArgs.WithVolumes == withVolumesBlock && volume.VolumeType != instance.VolumeServerVolumeTypeBSSD && volume.VolumeType != instance.VolumeServerVolumeTypeSbsVolume:
					continue
				case volume.VolumeType == instance.VolumeServerVolumeTypeScratch:
					continue
				}
				err := serverDeleteVolume(volume, api, block.NewAPI(client))
				if err != nil {
					return nil, err
				}
				humanSize, err := human.Marshal(volume.Size, nil)
				if err != nil {
					logger.Debugf("cannot marshal human size %v", volume.Size)
				}
				volumeName := ""
				if volume.Name != nil {
					volumeName = *volume.Name
				}
				deletedVolumeMessages = append(deletedVolumeMessages, [2]string{
					index,
					fmt.Sprintf("successfully deleted volume %s (%s %s)", volumeName, humanSize, volume.VolumeType),
				})
			}

			// Sort and print deleted volume messages
			sort.Slice(deletedVolumeMessages, func(i, j int) bool {
				return deletedVolumeMessages[i][0] < deletedVolumeMessages[j][0]
			})
			for _, message := range deletedVolumeMessages {
				_, _ = interactive.Println(message[1])
			}

			return &core.SuccessResult{
				TargetResource: server.Server,
			}, nil
		},
	}
}

func errorDeletingResource(err error) error {
	return &core.CliError{
		Err:  err,
		Hint: "Make sure this resource have been deleted or try to delete it manually.",
	}
}

func serverDeleteVolume(
	volume *instance.VolumeServer,
	instanceAPI *instance.API,
	blockAPI *block.API,
) error {
	var err error

	if volume.VolumeType == instance.VolumeServerVolumeTypeSbsVolume {
		volumeAvailable := block.VolumeStatusAvailable
		_, err = blockAPI.WaitForVolumeAndReferences(&block.WaitForVolumeAndReferencesRequest{
			Zone:                 volume.Zone,
			VolumeID:             volume.ID,
			VolumeTerminalStatus: &volumeAvailable,
		})
		if err != nil {
			return errorDeletingResource(err)
		}

		err = blockAPI.DeleteVolume(&block.DeleteVolumeRequest{
			Zone:     volume.Zone,
			VolumeID: volume.ID,
		})
	} else {
		_, err = instanceAPI.WaitForVolume(&instance.WaitForVolumeRequest{
			VolumeID: volume.ID,
			Zone:     volume.Zone,
		})
		if err != nil {
			return errorDeletingResource(err)
		}
		err = instanceAPI.DeleteVolume(&instance.DeleteVolumeRequest{
			Zone:     volume.Zone,
			VolumeID: volume.ID,
		})
	}
	if err != nil {
		return errorDeletingResource(err)
	}

	return nil
}
