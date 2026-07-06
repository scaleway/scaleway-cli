// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package search

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	search "github.com/scaleway/scaleway-sdk-go/api/search/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		searchRoot(),
		searchResource(),
		searchResourceSearch(),
	)
}

func searchRoot() *core.Command {
	return &core.Command{
		Short:     `Search API`,
		Long:      ``,
		Namespace: "search",
	}
}

func searchResource() *core.Command {
	return &core.Command{
		Short:     `Resource search commands`,
		Long:      `Resource search commands.`,
		Namespace: "search",
		Resource:  "resource",
	}
}

func searchResourceSearch() *core.Command {
	return &core.Command{
		Short:     `Search API`,
		Long:      `Search API.`,
		Namespace: "search",
		Resource:  "resource",
		Verb:      "search",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(search.SearchResourcesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "query",
				Short:      `Search query`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "project-ids.{index}",
				Short:      `List of Project IDs to filter the resources by`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "types.{index}",
				Short:      `List of resource types to filter the resources by`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"instance_server",
					"instance_volume",
					"instance_image",
					"instance_security_group",
					"instance_private_nic",
					"instance_snapshot",
					"instance_placement_group",
					"k8s_cluster",
					"k8s_pool",
					"k8s_node",
					"domain_domain",
					"dns_zone",
					"vpc_private_network",
					"vpc_vpc",
					"vpg_gateway",
					"apple_silicon_server",
					"rdb_instance",
					"rdb_snapshot",
					"rdb_backup",
					"baremetal_server",
					"tem_domain",
					"lb_server",
					"serverless_functions_function",
					"serverless_containers_container",
					"wbh_hosting",
					"redis_cluster",
					"sm_secret",
					"kms_key",
					"edg_pipeline",
					"mnq_nats_account",
					"sbs_volume",
					"sbs_snapshot",
					"serverless_job_definition",
					"serverless_sqldb_database",
					"serverless_sqldb_backup",
					"ddl_datalab",
					"mgdb_instance",
					"mgdb_snapshot",
					"ifr_deployment",
					"ifr_model",
					"gapi_batch",
					"dtwh_deployment",
					"obs_datasource",
					"obs_exporter",
					"svpn_vpn_gateway",
					"svpn_customer_gateway",
					"svpn_connection",
					"svpn_routing_policy",
					"kafk_cluster",
					"iam_api_key",
					"iam_application",
					"iam_user",
					"iam_group",
					"iam_policy",
					"sedb_cluster",
					"autoscaling_group",
				},
			},
			{
				Name:       "localities.{index}",
				Short:      `List of scopes (zones, regions, or global) to filter the resources by`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_locality",
					"global",
					"fr_rz",
					"fr_srr",
					"fr_srr_1",
					"fr_par",
					"fr_par_1",
					"fr_par_2",
					"fr_par_3",
					"fr_par_4",
					"nl_ams",
					"nl_ams_1",
					"nl_ams_2",
					"nl_ams_3",
					"pl_waw",
					"pl_waw_1",
					"pl_waw_2",
					"pl_waw_3",
					"fr_int",
					"fr_int_1",
					"fr_lab",
					"fr_lab_1",
					"it_mil",
					"it_mil_1",
				},
			},
			{
				Name:       "created-after",
				Short:      `Filter resources created after this timestamp`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "created-before",
				Short:      `Filter resources created before this timestamp`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "modified-after",
				Short:      `Filter resources modified after this timestamp`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "modified-before",
				Short:      `Filter resources modified before this timestamp`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*search.SearchResourcesRequest)

			client := core.ExtractClient(ctx)
			api := search.NewAPI(client)

			return api.SearchResources(request)
		},
	}
}
