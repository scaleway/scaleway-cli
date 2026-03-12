package domain

import (
	"context"
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Marshalers
//

var taskStatusMarshalSpecs = human.EnumMarshalSpecs{
	domain.TaskStatusSuccess:        &human.EnumMarshalSpec{Attribute: color.FgGreen},
	domain.TaskStatusError:          &human.EnumMarshalSpec{Attribute: color.FgRed},
	domain.TaskStatusPending:        &human.EnumMarshalSpec{Attribute: color.FgBlue},
	domain.TaskStatusNew:            &human.EnumMarshalSpec{Attribute: color.FgCyan},
	domain.TaskStatusWaitingPayment: &human.EnumMarshalSpec{Attribute: color.FgYellow},
}

//
// Commands
//

func dnsTask() *core.Command {
	return &core.Command{
		Short:     `DNS tasks management`,
		Long:      `DNS tasks management.`,
		Namespace: "dns",
		Resource:  "task",
	}
}

func dnsTaskListCommand() *core.Command {
	return &core.Command{
		Short: `List DNS tasks`,
		Long: `Retrieve a list of tasks for domain operations (creation, transfer, renewal, etc.).
This is useful for tracking operations and retrieving task IDs needed for Terraform imports.`,
		Namespace: "dns",
		Resource:  "task",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(domain.RegistrarAPIListTasksRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "organization-id",
				Short:      `Organization ID to filter on`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain",
				Short:      `Domain name to filter on`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "types.{index}",
				Short:      `Task types to filter on`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"create_domain",
					"create_external_domain",
					"renew_domain",
					"transfer_domain",
					"trade_domain",
					"lock_domain_transfer",
					"unlock_domain_transfer",
					"enable_dnssec",
					"disable_dnssec",
					"update_domain",
					"update_contact",
					"delete_domain",
					"cancel_task",
					"generate_ssl_certificate",
					"renew_ssl_certificate",
					"send_message",
					"delete_domain_expired",
					"delete_external_domain",
					"create_host",
					"update_host",
					"delete_host",
					"move_project",
					"transfer_online_domain",
				},
			},
			{
				Name:       "statuses.{index}",
				Short:      `Task statuses to filter on`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unavailable",
					"new",
					"waiting_payment",
					"pending",
					"success",
					"error",
				},
			},
			{
				Name:       "order-by",
				Short:      `Sort order of the returned tasks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("domain_desc"),
				EnumValues: []string{
					"domain_desc",
					"domain_asc",
				},
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIListTasksRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListTasks(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Tasks, nil
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Domain",
			},
			{
				FieldName: "Type",
			},
			{
				FieldName: "Status",
			},
			{
				FieldName: "StartedAt",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "Message",
			},
			{
				FieldName: "ProjectID",
			},
		}},
		Examples: []*core.Example{
			{
				Short:    "List all tasks",
				ArgsJSON: `{}`,
			},
			{
				Short:    "List tasks for a specific domain",
				ArgsJSON: `{"domain": "example.com"}`,
			},
			{
				Short:    "List tasks with create_domain type",
				ArgsJSON: `{"types": ["create_domain"]}`,
			},
		},
	}
}
