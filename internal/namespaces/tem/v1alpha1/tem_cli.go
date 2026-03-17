// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package tem

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	tem "github.com/scaleway/scaleway-sdk-go/api/tem/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		temRoot(),
		temEmail(),
		temDomain(),
		temWebhook(),
		temProjectSettings(),
		temBlocklists(),
		temOffers(),
		temProjectConsumption(),
		temEmailCreate(),
		temEmailGet(),
		temEmailList(),
		temEmailGetStatistics(),
		temEmailCancel(),
		temDomainCreate(),
		temDomainGet(),
		temDomainList(),
		temDomainRevoke(),
		temDomainCheck(),
		temDomainGetLastStatus(),
		temDomainUpdate(),
		temWebhookCreate(),
		temWebhookList(),
		temWebhookGet(),
		temWebhookUpdate(),
		temWebhookDelete(),
		temWebhookListEvents(),
		temBlocklistsList(),
		temBlocklistsCreate(),
		temBlocklistsDelete(),
		temOffersUpdate(),
		temOffersList(),
		temProjectConsumptionGet(),
	)
}

func temRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your Transactional Email services`,
		Long:      `This API allows you to manage your Transactional Email services.`,
		Namespace: "tem",
	}
}

func temEmail() *core.Command {
	return &core.Command{
		Short:     `Email management commands`,
		Long:      `This section lists your emails and shows you how to manage them.`,
		Namespace: "tem",
		Resource:  "email",
	}
}

func temDomain() *core.Command {
	return &core.Command{
		Short:     `Domain management commands`,
		Long:      `This section lists your domains, shows you to manage them, and gives you information about them.`,
		Namespace: "tem",
		Resource:  "domain",
	}
}

func temWebhook() *core.Command {
	return &core.Command{
		Short:     `Webhook management commands`,
		Long:      `Webhooks enable real-time communication and automation between systems by sending messages through all protocols supported by SNS, such as HTTP, HTTPS, and Serverless Functions, allowing for immediate updates and actions based on specific events. This feature is in beta. You can request quotas from the [Scaleway betas page](https://www.scaleway.com/fr/betas/#email-webhooks).`,
		Namespace: "tem",
		Resource:  "webhook",
	}
}

func temProjectSettings() *core.Command {
	return &core.Command{
		Short:     `Project settings management commands`,
		Long:      `Project settings allow you to manage the configuration of your projects.`,
		Namespace: "tem",
		Resource:  "project-settings",
	}
}

func temBlocklists() *core.Command {
	return &core.Command{
		Short:     `Blocklist`,
		Long:      `This section allows you to manage the blocklist of your emails.`,
		Namespace: "tem",
		Resource:  "blocklists",
	}
}

func temOffers() *core.Command {
	return &core.Command{
		Short:     `Project offers management commands`,
		Long:      `This section allows you to manage and get get subscribed information about your project email offer.`,
		Namespace: "tem",
		Resource:  "offers",
	}
}

func temProjectConsumption() *core.Command {
	return &core.Command{
		Short:     `Project consumption management commands`,
		Long:      `Project consumption allow you to see your project consumption.`,
		Namespace: "tem",
		Resource:  "project-consumption",
	}
}

