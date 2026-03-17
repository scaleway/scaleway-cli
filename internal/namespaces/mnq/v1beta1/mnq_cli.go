// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package mnq

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	mnq "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		mnqRoot(),
		mnqNats(),
		mnqSns(),
		mnqSqs(),
		mnqNatsCreateAccount(),
		mnqNatsDeleteAccount(),
		mnqNatsUpdateAccount(),
		mnqNatsGetAccount(),
		mnqNatsListAccounts(),
		mnqNatsCreateCredentials(),
		mnqNatsDeleteCredentials(),
		mnqNatsGetCredentials(),
		mnqNatsListCredentials(),
		mnqSnsActivate(),
		mnqSnsGetInfo(),
		mnqSnsDeactivate(),
		mnqSnsCreateCredentials(),
		mnqSnsDeleteCredentials(),
		mnqSnsUpdateCredentials(),
		mnqSnsGetCredentials(),
		mnqSnsListCredentials(),
		mnqSqsActivate(),
		mnqSqsGetInfo(),
		mnqSqsDeactivate(),
		mnqSqsCreateCredentials(),
		mnqSqsDeleteCredentials(),
		mnqSqsUpdateCredentials(),
		mnqSqsGetCredentials(),
		mnqSqsListCredentials(),
	)
}

func mnqRoot() *core.Command {
	return &core.Command{
		Short:     `These APIs allow you to manage your Messaging and Queuing NATS, Queues and Topics and Events services`,
		Long:      `These APIs allow you to manage your Messaging and Queuing NATS, Queues and Topics and Events services.`,
		Namespace: "mnq",
	}
}

func mnqNats() *core.Command {
	return &core.Command{
		Short:     `MnQ NATS commands`,
		Long:      `MnQ NATS commands.`,
		Namespace: "mnq",
		Resource:  "nats",
	}
}

func mnqSns() *core.Command {
	return &core.Command{
		Short:     `MnQ Topics and Events commands`,
		Long:      `MnQ Topics and Events commands.`,
		Namespace: "mnq",
		Resource:  "sns",
	}
}

func mnqSqs() *core.Command {
	return &core.Command{
		Short:     `MnQ Queues commands`,
		Long:      `MnQ Queues commands.`,
		Namespace: "mnq",
		Resource:  "sqs",
	}
}

func mnqNatsCreateAccount() *core.Command {
	return &core.Command{
		Short:     `Create a NATS account`,
		Long:      `Create a NATS account associated with a Project.`,
		Namespace: "mnq",
		Resource:  "nats",
		Verb:      "create-account",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.NatsAPICreateNatsAccountRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `NATS account name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("mnq"),
			},
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.NatsAPICreateNatsAccountRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)

			return api.CreateNatsAccount(request)
		},
	}
}

func mnqNatsDeleteAccount() *core.Command {
	return &core.Command{
		Short:     `Delete a NATS account`,
		Long:      `Delete a NATS account, specified by its NATS account ID. Note that deleting a NATS account is irreversible, and any credentials, streams, consumer and stored messages belonging to this NATS account will also be deleted.`,
		Namespace: "mnq",
		Resource:  "nats",
		Verb:      "delete-account",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.NatsAPIDeleteNatsAccountRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "nats-account-id",
				Short:      `ID of the NATS account to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.NatsAPIDeleteNatsAccountRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)
			e = api.DeleteNatsAccount(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "nats",
				Verb:     "delete-account",
			}, nil
		},
	}
}

func mnqNatsUpdateAccount() *core.Command {
	return &core.Command{
		Short:     `Update the name of a NATS account`,
		Long:      `Update the name of a NATS account, specified by its NATS account ID.`,
		Namespace: "mnq",
		Resource:  "nats",
		Verb:      "update-account",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.NatsAPIUpdateNatsAccountRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "nats-account-id",
				Short:      `ID of the NATS account to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `NATS account name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.NatsAPIUpdateNatsAccountRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)

			return api.UpdateNatsAccount(request)
		},
	}
}

