package instance

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
)

//
// Marshalers
//

var (
	securityGroupActionMarshalSpecs = human.EnumMarshalSpecs{
		instance.SecurityGroupActionDrop:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		instance.SecurityGroupActionAccept: &human.EnumMarshalSpec{Attribute: color.FgGreen},
	}

	securityGroupRuleActionMarshalSpecs = human.EnumMarshalSpecs{
		instance.SecurityGroupRuleActionDrop:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		instance.SecurityGroupRuleActionAccept: &human.EnumMarshalSpec{Attribute: color.FgGreen},
	}

	//securityGroupStateMarshalSpecs = human.EnumMarshalSpecs{
	//	instance.SecurityGroupStateAvailable:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
	//	instance.SecurityGroupStateSyncing:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	//	instance.SecurityGroupStateSyncingError: &human.EnumMarshalSpec{Attribute: color.FgRed},
	//}
)

func securityGroupMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type humanSecurityGroup instance.SecurityGroup
	securityGroup := humanSecurityGroup(i.(instance.SecurityGroup))

	// Sections
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "Tags",
			Title:     "Tags",
		},
		{
			FieldName: "Rules",
			Title:     "Rules",
		},
		{
			FieldName: "DefaultRules",
			Title:     "Default Rules",
		},
	}

	str, err := human.Marshal(securityGroup, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

//
// Builders
//

func securityGroupCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("name").Required = true
	c.ArgSpecs.GetByName("name").Default = core.RandomValueGenerator("sg")
	c.ArgSpecs.GetByName("inbound-default-action").Required = true
	c.ArgSpecs.GetByName("inbound-default-action").Default = core.DefaultValueSetter(
		instance.SecurityGroupActionAccept.String(),
	)
	c.ArgSpecs.GetByName("outbound-default-action").Required = true
	c.ArgSpecs.GetByName("outbound-default-action").Default = core.DefaultValueSetter(
		instance.SecurityGroupActionAccept.String(),
	)
	c.ArgSpecs.GetByName("stateless").Default = core.DefaultValueSetter("false")
	c.ArgSpecs.GetByName("disable-default-rules").Default = core.DefaultValueSetter("false")

	return c
}

func securityGroupGetBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("security-group-id").Positional = true

	return c
}

func securityGroupUpdateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("security-group-id").Positional = true

	return c
}

func securityGroupDeleteBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("security-group-id").Positional = true

	return c
}

func securityGroupListBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}

		sgList, ok := rawResp.(*instance.ListSecurityGroupsResponse)
		if !ok {
			return rawResp, fmt.Errorf(
				"expected response of type *instance.ListSecurityGroupsResponse, got %T",
				rawResp,
			)
		}

		return sgList.SecurityGroups, nil
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
				Label:     "DESCRIPTION",
				FieldName: "Description",
			},
			{
				Label:     "TAGS",
				FieldName: "Tags",
			},
			{
				Label:     "INBOUND DEF. ACT.",
				FieldName: "InboundDefaultAction",
			},
			{
				Label:     "OUTBOUND DEF. ACT.",
				FieldName: "OutboundDefaultAction",
			},
			{
				Label:     "STATELESS",
				FieldName: "Stateless",
			},
			{
				Label:     "PROJECT DEFAULT",
				FieldName: "ProjectDefault",
			},
			{
				Label:     "PROJECT ID",
				FieldName: "ProjectID",
			},
			{
				Label:     "ZONE",
				FieldName: "Zone",
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
