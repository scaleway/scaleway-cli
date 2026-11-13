package instance

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	serverActionTimeout = 10 * time.Minute
	volumeCreateTimeout = 10 * time.Minute
)

//
// Marshalers
//

var volumeStatusMarshalSpecs = human.EnumMarshalSpecs{
	instance.VolumeStatusAvailable:     &human.EnumMarshalSpec{Attribute: color.FgGreen},
	instance.VolumeStatusSnapshotting:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.VolumeStatusAttaching:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.VolumeStatusDetaching:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.VolumeStatusCreating:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.VolumeStatusMigrating:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.VolumeStatusError:         &human.EnumMarshalSpec{Attribute: color.FgRed},
	instance.VolumeStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
}

// serversMarshalerFunc marshals a VolumeSummary.
//func volumeSummaryMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
//	volumeSummary := i.(instance.VolumeSummary)
//
//	return human.Marshal(volumeSummary.ID, opt)
//}

// volumeMapMarshalerFunc returns the length of the map.
func volumeMapMarshalerFunc(i any, _ *human.MarshalOpt) (string, error) {
	volumes := i.(map[string]*instance.Volume)

	return strconv.Itoa(len(volumes)), nil
}

// Builders

func volumeCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("name").Default = core.RandomValueGenerator("vol")

	type customCreateVolumeRequest struct {
		*instance.VolumeAPICreateVolumeRequest
		BaseSnapshot *string `json:"base_snapshot"`
	}

	c.ArgsType = reflect.TypeFor[customCreateVolumeRequest]()

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
		args := argsI.(*customCreateVolumeRequest)

		if args.VolumeAPICreateVolumeRequest == nil {
			args.VolumeAPICreateVolumeRequest = &instance.VolumeAPICreateVolumeRequest{}
		}

		// If no 'base-snapshot-id' arg was passed, try deprecated 'base-snapshot' arg instead.
		request := args.VolumeAPICreateVolumeRequest
		if request.BaseSnapshotID == nil {
			request.BaseSnapshotID = args.BaseSnapshot
		}

		return runner(ctx, request)
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI any) (any, error) {
		api := instance.NewVolumeAPI(core.ExtractClient(ctx))

		return api.WaitForVolume(&instance.WaitForVolumeRequest{
			VolumeID:      respI.(*instance.Volume).ID,
			Zone:          argsI.(*customCreateVolumeRequest).Zone,
			Timeout:       new(volumeCreateTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func volumeListBuilder(c *core.Command) *core.Command {
	//	type customListVolumesRequest struct {
	//		*instance.ListVolumesRequest
	//		OrganizationID *string
	//		ProjectID      *string
	//	}
	//
	//	//renameOrganizationIDArgSpec(c.ArgSpecs)
	//	//renameProjectIDArgSpec(c.ArgSpecs)
	//
	//	c.ArgsType = reflect.TypeFor[customListVolumesRequest]()
	//
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}

		volList, ok := rawResp.(*instance.ListVolumesResponse)
		if !ok {
			return rawResp, fmt.Errorf(
				"expected response of type instance.ListVolumesResponse, got %T",
				rawResp,
			)
		}

		return volList.Volumes, nil
	}

	c.View = &core.View{
		Fields: []*core.ViewField{
			{Label: "ID", FieldName: "ID"},
			{Label: "NAME", FieldName: "Name"},
			{Label: "SIZE", FieldName: "Size"},
			{Label: "BASE SNAPSHOT ID", FieldName: "BaseSnapshotID"},
			{Label: "VOLUME TYPE", FieldName: "VolumeType"},
			{Label: "STATUS", FieldName: "Status"},
			{Label: "TAGS", FieldName: "Tags"},
			{Label: "SERVER ID", FieldName: "ServerID"},
			{Label: "ZONE", FieldName: "Zone"},
			{Label: "PROJECT ID", FieldName: "ProjectID"},
			{Label: "CREATED AT", FieldName: "CreatedAt"},
			{Label: "UPDATED AT", FieldName: "UpdatedAt"},
			{Label: "SRN", FieldName: "Srn"},
		},
	}

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
		ArgsType:  reflect.TypeFor[volumeWaitRequest](),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			args := argsI.(*volumeWaitRequest)

			return instance.NewVolumeAPI(core.ExtractClient(ctx)).
				WaitForVolume(&instance.WaitForVolumeRequest{
					Zone:          args.Zone,
					VolumeID:      args.VolumeID,
					Timeout:       new(args.Timeout),
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

func volumeDeleteBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("volume-id").Positional = true

	return c
}