func mnqNatsGetAccount() *core.Command {
	return &core.Command{
		Short:     `Get a NATS account`,
		Long:      `Retrieve information about an existing NATS account identified by its NATS account ID. Its full details, including name and endpoint, are returned in the response.`,
		Namespace: "mnq",
		Resource:  "nats",
		Verb:      "get-account",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.NatsAPIGetNatsAccountRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "nats-account-id",
				Short:      `ID of the NATS account to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.NatsAPIGetNatsAccountRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)

			return api.GetNatsAccount(request)
		},
	}
}

func mnqNatsListAccounts() *core.Command {
	return &core.Command{
		Short:     `List NATS accounts`,
		Long:      `List all NATS accounts in the specified region, for a Scaleway Organization or Project. By default, the NATS accounts returned in the list are ordered by creation date in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "mnq",
		Resource:  "nats",
		Verb:      "list-accounts",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.NatsAPIListNatsAccountsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Include only NATS accounts in this Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"updated_at_asc",
					"updated_at_desc",
					"name_asc",
					"name_desc",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.NatsAPIListNatsAccountsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListNatsAccounts(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.NatsAccounts, nil
		},
	}
}

func mnqNatsCreateCredentials() *core.Command {
	return &core.Command{
		Short:     `Create NATS credentials`,
		Long:      `Create a set of credentials for a NATS account, specified by its NATS account ID.`,
		Namespace: "mnq",
		Resource:  "nats",
		Verb:      "create-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.NatsAPICreateNatsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "nats-account-id",
				Short:      `NATS account containing the credentials`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the credentials`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("mnq"),
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.NatsAPICreateNatsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)

			return api.CreateNatsCredentials(request)
		},
	}
}

func mnqNatsDeleteCredentials() *core.Command {
	return &core.Command{
		Short:     `Delete NATS credentials`,
		Long:      `Delete a set of credentials, specified by their credentials ID. Deleting credentials is irreversible and cannot be undone. The credentials can no longer be used to access the NATS account, and active connections using this credentials will be closed.`,
		Namespace: "mnq",
		Resource:  "nats",
		Verb:      "delete-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.NatsAPIDeleteNatsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "nats-credentials-id",
				Short:      `ID of the credentials to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.NatsAPIDeleteNatsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)
			e = api.DeleteNatsCredentials(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "nats",
				Verb:     "delete-credentials",
			}, nil
		},
	}
}

func mnqNatsGetCredentials() *core.Command {
	return &core.Command{
		Short:     `Get NATS credentials`,
		Long:      `Retrieve an existing set of credentials, identified by the ` + "`" + `nats_credentials_id` + "`" + `. The credentials themselves are NOT returned, only their metadata (NATS account ID, credentials name, etc), are returned in the response.`,
		Namespace: "mnq",
		Resource:  "nats",
		Verb:      "get-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.NatsAPIGetNatsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "nats-credentials-id",
				Short:      `ID of the credentials to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.NatsAPIGetNatsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)

			return api.GetNatsCredentials(request)
		},
	}
}

func mnqNatsListCredentials() *core.Command {
	return &core.Command{
		Short:     `List NATS credentials`,
		Long:      `List existing credentials in the specified NATS account. The response contains only the metadata for the credentials, not the credentials themselves, which are only returned after a **Create Credentials** call.`,
		Namespace: "mnq",
		Resource:  "nats",
		Verb:      "list-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.NatsAPIListNatsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Include only NATS accounts in this Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nats-account-id",
				Short:      `Include only credentials for this NATS account`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"updated_at_asc",
					"updated_at_desc",
					"name_asc",
					"name_desc",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.NatsAPIListNatsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListNatsCredentials(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.NatsCredentials, nil
		},
	}
}

func mnqSnsActivate() *core.Command {
	return &core.Command{
		Short:     `Activate Topics and Events`,
		Long:      `Activate Topics and Events for the specified Project ID. Topics and Events must be activated before any usage. Activating Topics and Events does not trigger any billing, and you can deactivate at any time.`,
		Namespace: "mnq",
		Resource:  "sns",
		Verb:      "activate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIActivateSnsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SnsAPIActivateSnsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)

			return api.ActivateSns(request)
		},
	}
}

func mnqSnsGetInfo() *core.Command {
	return &core.Command{
		Short:     `Get Topics and Events info`,
		Long:      `Retrieve the Topics and Events information of the specified Project ID. information include the activation status and the Topics and Events API endpoint URL.`,
		Namespace: "mnq",
		Resource:  "sns",
		Verb:      "get-info",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIGetSnsInfoRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SnsAPIGetSnsInfoRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)

			return api.GetSnsInfo(request)
		},
	}
}

func mnqSnsDeactivate() *core.Command {
	return &core.Command{
		Short:     `Deactivate Topics and Events`,
		Long:      `Deactivate Topics and Events for the specified Project ID. You must delete all topics and credentials before this call or you need to set the force_delete parameter.`,
		Namespace: "mnq",
		Resource:  "sns",
		Verb:      "deactivate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIDeactivateSnsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SnsAPIDeactivateSnsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)

			return api.DeactivateSns(request)
		},
	}
}

func mnqSnsCreateCredentials() *core.Command {
	return &core.Command{
		Short:     `Create Topics and Events credentials`,
		Long:      `Create a set of credentials for Topics and Events, specified by a Project ID. Credentials give the bearer access to topics, and the level of permissions can be defined granularly.`,
		Namespace: "mnq",
		Resource:  "sns",
		Verb:      "create-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPICreateSnsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the credentials`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("mnq_sns"),
			},
			{
				Name:       "permissions.can-publish",
				Short:      `Defines whether the credentials bearer can publish messages to the service (publish to Topics and Events topics)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-receive",
				Short:      `Defines whether the credentials bearer can receive messages from the service (configure subscriptions)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-manage",
				Short:      `Defines whether the credentials bearer can manage the associated Topics and Events topics or subscriptions`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SnsAPICreateSnsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)

			return api.CreateSnsCredentials(request)
		},
	}
}