func temEmailCreate() *core.Command {
	return &core.Command{
		Short:     `Send an email`,
		Long:      `You must specify the ` + "`" + `region` + "`" + `, the sender and the recipient's information and the ` + "`" + `project_id` + "`" + ` to send an email from a checked domain.`,
		Namespace: "tem",
		Resource:  "email",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.CreateEmailRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "from.email",
				Short:      `Email address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "from.name",
				Short:      `(Optional) Name displayed`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "to.{index}.email",
				Short:      `Email address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "to.{index}.name",
				Short:      `(Optional) Name displayed`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cc.{index}.email",
				Short:      `Email address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cc.{index}.name",
				Short:      `(Optional) Name displayed`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bcc.{index}.email",
				Short:      `Email address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bcc.{index}.name",
				Short:      `(Optional) Name displayed`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "subject",
				Short:      `Subject of the email`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "text",
				Short:      `Text content`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "html",
				Short:      `HTML content`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "attachments.{index}.name",
				Short:      `Filename of the attachment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "attachments.{index}.type",
				Short:      `MIME type of the attachment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "attachments.{index}.content",
				Short:      `Content of the attachment encoded in base64`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "send-before",
				Short:      `Maximum date to deliver the email`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "additional-headers.{index}.key",
				Short:      `Email header key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "additional-headers.{index}.value",
				Short:      `Email header value`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.CreateEmailRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.CreateEmail(request)
		},
	}
}

func temEmailGet() *core.Command {
	return &core.Command{
		Short:     `Get an email`,
		Long:      `Retrieve information about a specific email using the ` + "`" + `email_id` + "`" + ` and ` + "`" + `region` + "`" + ` parameters.`,
		Namespace: "tem",
		Resource:  "email",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.GetEmailRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "email-id",
				Short:      `ID of the email to retrieve`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.GetEmailRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.GetEmail(request)
		},
	}
}

func temEmailList() *core.Command {
	return &core.Command{
		Short:     `List emails`,
		Long:      `Retrieve the list of emails sent from a specific domain or for a specific Project or Organization. You must specify the ` + "`" + `region` + "`" + `.`,
		Namespace: "tem",
		Resource:  "email",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.ListEmailsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `(Optional) ID of the Project in which to list the emails`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain-id",
				Short:      `(Optional) ID of the domain for which to list the emails`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "message-id",
				Short:      `(Optional) ID of the message for which to list the emails`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "since",
				Short:      `(Optional) List emails created after this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "until",
				Short:      `(Optional) List emails created before this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mail-from",
				Short:      `(Optional) List emails sent with this sender's email address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mail-to",
				Short:      `Deprecated. List emails sent to this recipient's email address`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "mail-rcpt",
				Short:      `(Optional) List emails sent to this recipient's email address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "statuses.{index}",
				Short:      `(Optional) List emails with any of these statuses`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"new",
					"sending",
					"sent",
					"failed",
					"canceled",
				},
			},
			{
				Name:       "subject",
				Short:      `(Optional) List emails with this subject`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "search",
				Short:      `(Optional) List emails by searching to all fields`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `(Optional) List emails corresponding to specific criteria`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
					"updated_at_desc",
					"updated_at_asc",
					"status_desc",
					"status_asc",
					"mail_from_desc",
					"mail_from_asc",
					"mail_rcpt_desc",
					"mail_rcpt_asc",
					"subject_desc",
					"subject_asc",
				},
			},
			{
				Name:       "flags.{index}",
				Short:      `(Optional) List emails containing only specific flags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_flag",
					"soft_bounce",
					"hard_bounce",
					"spam",
					"mailbox_full",
					"mailbox_not_found",
					"greylisted",
					"send_before_expiration",
					"blocklisted",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.ListEmailsRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListEmails(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Emails, nil
		},
	}
}

func temEmailGetStatistics() *core.Command {
	return &core.Command{
		Short:     `Email statuses`,
		Long:      `Get information on your emails' statuses.`,
		Namespace: "tem",
		Resource:  "email",
		Verb:      "get-statistics",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.GetStatisticsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `(Optional) Number of emails for this Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain-id",
				Short:      `(Optional) Number of emails sent from this domain (must be coherent with the ` + "`" + `project_id` + "`" + ` and the ` + "`" + `organization_id` + "`" + `)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "since",
				Short:      `(Optional) Number of emails created after this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "until",
				Short:      `(Optional) Number of emails created before this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mail-from",
				Short:      `(Optional) Number of emails sent with this sender's email address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.GetStatisticsRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.GetStatistics(request)
		},
	}
}

func temEmailCancel() *core.Command {
	return &core.Command{
		Short:     `Cancel an email`,
		Long:      `You can cancel the sending of an email if it has not been sent yet. You must specify the ` + "`" + `region` + "`" + ` and the ` + "`" + `email_id` + "`" + ` of the email you want to cancel.`,
		Namespace: "tem",
		Resource:  "email",
		Verb:      "cancel",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.CancelEmailRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "email-id",
				Short:      `ID of the email to cancel`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.CancelEmailRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.CancelEmail(request)
		},
	}
}

func temDomainCreate() *core.Command {
	return &core.Command{
		Short:     `Register a domain in a project`,
		Long:      `You must specify the ` + "`" + `region` + "`" + `, ` + "`" + `project_id` + "`" + ` and ` + "`" + `domain_name` + "`" + ` to register a domain in a specific Project.`,
		Namespace: "tem",
		Resource:  "domain",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.CreateDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "domain-name",
				Short:      `Fully qualified domain dame`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "accept-tos",
				Short:      `Deprecated. Accept Scaleway's Terms of Service`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "autoconfig",
				Short:      `Activate auto-configuration of the domain's DNS zone`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.CreateDomainRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.CreateDomain(request)
		},
	}
}

