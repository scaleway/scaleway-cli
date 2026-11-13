package instance

import (
	"context"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func volumeTypeListBuilder(cmd *core.Command) *core.Command {
	type customVolumeType struct {
		Name    string   `json:"name"`
		MinSize scw.Size `json:"min_size"`
		MaxSize scw.Size `json:"max_size"`
	}

	cmd.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
			res, err := runner(ctx, argsI)
			if err != nil {
				return res, err
			}

			volTypeList, ok := res.(*instance.ListVolumeTypesResponse)
			if !ok {
				return res, fmt.Errorf(
					"expected response of type instance.ListVolumeTypesResponse, got %T",
					res,
				)
			}

			customVolumeTypes := make([]customVolumeType, 0, volTypeList.TotalCount)
			for _, volumeType := range volTypeList.VolumeTypes {
				customVolumeTypes = append(customVolumeTypes, customVolumeType{
					Name:    volumeType.Name.String(),
					MinSize: scw.Size(volumeType.MinSize),
					MaxSize: scw.Size(volumeType.MaxSize),
				})
			}

			return customVolumeTypes, nil
		},
	)

	cmd.AllowAnonymousClient = true

	cmd.View = &core.View{
		Fields: []*core.ViewField{
			{FieldName: "Name", Label: "NAME"},
			{FieldName: "MinSize", Label: "MIN SIZE"},
			{FieldName: "MaxSize", Label: "MAX SIZE"},
		},
	}

	return cmd
}