func mnqSnsDeleteCredentials() *core.Command {
	return &core.Command{
		Short:     `Delete Topics and Events credentials`,
		Long:      `Delete a set of Topics and Events credentials, specified by their credentials ID. Deleting credentials is irreversible and cannot be undone. The credentials can then no longer be used to access Topics and Events.`,
		Namespace: "mnq",
		Resource:  "sns",
		Verb:      "delete-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIDeleteSnsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "sns-credentials-id",
				Short:      `ID of the credentials to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SnsAPIDeleteSnsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)
			e = api.DeleteSnsCredentials(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "sns",
				Verb:     "delete-credentials",
			}, nil
		},
	}
}

func mnqSnsUpdateCredentials() *core.Command {
	return &core.Command{
		Short:     `Update Topics and Events credentials`,
		Long:      `Update a set of Topics and Events credentials. You can update the credentials' name, or their permissions.`,
		Namespace: "mnq",
		Resource:  "sns",
		Verb:      "update-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIUpdateSnsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "sns-credentials-id",
				Short:      `ID of the Topics and Events credentials to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the credentials`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-publish",
				Short:      `Defines whether the credentials bearer can publish messages to the service (publish to Topics and Events topics)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-receive",
				Short:      `Defines whether the credentials bearer can receive messages from the service (configure subscriptions)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-manage",
				Short:      `Defines whether the credentials bearer can manage the associated Topics and Events topics or subscriptions`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SnsAPIUpdateSnsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)

			return api.UpdateSnsCredentials(request)
		},
	}
}

