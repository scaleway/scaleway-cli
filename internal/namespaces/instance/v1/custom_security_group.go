package instance

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Marshalers
//

var (
	securityGroupPolicyFormat = humanv2.EnumFormats[instance.SecurityGroupPolicy]{
		instance.SecurityGroupPolicyDrop:   humanv2.EnumFormat{Color: color.FgRed},
		instance.SecurityGroupPolicyAccept: humanv2.EnumFormat{Color: color.FgGreen},
	}

	securityGroupRuleActionFormat = humanv2.EnumFormats[instance.SecurityGroupRuleAction]{
		instance.SecurityGroupRuleActionDrop:   humanv2.EnumFormat{Color: color.FgRed},
		instance.SecurityGroupRuleActionAccept: humanv2.EnumFormat{Color: color.FgGreen},
	}

	securityGroupStateFormat = humanv2.EnumFormats[instance.SecurityGroupState]{
		instance.SecurityGroupStateAvailable:    humanv2.EnumFormat{Color: color.FgGreen},
		instance.SecurityGroupStateSyncing:      humanv2.EnumFormat{Color: color.FgBlue},
		instance.SecurityGroupStateSyncingError: humanv2.EnumFormat{Color: color.FgRed},
	}
)

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
