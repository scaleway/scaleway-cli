// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package autoscaling

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	autoscaling "github.com/scaleway/scaleway-sdk-go/api/autoscaling/v1alpha2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		autoscalingRoot(),
		autoscalingGroup(),
		autoscalingLogs(),
		autoscalingServers(),
		autoscalingAlerts(),
		autoscalingGroupList(),
		autoscalingGroupGet(),
		autoscalingGroupCreate(),
		autoscalingGroupUpdate(),
		autoscalingGroupDelete(),
		autoscalingLogsList(),
		autoscalingServersList(),
		autoscalingAlertsList(),
	)
}

func autoscalingRoot() *core.Command {
	return &core.Command{
		Short:     `Autoscaling Groups API`,
		Long:      ``,
		Namespace: "autoscaling",
	}
}

func autoscalingGroup() *core.Command {
	return &core.Command{
		Short: `Groups management commands`,
		Long: `Autoscaling groups automatically adjust the number of Instances based on metrics like CPU or memory usage.
This allows you to automatically scale your infrastructure up or down to meet demand while optimizing costs.`,
		Namespace: "autoscaling",
		Resource:  "group",
	}
}

func autoscalingLogs() *core.Command {
	return &core.Command{
		Short: `Logs management commands`,
		Long: `Logs provide visibility into the autoscaling group activities and events.
You can query logs to understand scaling decisions and troubleshoot issues.`,
		Namespace: "autoscaling",
		Resource:  "logs",
	}
}

func autoscalingServers() *core.Command {
	return &core.Command{
		Short:     `Servers management commands`,
		Long:      `List the Instances that belong to an autoscaling group.`,
		Namespace: "autoscaling",
		Resource:  "servers",
	}
}

func autoscalingAlerts() *core.Command {
	return &core.Command{
		Short:     `Alerts management commands`,
		Long:      `Alerts notify you of issues affecting your autoscaling groups, such as quota limits, stock issues, or configuration problems.`,
		Namespace: "autoscaling",
		Resource:  "alerts",
	}
}

