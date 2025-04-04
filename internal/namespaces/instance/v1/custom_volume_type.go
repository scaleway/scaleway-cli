package instance

import (
	"context"
	"sort"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

func volumeTypeListBuilder(cmd *core.Command) *core.Command {
	type customVolumeType struct {
		Type string `json:"type"`
		instance.VolumeType
	}

	cmd.AddInterceptors(
		func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
			res, err := runner(ctx, argsI)
			if err != nil {
				return res, err
			}

			volumeTypes := []*customVolumeType(nil)
			for typeName, volumeType := range res.(*instance.ListVolumesTypesResponse).Volumes {
				volumeTypes = append(volumeTypes, &customVolumeType{
					Type:       typeName,
					VolumeType: *volumeType,
				})
			}

			// sort for consistent order output
			sort.Slice(volumeTypes, func(i, j int) bool {
				return volumeTypes[i].Type < volumeTypes[j].Type
			})

			return volumeTypes, nil
		},
	)

	cmd.AllowAnonymousClient = true

	cmd.View = &core.View{
		Fields: []*core.ViewField{
			{FieldName: "Type", Label: "Type"},
			{FieldName: "DisplayName", Label: "Name"},
			{FieldName: "Capabilities.Snapshot", Label: "Snapshot"},
			{FieldName: "Constraints.Min", Label: "Min"},
			{FieldName: "Constraints.Max", Label: "Max"},
		},
	}

	return cmd
}