func temDomainGet() *core.Command {
	return &core.Command{
		Short:     `Get information about a domain`,
		Long:      `Retrieve information about a specific domain using the ` + "`" + `region` + "`" + ` and ` + "`" + `domain_id` + "`" + ` parameters. Monitor your domain's reputation and improve **average** and **bad** reputation statuses, using your domain's **Email activity** tab on the [Scaleway console](https://console.scaleway.com/transactional-email/domains) to get a more detailed report. Check out our [dedicated documentation](https://www.scaleway.com/en/docs/managed-services/transactional-email/reference-content/understanding-tem-reputation-score/) to improve your domain's reputation.`,
		Namespace: "tem",
		Resource:  "domain",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.GetDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain-id",
				Short:      `ID of the domain`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.GetDomainRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.GetDomain(request)
		},
	}
}

func temDomainList() *core.Command {
	return &core.Command{
		Short:     `List domains`,
		Long:      `Retrieve domains in a specific Project or in a specific Organization using the ` + "`" + `region` + "`" + ` parameter.`,
		Namespace: "tem",
		Resource:  "domain",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.ListDomainsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `(Optional) ID of the Project in which to list the domains`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status.{index}",
				Short:      `(Optional) List domains under specific statuses`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"checked",
					"unchecked",
					"invalid",
					"locked",
					"revoked",
					"pending",
					"autoconfiguring",
				},
			},
			{
				Name:       "name",
				Short:      `(Optional) Names of the domains to list`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `(Optional) ID of the Organization in which to list the domains`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.ListDomainsRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListDomains(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Domains, nil
		},
	}
}

func temDomainRevoke() *core.Command {
	return &core.Command{
		Short:     `Delete a domain`,
		Long:      `You must specify the domain you want to delete by the ` + "`" + `region` + "`" + ` and ` + "`" + `domain_id` + "`" + `. Deleting a domain is permanent and cannot be undone.`,
		Namespace: "tem",
		Resource:  "domain",
		Verb:      "revoke",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.RevokeDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain-id",
				Short:      `ID of the domain to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.RevokeDomainRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.RevokeDomain(request)
		},
	}
}

func temDomainCheck() *core.Command {
	return &core.Command{
		Short:     `Domain DNS check`,
		Long:      `Perform an immediate DNS check of a domain using the ` + "`" + `region` + "`" + ` and ` + "`" + `domain_id` + "`" + ` parameters.`,
		Namespace: "tem",
		Resource:  "domain",
		Verb:      "check",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.CheckDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain-id",
				Short:      `ID of the domain to check`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.CheckDomainRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.CheckDomain(request)
		},
	}
}

func temDomainGetLastStatus() *core.Command {
	return &core.Command{
		Short:     `Display SPF, DKIM, DMARC and MX records status and potential errors`,
		Long:      `Display SPF, DKIM, DMARC and MX records status and potential errors, including the found records to make debugging easier.`,
		Namespace: "tem",
		Resource:  "domain",
		Verb:      "get-last-status",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.GetDomainLastStatusRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain-id",
				Short:      `ID of the domain to get records status`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.GetDomainLastStatusRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.GetDomainLastStatus(request)
		},
	}
}

func temDomainUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a domain`,
		Long:      `Update a domain auto-configuration.`,
		Namespace: "tem",
		Resource:  "domain",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.UpdateDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain-id",
				Short:      `ID of the domain to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "autoconfig",
				Short:      `(Optional) If set to true, activate auto-configuration of the domain's DNS zone`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.UpdateDomainRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.UpdateDomain(request)
		},
	}
}

