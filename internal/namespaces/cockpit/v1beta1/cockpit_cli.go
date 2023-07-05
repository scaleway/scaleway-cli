// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package cockpit

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		cockpitRoot(),
		cockpitCockpit(),
		cockpitToken(),
		cockpitGrafanaUser(),
		cockpitPlan(),
		cockpitAlert(),
		cockpitContact(),
		cockpitCockpitActivate(),
		cockpitCockpitGet(),
		cockpitCockpitDeactivate(),
		cockpitCockpitResetGrafana(),
		cockpitTokenCreate(),
		cockpitTokenList(),
		cockpitTokenGet(),
		cockpitTokenDelete(),
		cockpitContactCreate(),
		cockpitContactList(),
		cockpitContactDelete(),
		cockpitAlertEnable(),
		cockpitAlertDisable(),
		cockpitAlertTest(),
		cockpitGrafanaUserCreate(),
		cockpitGrafanaUserList(),
		cockpitGrafanaUserDelete(),
		cockpitGrafanaUserResetPassword(),
		cockpitPlanList(),
		cockpitPlanSelect(),
	)
}
func cockpitRoot() *core.Command {
	return &core.Command{
		Short:     `Cockpit API`,
		Long:      `Cockpit's API allows you to activate your Cockpit on your Projects. Scaleway's Cockpit stores metrics and logs and provides a dedicated Grafana for dashboarding to visualize them.`,
		Namespace: "cockpit",
	}
}

func cockpitCockpit() *core.Command {
	return &core.Command{
		Short:     `Cockpit management commands`,
		Long:      `Cockpit management commands.`,
		Namespace: "cockpit",
		Resource:  "cockpit",
	}
}

func cockpitToken() *core.Command {
	return &core.Command{
		Short:     `Token management commands`,
		Long:      `Token management commands.`,
		Namespace: "cockpit",
		Resource:  "token",
	}
}

func cockpitGrafanaUser() *core.Command {
	return &core.Command{
		Short:     `Grafana user management commands`,
		Long:      `Grafana user management commands.`,
		Namespace: "cockpit",
		Resource:  "grafana-user",
	}
}

func cockpitPlan() *core.Command {
	return &core.Command{
		Short:     `Pricing plans management commands`,
		Long:      `Pricing plans management commands.`,
		Namespace: "cockpit",
		Resource:  "plan",
	}
}

func cockpitAlert() *core.Command {
	return &core.Command{
		Short:     `Managed alerts management commands`,
		Long:      `Managed alerts management commands.`,
		Namespace: "cockpit",
		Resource:  "alert",
	}
}

func cockpitContact() *core.Command {
	return &core.Command{
		Short:     `Contacts management commands`,
		Long:      `Contacts management commands.`,
		Namespace: "cockpit",
		Resource:  "contact",
	}
}

func cockpitCockpitActivate() *core.Command {
	return &core.Command{
		Short:     `Activate a Cockpit`,
		Long:      `Activate the Cockpit of the specified Project ID.`,
		Namespace: "cockpit",
		Resource:  "cockpit",
		Verb:      "activate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.ActivateCockpitRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.ActivateCockpitRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			return api.ActivateCockpit(request)

		},
	}
}

func cockpitCockpitGet() *core.Command {
	return &core.Command{
		Short:     `Get a Cockpit`,
		Long:      `Retrieve the Cockpit of the specified Project ID.`,
		Namespace: "cockpit",
		Resource:  "cockpit",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.GetCockpitRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.GetCockpitRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			return api.GetCockpit(request)

		},
	}
}

func cockpitCockpitDeactivate() *core.Command {
	return &core.Command{
		Short:     `Deactivate a Cockpit`,
		Long:      `Deactivate the Cockpit of the specified Project ID.`,
		Namespace: "cockpit",
		Resource:  "cockpit",
		Verb:      "deactivate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.DeactivateCockpitRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.DeactivateCockpitRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			return api.DeactivateCockpit(request)

		},
	}
}

