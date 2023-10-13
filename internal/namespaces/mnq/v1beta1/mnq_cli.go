// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package mnq

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		mnqRoot(),
		mnqNatsAccounts(),
		mnqNatsCredentials(),
		mnqSns(),
		mnqSnsCredentials(),
		mnqSqs(),
		mnqSqsCredentials(),
		mnqNatsAccountsCreate(),
		mnqNatsAccountsDelete(),
		mnqNatsAccountsUpdate(),
		mnqNatsAccountsGet(),
		mnqNatsAccountsList(),
		mnqNatsCredentialsCreate(),
		mnqNatsCredentialsDelete(),
		mnqNatsCredentialsGet(),
		mnqNatsCredentialsList(),
		mnqSnsActivate(),
		mnqSnsGet(),
		mnqSnsDeactivate(),
		mnqSnsCredentialsCreate(),
		mnqSnsCredentialsDelete(),
		mnqSnsCredentialsUpdate(),
		mnqSnsCredentialsGet(),
		mnqSnsCredentialsList(),
		mnqSqsActivate(),
		mnqSqsGet(),
		mnqSqsDeactivate(),
		mnqSqsCredentialsCreate(),
		mnqSqsCredentialsDelete(),
		mnqSqsCredentialsUpdate(),
		mnqSqsCredentialsGet(),
		mnqSqsCredentialsList(),
	)
}
func mnqRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage Scaleway Messaging and Queueing NATS accounts`,
		Long:      `Messaging and Queuing NATS API.`,
		Namespace: "mnq",
	}
}

func mnqNatsAccounts() *core.Command {
	return &core.Command{
		Short:     `MnQ NATS Accounts commands`,
		Long:      `MnQ NATS Accounts commands.`,
		Namespace: "mnq",
		Resource:  "nats-accounts",
	}
}

func mnqNatsCredentials() *core.Command {
	return &core.Command{
		Short:     `MnQ NATS Credentials commands`,
		Long:      `MnQ NATS Credentials commands.`,
		Namespace: "mnq",
		Resource:  "nats-credentials",
	}
}

func mnqSns() *core.Command {
	return &core.Command{
		Short:     `MnQ SNS commands`,
		Long:      `MnQ SNS commands.`,
		Namespace: "mnq",
		Resource:  "sns",
	}
}

func mnqSnsCredentials() *core.Command {
	return &core.Command{
		Short:     `MnQ SNS Credentials commands`,
		Long:      `MnQ SNS Credentials commands.`,
		Namespace: "mnq",
		Resource:  "sns-credentials",
	}
}

func mnqSqs() *core.Command {
	return &core.Command{
		Short:     `MnQ SQS commands`,
		Long:      `MnQ SQS commands.`,
		Namespace: "mnq",
		Resource:  "sqs",
	}
}

func mnqSqsCredentials() *core.Command {
	return &core.Command{
		Short:     `MnQ SQS Credentials commands`,
		Long:      `MnQ SQS Credentials commands.`,
		Namespace: "mnq",
		Resource:  "sqs-credentials",
	}
}

func mnqNatsAccountsCreate() *core.Command {
	return &core.Command{
		Short:     `Create a NATS account`,
		Long:      `Create a NATS account associated with a Project.`,
		Namespace: "mnq",
		Resource:  "nats-accounts",
		Verb:      "create",
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.NatsAPICreateNatsAccountRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)
			return api.CreateNatsAccount(request)

		},
	}
}

func mnqNatsAccountsDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a NATS account`,
		Long:      `Delete a NATS account, specified by its NATS account ID. Note that deleting a NATS account is irreversible, and any credentials, streams, consumer and stored messages belonging to this NATS account will also be deleted.`,
		Namespace: "mnq",
		Resource:  "nats-accounts",
		Verb:      "delete",
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.NatsAPIDeleteNatsAccountRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)
			e = api.DeleteNatsAccount(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "nats-accounts",
				Verb:     "delete",
			}, nil
		},
	}
}

