package instance_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v2alpha1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/stretchr/testify/assert"
)

func Test_SecurityGroup(t *testing.T) {
	t.Run("Create", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd: "scw instance security-group create name=cli-test-sg-create outbound-default-action=drop " +
			"stateless=true tags.0=cli-test tags.1=instance-sg tags.2=create description=\"Add desc here.\"",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				if res, ok := ctx.Result.(*instanceSDK.SecurityGroup); ok {
					assert.Equal(t, "cli-test-sg-create", res.Name)
					assert.Equal(t, "drop", res.OutboundDefaultAction.String())
					assert.Equal(t, "accept", res.InboundDefaultAction.String())
					assert.Equal(t, true, res.Stateless)
					assert.Equal(t, "cli-test", res.Tags[0])
					assert.Equal(t, "instance-sg", res.Tags[1])
					assert.Equal(t, "create", res.Tags[2])
					assert.Equal(t, "Add desc here.", res.Description)
				} else {
					t.Errorf(
						"expected result of type *instanceSDK.SecurityGroup, got %T",
						ctx.Result,
					)
				}
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.ExecAfterCmd("scw instance security-group delete {{ .CmdResult.ID }}"),
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
				if res, ok := ctx.Result.([]*instanceSDK.SecurityGroup); ok {
					assert.Len(t, res, 3)
				} else {
					t.Error(
						fmt.Sprintf(
							"expected result of type []*instanceSDK.SecurityGroup, got %T",
							ctx.Result,
						),
					)
				}
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
				if res, ok := ctx.Result.(*instanceSDK.SecurityGroup); ok {
					assert.Equal(t, "cli-test-sg-update", res.Name)
					assert.Equal(t, true, res.Stateless)
					assert.Equal(t, "accept", res.OutboundDefaultAction.String())
					assert.Equal(t, "drop", res.InboundDefaultAction.String())
					assert.Equal(t, false, res.DisableDefaultRules)
					assert.Len(t, res.DefaultRules, 6)
					assert.Equal(t, "cli-test", res.Tags[0])
					assert.Equal(t, "instance-sg", res.Tags[1])
					assert.Equal(t, "update", res.Tags[2])
					assert.Equal(t, "updated", res.Description)
				} else {
					t.Errorf(
						"expected result of type *instanceSDK.SecurityGroup, got %T",
						ctx.Result,
					)
				}
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.ExecAfterCmd("scw instance security-group delete {{ .SecurityGroup.ID }}"),
	}))
}
