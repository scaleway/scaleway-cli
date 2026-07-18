// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package billing

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/billing/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		billingRoot(),
		billingBudget(),
		billingBudgetAlert(),
		billingBudgetAlertNotification(),
		billingBudgetList(),
		billingBudgetGet(),
		billingBudgetCreate(),
		billingBudgetUpdate(),
		billingBudgetDelete(),
		billingBudgetAlertCreate(),
		billingBudgetAlertUpdate(),
		billingBudgetAlertDelete(),
		billingBudgetAlertNotificationCreate(),
		billingBudgetAlertNotificationUpdate(),
		billingBudgetAlertNotificationDelete(),
	)
}

func billingRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to query billing related objects`,
		Long:      `This API allows you to query billing related objects.`,
		Namespace: "billing",
	}
}

func billingBudget() *core.Command {
	return &core.Command{
		Short:     `Budget management commands`,
		Long:      `Budget management commands.`,
		Namespace: "billing",
		Resource:  "budget",
	}
}

func billingBudgetAlert() *core.Command {
	return &core.Command{
		Short:     `Budget alerts management commands`,
		Long:      `Budget alerts management commands.`,
		Namespace: "billing",
		Resource:  "budget-alert",
	}
}

func billingBudgetAlertNotification() *core.Command {
	return &core.Command{
		Short:     `Budget alert notification management commands`,
		Long:      `Budget alert notification management commands.`,
		Namespace: "billing",
		Resource:  "budget-alert-notification",
	}
}

func billingBudgetList() *core.Command {
	return &core.Command{
		Short:     `List your budgets, filtering by ` + "`" + `organization_id` + "`" + `.`,
		Long:      `List your budgets, filtering by ` + "`" + `organization_id` + "`" + `.`,
		Namespace: "billing",
		Resource:  "budget",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.ListBudgetsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "organization-id",
				Short:      `Filter by organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*billing.ListBudgetsRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListBudgets(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Budgets, nil
		},
	}
}

func billingBudgetGet() *core.Command {
	return &core.Command{
		Short:     `Fetch a budget.`,
		Long:      `Fetch a budget.`,
		Namespace: "billing",
		Resource:  "budget",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.GetBudgetRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "budget-id",
				Short:      `The ID of the budget`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*billing.GetBudgetRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)

			return api.GetBudget(request)
		},
	}
}

func billingBudgetCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new budget.`,
		Long:      `Create a new budget.`,
		Namespace: "billing",
		Resource:  "budget",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.CreateBudgetRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "consumption-limit",
				Short:      `Cost limit for the budget`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enabled",
				Short:      `Whether the budget is enabled or not`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*billing.CreateBudgetRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)

			return api.CreateBudget(request)
		},
	}
}

func billingBudgetUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a budget.`,
		Long:      `Update a budget.`,
		Namespace: "billing",
		Resource:  "budget",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.UpdateBudgetRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "budget-id",
				Short:      `The ID of the budget to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "consumption-limit",
				Short:      `Cost limit for the budget`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enabled",
				Short:      `Whether the budget will be enabled or not`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*billing.UpdateBudgetRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)

			return api.UpdateBudget(request)
		},
	}
}

func billingBudgetDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a budget.`,
		Long:      `Delete a budget.`,
		Namespace: "billing",
		Resource:  "budget",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.DeleteBudgetRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "budget-id",
				Short:      `The ID of the budget to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*billing.DeleteBudgetRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)
			e = api.DeleteBudget(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "budget",
				Verb:     "delete",
			}, nil
		},
	}
}

func billingBudgetAlertCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new budget alert.`,
		Long:      `Create a new budget alert.`,
		Namespace: "billing",
		Resource:  "budget-alert",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.CreateBudgetAlertRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "budget-id",
				Short:      `The ID of the budget to create alert for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "threshold",
				Short:      `Threshold above which the alert is sent`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*billing.CreateBudgetAlertRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)

			return api.CreateBudgetAlert(request)
		},
	}
}

func billingBudgetAlertUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a budget alert.`,
		Long:      `Update a budget alert.`,
		Namespace: "billing",
		Resource:  "budget-alert",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.UpdateBudgetAlertRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "budget-alert-id",
				Short:      `The ID of the budget alert to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "threshold",
				Short:      `Threshold above which the alert is sent`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*billing.UpdateBudgetAlertRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)

			return api.UpdateBudgetAlert(request)
		},
	}
}

func billingBudgetAlertDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a budget alert.`,
		Long:      `Delete a budget alert.`,
		Namespace: "billing",
		Resource:  "budget-alert",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.DeleteBudgetAlertRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "budget-alert-id",
				Short:      `The ID of the budget alert to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*billing.DeleteBudgetAlertRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)
			e = api.DeleteBudgetAlert(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "budget-alert",
				Verb:     "delete",
			}, nil
		},
	}
}

func billingBudgetAlertNotificationCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new budget alert notification.`,
		Long:      `Create a new budget alert notification.`,
		Namespace: "billing",
		Resource:  "budget-alert-notification",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.CreateBudgetAlertNotificationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "budget-alert-id",
				Short:      `The ID of the budget alert to create notification for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sms-phone-numbers.{index}",
				Short:      `List of phone numbers to receive sms notifications`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "email-addresses.{index}",
				Short:      `List of email addresses to receive email notifications`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "webhook-urls.{index}",
				Short:      `List of webhook url to receive webhook notifications`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*billing.CreateBudgetAlertNotificationRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)

			return api.CreateBudgetAlertNotification(request)
		},
	}
}

func billingBudgetAlertNotificationUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a budget alert notification.`,
		Long:      `Update a budget alert notification.`,
		Namespace: "billing",
		Resource:  "budget-alert-notification",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.UpdateBudgetAlertNotificationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "budget-alert-notification-id",
				Short:      `The ID of the budget alert notification to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sms-phone-numbers.{index}",
				Short:      `List of phone numbers to receive sms notifications`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "email-addresses.{index}",
				Short:      `List of email addresses to receive email notifications`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "webhook-urls.{index}",
				Short:      `List of webhook url to receive webhook notifications`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*billing.UpdateBudgetAlertNotificationRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)

			return api.UpdateBudgetAlertNotification(request)
		},
	}
}

func billingBudgetAlertNotificationDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a budget alert notification.`,
		Long:      `Delete a budget alert notification.`,
		Namespace: "billing",
		Resource:  "budget-alert-notification",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(billing.DeleteBudgetAlertNotificationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "budget-alert-notification-id",
				Short:      `The ID of the budget alert notification to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*billing.DeleteBudgetAlertNotificationRequest)

			client := core.ExtractClient(ctx)
			api := billing.NewAPI(client)
			e = api.DeleteBudgetAlertNotification(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "budget-alert-notification",
				Verb:     "delete",
			}, nil
		},
	}
}
