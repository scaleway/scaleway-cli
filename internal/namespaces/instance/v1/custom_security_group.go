package instance

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Marshalers
//

var (
	securityGroupPolicyMarshalSpecs = human.EnumMarshalSpecs{
		instance.SecurityGroupPolicyDrop:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		instance.SecurityGroupPolicyAccept: &human.EnumMarshalSpec{Attribute: color.FgGreen},
	}

	securityGroupRuleActionMarshalSpecs = human.EnumMarshalSpecs{
		instance.SecurityGroupRuleActionDrop:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		instance.SecurityGroupRuleActionAccept: &human.EnumMarshalSpec{Attribute: color.FgGreen},
	}

	securityGroupStateMarshalSpecs = human.EnumMarshalSpecs{
		instance.SecurityGroupStateAvailable:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		instance.SecurityGroupStateSyncing:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
		instance.SecurityGroupStateSyncingError: &human.EnumMarshalSpec{Attribute: color.FgRed},
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

// MarshalHuman marshals a customSecurityGroupResponse.
func (sg *customSecurityGroupResponse) MarshalHuman() (out string, err error) {
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

type customSecurityGroupResponse struct {
	instance.SecurityGroup

	Rules []*instance.SecurityGroupRule
}

func securityGroupCreateBuilder(c *core.Command) *core.Command {
	type customCreateSecurityGroupRequest struct {
		*instance.CreateSecurityGroupRequest
		OrganizationID *string
		ProjectID      *string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customCreateSecurityGroupRequest{})

	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
			args := argsI.(*customCreateSecurityGroupRequest)

			if args.CreateSecurityGroupRequest == nil {
				args.CreateSecurityGroupRequest = &instance.CreateSecurityGroupRequest{}
			}

			request := args.CreateSecurityGroupRequest
			request.Organization = args.OrganizationID
			request.Project = args.ProjectID

			return runner(ctx, request)
		},
	)

	return c
}

func securityGroupGetBuilder(c *core.Command) *core.Command {
	c.Run = func(ctx context.Context, argsI any) (any, error) {
		req := argsI.(*instance.GetSecurityGroupRequest)

		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)
		securityGroup, err := api.GetSecurityGroup(req)
		if err != nil {
			return nil, err
		}

		securityGroupRules, err := api.ListSecurityGroupRules(
			&instance.ListSecurityGroupRulesRequest{
				Zone:            req.Zone,
				SecurityGroupID: securityGroup.SecurityGroup.ID,
			},
			scw.WithAllPages(),
		)
		if err != nil {
			return nil, err
		}

		return &customSecurityGroupResponse{
			SecurityGroup: *securityGroup.SecurityGroup,
			Rules:         securityGroupRules.Rules,
		}, nil
	}

	return c
}

func securityGroupListBuilder(c *core.Command) *core.Command {
	type customListSecurityGroupsRequest struct {
		*instance.ListSecurityGroupsRequest
		OrganizationID *string
		ProjectID      *string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customListSecurityGroupsRequest{})

	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
			args := argsI.(*customListSecurityGroupsRequest)

			if args.ListSecurityGroupsRequest == nil {
				args.ListSecurityGroupsRequest = &instance.ListSecurityGroupsRequest{}
			}

			request := args.ListSecurityGroupsRequest
			request.Organization = args.OrganizationID
			request.Project = args.ProjectID

			return runner(ctx, request)
		},
	)

	return c
}

func securityGroupDeleteBuilder(c *core.Command) *core.Command {
	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
			res, originalErr := runner(ctx, argsI)
			if originalErr == nil {
				return res, nil
			}

			if strings.HasSuffix(originalErr.Error(), "group is in use. you cannot delete it.") {
				req := argsI.(*instance.DeleteSecurityGroupRequest)
				api := instance.NewAPI(core.ExtractClient(ctx))

				newError := &core.CliError{
					Err: errors.New("cannot delete security-group currently in use"),
				}

				// Get security-group.
				sg, err := api.GetSecurityGroup(&instance.GetSecurityGroupRequest{
					SecurityGroupID: req.SecurityGroupID,
				})
				if err != nil {
					// Ignore API error and return a minimal error.
					return nil, newError
				}

				// Create detail message.
				var hintBuilder strings.Builder
				hintBuilder.WriteString(
					"Attach all these instances to another security-group before deleting this one:",
				)
				for _, s := range sg.SecurityGroup.Servers {
					fmt.Fprintf(
						&hintBuilder,
						"\nscw instance server update %s security-group.id=$NEW_SECURITY_GROUP_ID",
						s.ID,
					)
				}

				newError.Hint = hintBuilder.String()

				return nil, newError
			}

			return nil, originalErr
		},
	)

	return c
}

