package instance

import (
	"context"
	"fmt"
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Marshalers
//

var (
	volumeStateMarshalSpecs = human.EnumMarshalSpecs{
		instance.VolumeStateAvailable:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		instance.VolumeStateError:        &human.EnumMarshalSpec{Attribute: color.FgRed},
		instance.VolumeStateFetching:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
		instance.VolumeStateHotsyncing:   &human.EnumMarshalSpec{Attribute: color.FgBlue},
		instance.VolumeStateResizing:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
		instance.VolumeStateSaving:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
		instance.VolumeStateSnapshotting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
	}
)

// serversMarshalerFunc marshals a VolumeSummary.
func volumeSummaryMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	volumeSummary := i.(instance.VolumeSummary)
	return human.Marshal(volumeSummary.ID, opt)
}

// volumeMapMarshalerFunc returns the length of the map.
func volumeMapMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	volumes := i.(map[string]*instance.Volume)
	return fmt.Sprintf("%v", len(volumes)), nil
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

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		args := argsI.(*customCreateVolumeRequest)

		if args.CreateVolumeRequest == nil {
			args.CreateVolumeRequest = &instance.CreateVolumeRequest{}
		}

		request := args.CreateVolumeRequest
		request.Organization = args.OrganizationID
		request.Project = args.ProjectID

		return runner(ctx, request)
	})
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

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		args := argsI.(*customListVolumesRequest)

		if args.ListVolumesRequest == nil {
			args.ListVolumesRequest = &instance.ListVolumesRequest{}
		}

		request := args.ListVolumesRequest
		request.Organization = args.OrganizationID
		request.Project = args.ProjectID

		return runner(ctx, request)
	})
	return c
}