func mnqNatsAccountsUpdate() *core.Command {
	return &core.Command{
		Short:     `Update the name of a NATS account`,
		Long:      `Update the name of a NATS account, specified by its NATS account ID.`,
		Namespace: "mnq",
		Resource:  "nats-accounts",
		Verb:      "update",
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.NatsAPIUpdateNatsAccountRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)
			return api.UpdateNatsAccount(request)

		},
	}
}

func mnqNatsAccountsGet() *core.Command {
	return &core.Command{
		Short:     `Get a NATS account`,
		Long:      `Retrieve information about an existing NATS account identified by its NATS account ID. Its full details, including name and endpoint, are returned in the response.`,
		Namespace: "mnq",
		Resource:  "nats-accounts",
		Verb:      "get",
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.NatsAPIGetNatsAccountRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)
			return api.GetNatsAccount(request)

		},
	}
}

func mnqNatsAccountsList() *core.Command {
	return &core.Command{
		Short:     `List NATS accounts`,
		Long:      `List all NATS accounts in the specified region, for a Scaleway Organization or Project. By default, the NATS accounts returned in the list are ordered by creation date in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "mnq",
		Resource:  "nats-accounts",
		Verb:      "list",
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
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "name_asc", "name_desc"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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

func mnqNatsCredentialsCreate() *core.Command {
	return &core.Command{
		Short:     `Create NATS credentials`,
		Long:      `Create a set of credentials for a NATS account, specified by its NATS account ID.`,
		Namespace: "mnq",
		Resource:  "nats-credentials",
		Verb:      "create",
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.NatsAPICreateNatsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)
			return api.CreateNatsCredentials(request)

		},
	}
}

func mnqNatsCredentialsDelete() *core.Command {
	return &core.Command{
		Short:     `Delete NATS credentials`,
		Long:      `Delete a set of credentials, specified by their credentials ID. Deleting credentials is irreversible and cannot be undone. The credentials can no longer be used to access the NATS account, and active connections using this credentials will be closed.`,
		Namespace: "mnq",
		Resource:  "nats-credentials",
		Verb:      "delete",
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.NatsAPIDeleteNatsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)
			e = api.DeleteNatsCredentials(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "nats-credentials",
				Verb:     "delete",
			}, nil
		},
	}
}

func mnqNatsCredentialsGet() *core.Command {
	return &core.Command{
		Short:     `Get NATS credentials`,
		Long:      `Retrieve an existing set of credentials, identified by the ` + "`" + `nats_credentials_id` + "`" + `. The credentials themselves are NOT returned, only their metadata (NATS account ID, credentials name, etc), are returned in the response.`,
		Namespace: "mnq",
		Resource:  "nats-credentials",
		Verb:      "get",
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.NatsAPIGetNatsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewNatsAPI(client)
			return api.GetNatsCredentials(request)

		},
	}
}

func mnqNatsCredentialsList() *core.Command {
	return &core.Command{
		Short:     `List NATS credentials`,
		Long:      `List existing credentials in the specified NATS account. The response contains only the metadata for the credentials, not the credentials themselves, which are only returned after a **Create Credentials** call.`,
		Namespace: "mnq",
		Resource:  "nats-credentials",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.NatsAPIListNatsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "nats-account-id",
				Short:      `Include only credentials for this NATS account`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "name_asc", "name_desc"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Short:     `Activate SNS`,
		Long:      `Activate SNS for the specified Project ID. SNS must be activated before any usage. Activating SNS does not trigger any billing, and you can deactivate at any time.`,
		Namespace: "mnq",
		Resource:  "sns",
		Verb:      "activate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIActivateSnsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SnsAPIActivateSnsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)
			return api.ActivateSns(request)

		},
	}
}

func mnqSnsGet() *core.Command {
	return &core.Command{
		Short:     `Get SNS info`,
		Long:      `Retrieve the SNS information of the specified Project ID. Informations include the activation status and the SNS API endpoint URL.`,
		Namespace: "mnq",
		Resource:  "sns",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIGetSnsInfoRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SnsAPIGetSnsInfoRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)
			return api.GetSnsInfo(request)

		},
	}
}

func mnqSnsDeactivate() *core.Command {
	return &core.Command{
		Short:     `Deactivate SNS`,
		Long:      `Deactivate SNS for the specified Project ID.You must delete all topics and credentials before this call or you need to set the force_delete parameter.`,
		Namespace: "mnq",
		Resource:  "sns",
		Verb:      "deactivate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIDeactivateSnsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SnsAPIDeactivateSnsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)
			return api.DeactivateSns(request)

		},
	}
}

func mnqSnsCredentialsCreate() *core.Command {
	return &core.Command{
		Short:     `Create SNS credentials`,
		Long:      `Create a set of credentials for SNS, specified by a Project ID. Credentials give the bearer access to topics, and the level of permissions can be defined granularly.`,
		Namespace: "mnq",
		Resource:  "sns-credentials",
		Verb:      "create",
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
				Short:      `Defines whether the credentials bearer can publish messages to the service (publish to SNS topics)`,
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
				Short:      `Defines whether the credentials bearer can manage the associated SNS topics or subscriptions`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SnsAPICreateSnsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)
			return api.CreateSnsCredentials(request)

		},
	}
}

func mnqSnsCredentialsDelete() *core.Command {
	return &core.Command{
		Short:     `Delete SNS credentials`,
		Long:      `Delete a set of SNS credentials, specified by their credentials ID. Deleting credentials is irreversible and cannot be undone. The credentials can then no longer be used to access SNS.`,
		Namespace: "mnq",
		Resource:  "sns-credentials",
		Verb:      "delete",
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SnsAPIDeleteSnsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)
			e = api.DeleteSnsCredentials(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "sns-credentials",
				Verb:     "delete",
			}, nil
		},
	}
}

func mnqSnsCredentialsUpdate() *core.Command {
	return &core.Command{
		Short:     `Update SNS credentials`,
		Long:      `Update a set of SNS credentials. You can update the credentials' name, or their permissions.`,
		Namespace: "mnq",
		Resource:  "sns-credentials",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIUpdateSnsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "sns-credentials-id",
				Short:      `ID of the SNS credentials to update`,
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
				Short:      `Defines whether the credentials bearer can publish messages to the service (publish to SNS topics)`,
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
				Short:      `Defines whether the credentials bearer can manage the associated SNS topics or subscriptions`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SnsAPIUpdateSnsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)
			return api.UpdateSnsCredentials(request)

		},
	}
}

