package iam_test

import (
	"os"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/account/v3"
	iam "github.com/scaleway/scaleway-cli/v2/internal/namespaces/iam/v1alpha1"
	iamSdk "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
)

func Test_iamAPIKeyGet(t *testing.T) {
	if isNightly := os.Getenv("SLACK_WEBHOOK_NIGHTLY"); isNightly != "" {
		t.Skip()
	}

	commands := iam.GetCommands()
	commands.Merge(account.GetCommands())

	/*
		In case you need to record new golden files
		Be aware that theses tests purpose is to check the output of the `scw iam api-key get` command and display
		some information about the user, the API key and the policies attached to the user.

		The member (`GetMemberAPIKey`) test relies on an activated member user having:
		- an API key
		- a few policies attached to be displayed in the output

		This member cannot be created and cleaned up dynamically: IAM users can only be invited to an Organization
		(activation happens in the Scaleway console), the organization owner cannot create an API key for another
		member (no impersonation), and members cannot be deleted through the IAM API. The member is therefore a
		manual, one-time setup documented in docs/developer.md. When the member is absent, the subtest is skipped
		instead of recording an empty cassette.

		The application (`GetApplicationAPIKey`) test is fully automated: the application and its API key are
		created through the CLI before the command and deleted after it.
	*/

	t.Run("GetOwnerAPIKey", func(t *testing.T) {
		userResulter := newSliceResulter(t, "*iamSdk.User", func(users []*iamSdk.User) any {
			return users[0].ID
		})

		apiKeyResulter := newSliceResulter(t, "*iamSdk.APIKey", func(keys []*iamSdk.APIKey) any {
			return keys[0].AccessKey
		})

		core.Test(&core.TestConfig{
			Commands: commands,
			BeforeFunc: core.BeforeFuncCombine(
				core.ExecStoreBeforeCmdWithResulter(
					"owner",
					"scw iam user list type=owner",
					userResulter,
				),

				core.ExecStoreBeforeCmdWithResulter(
					"ownerAPIKey",
					"scw iam api-key list bearer-id={{ .owner }}",
					apiKeyResulter,
				),
			),
			Cmd: `scw iam api-key get {{ .ownerAPIKey }}`,
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		})(t)
	})

	t.Run("GetMemberAPIKey", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands: commands,
			BeforeFunc: core.BeforeFuncCombine(
				func(ctx *core.BeforeFuncCtx) error {
					result := core.ExecBeforeCmdWithResult(
						ctx,
						"scw iam user list type=member",
					)
					users, ok := result.([]*iamSdk.User)
					if !ok {
						ctx.T.Skipf(
							"unexpected result type %T for IAM member user list, see docs/developer.md IAM setup section",
							result,
						)
					}
					if len(users) == 0 {
						ctx.T.Skipf(
							"no IAM member user present in this organization, see docs/developer.md IAM setup section",
						)
					}

					ctx.Meta["member"] = users[0].ID

					return nil
				},
				func(ctx *core.BeforeFuncCtx) error {
					result := core.ExecBeforeCmdWithResult(
						ctx,
						"scw iam api-key list bearer-id={{ .member }}",
					)
					keys, ok := result.([]*iamSdk.APIKey)
					if !ok {
						ctx.T.Skipf(
							"unexpected result type %T for IAM member API key list, see docs/developer.md IAM setup section",
							result,
						)
					}
					if len(keys) == 0 {
						ctx.T.Skipf(
							"IAM member user %s has no API key, see docs/developer.md IAM setup section",
							ctx.Meta["member"],
						)
					}

					ctx.Meta["memberAPIKey"] = keys[0].AccessKey

					return nil
				},
			),
			Cmd: `scw iam api-key get {{ .memberAPIKey }}`,
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		})(t)
	})

	t.Run("GetApplicationAPIKey", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands: commands,
			BeforeFunc: core.BeforeFuncCombine(
				core.ExecStoreBeforeCmd(
					"application",
					"scw iam application create name=test-cli-iam-api-key{{ randint }}",
				),
				core.ExecStoreBeforeCmd(
					"applicationAPIKey",
					"scw iam api-key create application-id={{ .application.ID }} description=test-cli-iam-api-key",
				),
			),
			Cmd: `scw iam api-key get {{ .applicationAPIKey.AccessKey }}`,
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.AfterFuncWhenUpdatingCassette(
				core.AfterFuncCombine(
					core.ExecAfterCmd("scw iam api-key delete {{ .applicationAPIKey.AccessKey }}"),
					core.ExecAfterCmd("scw iam application delete {{ .application.ID }}"),
				),
			),
		})(t)
	})
}