func temWebhookCreate() *core.Command {
	return &core.Command{
		Short:     `Create a Webhook`,
		Long:      `Create a new Webhook triggered by a list of event types and pushed to a Scaleway SNS ARN.`,
		Namespace: "tem",
		Resource:  "webhook",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.CreateWebhookRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain-id",
				Short:      `ID of the Domain to watch for triggering events`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the Webhook`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "event-types.{index}",
				Short:      `List of event types that will trigger an event`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"email_queued",
					"email_dropped",
					"email_deferred",
					"email_delivered",
					"email_spam",
					"email_mailbox_not_found",
					"email_blocklisted",
					"blocklist_created",
				},
			},
			{
				Name:       "sns-arn",
				Short:      `Scaleway SNS ARN topic to push the events to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.CreateWebhookRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.CreateWebhook(request)
		},
	}
}

func temWebhookList() *core.Command {
	return &core.Command{
		Short:     `List Webhooks`,
		Long:      `Retrieve Webhooks in a specific Project or in a specific Organization using the ` + "`" + `region` + "`" + ` parameter.`,
		Namespace: "tem",
		Resource:  "webhook",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.ListWebhooksRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `(Optional) List Webhooks corresponding to specific criteria`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
				},
			},
			{
				Name:       "project-id",
				Short:      `(Optional) ID of the Project for which to list the Webhooks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain-id",
				Short:      `(Optional) ID of the Domain for which to list the Webhooks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `(Optional) ID of the Organization for which to list the Webhooks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.ListWebhooksRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListWebhooks(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Webhooks, nil
		},
	}
}

func temWebhookGet() *core.Command {
	return &core.Command{
		Short:     `Get information about a Webhook`,
		Long:      `Retrieve information about a specific Webhook using the ` + "`" + `webhook_id` + "`" + ` and ` + "`" + `region` + "`" + ` parameters.`,
		Namespace: "tem",
		Resource:  "webhook",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.GetWebhookRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "webhook-id",
				Short:      `ID of the Webhook to check`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.GetWebhookRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.GetWebhook(request)
		},
	}
}

func temWebhookUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a Webhook`,
		Long:      `Update a Webhook events type, SNS ARN or name.`,
		Namespace: "tem",
		Resource:  "webhook",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.UpdateWebhookRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "webhook-id",
				Short:      `ID of the Webhook to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the Webhook to update`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "event-types.{index}",
				Short:      `List of event types to update`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"email_queued",
					"email_dropped",
					"email_deferred",
					"email_delivered",
					"email_spam",
					"email_mailbox_not_found",
					"email_blocklisted",
					"blocklist_created",
				},
			},
			{
				Name:       "sns-arn",
				Short:      `Scaleway SNS ARN topic to update`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.UpdateWebhookRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.UpdateWebhook(request)
		},
	}
}

func temWebhookDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a Webhook`,
		Long:      `You must specify the Webhook you want to delete by the ` + "`" + `region` + "`" + ` and ` + "`" + `webhook_id` + "`" + `. Deleting a Webhook is permanent and cannot be undone.`,
		Namespace: "tem",
		Resource:  "webhook",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.DeleteWebhookRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "webhook-id",
				Short:      `ID of the Webhook to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.DeleteWebhookRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			e = api.DeleteWebhook(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "webhook",
				Verb:     "delete",
			}, nil
		},
	}
}

