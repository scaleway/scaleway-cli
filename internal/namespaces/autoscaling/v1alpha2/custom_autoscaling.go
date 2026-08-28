package autoscaling

import (
	"context"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	autoscaling "github.com/scaleway/scaleway-sdk-go/api/autoscaling/v1alpha2"
)

func groupMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp autoscaling.Group
	template := tmp(i.(autoscaling.Group))

	opt.Sections = []*human.MarshalSection{
		{FieldName: "Tags", HideIfEmpty: true},
		{FieldName: "ScalingPolicy", Title: "Scaling Policy"},
		{FieldName: "OpenAlerts", HideIfEmpty: true, Title: "Open Alerts"},
		{
			FieldName:   "LoadBalancerConfiguration",
			HideIfEmpty: true,
			Title:       "LoadBalancer Configuration",
		},
	}

	str, err := human.Marshal(template, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func GroupCreateBuilder(c *core.Command) *core.Command {
	// Both minimum-size and maximum-size are required
	c.ArgSpecs.GetByName("scaling-policy-spec.minimum-size").Required = true
	c.ArgSpecs.GetByName("scaling-policy-spec.maximum-size").Required = true

	// One of fixed-size, cpu-target or memory-target must be set
	c.ArgSpecs.GetByName("scaling-policy-spec.fixed-size.size").OneOfGroup = "scaling-policy-target"
	c.ArgSpecs.GetByName("scaling-policy-spec.cpu-target.target-avg-percent").OneOfGroup = "scaling-policy-target"
	c.ArgSpecs.GetByName("scaling-policy-spec.memory-target.target-avg-percent").OneOfGroup = "scaling-policy-target"
	c.ArgSpecs.GetByName("scaling-policy-spec.fixed-size.size").Required = true
	c.ArgSpecs.GetByName("scaling-policy-spec.cpu-target.target-avg-percent").Required = true
	c.ArgSpecs.GetByName("scaling-policy-spec.memory-target.target-avg-percent").Required = true

	// Name is required by the API, so the CLI generates one if none is provided
	c.ArgSpecs.GetByName("name").Default = core.RandomValueGenerator("asg")

	return c
}

func GroupUpdateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("group-id").Positional = true

	return c
}

func GroupGetBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("group-id").Positional = true

	return c
}

func GroupDeleteBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("group-id").Positional = true

	return c
}

func GroupListBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}

		groupList, ok := rawResp.(*autoscaling.ListGroupsResponse)
		if !ok {
			return "", fmt.Errorf(
				"expected response of type *autoscaling.ListGroupsResponse, got %T",
				rawResp,
			)
		}

		return groupList.GroupSummaries, nil
	}

	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "ID",
				FieldName: "ID",
			},
			{
				Label:     "NAME",
				FieldName: "Name",
			},
			{
				Label:     "TEMPLATE ID",
				FieldName: "TemplateID",
			},
			{
				Label:     "TAGS",
				FieldName: "Tags",
			},
			{
				Label:     "STATUS",
				FieldName: "Status",
			},
			{
				Label:     "LATEST OPEN ALERT",
				FieldName: "LatestOpenAlert.Type",
			},
			{
				Label:     "SIZE",
				FieldName: "CurrentSize",
			},
			{
				Label:     "MIN SIZE",
				FieldName: "MinimumSize",
			},
			{
				Label:     "MAX SIZE",
				FieldName: "MaximumSize",
			},
			{
				Label:     "ZONE",
				FieldName: "Zone",
			},

			{
				Label:     "LOAD BALANCER ID",
				FieldName: "LoadBalancerID",
			},
			{
				Label:     "TARGET TYPE",
				FieldName: "ScalingPolicyTargetType",
			},
			{
				Label:     "PROJECT ID",
				FieldName: "ProjectID",
			},
			{
				Label:     "CREATED AT",
				FieldName: "CreatedAt",
			},
			{
				Label:     "UPDATED AT",
				FieldName: "UpdatedAt",
			},
		},
	}

	return c
}

func AlertsListBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}

		alertList, ok := rawResp.(*autoscaling.ListAlertsResponse)
		if !ok {
			return "", fmt.Errorf(
				"expected response of type *autoscaling.ListAlertsResponse, got %T",
				rawResp,
			)
		}

		return alertList.Alerts, nil
	}

	return c
}

func LogsListBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("group-id").Required = true
	c.ArgSpecs.GetByName("group-id").Positional = true

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}

		logList, ok := rawResp.(*autoscaling.ListLogsResponse)
		if !ok {
			return "", fmt.Errorf(
				"expected response of type *autoscaling.ListLogsResponse, got %T",
				rawResp,
			)
		}

		return logList.Logs, nil
	}

	return c
}

func ServersListBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("group-id").Required = true
	c.ArgSpecs.GetByName("group-id").Positional = true

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}

		serverList, ok := rawResp.(*autoscaling.ListServersResponse)
		if !ok {
			return "", fmt.Errorf(
				"expected response of type *autoscaling.ListServersResponse, got %T",
				rawResp,
			)
		}

		return serverList.Servers, nil
	}

	return c
}
