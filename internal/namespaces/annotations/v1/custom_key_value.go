package annotations

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	annotationsSDK "github.com/scaleway/scaleway-sdk-go/api/annotations/v1"
)

type keyValueRow struct {
	KeyID   string `json:"key_id"`
	Key     string `json:"key"`
	ValueID string `json:"value_id"`
	Value   string `json:"value"`
}

func keyValueListBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "Key ID",
				FieldName: "KeyID",
			},
			{
				Label:     "Key Name",
				FieldName: "Key",
			},
			{
				Label:     "Value ID",
				FieldName: "ValueID",
			},
			{
				Label:     "Value Name",
				FieldName: "Value",
			},
		},
	}

	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
			originalRes, err := runner(ctx, argsI)
			if err != nil {
				return nil, err
			}

			keyValueResponse := originalRes.(*annotationsSDK.ListAllKeysAndValuesResponse)

			rows := make([]keyValueRow, 0, len(keyValueResponse.Keys))
			for _, key := range keyValueResponse.Keys {
				if len(key.Values) == 0 {
					rows = append(rows, keyValueRow{KeyID: key.ID, Key: key.Name})

					continue
				}
				for _, value := range key.Values {
					rows = append(rows, keyValueRow{
						KeyID:   key.ID,
						Key:     key.Name,
						ValueID: value.ID,
						Value:   value.Name,
					})
				}
			}

			return rows, nil
		},
	)

	return c
}