func cockpitCockpitResetGrafana() *core.Command {
	return &core.Command{
		Short:     `Reset a Grafana`,
		Long:      `Reset your Cockpit's Grafana associated with the specified Project ID.`,
		Namespace: "cockpit",
		Resource:  "cockpit",
		Verb:      "reset-grafana",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.ResetCockpitGrafanaRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.ResetCockpitGrafanaRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			return api.ResetCockpitGrafana(request)

		},
	}
}

func cockpitTokenCreate() *core.Command {
	return &core.Command{
		Short:     `Create a token`,
		Long:      `Create a token associated with the specified Project ID.`,
		Namespace: "cockpit",
		Resource:  "token",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.CreateTokenRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the token`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("token"),
			},
			{
				Name:       "scopes.query-metrics",
				Short:      `Permission to fetch metrics`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scopes.write-metrics",
				Short:      `Permission to write metrics`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scopes.setup-metrics-rules",
				Short:      `Permission to setup metrics rules`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scopes.query-logs",
				Short:      `Permission to fetch logs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scopes.write-logs",
				Short:      `Permission to write logs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scopes.setup-logs-rules",
				Short:      `Permission to setup logs rules`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scopes.setup-alerts",
				Short:      `Permission to setup alerts`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.CreateTokenRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			return api.CreateToken(request)

		},
	}
}

func cockpitTokenList() *core.Command {
	return &core.Command{
		Short:     `List tokens`,
		Long:      `Get a list of tokens associated with the specified Project ID.`,
		Namespace: "cockpit",
		Resource:  "token",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.ListTokensRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.ListTokensRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListTokens(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Tokens, nil

		},
	}
}

func cockpitTokenGet() *core.Command {
	return &core.Command{
		Short:     `Get a token`,
		Long:      `Retrieve the token associated with the specified token ID.`,
		Namespace: "cockpit",
		Resource:  "token",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.GetTokenRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "token-id",
				Short:      `ID of the token`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.GetTokenRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			return api.GetToken(request)

		},
	}
}

func cockpitTokenDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a token`,
		Long:      `Delete the token associated with the specified token ID.`,
		Namespace: "cockpit",
		Resource:  "token",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.DeleteTokenRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "token-id",
				Short:      `ID of the token`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.DeleteTokenRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			e = api.DeleteToken(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "token",
				Verb:     "delete",
			}, nil
		},
	}
}

func cockpitContactCreate() *core.Command {
	return &core.Command{
		Short:     `Create a contact point`,
		Long:      `Create a contact point to receive alerts for the default receiver.`,
		Namespace: "cockpit",
		Resource:  "contact",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.CreateContactPointRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "contact-point.email.to",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.CreateContactPointRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			return api.CreateContactPoint(request)

		},
	}
}

func cockpitContactList() *core.Command {
	return &core.Command{
		Short:     `List contact points`,
		Long:      `Get a list of contact points for the Cockpit associated with the specified Project ID.`,
		Namespace: "cockpit",
		Resource:  "contact",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.ListContactPointsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.ListContactPointsRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListContactPoints(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.ContactPoints, nil

		},
	}
}

func cockpitContactDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an alert contact point`,
		Long:      `Delete a contact point for the default receiver.`,
		Namespace: "cockpit",
		Resource:  "contact",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.DeleteContactPointRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "contact-point.email.to",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.DeleteContactPointRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			e = api.DeleteContactPoint(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "contact",
				Verb:     "delete",
			}, nil
		},
	}
}

func cockpitAlertEnable() *core.Command {
	return &core.Command{
		Short:     `Enable managed alerts`,
		Long:      `Enable the sending of managed alerts for the specified Project's Cockpit.`,
		Namespace: "cockpit",
		Resource:  "alert",
		Verb:      "enable",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.EnableManagedAlertsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.EnableManagedAlertsRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			e = api.EnableManagedAlerts(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "alert",
				Verb:     "enable",
			}, nil
		},
	}
}

