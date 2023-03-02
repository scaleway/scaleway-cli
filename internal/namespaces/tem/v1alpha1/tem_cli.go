// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package tem

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/tem/v1alpha1"
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
	)
}
func temRoot() *core.Command {
	return &core.Command{
		Short:     `Tem`,
		Long:      `Tem.`,
		Namespace: "tem",
	}
}

func temEmail() *core.Command {
	return &core.Command{
		Short:     `Email management commands`,
		Long:      `Email management commands.`,
		Namespace: "tem",
		Resource:  "email",
	}
}

func temDomain() *core.Command {
	return &core.Command{
		Short:     `Domain management commands`,
		Long:      `Domain management commands.`,
		Namespace: "tem",
		Resource:  "domain",
	}
}

func temEmailCreate() *core.Command {
	return &core.Command{
		Short:     `Send an email`,
		Long:      `Send an email.`,
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
				Short:      `Optional display name`,
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
				Short:      `Optional display name`,
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
				Short:      `Optional display name`,
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
				Short:      `Optional display name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "subject",
				Short:      `Message subject`,
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
				Short:      `MIME type of the attachment (Currently only allow, text files, pdf and html files)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "attachments.{index}.content",
				Short:      `Content of the attachment, encoded in base64`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "send-before",
				Short:      `Maximum date to deliver mail`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*tem.CreateEmailRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			return api.CreateEmail(request)

		},
	}
}

func temEmailGet() *core.Command {
	return &core.Command{
		Short:     `Get information about an email`,
		Long:      `Get information about an email.`,
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
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*tem.GetEmailRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			return api.GetEmail(request)

		},
	}
}

func temEmailList() *core.Command {
	return &core.Command{
		Short:     `List emails sent from a domain and/or for a project and/or for an organization`,
		Long:      `List emails sent from a domain and/or for a project and/or for an organization.`,
		Namespace: "tem",
		Resource:  "email",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.ListEmailsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Optional ID of the project in which to list the emails`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain-id",
				Short:      `Optional ID of the domain for which to list the emails`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "message-id",
				Short:      `Optional ID of the message for which to list the emails`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "since",
				Short:      `Optional, list emails created after this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "until",
				Short:      `Optional, list emails created before this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mail-from",
				Short:      `Optional, list emails sent with this ` + "`" + `mail_from` + "`" + ` sender's address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mail-to",
				Short:      `Optional, list emails sent with this ` + "`" + `mail_to` + "`" + ` recipient's address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "statuses.{index}",
				Short:      `Optional, list emails having any of this status`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "new", "sending", "sent", "failed", "canceled"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Short:     `Get statistics on the email statuses`,
		Long:      `Get statistics on the email statuses.`,
		Namespace: "tem",
		Resource:  "email",
		Verb:      "get-statistics",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.GetStatisticsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Optional, count emails for this project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain-id",
				Short:      `Optional, count emails send from this domain (must be coherent with the ` + "`" + `project_id` + "`" + ` and the ` + "`" + `organization_id` + "`" + `)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "since",
				Short:      `Optional, count emails created after this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "until",
				Short:      `Optional, count emails created before this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mail-from",
				Short:      `Optional, count emails sent with this ` + "`" + `mail_from` + "`" + ` sender's address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*tem.GetStatisticsRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			return api.GetStatistics(request)

		},
	}
}

func temEmailCancel() *core.Command {
	return &core.Command{
		Short:     `Try to cancel an email if it has not yet been sent`,
		Long:      `Try to cancel an email if it has not yet been sent.`,
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
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Long:      `Register a domain in a project.`,
		Namespace: "tem",
		Resource:  "domain",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.CreateDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "domain-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Long:      `Get information about a domain.`,
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
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*tem.GetDomainRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			return api.GetDomain(request)

		},
	}
}

func temDomainList() *core.Command {
	return &core.Command{
		Short:     `List domains in a project and/or in an organization`,
		Long:      `List domains in a project and/or in an organization.`,
		Namespace: "tem",
		Resource:  "domain",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.ListDomainsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "checked", "unchecked", "invalid", "locked", "revoked", "pending"},
			},
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Short:     `Revoke a domain`,
		Long:      `Revoke a domain.`,
		Namespace: "tem",
		Resource:  "domain",
		Verb:      "revoke",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(tem.RevokeDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain-id",
				Short:      `ID of the domain to revoke`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*tem.RevokeDomainRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			return api.RevokeDomain(request)

		},
	}
}

func temDomainCheck() *core.Command {
	return &core.Command{
		Short:     `Ask for an immediate check of a domain (DNS check)`,
		Long:      `Ask for an immediate check of a domain (DNS check).`,
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
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*tem.CheckDomainRequest)

			client := core.ExtractClient(ctx)
			api := tem.NewAPI(client)
			return api.CheckDomain(request)

		},
	}
}