//
// Commands
//

type instanceResetSecurityGroupArgs struct {
	Zone            scw.Zone
	SecurityGroupID string
}

func securityGroupClearCommand() *core.Command {
	return &core.Command{
		Short:     `Remove all rules of a security group`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "clear",
		ArgsType:  reflect.TypeOf(instanceResetSecurityGroupArgs{}),
		Examples: []*core.Example{
			{
				Short:    "Remove all rules of the given security group",
				ArgsJSON: `{"security_group_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
		Run: func(ctx context.Context, argsI any) (i any, e error) {
			args := argsI.(*instanceResetSecurityGroupArgs)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			rules, err := api.ListSecurityGroupRules(&instance.ListSecurityGroupRulesRequest{
				Zone:            args.Zone,
				SecurityGroupID: args.SecurityGroupID,
			}, scw.WithAllPages())
			if err != nil {
				return nil, err
			}

			for _, rule := range rules.Rules {
				if !rule.Editable {
					continue
				}
				err = api.DeleteSecurityGroupRule(&instance.DeleteSecurityGroupRuleRequest{
					Zone:                args.Zone,
					SecurityGroupID:     args.SecurityGroupID,
					SecurityGroupRuleID: rule.ID,
				})
				if err != nil {
					return nil, err
				}
			}

			return &core.SuccessResult{Message: "Successful reset of the security group rules"}, err
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "security-group-id",
				Short:    `ID of the security group to reset.`,
				Required: true,
			},
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
	}
}

var instanceSecurityGroupEditYamlExample = `rules:
- action: drop
  dest_port_from: 1200
  dest_port_to: 1300
  direction: inbound
  ip_range: 192.168.0.0/24
  protocol: TCP
- action: drop
  direction: inbound
  protocol: ICMP
  ip_range: 0.0.0.0/0
- action: accept
  dest_port_from: 25565
  direction: outbound
  ip_range: 0.0.0.0/0
  protocol: UDP
`

type instanceSecurityGroupEditArgs struct {
	Zone            scw.Zone
	SecurityGroupID string
	Mode            editor.MarshalMode
}

func securityGroupEditCommand() *core.Command {
	return &core.Command{
		Short:     `Edit all rules of a security group`,
		Long:      editor.LongDescription,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "edit",
		ArgsType:  reflect.TypeOf(instanceSecurityGroupEditArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "security-group-id",
				Short:      `ID of the security group to reset.`,
				Required:   true,
				Positional: true,
			},
			editor.MarshalModeArgSpec(),
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
		Run: func(ctx context.Context, argsI any) (i any, e error) {
			args := argsI.(*instanceSecurityGroupEditArgs)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			rules, err := api.ListSecurityGroupRules(&instance.ListSecurityGroupRulesRequest{
				Zone:            args.Zone,
				SecurityGroupID: args.SecurityGroupID,
			}, scw.WithAllPages(), scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to list security-group rules: %w", err)
			}

			// Get only rules that can be edited
			editableRules := []*instance.SecurityGroupRule(nil)
			for _, rule := range rules.Rules {
				if rule.Editable {
					editableRules = append(editableRules, rule)
				}
			}
			rules.Rules = editableRules

			setRequest := &instance.SetSecurityGroupRulesRequest{
				Zone:            args.Zone,
				SecurityGroupID: args.SecurityGroupID,
			}

			editedSetRequest, err := editor.UpdateResourceEditor(rules, setRequest, &editor.Config{
				PutRequest:   true,
				MarshalMode:  args.Mode,
				Template:     instanceSecurityGroupEditYamlExample,
				IgnoreFields: []string{"editable"},
			})
			if err != nil {
				return nil, err
			}

			setRequest = editedSetRequest.(*instance.SetSecurityGroupRulesRequest)

			resp, err := api.SetSecurityGroupRules(setRequest, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			return resp.Rules, nil
		},
	}
}
