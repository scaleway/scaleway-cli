// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package partner

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/partner/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		partnerRoot(),
		partnerOrganization(),
		partnerOrganizationCreate(),
		partnerOrganizationGet(),
		partnerOrganizationList(),
		partnerOrganizationLock(),
		partnerOrganizationUnlock(),
		partnerOrganizationUpdate(),
	)
}

func partnerRoot() *core.Command {
	return &core.Command{
		Short:     `Scaleway Partner API ( for partner only )`,
		Long:      `Scaleway Partner API ( for partner only ).`,
		Namespace: "partner",
	}
}

func partnerOrganization() *core.Command {
	return &core.Command{
		Short:     `Organization management commands`,
		Long:      `Organization management commands.`,
		Namespace: "partner",
		Resource:  "organization",
	}
}

func partnerOrganizationCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new organization`,
		Long:      `Create a new organization.`,
		Namespace: "partner",
		Resource:  "organization",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(partner.CreateOrganizationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "partner-id",
				Short:      `Your personal ` + "`" + `partner_id` + "`" + `. This is the same as your Organization ID.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "email",
				Short:      `The email of the new organization owner`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-name",
				Short:      `The name of the organization you want to create. Usually the company name`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-firstname",
				Short:      `The first name of the new organization owner`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-lastname",
				Short:      `The last name of the new organization owner`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "phone-number",
				Short:      `The phone number of the new organization owner`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "customer-id",
				Short:      `A custom ID for the customer in your own infrastructure`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "siren-number",
				Short:      `A SIREN number for the customer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*partner.CreateOrganizationRequest)

			client := core.ExtractClient(ctx)
			api := partner.NewAPI(client)

			return api.CreateOrganization(request)
		},
	}
}

func partnerOrganizationGet() *core.Command {
	return &core.Command{
		Short:     `Get an organization`,
		Long:      `Get an organization.`,
		Namespace: "partner",
		Resource:  "organization",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(partner.GetOrganizationRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*partner.GetOrganizationRequest)

			client := core.ExtractClient(ctx)
			api := partner.NewAPI(client)

			return api.GetOrganization(request)
		},
	}
}

func partnerOrganizationList() *core.Command {
	return &core.Command{
		Short:     `List Organizations`,
		Long:      `List Organizations.`,
		Namespace: "partner",
		Resource:  "organization",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(partner.ListOrganizationsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "status",
				Short:      `Only list organizations with this status`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_status",
					"opened",
					"locked",
					"closed",
				},
			},
			{
				Name:       "email",
				Short:      `Only list organizations created with this email`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "customer-id",
				Short:      `Only list organizations attached to this Customer ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "locked-by",
				Short:      `Only list organizations locked by a certain entity`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_locked_by",
					"partner",
					"scaleway",
				},
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*partner.ListOrganizationsRequest)

			client := core.ExtractClient(ctx)
			api := partner.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListOrganizations(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Organizations, nil
		},
	}
}

func partnerOrganizationLock() *core.Command {
	return &core.Command{
		Short:     `Lock an organization`,
		Long:      `Lock an organization.`,
		Namespace: "partner",
		Resource:  "organization",
		Verb:      "lock",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(partner.LockOrganizationRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*partner.LockOrganizationRequest)

			client := core.ExtractClient(ctx)
			api := partner.NewAPI(client)

			return api.LockOrganization(request)
		},
	}
}

func partnerOrganizationUnlock() *core.Command {
	return &core.Command{
		Short:     `Unlock an organization`,
		Long:      `Unlock an organization.`,
		Namespace: "partner",
		Resource:  "organization",
		Verb:      "unlock",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(partner.UnlockOrganizationRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*partner.UnlockOrganizationRequest)

			client := core.ExtractClient(ctx)
			api := partner.NewAPI(client)

			return api.UnlockOrganization(request)
		},
	}
}

func partnerOrganizationUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an organization`,
		Long:      `Update an organization.`,
		Namespace: "partner",
		Resource:  "organization",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(partner.UpdateOrganizationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "email",
				Short:      `The new email`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `The new name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-firstname",
				Short:      `The first name of the new owner`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-lastname",
				Short:      `The last name of the new owner`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "phone-number",
				Short:      `The phone number of the new owner`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "customer-id",
				Short:      `Customer ID associated with this organization`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "comment",
				Short:      `A comment about the organization`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*partner.UpdateOrganizationRequest)

			client := core.ExtractClient(ctx)
			api := partner.NewAPI(client)

			return api.UpdateOrganization(request)
		},
	}
}