func autoscalingGroupList() *core.Command {
	return &core.Command{
		Short:     `List autoscaling groups`,
		Long:      `List all autoscaling groups in a project.`,
		Namespace: "autoscaling",
		Resource:  "group",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[autoscaling.ListGroupsRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order criteria for listing groups`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
				},
			},
			{
				Name:       "page-size",
				Short:      `Page size for pagination`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "page-token",
				Short:      `Token for pagination`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "template-id",
				Short:      `Template ID to filter groups`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "load-balancer-id",
				Short:      `Load balancer ID to filter groups`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*autoscaling.ListGroupsRequest)

			client := core.ExtractClient(ctx)
			api := autoscaling.NewAPI(client)

			return api.ListGroups(request, scw.WithContext(ctx))
		},
		Examples: []*core.Example{
			{
				Short:    "List all autoscaling groups in your default zone",
				ArgsJSON: `null`,
			},
			{
				Short:    "List autoscaling groups filtered by template",
				ArgsJSON: `{"template_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "List autoscaling groups filtered by load balancer",
				ArgsJSON: `{"load_balancer_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func autoscalingGroupGet() *core.Command {
	return &core.Command{
		Short: `Get an autoscaling group`,
		Long: `Get details of a specified autoscaling group including its
configuration, current size, and status.`,
		Namespace: "autoscaling",
		Resource:  "group",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[autoscaling.GetGroupRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `ID of the group to get`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*autoscaling.GetGroupRequest)

			client := core.ExtractClient(ctx)
			api := autoscaling.NewAPI(client)

			return api.GetGroup(request, scw.WithContext(ctx))
		},
		Examples: []*core.Example{
			{
				Short:    "Get a specified autoscaling group",
				ArgsJSON: `{"group_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func autoscalingGroupCreate() *core.Command {
	return &core.Command{
		Short: `Create an autoscaling group`,
		Long: `Create a new autoscaling group with the specified configuration
including template, scaling policy, and optional load balancer
settings.`,
		Namespace: "autoscaling",
		Resource:  "group",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[autoscaling.CreateGroupRequest](),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the autoscaling group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags associated with the group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "template-id",
				Short:      `Template ID for instances in this group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.minimum-size",
				Short:      `Minimum number of instances in the group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.maximum-size",
				Short:      `Maximum number of instances in the group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.scale-out-cooldown",
				Short:      `Cooldown period after a scale out event`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.scale-in-cooldown",
				Short:      `Cooldown period after a scale in event`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.scale-in-step",
				Short:      `Number of instances to remove in a single scale in event`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.scale-out-step",
				Short:      `Number of instances to add in a single scale out event`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.fixed-size.size",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.cpu-target.target-avg-percent",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.memory-target.target-avg-percent",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "load-balancer-configuration-spec.load-balancer-id",
				Short:      `ID of the load balancer (set to empty to disable)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "load-balancer-configuration-spec.backends.{index}.backend-id",
				Short:      `ID of the load balancer backend`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "load-balancer-configuration-spec.backends.{index}.address-family",
				Short:      `IP address family (IPv4 or IPv6)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_address_family",
					"ipv4",
					"ipv6",
				},
			},
			{
				Name:       "load-balancer-configuration-spec.backends.{index}.private-network-id",
				Short:      `Optional private network ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "load-balancer-configuration-spec.auto-healing.enabled",
				Short:      `Whether auto-healing is enabled`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "load-balancer-configuration-spec.auto-healing.grace-period",
				Short:      `Grace period for health checks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*autoscaling.CreateGroupRequest)

			client := core.ExtractClient(ctx)
			api := autoscaling.NewAPI(client)

			return api.CreateGroup(request, scw.WithContext(ctx))
		},
		Examples: []*core.Example{
			{
				Short: "Create an autoscaling group with a fixed size",
				Raw:   `scw autoscaling group create scaling-policy-spec.fixed-size.size=3 scaling-policy-spec.maximum-size=5 scaling-policy-spec.minimum-size=1 name=my-autoscaling-group template-id=11111111-1111-1111-1111-111111111111 project-id=11111111-1111-1111-1111-111111111111`,
			},
		},
	}
}

func autoscalingGroupUpdate() *core.Command {
	return &core.Command{
		Short: `Update an autoscaling group`,
		Long: `Update the configuration of a specified autoscaling group including
name, tags, template, scaling policy, and load balancer settings.`,
		Namespace: "autoscaling",
		Resource:  "group",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[autoscaling.UpdateGroupRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `ID of the group to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `New name for the group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `New tags for the group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "template-id",
				Short:      `New template ID for instances`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.minimum-size",
				Short:      `Minimum number of instances in the group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.maximum-size",
				Short:      `Maximum number of instances in the group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.scale-out-cooldown",
				Short:      `Cooldown period after a scale out event`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.scale-in-cooldown",
				Short:      `Cooldown period after a scale in event`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.scale-in-step",
				Short:      `Number of instances to remove in a single scale in event`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.scale-out-step",
				Short:      `Number of instances to add in a single scale out event`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.fixed-size.size",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.cpu-target.target-avg-percent",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-policy-spec.memory-target.target-avg-percent",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "load-balancer-configuration-spec.load-balancer-id",
				Short:      `ID of the load balancer (set to empty to disable)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "load-balancer-configuration-spec.backends.{index}.backend-id",
				Short:      `ID of the load balancer backend`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "load-balancer-configuration-spec.backends.{index}.address-family",
				Short:      `IP address family (IPv4 or IPv6)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_address_family",
					"ipv4",
					"ipv6",
				},
			},
			{
				Name:       "load-balancer-configuration-spec.backends.{index}.private-network-id",
				Short:      `Optional private network ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "load-balancer-configuration-spec.auto-healing.enabled",
				Short:      `Whether auto-healing is enabled`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "load-balancer-configuration-spec.auto-healing.grace-period",
				Short:      `Grace period for health checks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*autoscaling.UpdateGroupRequest)

			client := core.ExtractClient(ctx)
			api := autoscaling.NewAPI(client)

			return api.UpdateGroup(request, scw.WithContext(ctx))
		},
		Examples: []*core.Example{
			{
				Short:    "Update the name of an autoscaling group",
				ArgsJSON: `{"group_id":"11111111-1111-1111-1111-111111111111","name":"new-name"}`,
			},
			{
				Short: "Update the scaling policy of an autoscaling group",
				Raw:   `scw autoscaling group update group-id=11111111-1111-1111-1111-111111111111 scaling_policy_spec.cpu_target.target_avg_percent=70`,
			},
		},
	}
}

func autoscalingGroupDelete() *core.Command {
	return &core.Command{
		Short: `Delete an autoscaling group`,
		Long: `Delete a specified autoscaling group and all its associated
resources.`,
		Namespace: "autoscaling",
		Resource:  "group",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[autoscaling.DeleteGroupRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `ID of the group to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*autoscaling.DeleteGroupRequest)

			client := core.ExtractClient(ctx)
			api := autoscaling.NewAPI(client)

			return api.DeleteGroup(request, scw.WithContext(ctx))
		},
		Examples: []*core.Example{
			{
				Short:    "Delete a specified autoscaling group",
				ArgsJSON: `{"group_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func autoscalingLogsList() *core.Command {
	return &core.Command{
		Short: `List autoscaling group logs`,
		Long: `List logs for a specified autoscaling group to view scaling events
and activities.`,
		Namespace: "autoscaling",
		Resource:  "logs",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[autoscaling.ListLogsRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
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
				Name:       "page-size",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "start-time",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "end-time",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*autoscaling.ListLogsRequest)

			client := core.ExtractClient(ctx)
			api := autoscaling.NewAPI(client)

			return api.ListLogs(request, scw.WithContext(ctx))
		},
		Examples: []*core.Example{
			{
				Short:    "List all logs for an autoscaling group",
				ArgsJSON: `{"group_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func autoscalingServersList() *core.Command {
	return &core.Command{
		Short:     `List autoscaling group servers`,
		Long:      `List all Instances belonging to a specified autoscaling group.`,
		Namespace: "autoscaling",
		Resource:  "servers",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[autoscaling.ListServersRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
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
				Name:       "page-size",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*autoscaling.ListServersRequest)

			client := core.ExtractClient(ctx)
			api := autoscaling.NewAPI(client)

			return api.ListServers(request, scw.WithContext(ctx))
		},
		Examples: []*core.Example{
			{
				Short:    "List all servers in an autoscaling group",
				ArgsJSON: `{"group_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func autoscalingAlertsList() *core.Command {
	return &core.Command{
		Short:     `List autoscaling group alerts`,
		Long:      `List active and historical alerts for a specified autoscaling group.`,
		Namespace: "autoscaling",
		Resource:  "alerts",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[autoscaling.ListAlertsRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
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
				Name:       "page-size",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*autoscaling.ListAlertsRequest)

			client := core.ExtractClient(ctx)
			api := autoscaling.NewAPI(client)

			return api.ListAlerts(request, scw.WithContext(ctx))
		},
		Examples: []*core.Example{
			{
				Short:    "List all alerts for an autoscaling group",
				ArgsJSON: `{"group_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