func mnqSnsCredentialsGet() *core.Command {
	return &core.Command{
		Short:     `Get SNS credentials`,
		Long:      `Retrieve an existing set of credentials, identified by the ` + "`" + `credentials_id` + "`" + `. The credentials themselves, as well as their metadata (name, project ID etc), are returned in the response.`,
		Namespace: "mnq",
		Resource:  "sns-credentials",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIGetSnsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "sns-credentials-id",
				Short:      `ID of the SNS credentials to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SnsAPIGetSnsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSnsAPI(client)
			return api.GetSnsCredentials(request)

		},
	}
}

func mnqSnsCredentialsList() *core.Command {
	return &core.Command{
		Short:     `List SNS credentials`,
		Long:      `List existing SNS credentials in the specified region. The response contains only the metadata for the credentials, not the credentials themselves.`,
		Namespace: "mnq",
		Resource:  "sns-credentials",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SnsAPIListSnsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Include only SNS credentials in this Project`,
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
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "name_asc", "name_desc"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Short:     `Activate SQS`,
		Long:      `Activate SQS for the specified Project ID. SQS must be activated before any usage such as creating credentials and queues. Activating SQS does not trigger any billing, and you can deactivate at any time.`,
		Namespace: "mnq",
		Resource:  "sqs",
		Verb:      "activate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIActivateSqsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SqsAPIActivateSqsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)
			return api.ActivateSqs(request)

		},
	}
}

