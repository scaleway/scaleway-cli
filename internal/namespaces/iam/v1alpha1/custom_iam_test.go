package iam_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/account/v3"
	iam "github.com/scaleway/scaleway-cli/v2/internal/namespaces/iam/v1alpha1"
	iamSdk "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
)

func Test_iamAPIKeyGet(t *testing.T) {
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

	userResulter := func(result any) any {
		users := result.([]*iamSdk.User)
		if users == nil {
			panic("users is nil")
		}
		if len(users) == 0 {
			panic("no user found")
		}

		return users[0].ID
	}

	apiKeyResulter := func(result any) any {
		apiKeys := result.([]*iamSdk.APIKey)
		if apiKeys == nil {
			panic("apiKeys is nil")
		}
		if len(apiKeys) == 0 {
			panic("no api key found")
		}

		return apiKeys[0].AccessKey
	}

	t.Run("GetOwnerAPIKey", core.Test(&core.TestConfig{
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
	}))

	t.Run("GetMemberAPIKey", core.Test(&core.TestConfig{
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
	}))
}
