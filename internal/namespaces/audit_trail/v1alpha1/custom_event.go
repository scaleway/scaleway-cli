package audit_trail

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	audit_trail "github.com/scaleway/scaleway-sdk-go/api/audit_trail/v1alpha1"
)

func eventListBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "Recorded At",
				FieldName: "RecordedAt",
			},
			{
				Label:     "Name",
				FieldName: "Resource.Name",
			},
			{
				Label:     "StatusCode",
				FieldName: "StatusCode",
			},
			{
				Label:     "MethodName",
				FieldName: "MethodName",
			},
			{
				Label:     "Principal",
				FieldName: "Principal.ID",
			},
			{
				Label:     "SourceIP",
				FieldName: "SourceIP",
			},
			{
				Label:     "ProjectID",
				FieldName: "ProjectID",
			},
			{
				Label:     "ProductName",
				FieldName: "ProductName",
			},
			{
				Label:     "RequestBody",
				FieldName: "RequestBody",
			},
			{
				Label:     "RequestID",
				FieldName: "RequestID",
			},
			{
				Label:     "User-Agent",
				FieldName: "UserAgent",
			},
			{
				Label:     "ID",
				FieldName: "ID",
			},
			{
				Label:     "ServiceName",
				FieldName: "ServiceName",
			},
			{
				Label:     "Type",
				FieldName: "Resource.Type",
			},
			{
				Label:     "ResourceID",
				FieldName: "Resource.ID",
			},
			{
				Label:     "Resource Created At",
				FieldName: "Resource.CreatedAt",
			},
			{
				Label:     "Resource Updated At",
				FieldName: "Resource.UpdatedAt",
			},
		},
	}

	c.AddInterceptors(
		func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
			originalRes, err := runner(ctx, argsI)
			if err != nil {
				return nil, err
			}

			eventsResponse := originalRes.(*audit_trail.ListEventsResponse)

			return eventsResponse.Events, nil
		},
	)

	return c
}
