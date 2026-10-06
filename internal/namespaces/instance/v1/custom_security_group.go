package instance

import (
	"strconv"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Marshalers
//

var securityGroupRuleActionMarshalSpecs = human.EnumMarshalSpecs{
	instance.SecurityGroupRuleActionDrop:   &human.EnumMarshalSpec{Attribute: color.FgRed},
	instance.SecurityGroupRuleActionAccept: &human.EnumMarshalSpec{Attribute: color.FgGreen},
}

func marshalSecurityGroupRules(i any, _ *human.MarshalOpt) (out string, err error) {
	rules := i.([]*instance.SecurityGroupRule)

	type humanRule struct {
		ID        string
		Direction string
		Protocol  instance.SecurityGroupRuleProtocol
		Action    instance.SecurityGroupRuleAction
		IPRange   string
		Dest      string
	}

	toHumanRule := func(rule *instance.SecurityGroupRule) *humanRule {
		dest := "ALL"
		if rule.DestPortFrom != nil {
			dest = strconv.Itoa(int(*rule.DestPortFrom))
		}
		if rule.DestPortTo != nil {
			dest += "-" + strconv.Itoa(int(*rule.DestPortTo))
		}

		return &humanRule{
			ID:        rule.ID,
			Direction: string(rule.Direction),
			Protocol:  rule.Protocol,
			Action:    rule.Action,
			IPRange:   rule.IPRange.String(),
			Dest:      dest,
		}
	}
	humanRules := make([]*humanRule, len(rules))

	for i, rule := range rules {
		humanRules[i] = toHumanRule(rule)
	}

	return human.Marshal(humanRules, nil)
}

//
// Builders
//

func securityGroupListDefaultRulesBuilder(c *core.Command) *core.Command {
	deprecateV1OnlyCommand(c)
	c.DeprecationMessage += " Default rules can be listed using \"scw instance security-group get <security-group-id>\"."

	return c
}

func securityGroupListRulesBuilder(c *core.Command) *core.Command {
	deprecateV1OnlyCommand(c)
	c.DeprecationMessage += " Rules can be listed using \"scw instance security-group get <security-group-id>\"."

	return c
}

func securityGroupCreateRuleBuilder(c *core.Command) *core.Command {
	deprecateV1OnlyCommand(c)
	c.DeprecationMessage += " Rules can be created using \"scw instance security-group add-rules <security-group-id ...> [arg=value ...]\"." // check usage

	return c
}

func securityGroupDeleteRuleBuilder(c *core.Command) *core.Command {
	deprecateV1OnlyCommand(c)
	c.DeprecationMessage += " Rules can be deleted using \"scw instance security-group delete-rules security-group-rule-ids.0=<security-group-rule-id>\"." // check usage

	return c
}

func securityGroupGetRuleBuilder(c *core.Command) *core.Command {
	deprecateV1OnlyCommand(c)
	c.DeprecationMessage += " Rules can be listed using \"scw instance security-group get <security-group-id>\"."

	return c
}
