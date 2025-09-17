// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package audit_trail

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	audit_trail "github.com/scaleway/scaleway-sdk-go/api/audit_trail/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		auditTrailRoot(),
		auditTrailEvent(),
		auditTrailProduct(),
		auditTrailEventList(),
		auditTrailProductList(),
	)
}

func auditTrailRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to ensure accountability and security by recording events and changes performed within your Scaleway Organization.`,
		Long:      `This API allows you to ensure accountability and security by recording events and changes performed within your Scaleway Organization.`,
		Namespace: "audit-trail",
	}
}

func auditTrailEvent() *core.Command {
	return &core.Command{
		Short:     `Represent an entry in the Audit Trail`,
		Long:      `Represent an entry in the Audit Trail.`,
		Namespace: "audit-trail",
		Resource:  "event",
	}
}

func auditTrailProduct() *core.Command {
	return &core.Command{
		Short:     `Product integrated with Audit Trail`,
		Long:      `Product integrated with Audit Trail.`,
		Namespace: "audit-trail",
		Resource:  "product",
	}
}

func auditTrailEventList() *core.Command {
	return &core.Command{
		Short:     `List events`,
		Long:      `Retrieve the list of Audit Trail events for a Scaleway Organization and/or Project. You must specify the ` + "`" + `organization_id` + "`" + ` and optionally, the ` + "`" + `project_id` + "`" + `.`,
		Namespace: "audit-trail",
		Resource:  "event",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(audit_trail.ListEventsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `(Optional) ID of the Project containing the Audit Trail events`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "resource-type",
				Short:      `(Optional) Type of the Scaleway resource`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"secm_secret",
					"secm_secret_version",
					"kube_cluster",
					"kube_pool",
					"kube_node",
					"kube_acl",
					"keym_key",
					"iam_user",
					"iam_application",
					"iam_group",
					"iam_policy",
					"iam_api_key",
					"iam_ssh_key",
					"iam_rule",
					"iam_saml",
					"iam_saml_certificate",
					"secret_manager_secret",
					"secret_manager_version",
					"key_manager_key",
					"account_user",
					"account_organization",
					"account_project",
					"instance_server",
					"instance_placement_group",
					"instance_security_group",
					"instance_volume",
					"instance_snapshot",
					"instance_image",
					"apple_silicon_server",
					"baremetal_server",
					"baremetal_setting",
					"ipam_ip",
					"sbs_volume",
					"sbs_snapshot",
					"load_balancer_lb",
					"load_balancer_ip",
					"load_balancer_frontend",
					"load_balancer_backend",
					"load_balancer_route",
					"load_balancer_acl",
					"load_balancer_certificate",
					"sfs_filesystem",
					"vpc_private_network",
				},
			},
			{
				Name:       "method-name",
				Short:      `(Optional) Name of the method of the API call performed`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status",
				Short:      `(Optional) HTTP status code of the request. Returns either ` + "`" + `200` + "`" + ` if the request was successful or ` + "`" + `403` + "`" + ` if the permission was denied`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "recorded-after",
				Short:      `(Optional) The ` + "`" + `recorded_after` + "`" + ` parameter defines the earliest timestamp from which Audit Trail events are retrieved. Returns ` + "`" + `one hour ago` + "`" + ` by default`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "recorded-before",
				Short:      `(Optional) The ` + "`" + `recorded_before` + "`" + ` parameter defines the latest timestamp up to which Audit Trail events are retrieved. Returns ` + "`" + `now` + "`" + ` by default`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"recorded_at_desc",
					"recorded_at_asc",
				},
			},
			{
				Name:       "page-size",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "page-token",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "product-name",
				Short:      `(Optional) Name of the Scaleway product in a hyphenated format`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "service-name",
				Short:      `(Optional) Name of the service of the API call performed`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "resource-id",
				Short:      `(Optional) ID of the Scaleway resource`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "principal-id",
				Short:      `(Optional) ID of the User or IAM application at the origin of the event`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "source-ip",
				Short:      `(Optional) IP address at the origin of the event`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*audit_trail.ListEventsRequest)

			client := core.ExtractClient(ctx)
			api := audit_trail.NewAPI(client)

			return api.ListEvents(request)
		},
	}
}

func auditTrailProductList() *core.Command {
	return &core.Command{
		Short:     `Retrieve the list of Scaleway resources for which you have Audit Trail events`,
		Long:      `Retrieve the list of Scaleway resources for which you have Audit Trail events.`,
		Namespace: "audit-trail",
		Resource:  "product",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(audit_trail.ListProductsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*audit_trail.ListProductsRequest)

			client := core.ExtractClient(ctx)
			api := audit_trail.NewAPI(client)

			return api.ListProducts(request)
		},
	}
}