func mnqSnsGetCredentials() *core.Command {
	return &core.Command{
		Short:     `Get Topics and Events credentials`,
		Long:      `Retrieve an existing set of credentials, identified by the ` + "`" + `credentials_id` + "`" + `. The credentials themselves, as well as their metadata (name, project ID etc), are returned in the response.`,
		Namespace: "mnq",
		Resource:  "sns",
		Verb:      "get-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIGetSnsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "sns-credentials-id",
				Short:      `ID of the Topics and Events credentials to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SnsAPIGetSnsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)

			return api.GetSnsCredentials(request)
		},
	}
}

func mnqSnsListCredentials() *core.Command {
	return &core.Command{
		Short:     `List Topics and Events credentials`,
		Long:      `List existing Topics and Events credentials in the specified region. The response contains only the metadata for the credentials, not the credentials themselves.`,
		Namespace: "mnq",
		Resource:  "sns",
		Verb:      "list-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIListSnsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Include only Topics and Events credentials in this Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"updated_at_asc",
					"updated_at_desc",
					"name_asc",
					"name_desc",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SnsAPIListSnsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListSnsCredentials(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.SnsCredentials, nil
		},
	}
}

func mnqSqsActivate() *core.Command {
	return &core.Command{
		Short:     `Activate Queues`,
		Long:      `Activate Queues for the specified Project ID. Queues must be activated before any usage such as creating credentials and queues. Activating Queues does not trigger any billing, and you can deactivate at any time.`,
		Namespace: "mnq",
		Resource:  "sqs",
		Verb:      "activate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIActivateSqsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SqsAPIActivateSqsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)

			return api.ActivateSqs(request)
		},
	}
}

func mnqSqsGetInfo() *core.Command {
	return &core.Command{
		Short:     `Get Queues info`,
		Long:      `Retrieve the Queues information of the specified Project ID. information include the activation status and the Queues API endpoint URL.`,
		Namespace: "mnq",
		Resource:  "sqs",
		Verb:      "get-info",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIGetSqsInfoRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SqsAPIGetSqsInfoRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)

			return api.GetSqsInfo(request)
		},
	}
}

func mnqSqsDeactivate() *core.Command {
	return &core.Command{
		Short:     `Deactivate Queues`,
		Long:      `Deactivate Queues for the specified Project ID. You must delete all queues and credentials before this call or you need to set the force_delete parameter.`,
		Namespace: "mnq",
		Resource:  "sqs",
		Verb:      "deactivate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIDeactivateSqsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SqsAPIDeactivateSqsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)

			return api.DeactivateSqs(request)
		},
	}
}

func mnqSqsCreateCredentials() *core.Command {
	return &core.Command{
		Short:     `Create Queues credentials`,
		Long:      `Create a set of credentials for Queues, specified by a Project ID. Credentials give the bearer access to queues, and the level of permissions can be defined granularly.`,
		Namespace: "mnq",
		Resource:  "sqs",
		Verb:      "create-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPICreateSqsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the credentials`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("mnq_sqs"),
			},
			{
				Name:       "permissions.can-publish",
				Short:      `Defines whether the credentials bearer can publish messages to the service (send messages to Queues queues)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-receive",
				Short:      `Defines whether the credentials bearer can receive messages from Queues queues`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-manage",
				Short:      `Defines whether the credentials bearer can manage the associated Queues queues`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SqsAPICreateSqsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)

			return api.CreateSqsCredentials(request)
		},
	}
}

func mnqSqsDeleteCredentials() *core.Command {
	return &core.Command{
		Short:     `Delete Queues credentials`,
		Long:      `Delete a set of Queues credentials, specified by their credentials ID. Deleting credentials is irreversible and cannot be undone. The credentials can then no longer be used to access Queues.`,
		Namespace: "mnq",
		Resource:  "sqs",
		Verb:      "delete-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIDeleteSqsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "sqs-credentials-id",
				Short:      `ID of the credentials to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SqsAPIDeleteSqsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)
			e = api.DeleteSqsCredentials(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "sqs",
				Verb:     "delete-credentials",
			}, nil
		},
	}
}

func mnqSqsUpdateCredentials() *core.Command {
	return &core.Command{
		Short:     `Update Queues credentials`,
		Long:      `Update a set of Queues credentials. You can update the credentials' name, or their permissions.`,
		Namespace: "mnq",
		Resource:  "sqs",
		Verb:      "update-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIUpdateSqsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "sqs-credentials-id",
				Short:      `ID of the Queues credentials to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the credentials`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-publish",
				Short:      `Defines whether the credentials bearer can publish messages to the service (send messages to Queues queues)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-receive",
				Short:      `Defines whether the credentials bearer can receive messages from Queues queues`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-manage",
				Short:      `Defines whether the credentials bearer can manage the associated Queues queues`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SqsAPIUpdateSqsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)

			return api.UpdateSqsCredentials(request)
		},
	}
}

func mnqSqsGetCredentials() *core.Command {
	return &core.Command{
		Short:     `Get Queues credentials`,
		Long:      `Retrieve an existing set of credentials, identified by the ` + "`" + `credentials_id` + "`" + `. The credentials themselves, as well as their metadata (name, project ID etc), are returned in the response.`,
		Namespace: "mnq",
		Resource:  "sqs",
		Verb:      "get-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIGetSqsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "sqs-credentials-id",
				Short:      `ID of the Queues credentials to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SqsAPIGetSqsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)

			return api.GetSqsCredentials(request)
		},
	}
}

func mnqSqsListCredentials() *core.Command {
	return &core.Command{
		Short:     `List Queues credentials`,
		Long:      `List existing Queues credentials in the specified region. The response contains only the metadata for the credentials, not the credentials themselves.`,
		Namespace: "mnq",
		Resource:  "sqs",
		Verb:      "list-credentials",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIListSqsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Include only Queues credentials in this Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"updated_at_asc",
					"updated_at_desc",
					"name_asc",
					"name_desc",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*mnq.SqsAPIListSqsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListSqsCredentials(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.SqsCredentials, nil
		},
	}
}
