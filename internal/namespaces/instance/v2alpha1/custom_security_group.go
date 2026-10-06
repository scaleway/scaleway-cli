package instance

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	securityGroupRulePortMin = 1
	securityGroupRulePortMax = 65536
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
)

func securityGroupMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmpSecurityGroup instance.SecurityGroup
	securityGroup := tmpSecurityGroup(i.(instance.SecurityGroup))

	// Sections
	opt.Sections = []*human.MarshalSection{
		{
			FieldName:   "Tags",
			Title:       "Tags",
			HideIfEmpty: true,
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

func securityGroupRulesMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmpSecurityGroupRules []*instance.SecurityGroupRule
	securityGroupRules := tmpSecurityGroupRules(i.([]*instance.SecurityGroupRule))

	opt.Fields = []*human.MarshalFieldOpt{
		{
			Label:     "ID",
			FieldName: "ID",
		},
		{
			Label:     "PROTOCOL",
			FieldName: "Protocol",
		},
		{
			Label:     "DIRECTION",
			FieldName: "Direction",
		},
		{
			Label:     "ACTION",
			FieldName: "Action",
		},
		{
			Label:     "SOURCE IP RANGE",
			FieldName: "SourceIPRange",
		},
		{
			Label:     "SOURCE PORTS",
			FieldName: "SourcePorts",
		},
		{
			Label:     "DEST IP RANGE",
			FieldName: "DestinationIPRange",
		},
		{
			Label:     "DEST PORTS",
			FieldName: "DestinationPorts",
		},
	}

	str, err := human.Marshal(securityGroupRules, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func securityGroupRulePortRangeMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmpSecurityGroupRulePortRange instance.SecurityGroupRulePortRange
	securityGroupRulePortRange := tmpSecurityGroupRulePortRange(
		i.(instance.SecurityGroupRulePortRange),
	)

	humanReadablePortRange := fmt.Sprintf(
		"%d-%d",
		securityGroupRulePortRange.Start,
		securityGroupRulePortRange.End,
	)

	if securityGroupRulePortRange.Start == securityGroupRulePortMin &&
		securityGroupRulePortRange.End == securityGroupRulePortMax {
		humanReadablePortRange = "ALL"
	}

	if securityGroupRulePortRange.Start == securityGroupRulePortRange.End {
		humanReadablePortRange = strconv.Itoa(int(securityGroupRulePortRange.Start))
	}

	str, err := human.Marshal(humanReadablePortRange, opt)
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
	c.ArgSpecs.GetByName("inbound-default-action").Default = core.DefaultValueSetter(
		instance.SecurityGroupActionAccept.String(),
	)
	c.ArgSpecs.GetByName("outbound-default-action").Default = core.DefaultValueSetter(
		instance.SecurityGroupActionAccept.String(),
	)
	c.ArgSpecs.GetByName("stateless").Default = core.DefaultValueSetter("false")
	c.ArgSpecs.GetByName("disable-default-rules").Default = core.DefaultValueSetter("false")

	examples := []*core.Example{
		{
			Short: "Create a security group with a specified name and description",
			Raw:   `scw instance security-group create name=foobar description="foobar foobar"`,
		},
		{
			Short: "Create a security group that will be applied as default on all Instances of this Project",
			Raw:   `scw instance security-group create project-default=true`,
		},
		{
			Short: "Create a security group that will have a default drop inbound policy (traffic your Instance receives)",
			Raw:   "scw instance security-group create inbound-default-action=drop",
		},
		{
			Short: "Create a security group that will have a default drop outbound policy (traffic your Instance transmits)",
			Raw:   "scw instance security-group create outbound-default-action=drop",
		},
		{
			Short: "Create a stateless security group",
			Raw:   "scw instance security-group create stateless=true",
		},
	}

	c.Examples = append(c.Examples, examples...)

	return c
}

func securityGroupGetBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("security-group-id").Positional = true

	c.Examples = append(c.Examples, &core.Example{
		Short: "Get a security group with the specified ID",
		Raw:   "scw instance security-group get 11111111-1111-1111-1111-111111111111",
	})

	return c
}

func securityGroupUpdateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("security-group-id").Positional = true

	return c
}

func securityGroupDeleteBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("security-group-id").Positional = true

	c.Examples = append(c.Examples, &core.Example{
		Short: "Delete the security group with the specified ID",
		Raw:   "scw instance security-group delete 11111111-1111-1111-1111-111111111111",
	})

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

	c.Examples = append(c.Examples, &core.Example{
		Short: "List all security groups that match the specified name",
		Raw:   "scw instance security-group list name=foobar",
	})

	return c
}

func securityGroupAddRulesBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("security-group-id").Positional = true
	c.ArgSpecs.GetByName("security-group-id").Required = true
	c.ArgSpecs.GetByName("security-group-rules.{index}.protocol").Name = "rules.{index}.protocol"
	c.ArgSpecs.GetByName("security-group-rules.{index}.direction").Name = "rules.{index}.direction"
	c.ArgSpecs.GetByName("security-group-rules.{index}.action").Name = "rules.{index}.action"
	c.ArgSpecs.GetByName("security-group-rules.{index}.source-ip-range").Name = "rules.{index}.source-ip-range"
	c.ArgSpecs.GetByName("security-group-rules.{index}.destination-ip-range").Name = "rules.{index}.destination-ip-range"
	c.ArgSpecs.GetByName("security-group-rules.{index}.source-ports.start").Name = "rules.{index}.source-ports.start"
	c.ArgSpecs.GetByName("security-group-rules.{index}.source-ports.end").Name = "rules.{index}.source-ports.end"
	c.ArgSpecs.GetByName("security-group-rules.{index}.destination-ports.start").Name = "rules.{index}.destination-ports.start"
	c.ArgSpecs.GetByName("security-group-rules.{index}.destination-ports.end").Name = "rules.{index}.destination-ports.end"
	c.ArgSpecs.GetByName("security-group-rules.{index}.position").Name = "rules.{index}.position"

	type customSecurityGroupAddRulesRequest struct {
		instance.AddSecurityGroupRulesRequest
		Rules []*instance.SecurityGroupRuleConfig
	}

	c.ArgsType = reflect.TypeFor[customSecurityGroupAddRulesRequest]()

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		args := argsI.(*customSecurityGroupAddRulesRequest)
		req := args.AddSecurityGroupRulesRequest
		req.SecurityGroupRules = args.Rules

		resp, err := runner(ctx, &req)
		if err != nil {
			return nil, err
		}

		return resp.(*instance.AddSecurityGroupRulesResponse).SecurityGroup, nil
	}

	return c
}

func securityGroupSetRulesBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("security-group-id").Positional = true
	c.ArgSpecs.GetByName("security-group-id").Required = true
	c.ArgSpecs.GetByName("security-group-rules.{index}.protocol").Name = "rules.{index}.protocol"
	c.ArgSpecs.GetByName("security-group-rules.{index}.direction").Name = "rules.{index}.direction"
	c.ArgSpecs.GetByName("security-group-rules.{index}.action").Name = "rules.{index}.action"
	c.ArgSpecs.GetByName("security-group-rules.{index}.source-ip-range").Name = "rules.{index}.source-ip-range"
	c.ArgSpecs.GetByName("security-group-rules.{index}.destination-ip-range").Name = "rules.{index}.destination-ip-range"
	c.ArgSpecs.GetByName("security-group-rules.{index}.source-ports.start").Name = "rules.{index}.source-ports.start"
	c.ArgSpecs.GetByName("security-group-rules.{index}.source-ports.end").Name = "rules.{index}.source-ports.end"
	c.ArgSpecs.GetByName("security-group-rules.{index}.destination-ports.start").Name = "rules.{index}.destination-ports.start"
	c.ArgSpecs.GetByName("security-group-rules.{index}.destination-ports.end").Name = "rules.{index}.destination-ports.end"
	c.ArgSpecs.GetByName("security-group-rules.{index}.position").Name = "rules.{index}.position"

	type customSecurityGroupSetRulesRequest struct {
		instance.SetSecurityGroupRulesRequest
		Rules []*instance.SecurityGroupRuleConfig
	}

	c.ArgsType = reflect.TypeFor[customSecurityGroupSetRulesRequest]()

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		args := argsI.(*customSecurityGroupSetRulesRequest)
		req := args.SetSecurityGroupRulesRequest
		req.SecurityGroupRules = args.Rules

		return runner(ctx, &req)
	}

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
		ArgsType:  reflect.TypeFor[instanceResetSecurityGroupArgs](),
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

			_, err := api.SetSecurityGroupRules(&instance.SetSecurityGroupRulesRequest{
				Zone:               args.Zone,
				SecurityGroupID:    args.SecurityGroupID,
				SecurityGroupRules: []*instance.SecurityGroupRuleConfig{},
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{Message: "Successful reset of the security group rules"}, err
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "security-group-id",
				Short:      `ID of the security group to reset.`,
				Required:   true,
				Positional: true,
			},
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
	}
}

var instanceSecurityGroupEditYamlExample = `security_group_rules:
- action: drop
  destination_ports:
    start: 1200
    end: 1300
  direction: inbound
  source_ip_range: 192.168.0.0/24
  protocol: TCP
- action: drop
  direction: inbound
  protocol: ICMP
  source_ip_range: 0.0.0.0/0
- action: accept
  destination_ports:
    start: 25565
  direction: outbound
  source_ip_range: 0.0.0.0/0
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
		ArgsType:  reflect.TypeFor[instanceSecurityGroupEditArgs](),
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

			sg, err := api.GetSecurityGroup(&instance.GetSecurityGroupRequest{
				Zone:            args.Zone,
				SecurityGroupID: args.SecurityGroupID,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to get security group: %w", err)
			}

			// Build resource with matching field names for editor
			type securityGroupRulesResource struct {
				SecurityGroupRules []*instance.SecurityGroupRule
			}
			resource := &securityGroupRulesResource{
				SecurityGroupRules: sg.Rules,
			}

			setRequest := &instance.SetSecurityGroupRulesRequest{
				Zone:            args.Zone,
				SecurityGroupID: args.SecurityGroupID,
			}

			editedSetRequest, err := editor.UpdateResourceEditor(
				resource,
				setRequest,
				&editor.Config{
					PutRequest:   true,
					MarshalMode:  args.Mode,
					Template:     instanceSecurityGroupEditYamlExample,
					IgnoreFields: []string{"security_group_id", "zone", "id", "position"},
				},
			)
			if err != nil {
				return nil, err
			}

			setRequest = editedSetRequest.(*instance.SetSecurityGroupRulesRequest)

			// Preserve path parameters that were filtered out by IgnoreFields
			setRequest.Zone = args.Zone
			setRequest.SecurityGroupID = args.SecurityGroupID

			for _, rule := range setRequest.SecurityGroupRules {
				rule.Protocol = instance.SecurityGroupRuleProtocol(
					strings.ToLower(rule.Protocol.String()),
				)

				removePortRangeDefaultValues(rule.SourcePorts)
				removePortRangeDefaultValues(rule.DestinationPorts)
			}

			resp, err := api.SetSecurityGroupRules(setRequest, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			return resp.Rules, nil
		},
	}
}

func removePortRangeDefaultValues(ports *instance.SecurityGroupRulePortRange) {
	if ports == nil || ports.End != securityGroupRulePortMax {
		return
	}

	if ports.Start == securityGroupRulePortMin {
		*ports = instance.SecurityGroupRulePortRange{}
	} else {
		ports.End = ports.Start
	}
}
