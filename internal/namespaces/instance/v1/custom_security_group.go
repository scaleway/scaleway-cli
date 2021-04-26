package instance

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-cli/internal/terminal"
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
		OrganizationDefault   bool
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

	type humanRule struct {
		ID       string
		Protocol instance.SecurityGroupRuleProtocol
		Action   instance.SecurityGroupRuleAction
		IPRange  string
		Dest     string
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
			ID:       rule.ID,
			Protocol: rule.Protocol,
			Action:   rule.Action,
			IPRange:  rule.IPRange.String(),
			Dest:     dest,
		}
	}

	inboundRules := []*humanRule(nil)
	outboundRules := []*humanRule(nil)
	for _, rule := range sg.Rules {
		switch rule.Direction {
		case instance.SecurityGroupRuleDirectionInbound:
			inboundRules = append(inboundRules, toHumanRule(rule))
		case instance.SecurityGroupRuleDirectionOutbound:
			outboundRules = append(outboundRules, toHumanRule(rule))
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
	inboundRulesView := b("Inbound Rules (default policy ") + b(defaultInboundPolicy) + b("):\n") + inboundRulesContent

	outboundRulesContent, err := human.Marshal(outboundRules, nil)
	if err != nil {
		return "", err
	}
	outboundRulesView := b("Outbound Rules (default policy ") + b(defaultOutboundPolicy) + b("):\n") + outboundRulesContent

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

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		args := argsI.(*customCreateSecurityGroupRequest)

		if args.CreateSecurityGroupRequest == nil {
			args.CreateSecurityGroupRequest = &instance.CreateSecurityGroupRequest{}
		}

		request := args.CreateSecurityGroupRequest
		request.Organization = args.OrganizationID
		request.Project = args.ProjectID

		return runner(ctx, request)
	})

	return c
}

func securityGroupGetBuilder(c *core.Command) *core.Command {
	c.Run = func(ctx context.Context, argsI interface{}) (interface{}, error) {
		req := argsI.(*instance.GetSecurityGroupRequest)

		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)
		securityGroup, err := api.GetSecurityGroup(req)
		if err != nil {
			return nil, err
		}

		securityGroupRules, err := api.ListSecurityGroupRules(&instance.ListSecurityGroupRulesRequest{
			Zone:            req.Zone,
			SecurityGroupID: securityGroup.SecurityGroup.ID,
		}, scw.WithAllPages())
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

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		args := argsI.(*customListSecurityGroupsRequest)

		if args.ListSecurityGroupsRequest == nil {
			args.ListSecurityGroupsRequest = &instance.ListSecurityGroupsRequest{}
		}

		request := args.ListSecurityGroupsRequest
		request.Organization = args.OrganizationID
		request.Project = args.ProjectID

		return runner(ctx, request)
	})

	return c
}

func securityGroupDeleteBuilder(c *core.Command) *core.Command {
	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		res, originalErr := runner(ctx, argsI)
		if originalErr == nil {
			return res, nil
		}

		if strings.HasSuffix(originalErr.Error(), "group is in use. you cannot delete it.") {
			req := argsI.(*instance.DeleteSecurityGroupRequest)
			api := instance.NewAPI(core.ExtractClient(ctx))

			newError := &core.CliError{
				Err: fmt.Errorf("cannot delete security-group currently in use"),
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
			hint := "Attach all these instances to another security-group before deleting this one:"
			for _, s := range sg.SecurityGroup.Servers {
				hint += "\nscw instance server update " + s.ID + " security-group.id=$NEW_SECURITY_GROUP_ID"
			}

			newError.Hint = hint
			return nil, newError
		}

		return nil, originalErr
	})
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
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
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
			core.ZoneArgSpec(),
		},
	}
}