func temWebhookListEvents() *core.Command {
	return &core.Command{
		Short:     `List Webhook triggered events`,
		Long:      `Retrieve the list of Webhook events triggered from a specific Webhook or for a specific Project or Organization. You must specify the ` + "`" + `region` + "`" + `.`,
		Namespace: "tem",
		Resource:  "webhook",
		Verb:      "list-events",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.ListWebhookEventsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `(Optional) List Webhook events corresponding to specific criteria`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
				},
			},
			{
				Name:       "webhook-id",
				Short:      `ID of the Webhook linked to the events`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "email-id",
				Short:      `ID of the email linked to the events`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "event-types.{index}",
				Short:      `List of event types linked to the events`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"email_queued",
					"email_dropped",
					"email_deferred",
					"email_delivered",
					"email_spam",
					"email_mailbox_not_found",
					"email_blocklisted",
					"blocklist_created",
				},
			},
			{
				Name:       "statuses.{index}",
				Short:      `List of event statuses`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_status",
					"sending",
					"sent",
					"failed",
				},
			},
			{
				Name:       "project-id",
				Short:      `ID of the webhook Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain-id",
				Short:      `ID of the domain to watch for triggering events`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `ID of the webhook Organization`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.ListWebhookEventsRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListWebhookEvents(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.WebhookEvents, nil
		},
	}
}

func temBlocklistsList() *core.Command {
	return &core.Command{
		Short:     `List blocklists`,
		Long:      `Retrieve the list of blocklists.`,
		Namespace: "tem",
		Resource:  "blocklists",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.ListBlocklistsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `(Optional) List blocklist corresponding to specific criteria`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
					"ends_at_desc",
					"ends_at_asc",
				},
			},
			{
				Name:       "domain-id",
				Short:      `(Optional) Filter by a domain ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "email",
				Short:      `(Optional) Filter by an email address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `(Optional) Filter by a blocklist type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"mailbox_full",
					"mailbox_not_found",
				},
			},
			{
				Name:       "custom",
				Short:      `(Optional) Filter by custom blocklist (true) or automatic Transactional Email blocklist (false)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.ListBlocklistsRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListBlocklists(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Blocklists, nil
		},
	}
}

func temBlocklistsCreate() *core.Command {
	return &core.Command{
		Short:     `Bulk create blocklists`,
		Long:      `Create multiple blocklists in a specific Project or Organization using the ` + "`" + `region` + "`" + ` parameter.`,
		Namespace: "tem",
		Resource:  "blocklists",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.BulkCreateBlocklistsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain-id",
				Short:      `Domain ID linked to the blocklist`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "emails.{index}",
				Short:      `Email blocked by the blocklist`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Type of blocklist`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"mailbox_full",
					"mailbox_not_found",
				},
			},
			{
				Name:       "reason",
				Short:      `Reason to block the email`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.BulkCreateBlocklistsRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.BulkCreateBlocklists(request)
		},
	}
}

func temBlocklistsDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a blocklist`,
		Long:      `You must specify the blocklist you want to delete by the ` + "`" + `region` + "`" + ` and ` + "`" + `blocklist_id` + "`" + `.`,
		Namespace: "tem",
		Resource:  "blocklists",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.DeleteBlocklistRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "blocklist-id",
				Short:      `ID of the blocklist to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.DeleteBlocklistRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			e = api.DeleteBlocklist(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "blocklists",
				Verb:     "delete",
			}, nil
		},
	}
}

func temOffersUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a subscribed offer`,
		Long:      `Update a subscribed offer.`,
		Namespace: "tem",
		Resource:  "offers",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.UpdateOfferSubscriptionRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the offer-subscription`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_name",
					"essential",
					"scale",
				},
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.UpdateOfferSubscriptionRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.UpdateOfferSubscription(request)
		},
	}
}

func temOffersList() *core.Command {
	return &core.Command{
		Short:     `List the available offers.`,
		Long:      `Retrieve the list of the available and free-of-charge offers you can subscribe to.`,
		Namespace: "tem",
		Resource:  "offers",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.ListOffersRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.ListOffersRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.ListOffers(request)
		},
	}
}

func temProjectConsumptionGet() *core.Command {
	return &core.Command{
		Short:     `Get project resource consumption.`,
		Long:      `Get project resource consumption.`,
		Namespace: "tem",
		Resource:  "project-consumption",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.GetProjectConsumptionRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*tem.GetProjectConsumptionRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)

			return api.GetProjectConsumption(request)
		},
	}
}
