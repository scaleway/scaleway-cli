package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v2alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func Test_SecurityGroup(t *testing.T) {
	t.Run("Create", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd: "scw instance security-group create name=cli-test-sg-create outbound-default-action=drop " +
			"stateless=true tags.0=cli-test tags.1=instance-sg tags.2=create description=\"Add desc here.\"",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				res := ctx.Result.(*instanceSDK.SecurityGroup)
				assert.Equal(t, "cli-test-sg-create", res.Name)
				assert.Equal(t, "drop", res.OutboundDefaultAction.String())
				assert.Equal(t, "accept", res.InboundDefaultAction.String())
				assert.True(t, res.Stateless)
				assert.Equal(t, "cli-test", res.Tags[0])
				assert.Equal(t, "instance-sg", res.Tags[1])
				assert.Equal(t, "create", res.Tags[2])
				assert.Equal(t, "Add desc here.", res.Description)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: instance.DeleteSecurityGroup("CmdResult"),
	}))

	t.Run("List", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("SG0", "scw instance security-group create tags.0=foo "),
			core.ExecStoreBeforeCmd("SG1", "scw instance security-group create"),
			core.ExecStoreBeforeCmd("SG2", "scw instance security-group create name=named-sg "+
				"inbound-default-action=drop outbound-default-action=drop stateless=true tags.0=foo"),
			core.ExecStoreBeforeCmd(
				"SG3",
				"scw instance security-group create tags.0=foo tags.1=bar "+
					"inbound-default-action=drop description=\"Descriptive text.\"",
			),
		),
		Cmd: "scw instance security-group list tags.0=foo",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				res := ctx.Result.([]*instanceSDK.SecurityGroupSummary)
				assert.Len(t, res, 3)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw instance security-group delete {{ .SG0.ID }}"),
			core.ExecAfterCmd("scw instance security-group delete {{ .SG1.ID }}"),
			core.ExecAfterCmd("scw instance security-group delete {{ .SG2.ID }}"),
			core.ExecAfterCmd("scw instance security-group delete {{ .SG3.ID }}"),
		),
	}))

	t.Run("Update", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("SecurityGroup", "scw instance security-group create "+
			"outbound-default-action=drop tags.0=cli-test tags.1=instance-sg tags.2=update description=tmp disable-default-rules=true"),
		Cmd: "scw instance security-group update {{ .SecurityGroup.ID }} name=cli-test-sg-update stateless=true " +
			"outbound-default-action=accept inbound-default-action=drop disable-default-rules=false " +
			"tags.0=cli-test tags.1=instance-sg tags.2=update description=updated",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				res := ctx.Result.(*instanceSDK.SecurityGroup)
				assert.Equal(t, "cli-test-sg-update", res.Name)
				assert.True(t, res.Stateless)
				assert.Equal(t, "accept", res.OutboundDefaultAction.String())
				assert.Equal(t, "drop", res.InboundDefaultAction.String())
				assert.False(t, res.DisableDefaultRules)
				assert.Len(t, res.DefaultRules, 6)
				assert.Equal(t, "cli-test", res.Tags[0])
				assert.Equal(t, "instance-sg", res.Tags[1])
				assert.Equal(t, "update", res.Tags[2])
				assert.Equal(t, "updated", res.Description)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: instance.DeleteSecurityGroup("SecurityGroup"),
	}))

	t.Run("AddRules", core.Test(&core.TestConfig{
		Commands:   instance.GetCommands(),
		BeforeFunc: instance.CreateSecurityGroup("SecurityGroup"),
		Cmd: "scw instance security-group add-rules {{ .SecurityGroup.ID }} " +
			"rules.0.protocol=any rules.0.direction=inbound rules.0.action=accept rules.0.source-ip-range=1.1.1.1/32 " +
			"rules.1.protocol=udp rules.1.direction=outbound rules.1.action=drop rules.1.source-ip-range=2.2.2.2/32",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				sg := ctx.Result.(*instanceSDK.SecurityGroup)
				assert.Len(t, sg.Rules, 2)
				assert.Equal(t, "any", sg.Rules[0].Protocol.String())
				assert.Equal(t, "udp", sg.Rules[1].Protocol.String())
				assert.Equal(t, "inbound", sg.Rules[0].Direction.String())
				assert.Equal(t, "outbound", sg.Rules[1].Direction.String())
				assert.Equal(t, "accept", sg.Rules[0].Action.String())
				assert.Equal(t, "drop", sg.Rules[1].Action.String())
				assert.Equal(t, "1.1.1.1/32", sg.Rules[0].SourceIPRange.String())
				assert.Equal(t, "2.2.2.2/32", sg.Rules[1].SourceIPRange.String())
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: instance.DeleteSecurityGroup("SecurityGroup"),
	}))

	t.Run("SetRules", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"SecurityGroup",
			"scw instance security-group create zone=fr-par-2",
		),
		Cmd: "scw instance security-group set-rules {{ .SecurityGroup.ID }} zone=fr-par-2 " +
			"rules.0.protocol=tcp rules.0.direction=inbound rules.0.action=drop rules.0.source-ip-range=1.1.1.1/32 " +
			"rules.1.protocol=udp rules.1.direction=outbound rules.1.action=drop rules.1.source-ip-range=2.2.2.2/32 " +
			"rules.2.protocol=icmp rules.2.direction=inbound rules.2.action=accept rules.2.source-ip-range=3.3.3.3/32",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				sg := ctx.Result.(*instanceSDK.SecurityGroup)
				assert.Len(t, sg.Rules, 3)
				assert.Equal(t, "tcp", sg.Rules[0].Protocol.String())
				assert.Equal(t, "udp", sg.Rules[1].Protocol.String())
				assert.Equal(t, "icmp", sg.Rules[2].Protocol.String())
				assert.Equal(t, "inbound", sg.Rules[0].Direction.String())
				assert.Equal(t, "outbound", sg.Rules[1].Direction.String())
				assert.Equal(t, "inbound", sg.Rules[2].Direction.String())
				assert.Equal(t, "drop", sg.Rules[0].Action.String())
				assert.Equal(t, "drop", sg.Rules[1].Action.String())
				assert.Equal(t, "accept", sg.Rules[2].Action.String())
				assert.Equal(t, "1.1.1.1/32", sg.Rules[0].SourceIPRange.String())
				assert.Equal(t, "2.2.2.2/32", sg.Rules[1].SourceIPRange.String())
				assert.Equal(t, "3.3.3.3/32", sg.Rules[2].SourceIPRange.String())
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.ExecAfterCmd(
			"scw instance security-group delete {{ .SecurityGroup.ID }} zone=fr-par-2",
		),
	}))

	t.Run("UpdateRule", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			instance.CreateSecurityGroup("SecurityGroup"),
			core.ExecStoreBeforeCmd(
				"SecurityGroupWithRules",
				"scw instance security-group add-rules {{ .SecurityGroup.ID }} "+
					"rules.0.protocol=udp rules.0.direction=inbound rules.0.action=drop rules.0.source-ip-range=1.1.1.1/32",
			),
		),
		Cmd: "scw instance security-group update-rule security-group-rule-id={{ (index .SecurityGroupWithRules.Rules 0).ID }} " +
			"protocol=icmp direction=outbound action=drop source-ip-range=2.2.2.2/32",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				sg := ctx.Result.(*instanceSDK.SecurityGroup)
				assert.Len(t, sg.Rules, 1)
				assert.Equal(t, "icmp", sg.Rules[0].Protocol.String())
				assert.Equal(t, "outbound", sg.Rules[0].Direction.String())
				assert.Equal(t, "drop", sg.Rules[0].Action.String())
				assert.Equal(t, "2.2.2.2/32", sg.Rules[0].SourceIPRange.String())
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: instance.DeleteSecurityGroup("SecurityGroup"),
	}))

	t.Run("DeleteRules", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			instance.CreateSecurityGroup("SecurityGroup"),
			core.ExecStoreBeforeCmd(
				"SecurityGroupWithRules",
				"scw instance security-group set-rules {{ .SecurityGroup.ID }} "+
					"rules.0.protocol=tcp rules.0.direction=inbound rules.0.action=drop rules.0.source-ip-range=1.1.1.1/32 "+
					"rules.1.protocol=udp rules.1.direction=outbound rules.1.action=drop rules.1.source-ip-range=2.2.2.2/32 "+
					"rules.2.protocol=icmp rules.2.direction=inbound rules.2.action=accept rules.2.source-ip-range=3.3.3.3/32 "+
					"rules.3.protocol=any rules.3.direction=outbound rules.3.action=accept rules.3.source-ip-range=4.4.4.4/32",
			),
		),
		Cmd: "scw instance security-group delete-rules security-group-rule-ids.0={{ (index .SecurityGroupWithRules.Rules 0).ID }} " +
			"security-group-rule-ids.1={{ (index .SecurityGroupWithRules.Rules 1).ID }} " +
			"security-group-rule-ids.2={{ (index .SecurityGroupWithRules.Rules 3).ID }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				sg, err := refreshSecurityGroupInfo(t, ctx, "SecurityGroup")
				if err != nil {
					t.Errorf("could not get security group: %v", err)
				}

				assert.Len(t, sg.Rules, 1)
				assert.Equal(t, "icmp", sg.Rules[0].Protocol.String())
				assert.Equal(t, "inbound", sg.Rules[0].Direction.String())
				assert.Equal(t, "accept", sg.Rules[0].Action.String())
				assert.Equal(t, "3.3.3.3/32", sg.Rules[0].SourceIPRange.String())
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: instance.DeleteSecurityGroup("SecurityGroup"),
	}))

	t.Run("Clear", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			instance.CreateSecurityGroup("SecurityGroup"),
			core.ExecBeforeCmd(
				"scw instance security-group set-rules {{ .SecurityGroup.ID }} "+
					"rules.0.protocol=tcp rules.0.direction=inbound rules.0.action=drop rules.0.source-ip-range=1.1.1.1/32 "+
					"rules.1.protocol=udp rules.1.direction=outbound rules.1.action=drop rules.1.source-ip-range=2.2.2.2/32 "+
					"rules.2.protocol=icmp rules.2.direction=inbound rules.2.action=accept rules.2.source-ip-range=3.3.3.3/32 "+
					"rules.3.protocol=any rules.3.direction=outbound rules.3.action=accept rules.3.source-ip-range=4.4.4.4/32",
			),
		),
		Cmd: "scw instance security-group clear {{ .SecurityGroup.ID }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				sg, err := refreshSecurityGroupInfo(t, ctx, "SecurityGroup")
				if err != nil {
					t.Errorf("could not get security group: %v", err)
				}

				assert.Empty(t, sg.Rules)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: instance.DeleteSecurityGroup("SecurityGroup"),
	}))

	t.Run("Edit", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			instance.CreateSecurityGroup("SecurityGroup"),
			core.ExecBeforeCmd(
				"scw instance security-group set-rules {{ .SecurityGroup.ID }} "+
					"rules.0.protocol=tcp rules.0.direction=inbound rules.0.action=drop rules.0.source-ip-range=1.1.1.1/32 "+
					"rules.1.protocol=icmp rules.1.direction=outbound rules.1.action=drop rules.1.source-ip-range=2.2.2.2/32",
			),
			func(ctx *core.BeforeFuncCtx) error {
				editor.SkipEditor = true
				editor.EditorTestContent = `security_group_rules:
- action: accept
  direction: inbound
  protocol: tcp
  source_ip_range: 1.1.1.1/32
  destination_ports:
    start: 1
    end: 81
- action: drop
  direction: inbound
  protocol: udp
  source_ip_range: 3.3.3.3/32
  source_ports:
    start: 8080
- action: drop
  direction: outbound
  protocol: icmp
  source_ip_range: 4.4.4.4/32
`

				return nil
			},
		),
		Cmd: "scw instance security-group edit {{ .SecurityGroup.ID }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				sgRules := ctx.Result.([]*instanceSDK.SecurityGroupRule)
				assert.Len(t, sgRules, 3)

				assert.Equal(t, "accept", sgRules[0].Action.String())
				assert.Equal(t, "inbound", sgRules[0].Direction.String())
				assert.Equal(t, "tcp", sgRules[0].Protocol.String())
				assert.Equal(t, "1.1.1.1/32", sgRules[0].SourceIPRange.String())
				assert.Equal(t, uint32(1), sgRules[0].DestinationPorts.Start)
				assert.Equal(t, uint32(81), sgRules[0].DestinationPorts.End)

				assert.Equal(t, "drop", sgRules[1].Action.String())
				assert.Equal(t, "inbound", sgRules[1].Direction.String())
				assert.Equal(t, "udp", sgRules[1].Protocol.String())
				assert.Equal(t, "3.3.3.3/32", sgRules[1].SourceIPRange.String())
				assert.Equal(t, uint32(8080), sgRules[1].SourcePorts.Start)
				assert.Equal(t, uint32(8080), sgRules[1].SourcePorts.End)

				assert.Equal(t, "drop", sgRules[2].Action.String())
				assert.Equal(t, "outbound", sgRules[2].Direction.String())
				assert.Equal(t, "icmp", sgRules[2].Protocol.String())
				assert.Equal(t, "4.4.4.4/32", sgRules[2].SourceIPRange.String())
				assert.Equal(
					t,
					uint32(instance.SecurityGroupRulePortMin),
					sgRules[2].SourcePorts.Start,
				)
				assert.Equal(
					t,
					uint32(instance.SecurityGroupRulePortMax),
					sgRules[2].SourcePorts.End,
				)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			instance.DeleteSecurityGroup("SecurityGroup"),
			func(ctx *core.AfterFuncCtx) error {
				editor.SkipEditor = false
				editor.EditorTestContent = ""

				return nil
			},
		),
	}))
}

func refreshSecurityGroupInfo(
	t *testing.T,
	ctx *core.CheckFuncCtx,
	metaKey string,
) (*instanceSDK.SecurityGroup, error) {
	t.Helper()
	api := instanceSDK.NewAPI(ctx.Client)
	sgMeta := testhelpers.MapValue[*instanceSDK.SecurityGroup](
		t,
		ctx.Meta,
		metaKey,
	)

	return api.GetSecurityGroup(&instanceSDK.GetSecurityGroupRequest{
		Zone:            "fr-par-1",
		SecurityGroupID: sgMeta.ID,
	}, scw.WithContext(t.Context()))
}