func mnqSqsGet() *core.Command {
	return &core.Command{
		Short:     `Get SQS info`,
		Long:      `Retrieve the SQS information of the specified Project ID. Informations include the activation status and the SQS API endpoint URL.`,
		Namespace: "mnq",
		Resource:  "sqs",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIGetSqsInfoRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SqsAPIGetSqsInfoRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)
			return api.GetSqsInfo(request)

		},
	}
}

func mnqSqsDeactivate() *core.Command {
	return &core.Command{
		Short:     `Deactivate SQS`,
		Long:      `Deactivate SQS for the specified Project ID. You must delete all queues and credentials before this call or you need to set the force_delete parameter.`,
		Namespace: "mnq",
		Resource:  "sqs",
		Verb:      "deactivate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIDeactivateSqsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SqsAPIDeactivateSqsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)
			return api.DeactivateSqs(request)

		},
	}
}

func mnqSqsCredentialsCreate() *core.Command {
	return &core.Command{
		Short:     `Create SQS credentials`,
		Long:      `Create a set of credentials for SQS, specified by a Project ID. Credentials give the bearer access to queues, and the level of permissions can be defined granularly.`,
		Namespace: "mnq",
		Resource:  "sqs-credentials",
		Verb:      "create",
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
				Short:      `Defines whether the credentials bearer can publish messages to the service (send messages to SQS queues)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-receive",
				Short:      `Defines whether the credentials bearer can receive messages from SQS queues`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-manage",
				Short:      `Defines whether the credentials bearer can manage the associated SQS queues`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SqsAPICreateSqsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)
			return api.CreateSqsCredentials(request)

		},
	}
}

func mnqSqsCredentialsDelete() *core.Command {
	return &core.Command{
		Short:     `Delete SQS credentials`,
		Long:      `Delete a set of SQS credentials, specified by their credentials ID. Deleting credentials is irreversible and cannot be undone. The credentials can then no longer be used to access SQS.`,
		Namespace: "mnq",
		Resource:  "sqs-credentials",
		Verb:      "delete",
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SqsAPIDeleteSqsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)
			e = api.DeleteSqsCredentials(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "sqs-credentials",
				Verb:     "delete",
			}, nil
		},
	}
}

func mnqSqsCredentialsUpdate() *core.Command {
	return &core.Command{
		Short:     `Update SQS credentials`,
		Long:      `Update a set of SQS credentials. You can update the credentials' name, or their permissions.`,
		Namespace: "mnq",
		Resource:  "sqs-credentials",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIUpdateSqsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "sqs-credentials-id",
				Short:      `ID of the SQS credentials to update`,
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
				Short:      `Defines whether the credentials bearer can publish messages to the service (send messages to SQS queues)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-receive",
				Short:      `Defines whether the credentials bearer can receive messages from SQS queues`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "permissions.can-manage",
				Short:      `Defines whether the credentials bearer can manage the associated SQS queues`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SqsAPIUpdateSqsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)
			return api.UpdateSqsCredentials(request)

		},
	}
}

func mnqSqsCredentialsGet() *core.Command {
	return &core.Command{
		Short:     `Get SQS credentials`,
		Long:      `Retrieve an existing set of credentials, identified by the ` + "`" + `credentials_id` + "`" + `. The credentials themselves, as well as their metadata (name, project ID etc), are returned in the response.`,
		Namespace: "mnq",
		Resource:  "sqs-credentials",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIGetSqsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "sqs-credentials-id",
				Short:      `ID of the SQS credentials to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*mnq.SqsAPIGetSqsCredentialsRequest)

			client := core.ExtractClient(ctx)
			api := mnq.NewSqsAPI(client)
			return api.GetSqsCredentials(request)

		},
	}
}

func mnqSqsCredentialsList() *core.Command {
	return &core.Command{
		Short:     `List SQS credentials`,
		Long:      `List existing SQS credentials in the specified region. The response contains only the metadata for the credentials, not the credentials themselves.`,
		Namespace: "mnq",
		Resource:  "sqs-credentials",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(mnq.SqsAPIListSqsCredentialsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Include only SQS credentials in this Project`,
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
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "name_asc", "name_desc"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
