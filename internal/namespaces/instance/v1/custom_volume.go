package instance

import (
	"context"
	"reflect"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Marshalers
//

var volumeStateMarshalSpecs = human.EnumMarshalSpecs{
	instance.VolumeStateAvailable:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
	instance.VolumeStateError:        &human.EnumMarshalSpec{Attribute: color.FgRed},
	instance.VolumeStateFetching:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.VolumeStateHotsyncing:   &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.VolumeStateResizing:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.VolumeStateSaving:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.VolumeStateSnapshotting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
}

// serversMarshalerFunc marshals a VolumeSummary.
func volumeSummaryMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	volumeSummary := i.(instance.VolumeSummary)

	return human.Marshal(volumeSummary.ID, opt)
}

// volumeMapMarshalerFunc returns the length of the map.
func volumeMapMarshalerFunc(i interface{}, _ *human.MarshalOpt) (string, error) {
	volumes := i.(map[string]*instance.Volume)

	return strconv.Itoa(len(volumes)), nil
}

// Builders

func volumeCreateBuilder(c *core.Command) *core.Command {
	type customCreateVolumeRequest struct {
		*instance.CreateVolumeRequest
		OrganizationID *string
		ProjectID      *string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customCreateVolumeRequest{})

	c.AddInterceptors(
		func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
			args := argsI.(*customCreateVolumeRequest)

			if args.CreateVolumeRequest == nil {
				args.CreateVolumeRequest = &instance.CreateVolumeRequest{}
			}

			request := args.CreateVolumeRequest
			request.Organization = args.OrganizationID
			request.Project = args.ProjectID

			return runner(ctx, request)
		},
	)

	return c
}

func volumeListBuilder(c *core.Command) *core.Command {
	type customListVolumesRequest struct {
		*instance.ListVolumesRequest
		OrganizationID *string
		ProjectID      *string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customListVolumesRequest{})

	c.AddInterceptors(
		func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
			args := argsI.(*customListVolumesRequest)

			if args.ListVolumesRequest == nil {
				args.ListVolumesRequest = &instance.ListVolumesRequest{}
			}

			request := args.ListVolumesRequest
			request.Organization = args.OrganizationID
			request.Project = args.ProjectID

			return runner(ctx, request)
		},
	)

	return c
}

type volumeWaitRequest struct {
	Zone     scw.Zone
	VolumeID string
	Timeout  time.Duration
}

func volumeWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for volume to reach a stable state`,
		Long:      `Wait for volume to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the volume.`,
		Namespace: "instance",
		Resource:  "volume",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(volumeWaitRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			args := argsI.(*volumeWaitRequest)

			return instance.NewAPI(core.ExtractClient(ctx)).
				WaitForVolume(&instance.WaitForVolumeRequest{
					Zone:          args.Zone,
					VolumeID:      args.VolumeID,
					Timeout:       scw.TimeDurationPtr(args.Timeout),
					RetryInterval: core.DefaultRetryInterval,
				})
		},
		ArgSpecs: core.ArgSpecs{
			core.WaitTimeoutArgSpec(serverActionTimeout),
			{
				Name:       "volume-id",
				Short:      `ID of the volume affected by the action.`,
				Required:   true,
				Positional: true,
			},
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a volume to reach a stable state",
				ArgsJSON: `{"volume_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func volumeMigrationBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.DeleteByName("snapshot-id")
	c.ArgSpecs.GetByName("volume-id").Positional = true

	return c
}
