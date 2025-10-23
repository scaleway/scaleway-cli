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

		They rely on the fact that a user member exist and has:
		- an API key
		- a few policies attached to be displayed in the output

		The user member and API Key creation cannot be automated with BeforeFunc because it is impossible to
		create a member API Key from the organization owner (no impersonation).
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
					"member",
					"scw iam user list type=member",
					userResulter,
				),
				core.ExecStoreBeforeCmdWithResulter(
					"memberAPIKey",
					"scw iam api-key list bearer-id={{ .member }}",
					apiKeyResulter,
				),
			),
			Cmd: `scw iam api-key get {{ .memberAPIKey }}`,
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		})(t)
	})

	t.Run("GetApplicationAPIKey", func(t *testing.T) {
		appResulter := newSliceResulter(
			t,
			"*iamSdk.Application",
			func(apps []*iamSdk.Application) any {
				return apps[0].ID
			},
		)

		apiKeyResulter := newSliceResulter(t, "*iamSdk.APIKey", func(keys []*iamSdk.APIKey) any {
			return keys[0].AccessKey
		})

		core.Test(&core.TestConfig{
			Commands: commands,
			BeforeFunc: core.BeforeFuncCombine(
				core.ExecStoreBeforeCmdWithResulter(
					"application",
					"scw iam application list",
					appResulter,
				),
				core.ExecStoreBeforeCmdWithResulter(
					"applicationAPIKey",
					"scw iam api-key list bearer-id={{ .application }}",
					apiKeyResulter,
				),
			),
			Cmd: `scw iam api-key get {{ .applicationAPIKey }}`,
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		})(t)
	})
}