// securityGroupUpdateCommand updates a security-group.
// Because the API for SecurityGroup works with a PUT but not a PATCH,
// the method UpdateSecurityGroup() is not generated.
// Instead, setSecurityGroup() is generated, and a custom UpdateSecurityGroup() method is handwritten the SDK.
// This is why 'scw instance security-group update' needs to be written by hand.
func securityGroupUpdateCommand() *core.Command {
	return &core.Command{
		Short:     `Update security group`,
		Long:      `Update security group.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(instance.UpdateSecurityGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "security-group-id",
				Short:    `ID of the security group to update`,
				Required: true,
			},
			{
				Name: "name",
			},
			{
				Name: "description",
			},
			{
				Name: "stateful",
			},
			{
				Name:       "inbound-default-policy",
				EnumValues: []string{"accept", "drop"},
			},
			{
				Name:       "outbound-default-policy",
				EnumValues: []string{"accept", "drop"},
			},
			{
				Name: "organization-default",
			},
			{
				Name: "project-default",
			},
			core.ZoneArgSpec(),
		},
		Examples: []*core.Example{
			{
				Short:    "Set the default outbound policy as drop",
				ArgsJSON: `{"security_group_id": "11111111-1111-1111-1111-111111111111", "outbound_default_policy": "drop"}`,
			},
			{
				Short:    "Set the given security group as the default for the project",
				ArgsJSON: `{"security_group_id": "11111111-1111-1111-1111-111111111111", "project_default": true}`,
			},
			{
				Short:    "Change the name of the given security group",
				ArgsJSON: `{"security_group_id": "11111111-1111-1111-1111-111111111111", "name": "foobar"}`,
			},
			{
				Short:    "Change the description of the given security group",
				ArgsJSON: `{"security_group_id": "11111111-1111-1111-1111-111111111111", "description": "foobar"}`,
			},
			{
				Short:    "Enable stateful security group",
				ArgsJSON: `{"security_group_id": "11111111-1111-1111-1111-111111111111", "stateful": true}`,
			},
			{
				Short:    "Set the default inbound policy as drop",
				ArgsJSON: `{"security_group_id": "11111111-1111-1111-1111-111111111111", "inbound_default_policy": "drop"}`,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			req := argsI.(*instance.UpdateSecurityGroupRequest)

			api := instance.NewAPI(core.ExtractClient(ctx))
			res, err := api.UpdateSecurityGroup(req)
			if err == nil {
				return res, nil
			}

			resErr, isResErr := err.(*scw.ResponseError)
			if !isResErr {
				return nil, err
			}

			// Try to find the error type and create a more user friendly one.
			switch resErr.Message {
			case "default security group can't be stateful":
				return nil, &core.CliError{
					Err: fmt.Errorf("your default security group cannot be stateful"),
					Details: interactive.RemoveIndent(`
						You have to make this security group stateless to use it as a project default.
						More info: https://www.scaleway.com/en/docs/how-to-activate-a-stateful-cloud-firewall
					`),
					Hint: "scw instance security-group update " + req.SecurityGroupID + " project-default=true stateful=false",
				}

			case "cannot have more than one project default", "cannot have more than one project default group":
				defaultSG, err := getDefaultProjectSecurityGroup(ctx, req.Zone)
				if err != nil {
					// Abort and return the original error.
					return nil, resErr
				}

				return nil, &core.CliError{
					Err: fmt.Errorf("you cannot have more than one project default"),
					Details: interactive.RemoveIndent(`
						You already have a project default security-group (` + defaultSG.ID + `).

						First, you need to set your current project default security-group as non-default with:
						scw instance security-group update ` + defaultSG.ID + ` project-default=false

						Then, retry this command:
						scw instance security-group update ` + req.SecurityGroupID + ` project-default=true stateful=false
					`),
				}
			default:
				// Unknown error, use default behavior.
				return nil, resErr
			}
		},
	}
}

func getDefaultProjectSecurityGroup(ctx context.Context, zone scw.Zone) (*instance.SecurityGroup, error) {
	api := instance.NewAPI(core.ExtractClient(ctx))

	projectID := core.GetProjectIDFromContext(ctx)
	sgList, err := api.ListSecurityGroups(&instance.ListSecurityGroupsRequest{
		Zone:    zone,
		Project: scw.StringPtr(projectID),
	}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}

	for _, sg := range sgList.SecurityGroups {
		if sg.ProjectDefault {
			return sg, nil
		}
	}

	return nil, fmt.Errorf("%s project does not have a default security group", projectID)
}