func cockpitAlertDisable() *core.Command {
	return &core.Command{
		Short:     `Disable managed alerts`,
		Long:      `Disable the sending of managed alerts for the specified Project's Cockpit.`,
		Namespace: "cockpit",
		Resource:  "alert",
		Verb:      "disable",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.DisableManagedAlertsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.DisableManagedAlertsRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			e = api.DisableManagedAlerts(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "alert",
				Verb:     "disable",
			}, nil
		},
	}
}

func cockpitAlertTest() *core.Command {
	return &core.Command{
		Short:     `Trigger a test alert`,
		Long:      `Trigger a test alert to all of the Cockpit's receivers.`,
		Namespace: "cockpit",
		Resource:  "alert",
		Verb:      "test",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.TriggerTestAlertRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.TriggerTestAlertRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			e = api.TriggerTestAlert(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "alert",
				Verb:     "test",
			}, nil
		},
	}
}

func cockpitGrafanaUserCreate() *core.Command {
	return &core.Command{
		Short:     `Create a Grafana user`,
		Long:      `Create a Grafana user for your Cockpit's Grafana instance. Make sure you save the automatically-generated password and the Grafana user ID.`,
		Namespace: "cockpit",
		Resource:  "grafana-user",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.CreateGrafanaUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "login",
				Short:      `Username of the Grafana user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "role",
				Short:      `Role assigned to the Grafana user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_role", "editor", "viewer"},
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.CreateGrafanaUserRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			return api.CreateGrafanaUser(request)

		},
	}
}

func cockpitGrafanaUserList() *core.Command {
	return &core.Command{
		Short:     `List Grafana users`,
		Long:      `Get a list of Grafana users who are able to connect to the Cockpit's Grafana instance.`,
		Namespace: "cockpit",
		Resource:  "grafana-user",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.ListGrafanaUsersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"login_asc", "login_desc"},
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.ListGrafanaUsersRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListGrafanaUsers(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.GrafanaUsers, nil

		},
	}
}

func cockpitGrafanaUserDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a Grafana user`,
		Long:      `Delete a Grafana user from a Grafana instance, specified by the Cockpit's Project ID and the Grafana user ID.`,
		Namespace: "cockpit",
		Resource:  "grafana-user",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.DeleteGrafanaUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "grafana-user-id",
				Short:      `ID of the Grafana user`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.DeleteGrafanaUserRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			e = api.DeleteGrafanaUser(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "grafana-user",
				Verb:     "delete",
			}, nil
		},
	}
}

func cockpitGrafanaUserResetPassword() *core.Command {
	return &core.Command{
		Short:     `Reset a Grafana user's password`,
		Long:      `Reset a Grafana user's password specified by the Cockpit's Project ID and the Grafana user ID.`,
		Namespace: "cockpit",
		Resource:  "grafana-user",
		Verb:      "reset-password",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.ResetGrafanaUserPasswordRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "grafana-user-id",
				Short:      `ID of the Grafana user`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.ResetGrafanaUserPasswordRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			return api.ResetGrafanaUserPassword(request)

		},
	}
}

func cockpitPlanList() *core.Command {
	return &core.Command{
		Short:     `List pricing plans`,
		Long:      `Get a list of all pricing plans available.`,
		Namespace: "cockpit",
		Resource:  "plan",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.ListPlansRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"name_asc", "name_desc"},
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.ListPlansRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListPlans(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Plans, nil

		},
	}
}

func cockpitPlanSelect() *core.Command {
	return &core.Command{
		Short:     `Select pricing plan`,
		Long:      `Select your chosen pricing plan for your Cockpit, specifying the Cockpit's Project ID and the pricing plan's ID in the request.`,
		Namespace: "cockpit",
		Resource:  "plan",
		Verb:      "select",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(cockpit.SelectPlanRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "plan-id",
				Short:      `ID of the pricing plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*cockpit.SelectPlanRequest)

			client := core.ExtractClient(ctx)
			api := cockpit.NewAPI(client)
			return api.SelectPlan(request)

		},
	}
}
