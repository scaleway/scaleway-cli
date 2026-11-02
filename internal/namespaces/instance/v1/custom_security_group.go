package instance

import (
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
)

//
// Marshalers
//

var securityGroupRuleActionMarshalSpecs = human.EnumMarshalSpecs{
	instance.SecurityGroupRuleActionDrop:   &human.EnumMarshalSpec{Attribute: color.FgRed},
	instance.SecurityGroupRuleActionAccept: &human.EnumMarshalSpec{Attribute: color.FgGreen},
}

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

// FormatHuman marshals a customSecurityGroupResponse as a single title-string
func (sg *customSecurityGroupResponse) FormatHuman() (out string, err error) {
	humanSecurityGroup := struct {
		ID                    string
		Name                  string
		State                 instance.SecurityGroupState
		Description           string
		EnableDefaultSecurity bool
		OrganizationID        string
		ProjectID             string
		OrganizationDefault   *bool
		ProjectDefault        bool
		CreationDate          *time.Time
		ModificationDate      *time.Time
		Stateful              bool
	}{
		ID:                    sg.ID,
		Name:                  sg.Name,
		State:                 sg.State,
		Description:           sg.Description,
		EnableDefaultSecurity: sg.EnableDefaultSecurity,
		OrganizationID:        sg.Organization,
		ProjectID:             sg.Project,
		OrganizationDefault:   sg.OrganizationDefault,
		ProjectDefault:        sg.ProjectDefault,
		CreationDate:          sg.CreationDate,
		ModificationDate:      sg.ModificationDate,
		Stateful:              sg.Stateful,
	}

	securityGroupView, err := human.Marshal(humanSecurityGroup, nil)
	if err != nil {
		return "", err
	}
	securityGroupView = terminal.Style("Security Group:\n", color.Bold) + securityGroupView

	inboundRules := []*instance.SecurityGroupRule(nil)
	outboundRules := []*instance.SecurityGroupRule(nil)
	for _, rule := range sg.Rules {
		switch rule.Direction {
		case instance.SecurityGroupRuleDirectionInbound:
			inboundRules = append(inboundRules, rule)
		case instance.SecurityGroupRuleDirectionOutbound:
			outboundRules = append(outboundRules, rule)
		default:
			logger.Warningf("invalid security group rule direction: %v", rule.Direction)
		}
	}

	// defaultInboundPolicy will already be colored in green or red by the marshaler.
	defaultInboundPolicy, err := human.Marshal(sg.InboundDefaultPolicy, nil)
	if err != nil {
		return "", err
	}

	// defaultOutboundPolicy will already be colored in green or red by the marshaler.
	defaultOutboundPolicy, err := human.Marshal(sg.OutboundDefaultPolicy, nil)
	if err != nil {
		return "", err
	}

	// b returns the given string in bold.
	// For inboundRulesView and outboundRulesView, this function must be called for every
	// concatenated part of the string because of the color package escaping at the end of
	// a color resulting in a non-bold format after the default{In|Out}boundPolicy.
	b := color.New(color.Bold).SprintFunc()

	inboundRulesContent, err := human.Marshal(inboundRules, nil)
	if err != nil {
		return "", err
	}
	inboundRulesView := b(
		"Inbound Rules (default policy ",
	) + b(
		defaultInboundPolicy,
	) + b(
		"):\n",
	) + inboundRulesContent

	outboundRulesContent, err := human.Marshal(outboundRules, nil)
	if err != nil {
		return "", err
	}
	outboundRulesView := b(
		"Outbound Rules (default policy ",
	) + b(
		defaultOutboundPolicy,
	) + b(
		"):\n",
	) + outboundRulesContent

	serversContent, err := human.Marshal(sg.Servers, nil)
	if err != nil {
		return "", err
	}
	serversView := terminal.Style("Servers:\n", color.Bold) + serversContent

	return strings.Join([]string{
		securityGroupView,
		inboundRulesView,
		outboundRulesView,
		serversView,
	}, "\n\n"), nil
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
